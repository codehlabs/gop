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

type Config struct {
	UseBuiltInSaveLogic bool // when calling [User] to save it will use the custom built in logic
	UniqueIDLength      int  // defines built int unique id seed length
}

var config = &Config{UseBuiltInSaveLogic: true}

// Set Custom configuration
func SetConfig(c Config) {
	config = &c
}

type User struct {
	Id        string    `bson:"_id" sql:"id,text,unique"`
	Username  string    `bson:"username" form:"username" sql:"username,text"`
	FirsName  string    `bson:"firstname" form:"firstname" sql:"firstname,text"`
	LastName  string    `bson:"lastname" form:"lastname" sql:"lastname,text"`
	Password  string    `bson:"password" form:"password" sql:"password,text"`
	Salt      string    `bson:"salt" sql:"salt,text"`
	Email     string    `bson:"email" form:"email" sql:"email,text"`
	Phone     string    `bson:"phone" form:"phone" sql:"phone,text"`
	Age       int32     `bson:"age" sql:"age,integer"`
	DOB       time.Time `bson:"dob" form:"dob" sql:"dob,integer"` // Date of Birth
	Address   Address   `bson:"address,omitempty"`
	Profile   any       `bson:"profile" sql:"omit"` // Platform related profile
	CreatedAt time.Time `bson:"created_at" sql:"created_at,integer"`
	UpdatedAt time.Time `bson:"updated_at,omitempty" sql:"updated_at,integer"`
	DeletedAt time.Time `bson:"deleted_at,omitempty" sql:"deleted_at,integer"`
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
	Address  string `bson:"address" form:"address" sql:"address,text"`
	Address2 string `bson:"address2,omitempty" form:"address2" sql:"address2,text"`
	City     string `bson:"city" form:"city" sql:"city,text"`
	State    string `bson:"state,omitempty" form:"state" sql:"state,text"`
	ZipCode  string `bson:"zipcode" form:"zipcode" sql:"zipcode,text"`
}

func (a Address) String() string {
	return fmt.Sprintf("%s %s %s, %s, %s", a.Address, a.Address2, a.City, a.State, a.ZipCode)
}

func (u User) Save(db Db) error {

	if config.UseBuiltInSaveLogic {
		if err := u.HashAndSalt(); err != nil {
			return err
		}
		u.CreatedAt = time.Now()
		age := time.Now().Sub(u.DOB)
		u.Age = int32(age.Hours() / 24 / 365)
		seed_length := config.UniqueIDLength
		if seed_length == 0 {
			seed_length = 16
		}
		uid, err := unique_id(seed_length)
		if err != nil {
			return err
		}
		u.Id = uid
	}

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

type Db interface {
	Save(u User) error
	Update(u User) error
	Delete(id string) error
	Read(id string) (User, error)
}
