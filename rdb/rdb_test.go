package rdb

import (
	"testing"
	"time"
)

type Sample struct {
	Id       int       `sql:"id,integer,primary key,autoincrement"`
	Name     string    `sql:"name,TEXT"`
	Foo      Foo       `sql:"omit"`
	Content  string    `sql:"content,TEXT"`
	CratedAt time.Time `sql:"created_at,integer,default current_timestamp"`
}

type Foo struct {
	Set  int `sql:"[set],integer"`
	Rank int `sql:"rank,integer"`
}

func TestRDB(t *testing.T) {
	s := Sample{
		Name:    "dumbo",
		Content: "Hello dumbo how are today ? In this crazy world",
	}
	orm, err := Open("libsql", "file:./test.db")
	if err != nil {
		t.Error(err)
	}
	defer orm.db.Close()

	if err := orm.CreateTable(s, ""); err != nil {
		t.Error(err)
	}

	if err := orm.Save(&s, "samples"); err != nil {
		t.Error(err)
	}
}
