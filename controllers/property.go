package controllers

import (
	"Beego_API/services"
	"Beego_API/utils"
	"errors"
	"strings"

	beego "github.com/beego/beego/v2/server/web"
)

type PropertyController struct {
	beego.Controller
}

func (c *PropertyController) GetByID() {

	id := c.Ctx.Input.Param(":id")
	//Avoiding the error for string Like "  "
	id = strings.TrimSpace(id)

	if id == "" {

		c.Data["json"] = map[string]string{
			"error": "invalid id",
		}
		c.Ctx.ResponseWriter.WriteHeader(400)
		c.ServeJSON()
		return

	}

	property, err := services.PropertyServices.GetPropertyByID(id)

	var apiError *utils.APIError

	//We are using it here so that user easily understand that property not found

	if errors.As(err, &apiError) {

		c.Ctx.ResponseWriter.WriteHeader(apiError.Status)

		c.Data["json"] = map[string]string{
			"Error": apiError.Error(),
		}
		 

		c.ServeJSON()

		return
	}

	if err != nil {
		c.Data["json"] = map[string]string{
			"error": err.Error(),
		}
		c.Ctx.ResponseWriter.WriteHeader(500)
		c.ServeJSON()
		return
	}
    
	c.Data["json"] = property
	c.Ctx.ResponseWriter.WriteHeader(200)
	c.ServeJSON()
}

func (c *PropertyController) Get() {

	//Parsing Query
	query, err := utils.ParsePropertyQueryWithValidation(&c.Controller)

	if err != nil {
		c.Data["json"] = map[string]string{
			"error": err.Error(),
		}
		c.Ctx.ResponseWriter.WriteHeader(400)
		c.ServeJSON()
		return
	}
   
	 //Getting required properties with query
	properties, err := services.PropertyServices.GetProperties(query)

	if err != nil {
		c.Data["json"] = map[string]string{
			"error": err.Error(),
		}
		c.Ctx.ResponseWriter.WriteHeader(500)
		c.ServeJSON()
		return
	}

	c.Data["json"] = properties
	c.Ctx.ResponseWriter.WriteHeader(200)
	c.ServeJSON()

}
