package routers

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/galamarv/test_backend_mnc/controllers"
)

func init() {
	web.Router("/register", &controllers.UserController{}, "post:Register")
	web.Router("/login", &controllers.UserController{}, "post:Login")
	web.Router("/topup", &controllers.TransactionController{}, "post:TopUp")
	web.Router("/pay", &controllers.TransactionController{}, "post:Pay")
	web.Router("/transfer", &controllers.TransactionController{}, "post:Transfer")
	web.Router("/transactions", &controllers.TransactionController{}, "get:TransactionReport")
}
