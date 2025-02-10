package gop

import (
	"testing"

	"github.com/racg0092/gop/rdb"
	_ "github.com/tursodatabase/go-libsql"
)

func TestSqlUser(t *testing.T) {
	orm, err := rdb.Open("libsql", "file:./rdb/test.db")
	if err != nil {
		t.Error(err)
	}

	u := User{}

	err = orm.CreateTable(u, "")
	if err != nil {
		t.Error(err)
	}
}
