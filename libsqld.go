package gop

import (
	"database/sql"
	"fmt"

	_ "github.com/tursodatabase/go-libsql"
)

type LibSqlADriver struct {
	dbpath string
	db     *sql.DB
}

// Creates a new LibSql Driver where dbpath is the database and usertable is the table
// to save the users to
func NewLibSqlADriver(dbpath, usertable string) (LibSqlADriver, error) {
	driver := LibSqlADriver{dbpath: dbpath}
	db, err := sql.Open("libsql", dbpath)
	if err != nil {
		return driver, err
	}
	driver.db = db
	return driver, nil
}

// Check is table exists in the database
func tableExists(driver LibSqlADriver, table string) (bool, error) {
	var name string
	err := driver.db.QueryRow(`
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
	query := fmt.Sprintf(`
  CREATE TABLE IF NOT EXISTS %s (
    id TEXT PRIMARY KEY,
    username TEXT NOT NULL,
    firstname TEXT NOT NULL,
    lastname TEXT NOT NULL,
    password TEXT NOT NULL,
    salt TEXT NOT NULL,
    email TEXT NULL,
    phone TEXT NULL,
    age INTEGER NULL,
    dob INTEGER NULL,
    created_at INTEGER DEFAULT CURRENT_TIMESTAMP,
    updated_at INTEGER NULL
  );
  `, tablename)

	_, err := driver.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
