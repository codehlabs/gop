package gop

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"time"
)

type User struct {
	Id        string    `bson:"_id" sql:"id,text,unique"`
	Username  string    `bson:"username" form:"username" sql:"username,text"`
	FirsName  string    `bson:"firstname" form:"firstname" sql:"firstname,text"`
	LastName  string    `bson:"lastname" form:"lastname" sql:"lastname,text"`
	Password  string    `bson:"password" form:"password" sql:"password,text"`
	Salt      string    `bson:"salt" sql:"salt,text"`
	Email     string    `bson:"email" form:"email" sql:"email,text"`
	Phone     string    `bson:"phone,omitempty" form:"phone" sql:"phone,text"`
	Age       int32     `bson:"age" sql:"age,integer"`
	DOB       time.Time `bson:"dob" form:"dob" sql:"dob,integer"` // Date of Birth
	Address   Address   `bson:"address,omitempty"`
	Profile   any       `bson:"profile" sql:"omit"` // Domain defined profile
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
