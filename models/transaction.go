package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Transaction struct {
	Id            string `orm:"pk;size(36)"` // Store UUID as a string
	UserId        string `orm:"size(36)"`    // Foreign key to User's UUID
	Type          string `orm:"size(10)"`    // "CREDIT" or "DEBIT"
	Amount        float64
	Remarks       string `orm:"size(255)"`
	BalanceBefore float64
	BalanceAfter  float64
	CreatedDate   time.Time `orm:"auto_now_add;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(Transaction))
}
