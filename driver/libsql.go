package driver

import (
	"database/sql"

	"github.com/racg0092/gop"
	"github.com/racg0092/gop/rdb"
	_ "github.com/tursodatabase/go-libsql"
)

type LibSqlADriver struct {
	dbpath string
	orm    rdb.ORM
}

// Closes the connection
func (d LibSqlADriver) Close() error {
	return d.orm.Close()
}

// Creates a new LibSql Driver where dbpath is the database and usertable is the table
// to save the users to
func NewLibSqlADriver(dbpath string) (LibSqlADriver, error) {
	driver := LibSqlADriver{dbpath: dbpath}
	orm, err := rdb.Open("libsql", dbpath)
	if err != nil {
		return driver, err
	}
	driver.orm = orm
	return driver, nil
}

func (d LibSqlADriver) Login(username, email, phone, password string) (id string, err error) {

	return id, err
}

func (d LibSqlADriver) Save(u gop.User) error {
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

func (d LibSqlADriver) Update(u gop.User) error {
	return nil
}

func (d LibSqlADriver) Delete(id string) error {
	return nil
}

func (d LibSqlADriver) Read(id string) (gop.User, error) {
	return gop.User{}, nil
}

// Check is table exists in the database
func tableExists(driver LibSqlADriver, table string) (bool, error) {
	var name string
	var err error
	driver.orm.Raw(func(db *sql.DB) {
		err = db.QueryRow(`
      select name
      from sqlite_master
      where type='table'
      and name =?
      `, table).Scan(&name)
		if err == sql.ErrNoRows {
			return
		}
		if err != nil {
			return
		}
		return
	})
	if err != nil {
		return false, err
	}
	return true, err
}

// Creates user table
func createUserTable(driver LibSqlADriver, tablename string) error {
	u := gop.User{}
	err := driver.orm.CreateTable(u, tablename)
	if err != nil {
		return err
	}
	return nil
}

// Chekcs if there are any duplicates in the database
func checkIfDup(driver LibSqlADriver, config *DriverConfig) error {
	var err error
	driver.orm.Raw(func(db *sql.DB) {
		//TODO: finish business logic
	})
	return err
}
