package gop

import (
	"database/sql"
	"fmt"

	"github.com/racg0092/gop/rdb"
	// _ "github.com/tursodatabase/go-libsql"
	"go.mongodb.org/mongo-driver/mongo"
	_ "modernc.org/sqlite"
)

type LibSqlADriver struct {
	dbpath string
	orm    rdb.ORM
}

func (d LibSqlADriver) Db() *sql.DB {
	return d.orm.Db()
}

func (d LibSqlADriver) MongoDb() *mongo.Database {
	return nil
}

// Closes the connection
func (d LibSqlADriver) Close() error {
	return d.orm.Close()
}

// Creates a new LibSql Driver where dbpath is the database and usertable is the table
// to save the users to
func NewLibSqlADriver(dbpath string) (*LibSqlADriver, error) {
	driver := &LibSqlADriver{dbpath: dbpath}
	orm, err := rdb.Open("libsql", dbpath)
	if err != nil {
		return nil, err
	}
	driver.orm = orm
	return driver, nil
}

func (d LibSqlADriver) Login(username, email, phone, password string) (id string, err error) {

	return id, err
}

func (d LibSqlADriver) Save(u User) error {
	err := createUserTable(d, "")
	if err != nil {
		return err
	}

	//TODO: double check fro duplicates

	err = d.orm.Save(&u, "users")
	if err != nil {
		return err
	}

	return nil
}

func (d LibSqlADriver) Update(u User) error {
	return nil
}

func (d LibSqlADriver) Delete(id string) error {
	return nil
}

func (d LibSqlADriver) Read(id string, includeProfile bool) (User, error) {
	return User{}, nil
}

func (d LibSqlADriver) ReadNonCritical(id string, includeProfile bool) (User, error) {
	return User{}, nil
}

// Check is table exists in the database
func tableExists(driver ActionDriver, table string) (bool, error) {
	var name string
	var err error
	db := driver.Db()

	if db == nil {
		return false, ErrDbIsNil
	}

	err = db.QueryRow(`
      select name
      from sqlite_master
      where type='table'
      and name =?
      `, table).Scan(&name)

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

// Creates user table
func createUserTable(driver LibSqlADriver, tablename string) error {
	u := User{}
	err := driver.orm.CreateTable(u, tablename)
	if err != nil {
		return err
	}
	return nil
}

// Checks if there are any duplicates in the database
func checkIfDup(driver ActionDriver, config *DriverBehavior, u User) error {
	//NOTE: may need to pass in table to make it explicit instead of implicit
	db := driver.Db()

	if db == nil {
		return ErrDbIsNil
	}

	query := "select count(id) from users where"

	if config.UniqueEmail {
		query = fmt.Sprintf("%s email = '%s'", query, u.Email)
	}

	if config.UniquePhone {
		query = fmt.Sprintf("%s OR phone = '%s'", query, u.Phone)
	}

	if config.UniqueUsername {
		query = fmt.Sprintf("%s OR username = '%s'", query, u.Username)
	}

	result := db.QueryRow(query)
	err := result.Err()

	if err == sql.ErrNoRows {
		return nil
	}

	if err != nil {
		return err
	}

	var affected int
	err = result.Scan(&affected)
	if err != nil {
		return err
	}

	if affected > 0 {
		return ErrDupUser
	}

	return nil
}
