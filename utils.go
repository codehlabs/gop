package gop

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// Authenticates into the system using email as the identifier
func LoginViaEmail(email, password string, driver ActionDriver) (string, error) {
	return driver.Login("", email, "", password)
}

// Authenti into the system using username as the identifier
func LoginViaUsername(username, password string, driver ActionDriver) (string, error) {
	return driver.Login(username, "", "", password)
}

// Authenticate into the system using the phone number as the identifier
// NOTE: may remove this one
func LoginViaPhone(phone, password string, driver ActionDriver) (string, error) {
	return driver.Login("", "", phone, password)
}

// Takes salt and password and returns a hash
func ValidateHash(salt, password string) (string, error) {
	bslice, err := hex.DecodeString(salt)
	if err != nil {
		return "", err
	}
	bslice = append(bslice, []byte(password)...)
	sha := sha256.New()
	_, err = sha.Write(bslice)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", sha.Sum(nil)), nil
}
