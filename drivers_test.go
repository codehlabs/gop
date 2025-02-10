package gop

import (
	"fmt"
	"testing"
)

func TestTableExists(t *testing.T) {
	driver, err := NewLibSqlADriver("file:./local.db", "user")
	if err != nil {
		t.Error(err)
	}
	exist, err := tableExists(driver, "user")
	if err != nil {
		t.Error(err)
	}
	if !exist {
		if err := createUserTable(driver, "user"); err != nil {
			t.Error(err)
		}
	} else {
		fmt.Println("user table exists:", exist)
	}
}
