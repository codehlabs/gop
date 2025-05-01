package gop

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"time"
)

type Db interface {
	Save(u User) error
	Update(u User) error
	Delete(id string) error
	Read(id string) (User, error)
}

type Config struct {
	UseBuiltInSaveLogic bool // when calling [User] to save it will use the custom built in logic
	UniqueIDLength      int  // defines built in unique id seed length
}

var config = &Config{UseBuiltInSaveLogic: true, UniqueIDLength: 32}

func SetConfig(c Config) {
	config = &c
}

func GetConfig() *Config {
	return config
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
