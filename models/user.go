package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type User struct {
	Id          string    `orm:"pk;size(36)"` // Store UUID as a string
	FirstName   string    `orm:"size(100)"`
	LastName    string    `orm:"size(100)"`
	PhoneNumber string    `orm:"unique;size(15)"`
	Address     string    `orm:"size(255)"`
	Pin         string    `orm:"size(6)"`
	Balance     float64   `orm:"default(0)"`
	CreatedDate time.Time `orm:"auto_now_add;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(User))
}
