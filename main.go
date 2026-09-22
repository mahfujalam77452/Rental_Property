package main

import (
	_ "Beego_API/routers"
	"Beego_API/services"
	_ "fmt"
	

	"github.com/beego/beego/v2/core/logs"
	bee "github.com/beego/beego/v2/server/web"
)

func main() {
    path,err := bee.AppConfig.String("data_file")
    
	if err != nil {
		logs.Error("Extracting path variable from app.conf file failed !")
	}
	
	propertiesObject,err := services.NewPropertyServices(path)
	
	if err != nil {
		logs.Error("Property load failed : ",err)
	}
	services.PropertyServices = *propertiesObject
	bee.Run()
}
