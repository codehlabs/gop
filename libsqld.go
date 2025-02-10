package gop

import (
	"database/sql"

	"github.com/racg0092/gop/rdb"
	_ "github.com/tursodatabase/go-libsql"
)

type LibSqlADriver struct {
	dbpath string
	orm    rdb.ORM
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
	u := User{}
	err := driver.orm.CreateTable(u, tablename)
	if err != nil {
		return err
	}
	return nil
}
