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
    property,err := services.PropertyServices.GetPropertyByID(id)

	if err != nil {
		c.Data["json"] = map[string]string{
			"error":err.Error(),
		}
		c.Ctx.ResponseWriter.WriteHeader(404)
		c.ServeJSON()
		return
	}

	c.Data["json"] = property
	c.ServeJSON()
}

func (c *PropertyController) Get() {

	properties,err := services.PropertyServices.GetProperties()

	if err != nil {
		c.Data["json"] = map[string]string{
			"error":err.Error(),
		}
		c.Ctx.ResponseWriter.WriteHeader(404)
		c.ServeJSON()
		return
	}

	c.Data["json"] = properties
	c.ServeJSON()

	
}