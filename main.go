package main

import (
	"github.com/beego/beego/v2/server/web"
	_ "github.com/galamarv/test_backend_mnc/routers"
)

func main() {
	web.Run()
}
