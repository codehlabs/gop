package gop

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
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

// Creates unique id
func unique_id(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	now := make([]byte, 8)
	binary.BigEndian.PutUint64(now, uint64(time.Now().Unix()))
	b = append(b, now...)
	return base64.RawURLEncoding.EncodeToString(b)[:length+8], nil
}

// Checks if the password is a secure password. If checkifpwned is true then it will check against
// hundred of millions of passwoed previously exposed in data breaches. If checkifbad is true then it will check
// against 1 million commonly used bad passwords
func SecurePassword(password string, checkifpwned, checkifbad bool) error {
	if len(password) < 20 {
		return ErrShortPassword
	}

	if checkifpwned {
		s := sha1.New()
		_, err := s.Write([]byte(password))
		if err != nil {
			return err
		}
		hashed := strings.ToUpper(fmt.Sprintf("%x", s.Sum(nil)))
		res, err := http.Get("https://api.pwnedpasswords.com/range/" + string(hashed[:5]))
		if err != nil {
			return err
		}
		defer res.Body.Close()
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		if idx := bytes.Index(data, []byte(hashed[5:])); idx != -1 {
			data = data[idx:]
			newline := bytes.Index(data, []byte("\n"))
			colon := bytes.Index(data, []byte(":"))
			data = data[colon+1 : newline]
			data = bytes.Replace(data, []byte("\r"), []byte(""), 1)
			breaches, err := strconv.Atoi(string(data))
			if err == nil {
				if breaches > 500 {
					ErrPwnedPassword.append("highly insecure password it has been breached " + string(data))
				}
			}
			return ErrPwnedPassword
		}
	}

	if checkifbad {

	}

	return nil
}
