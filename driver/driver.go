package driver

import (
	"database/sql"
	"errors"
	"github.com/racg0092/gop/core"
)

type Type int

func (t Type) String() string {
	switch t {
	case MONGO:
		return "MONGO"
	case LIBSQL:
		return "LIBSQL"
	case SQLITE:
		return "SQLITE"
	default:
		return "UNKNOW"
	}
}

const (
	MONGO Type = iota + 1
	LIBSQL
	SQLITE
)

// Configures driver behavior
type DriverConfig struct {
	UniqueUsername bool // unique database username defaults to true
	UniqueEmail    bool // unique database email defaults to true
	UniquePhone    bool // unique database phone defaults to true
}

// Driver confifuration
type InitConfig struct {
	Conn       string
	Database   string
	TableName  string
	Collection string // when using mongo action driver only
}

// Check if any of the struct fields has the default initialization value
func (i InitConfig) IsDefault() bool {

	if i.Conn == "" {
		return true
	}

	if i.Database == "" {
		return true
	}

	if i.TableName == "" && i.Collection == "" {
		return true
	}

	return false
}

var config *InitConfig

// Set driver configuration
func Config(i InitConfig) (*InitConfig, error) {
	if i.IsDefault() {
		return nil, errors.New("all values must be set for driver configuration")
	}
	config = &i
	return config, nil
}

// Driver configuration
func GetConfig() *InitConfig {
	return config
}

var driver_config = &DriverConfig{true, true, true}

type ActionDriver interface {
	core.Db
	Login(username, email, phone string, password string) (id string, err error)
	Db() *sql.DB
}

// Sets driver configuration
func SetDriverConfig(c DriverConfig) {
	driver_config = &c
}

// Returns new driver based on dt and configuration
func New(dt Type, config InitConfig) (ActionDriver, error) {
	var driver ActionDriver
	var err error
	switch dt {
	case MONGO:
		driver, err = NewMongoADriver(config.Conn, config.Database, config.Collection)
	case LIBSQL, SQLITE:
		driver, err = NewLibSqlADriver(config.Conn)
	default:
		err = ErrUnknowDriver
	}
	return driver, err
}

var (
	ErrUnknowDriver = errors.New("unknow action driver type")
	ErrDbIsNil      = errors.New("driver sql db is <nil>")
	ErrDupUser      = errors.New("duplicate user email, phone or username is already in the database")
)
