package gop

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

type User struct {
	Id        string    `bson:"_id"`
	Username  string    `bson:"username" form:"username"`
	FirsName  string    `bson:"firstname" form:"firstname"`
	LastName  string    `bson:"lastname" form:"lastname"`
	Password  string    `bson:"password" form:"password"`
	Salt      string    `bson:"salt"`
	Email     string    `bson:"email" form:"email"`
	Phone     string    `bson:"phone" form:"phone"`
	Age       int32     `bson:"age"`
	DOB       time.Time `bson:"dob" form:"dob"` // Date of Birth
	Address   Address   `bson:"address"`
	Profile   any       `bson:"profile"` // Platform related profile
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
	DeletedAt time.Time `bson:"deleted_at"`
}

func (u *User) HashAndSalt() (err error) {
	bytes := make([]byte, 32)
	read_bytes, err := rand.Read(bytes)
	if err != nil {
		return
	}
	salt := fmt.Sprintf("%x", bytes)
	bytes = bytes[:read_bytes]
	sha := sha256.New()
	bytes = append(bytes, []byte(u.Password)...)
	sha.Write(bytes)
	hash := fmt.Sprintf("%x", sha.Sum(nil))
	u.Password = hash
	u.Salt = salt
	return nil
}

type Address struct {
	Address  string `bson:"address" form:"address"`
	Address2 string `bson:"address2,omitempty" form:"address2"`
	City     string `bson:"city" form:"city"`
	State    string `bson:"state,omitempty" form:"state"`
	ZipCode  string `bson:"zipcode" form:"zipcode"`
}

func (a Address) String() string {
	return fmt.Sprintf("%s %s %s, %s, %s", a.Address, a.Address2, a.City, a.State, a.ZipCode)
}

func (u User) Save(db Db) error {

	if err := u.HashAndSalt(); err != nil {
		return err
	}

	u.CreatedAt = time.Now()

	age := u.DOB.Sub(time.Now())
	u.Age = int32(age.Hours() / 24 / 365)

	return db.Save(u)
}

func (u User) Delete(db Db) error {
	return db.Delete(u.Id)
}

func (u User) Update(db Db) error {
	return db.Update(u)
}

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
					parsedTime, err := time.Parse(time.RFC3339, formValue)
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

type Db interface {
	Save(u User) error
	Update(u User) error
	Delete(id string) error
	Read(id string) (User, error)
}
