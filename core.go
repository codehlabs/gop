package gop

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"golang.org/x/crypto/argon2"
	"strings"
	"time"
)

type HAlgo int

const (
	SHA256 HAlgo = iota + 1
	ARGON2
)

type Db interface {
	Save(u User) (string, error)

	Update(u User) error

	// Delete user
	Delete(id string) error

	// Read everything excluding only the password
	Read(id string, includeProfile bool) (User, error)

	// Read user data excluding secure data
	ReadNonCritical(id string, includeProfile bool) (UserNonConfidential, error)
}

type Config struct {
	UseBuiltInSaveLogic     bool  // when calling [User] to save it will use the custom built in logic
	UniqueIDLength          int   // defines built in unique id seed length
	HashAlgo                HAlgo // Hashing algorithm used defaults to SHA256
	SaltAndPepper           bool  // Defines if you want to salt and pepper the paswword
	UseEmailIfUsernameBlank bool  // Defines if you want to use an email as an username if none set
	IsPawnedPassword        bool  // Check is password is pawned
	IsBadPassword           bool  // Checks if passwored is bad
}

// Sets email as username if username is blank
func (c *Config) IfUserNameBlankUseEmail() *Config {
	c.UseEmailIfUsernameBlank = true
	return c
}

func (c *Config) CheckIfPawnedPassword() *Config {
	c.IsPawnedPassword = true
	return c
}

func (c *Config) CheckIfBadPassword() *Config {
	c.IsBadPassword = true
	return c
}

var config = &Config{UseBuiltInSaveLogic: true, UniqueIDLength: 32, HashAlgo: SHA256}

// Sets configuration to max security
func SetConfigMaxSec() *Config {
	config = &Config{
		UseBuiltInSaveLogic: true,
		UniqueIDLength:      32,
		HashAlgo:            ARGON2,
		SaltAndPepper:       true,
		IsPawnedPassword:    true,
		IsBadPassword:       true,
	}
	return config
}

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

// Takes password from user and store password and compares them
func ValidateHash(password, storepassword string) (bool, error) {
	parts := strings.Split(storepassword, "$")
	algo := parts[0]
	parts = parts[1:]

	switch algo {
	case "sha256":
		if len(parts) < 2 {
			return false, ErrMalformedEncodedPassword
		}
		salt, e := hex.DecodeString(parts[0])
		if e != nil {
			return false, e
		}
		return validate_sha256(password, storepassword, salt), nil
	case "argon2i":
		if len(parts) < 5 {
			return false, ErrMalformedEncodedPassword
		}
		salt, e := hex.DecodeString(parts[3])
		if e != nil {
			return false, e
		}
		return validate_argon2i(password, storepassword, salt)
	default:
		return false, errors.New("unknow hashing algorithm")
	}
}

func validate_argon2i(password, storepassword string, salt []byte) (bool, error) {
	hash, _, e := argon2i_hash(password, salt)
	if e != nil {
		return false, e
	}
	return hash == storepassword, nil
}

func validate_sha256(password, oldpassword string, salt []byte) bool {
	sha := sha256.New()
	sha.Write(salt)
	sha.Write([]byte(password))
	hash := fmt.Sprintf("%x", sha.Sum(nil))
	hexsalt := fmt.Sprintf("%x", salt)
	passwordhash := fmt.Sprintf("sha256$%s$%s", hexsalt, hash)
	return passwordhash == oldpassword
}

func argon2i_hash(password string, salt []byte) (hash string, version int, err error) {
	if salt == nil {
		salt = make([]byte, SaltLength)
		_, err = rand.Read(salt)
		if err != nil {
			return "", -1, err
		}
	}

	hashbuff := argon2.Key([]byte(password), salt, Iterations, Memory, Parallelism, KeyLength)

	encodedSalt := fmt.Sprintf("%x", salt)
	encodedHash := fmt.Sprintf("%x", hashbuff)

	version = argon2.Version

	hash = fmt.Sprintf("argon2i$%d$%d$%d$%s$%s", Iterations, Memory, Parallelism, encodedSalt, encodedHash)
	return hash, version, nil
}
