package gop

import (
	"testing"
)

func TestHashing(t *testing.T) {

	t.Run("hashing 256", func(t *testing.T) {
		u := User{Password: "shaula"}
		u.HashAndSalt()

		valid, e := ValidateHash("shaula", u.Password)
		if e != nil {
			t.Error(e)
		}

		if !valid {
			t.Error("hashes do not match")
		}

	})

	t.Run("hashes argon2i", func(t *testing.T) {
		u := User{Password: "shaula"}
		e := u.ArgonHash()
		if e != nil {
			t.Error(e)
		}
		valid, e := ValidateHash("shaula", u.Password)
		if e != nil {
			t.Error(e)
		}
		if !valid {
			t.Error("hashes do not match")
		}
	})
}
