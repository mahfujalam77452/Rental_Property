package controllers

import (
	"Beego_API/services"
	

	beego "github.com/beego/beego/v2/server/web"
)
type PropertyController struct {
	beego.Controller
}

func (c *PropertyController) GetByID(){

	id := c.Ctx.Input.Param(":id")
    properyty,err := services.PropertyServices.GetPropertyByID(id)

	if err != nil {
		c.Data["json"] = map[string]string{
			"error":err.Error(),
		}
		c.Ctx.ResponseWriter.WriteHeader(404)
		c.ServeJSON()
		return
	}

	c.Data["json"] = properyty
	c.ServeJSON()
}