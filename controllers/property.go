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
// @Title Get Property By ID
// @Description Get a rental property by its ID
// @Param id path string true "Property ID"
// @Success 200 {object} models.PropertyResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @router /:id [get]
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


// @Title Get Properties
// @Description Get rental properties with optional filters
//
// @Param limit query int false "Maximum number of results (1-100)"
// @Param min_price query number false "Minimum price"
// @Param max_price query number false "Maximum price"
// @Param min_star_rating query int false "Minimum star rating (0-5)"
// @Param min_review_score query number false "Minimum review score"
// @Param min_reviews query int false "Minimum number of reviews"
// @Param published query boolean false "Published status"
// @Param property_type query string false "Property type"
// @Param feed query int false "Feed ID (11, 12, 22, 24)"
// @Param min_bedroom query int false "Minimum number of bedrooms"
// @Param amenities query string false "Comma-separated amenities; any match is accepted"
//
// @Success 200 {object} models.PropertyListResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
//
// @router / [get]
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
