package gop

import (
	"fmt"
	"os"
	"testing"
)

func TestMDBDriver(t *testing.T) {
	mdbd, err := NewMongoADriver(os.Getenv("mdb"), "d", "user")
	if err != nil {
		t.Error(err)
	}
	id, err := LoginViaUsername("cgrichard", "shaula", mdbd)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(id)
}

func TestAddingUser(t *testing.T) {
	mad, err := NewMongoADriver(os.Getenv("mdb"), "d", "user")
	if err != nil {
		t.Error(err)
	}

	u := User{}
	u.Username = "cgrichard"
	u.Password = "shaula"
	u.Email = "richard@test.com"
	e := u.HashAndSalt()
	if e != nil {
		t.Error(e)
	}

	e = u.Save(mad)
	if e != nil {
		t.Error(e)
	}
}
