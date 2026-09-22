package main

import (
	_ "Beego_API/routers"

	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	
	beego.Run()
}
