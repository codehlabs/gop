package gop

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"time"
)

type User struct {
	Id        string    `bson:"_id"`
	Username  string    `bson:"username"`
	FirsName  string    `bson:"firstname"`
	LastName  string    `bson:"lastname"`
	Password  string    `bson:"password"`
	Salt      string    `bson:"salt"`
	Email     string    `bson:"email"`
	Phone     string    `bson:"phone"`
	Age       int32     `bson:"age"`
	DOB       time.Time `bson:"dob"` // Date of Birth
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
	Address  string `bson:"address"`
	Address2 string `bson:"address2,omitempty"`
	City     string `bson:"city"`
	State    string `bson:"state,omitempty"`
	ZipCode  string `bson:"zipcode"`
}

func (a Address) String() string {
	return fmt.Sprintf("%s %s %s, %s, %s", a.Address, a.Address2, a.City, a.State, a.ZipCode)
}

func (u User) Save(db Db) error {
	u.Id = u.Email
	u.CreatedAt = time.Now()
	return db.Save(u)
}

func (u User) Delete(db Db) error {
	return db.Delete(u.Id)
}

func (u User) Update(db Db) error {
	return db.Update(u)
}

type Db interface {
	Save(u User) error
	Update(u User) error
	Delete(id string) error
	Read(id string) (User, error)
}
