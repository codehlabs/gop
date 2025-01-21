package gop

import (
	"context"
	"fmt"
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UseDb struct {
}

func (udb UseDb) Save(u User) error {
	ctx := context.Background()
	client, e := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("mdb")))
	if e != nil {
		return e
	}
	defer client.Disconnect(ctx)

	col := client.Database("d").Collection("user")

	_, e = col.InsertOne(ctx, u)
	if e != nil {
		return e
	}

	return nil
}

func (ubd UseDb) Delete(id string) error {
	return nil
}

func (udb UseDb) Update(u User) error {
	return nil
}

func (udb UseDb) Read(id string) (User, error) {
	return User{}, nil
}

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
	u := User{}
	u.Username = "cgrichard"
	u.Password = "shaula"
	u.Email = "richard@test.com"
	e := u.HashAndSalt()
	if e != nil {
		t.Error(e)
	}
	dbi := UseDb{}

	e = u.Save(dbi)
	if e != nil {
		t.Error(e)
	}
}
