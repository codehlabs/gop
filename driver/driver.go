package driver

import (
	"database/sql"
	"errors"

	"github.com/racg0092/gop"
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

var driver_config = &DriverConfig{true, true, true}

type ActionDriver interface {
	gop.Db
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
