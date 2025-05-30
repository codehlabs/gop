package gop

import (
	"os"
	"testing"
	"time"
)

func TestDriver(t *testing.T) {

	t.Run("mongo_driver", func(t *testing.T) {
		d, err := NewDriver(MONGO, DriverConfig{
			Conn:       os.Getenv("mdb"),
			Database:   "d",
			Collection: "user",
		})

		if err != nil {
			t.Error(err)
		}

		dob, err := time.Parse("2006-01-02", "1992-01-07")
		if err != nil {
			t.Error(err)
		}

		u := User{
			Username:  "jdoes00",
			FirstName: "Jon",
			LastName:  "Doe",
			Password:  "running with the lions in the jungle 0033",
			Email:     "jdoe@email.com",
			DOB:       dob,
			Phone:     "+19999999999",
		}

		if _, err := u.Save(d); err != nil {
			t.Error(err)
		}

	})

	t.Run("libsql_driver", func(t *testing.T) {
		d, err := NewDriver(LIBSQL, DriverConfig{Conn: "file:local.db"})
		if err != nil {
			t.Error(err)
		}

		dob, err := time.Parse("2006-01-02", "1992-01-07")
		if err != nil {
			t.Error(err)
		}

		u := User{
			Username:  "jdoes00",
			FirstName: "Jon",
			LastName:  "Doe",
			Password:  "running with the lions in the jungle 0033",
			Email:     "jdoe@email.com",
			DOB:       dob,
			Phone:     "+19999999999",
		}

		_, err = u.Save(d)
		if err != nil {
			t.Error(err)
		}
	})

}

func TestUtilsLibSql(t *testing.T) {
	config := DriverConfig{Conn: "file:local.db"}
	driver, err := NewDriver(LIBSQL, config)
	if err != nil {
		t.Error(err)
	}

	t.Run("table not found", func(t *testing.T) {
		exist, err := tableExists(driver, "foos")

		if err != nil {
			t.Error(err)
		}

		if exist == true {
			t.Errorf("expected foos to be false but got %v", exist)
		}

	})

	t.Run("table found", func(t *testing.T) {
		exist, err := tableExists(driver, "bar")
		if err != nil {
			t.Error(err)
		}

		if exist == false {
			t.Errorf("expected bar to extist but got %v", exist)
		}
	})

	t.Run("duplicate account", func(t *testing.T) {
		u := User{
			Username: "jdoes00",
			Email:    "jdoe@emai.com",
			Phone:    "+19999999999",
		}
		err := checkIfDup(driver, &DriverBehavior{true, true, true}, u)

		if err != ErrDupUser {
			t.Errorf("expected %q got %q", ErrDupUser, err)
		}
	})

	t.Run("no duplicate account", func(t *testing.T) {
		u := User{
			Username: "odinson",
			Email:    "o@son.com",
			Phone:    "+16767776655",
		}

		err := checkIfDup(driver, &DriverBehavior{true, true, true}, u)

		if err != nil {
			t.Error(err)
		}
	})

}
