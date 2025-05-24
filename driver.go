package gop

import (
	"database/sql"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
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
type DriverBehavior struct {
	UniqueUsername bool // unique database username defaults to true
	UniqueEmail    bool // unique database email defaults to true
	UniquePhone    bool // unique database phone defaults to true
}

// Driver confifuration
type DriverConfig struct {
	Conn       string
	Database   string
	TableName  string
	Collection string // when using mongo action driver only
}

// Check if any of the struct fields has the default initialization value
func (i DriverConfig) IsDefault() bool {

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

var driver_config *DriverConfig

// Set driver configuration
func SetDriverConfig(i DriverConfig) (*DriverConfig, error) {
	if i.IsDefault() {
		return nil, errors.New("all values must be set for driver configuration")
	}
	driver_config = &i
	return driver_config, nil
}

// Driver configuration
func GetDriverConfig() *DriverConfig {
	return driver_config
}

var driverbehavior = &DriverBehavior{true, true, true}

type ActionDriver interface {
	Db
	Login(username, email, phone string, password string) (id string, err error)
	Db() *sql.DB
	MongoDb() *mongo.Database
}

var driver ActionDriver

// Returns new driver based on dt and configuration
func NewDriver(dt Type, config DriverConfig) (ActionDriver, error) {
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

// Get driver
func GetDriver() ActionDriver {
	if driver == nil {
		panic("driver is not set")
	}
	return driver
}

var (
	ErrUnknowDriver = errors.New("unknow action driver type")
	ErrDbIsNil      = errors.New("driver sql db is <nil>")
	ErrDupUser      = errors.New("duplicate user email, phone or username is already in the database")
)
