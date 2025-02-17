package driver

import (
	"github.com/racg0092/gop"
	"os"
	"testing"
	"time"
)

func TestDriver(t *testing.T) {

	t.Run("mongo_driver", func(t *testing.T) {
		d, err := New(MONGO, InitConfig{
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

		u := gop.User{
			Username: "jdoes00",
			FirsName: "Jon",
			LastName: "Doe",
			Password: "running with the lions in the jungle 0033",
			Email:    "jdoe@email.com",
			DOB:      dob,
			Phone:    "+19999999999",
		}

		if err := u.Save(d); err != nil {
			t.Error(err)
		}

	})

	t.Run("libsql_driver", func(t *testing.T) {
		d, err := New(LIBSQL, InitConfig{Conn: "file:../local.db"})
		if err != nil {
			t.Error(err)
		}

		dob, err := time.Parse("2006-01-02", "1992-01-07")
		if err != nil {
			t.Error(err)
		}

		u := gop.User{
			Username: "jdoes00",
			FirsName: "Jon",
			LastName: "Doe",
			Password: "running with the lions in the jungle 0033",
			Email:    "jdoe@email.com",
			DOB:      dob,
			Phone:    "+19999999999",
		}

		err = u.Save(d)
		if err != nil {
			t.Error(err)
		}
	})

}
