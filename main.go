package main

import (
	_ "Beego_API/routers"
	"Beego_API/services"
	_ "fmt"
	

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {

	properties,err := services.NewProperties("data/rental_properties.json")
	
	if err != nil {
		logs.Error("Property load failed : ",err)
	}
	services.PropertyService = *properties
	beego.Run()
}
