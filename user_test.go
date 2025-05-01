package gop

import (
	"github.com/racg0092/gop/rdb"
	_ "github.com/tursodatabase/go-libsql"
	"testing"
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
