package services

import (
	"Beego_API/models"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	_ "github.com/beego/beego/v2/core/logs"
	_ "github.com/beego/beego/v2/core/logs"
)


func transformProperty(property models.Property)(*models.PropertyResponse,error) {
   var breadcrumbs []models.Breadcrumb

   err := json.Unmarshal([]byte(property.Categories),&breadcrumbs)
   
   if err != nil {
	return  nil,fmt.Errorf("Failed to parse the categories : %w",err)
   }

   if len(property.LonLat.Coordinates) < 2 {
	return  nil,errors.New("Invalid property condinates !")
   }

   response := &models.PropertyResponse{
        ID : property.ID,
		Feed: property.Feed,
		Published: property.Published,

		GeoInfo: models.GeoInfo{
			Breadcrumbs: breadcrumbs,
			City: property.City,
			Country: property.Country,
			CountryCode: property.CountryCode,
			Name: property.Display,
			LocationID: property.LocationID,
			Lat: property.LonLat.Coordinates[1],
			Lon:property.LonLat.Coordinates[0],
			State: property.State,
			StateAbbr: property.StateAbbr,

		},
		Property: models.PropertyInfo{
			Amenities: property.AmenityCategories,
			Name: property.PropertyName,
			Slug: property.PropertySlug,
			PropertyType: property.PropertyTypeCategory,
			Price: property.USDPrice,
			ReviewScore: property.ReviewScoreGeneral,
			StarRating: property.StarRating,

			Counts: models.Counts{
				Bathroom: property.BathroomCount,
				Bedroom: property.BedroomCount,
				Reviews: property.NumberOfReview,
				Occupancy: property.Occupancy,
			},
			Image: models.Image{
				Count: len(property.Images),
				Images: property.Images,
			},
		},

   }

   return  response,nil
}
func loadProperties(path string)([]models.Property,error){
      
	  data,err := os.ReadFile(path)
	  

	  if err != nil {
		return nil,fmt.Errorf("JSON file read failed : %w",err)
	  }

	  var properties []models.Property

	  err = json.Unmarshal(data,&properties)

	  if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("Converstion of Propeties data json -> struct failed : %w",err)
	  }
	  
	 
	  

	  return properties,nil
	  
}

type AllPropertyServices struct {
	Properties []models.Property
}

func NewPropertyServices(filePath string)(*AllPropertyServices,error) {
    
	 filePath = strings.TrimSpace(filePath)

	 if(filePath == "") {
             return  nil,errors.New("invalid empty file path")
	 }

	 properties,err := loadProperties(filePath)

	 if err != nil {
		return  nil,err
	 }

	 

	 return  &AllPropertyServices{Properties: properties},nil

}

var PropertyServices AllPropertyServices

func (s *AllPropertyServices) GetPropertyByID(id string)(*models.PropertyResponse,error) {

	for _,property := range s.Properties {
         if property.ID == id {
			
			transformProperty,err := transformProperty(property)

			if err != nil {
				return  nil,fmt.Errorf("can't transform data !")

			}

			return transformProperty,nil
			
		 }
	}

	return nil,errors.New("Property Not found !")
}

func (s *AllPropertyServices) GetProperties()(*models.PropertyListResponse,error) {

	var propertyList []models.PropertyResponse

	for _,property :=range(s.Properties) {
        
		transformedProperty,err := transformProperty(property)

		if err != nil {
			return  nil,fmt.Errorf("Tranfomation of data failed !")
		}

		propertyList = append(propertyList, *transformedProperty)


	}

	return &models.PropertyListResponse{
         Result: models.Result{
			Count: len(propertyList),
			Items: propertyList,
		 },
	},nil
}