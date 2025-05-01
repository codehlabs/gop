package gop

import "fmt"

type Address struct {
	Address  string `bson:"address" form:"address" sql:"address,text"`
	Address2 string `bson:"address2,omitempty" form:"address2" sql:"address2,text"`
	City     string `bson:"city" form:"city" sql:"city,text"`
	State    string `bson:"state,omitempty" form:"state" sql:"state,text"`
	ZipCode  string `bson:"zipcode" form:"zipcode" sql:"zipcode,text"`
}

func (a Address) String() string {
	return fmt.Sprintf("%s %s %s, %s, %s", a.Address, a.Address2, a.City, a.State, a.ZipCode)
}
