package controllers

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
	"github.com/dgrijalva/jwt-go"
	"github.com/galamarv/test_backend_mnc/models"
	"github.com/google/uuid"
)

type UserController struct {
	web.Controller
}

type RegisterInput struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
	Pin         string `json:"pin"`
}

func (u *UserController) Register() {
	var input RegisterInput
	if err := u.ParseForm(&input); err != nil {
		u.CustomAbort(400, "Invalid input")
	}

	user := models.User{
		Id:          uuid.New().String(), // Convert UUID to string
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		PhoneNumber: input.PhoneNumber,
		Address:     input.Address,
		Pin:         input.Pin,
		CreatedDate: time.Now(),
	}

	o := orm.NewOrm()
	_, err := o.Insert(&user)
	if err != nil {
		u.CustomAbort(400, "Phone Number already registered")
	}

	u.Data["json"] = map[string]interface{}{
		"status": "SUCCESS",
		"result": user,
	}
	u.ServeJSON()
}

type LoginInput struct {
	PhoneNumber string `json:"phone_number"`
	Pin         string `json:"pin"`
}

func (u *UserController) Login() {
	var input LoginInput
	if err := u.ParseForm(&input); err != nil {
		u.CustomAbort(400, "Invalid input")
	}

	o := orm.NewOrm()
	user := models.User{PhoneNumber: input.PhoneNumber}
	if err := o.Read(&user, "PhoneNumber"); err != nil {
		u.CustomAbort(400, "Phone Number not registered")
	}

	if user.Pin != input.Pin {
		u.CustomAbort(400, "Phone Number and PIN don’t match")
	}

	// Generate JWT tokens
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.Id,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, _ := token.SignedString([]byte("your_jwt_secret"))

	u.Data["json"] = map[string]interface{}{
		"status": "SUCCESS",
		"result": map[string]string{
			"access_token": tokenString,
		},
	}
	u.ServeJSON()
}
