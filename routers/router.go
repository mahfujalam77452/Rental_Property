// @APIVersion 1.0.0
// @Title Travel API
// @Description API for searching and booking travel properties
// @Contact admin@w3engineers.com
package routers

import (
	"Beego_API/controllers"
	beego "github.com/beego/beego/v2/server/web"
	
)

func init() {
	ns := beego.NewNamespace("/v1",
	   beego.NSNamespace("/properties",
		beego.NSInclude(&controllers.PropertyController{}),
	   ))

		

	beego.AddNamespace(ns)
    // Swagger UI + swagger.json/yml
	beego.SetStaticPath("/swagger", "swagger")
	
}


