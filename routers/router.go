package routers

import (
	"Beego_API/controllers"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSRouter("/properties/:id", &controllers.PropertyController{}, "get:GetByID"),
		beego.NSRouter("/properties", &controllers.PropertyController{}, "get:Get"),
	)

	beego.AddNamespace(ns)
}
