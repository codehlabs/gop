package gop

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Bind form data to data
func bind_form_data(r *http.Request, data interface{}) error {

	v := reflect.ValueOf(data).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i += 1 {
		field := t.Field(i)
		//TODO: implement split tag and find omit empty fields
		formKey := field.Tag.Get("form")
		if formKey == "" {
			continue
		}

		formValue := r.Form.Get(formKey)
		if formValue == "" {
			continue
		}

		structField := v.Field(i)
		if structField.CanSet() {
			switch structField.Kind() {
			case reflect.String:
				structField.SetString(formValue)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				intval, err := strconv.Atoi(formValue)
				if err == nil {
					structField.SetInt(int64(intval))
				}
			case reflect.Float32, reflect.Float64:
				floatval, err := strconv.ParseFloat(formValue, 64)
				if err == nil {
					structField.SetFloat(floatval)
				}
			case reflect.Bool:
				boolval, err := strconv.ParseBool(formValue)
				if err == nil {
					structField.SetBool(boolval)
				}
			case reflect.Struct:
				if structField.Type() == reflect.TypeOf(time.Time{}) {
					parsedTime, err := time.Parse("2006-01-02", formValue)
					if err == nil {
						structField.Set(reflect.ValueOf(parsedTime))
					}
				} else {
					nestedStruct := reflect.New(structField.Type()).Elem()
					if err := bind_form_data(r, &nestedStruct); err == nil {
						structField.Set(nestedStruct)
					}
				}

			}
		}
	}

	return nil
}

// Creates an user struct from the form data
func UserFromForm(r *http.Request) (User, error) {
	u := User{}
	err := r.ParseForm()
	if err != nil {
		return u, err
	}

	err = bind_form_data(r, &u)
	if err != nil {
		return u, err
	}

	return u, nil
}

// Checks if the password is a secure password. If checkifpwned is true then it will check against
// hundred of millions of passwoed previously exposed in data breaches. If checkifbad is true then it will check
// against 1 million commonly used bad passwords
func SecurePassword(password string, checkifpwned, checkifbad bool) error {
	if len(password) < 20 {
		return ErrShortPassword
	}

	if checkifpwned {
		s := sha1.New()
		_, err := s.Write([]byte(password))
		if err != nil {
			return err
		}
		hashed := strings.ToUpper(fmt.Sprintf("%x", s.Sum(nil)))
		res, err := http.Get("https://api.pwnedpasswords.com/range/" + string(hashed[:5]))
		if err != nil {
			return err
		}
		defer res.Body.Close()
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		if idx := bytes.Index(data, []byte(hashed[5:])); idx != -1 {
			data = data[idx:]
			newline := bytes.Index(data, []byte("\n"))
			colon := bytes.Index(data, []byte(":"))
			data = data[colon+1 : newline]
			data = bytes.Replace(data, []byte("\r"), []byte(""), 1)
			breaches, err := strconv.Atoi(string(data))
			if err == nil {
				if breaches > 500 {
					ErrPwnedPassword.Append("highly insecure password it has been breached " + string(data))
				}
			}
			return ErrPwnedPassword
		}
	}

	if checkifbad {

	}

	return nil
}
