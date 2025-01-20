package gop

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func LoginViaEmail(email, password string, driver ActionDriver) (string, error) {
	return driver.Login("", email, "", password)
}

func LoginViaUsername(username, password string, driver ActionDriver) (string, error) {
	return driver.Login(username, "", "", password)
}

func LoginViaPhone(phone, password string, driver ActionDriver) (string, error) {
	return driver.Login("", "", phone, password)
}

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
