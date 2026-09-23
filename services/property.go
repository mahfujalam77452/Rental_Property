package services

import (
	"Beego_API/models"
	"Beego_API/utils"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	_ "github.com/beego/beego/v2/core/logs"
)

// Helper function for Transform the property object to the Response Formate
func transformProperty(property models.Property) (*models.PropertyResponse, error) {

	//For store the unmarshalled value of Category field stored before as String
	var breadcrumbs []models.Breadcrumb

	err := json.Unmarshal([]byte(property.Categories), &breadcrumbs)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse the categories : %w", err)
	}

	if len(property.LonLat.Coordinates) < 2 {
		return nil, errors.New("Invalid property condinates !")
	}

	response := &models.PropertyResponse{
		ID:        property.ID,
		Feed:      property.Feed,
		Published: property.Published,

		GeoInfo: models.GeoInfo{
			Breadcrumbs: breadcrumbs,
			City:        property.City,
			Country:     property.Country,
			CountryCode: property.CountryCode,
			Name:        property.Display,
			LocationID:  property.LocationID,
			Lat:         property.LonLat.Coordinates[1],
			Lon:         property.LonLat.Coordinates[0],
			State:       property.State,
			StateAbbr:   property.StateAbbr,
		},
		Property: models.PropertyInfo{
			Amenities:    property.AmenityCategories,
			Name:         property.PropertyName,
			Slug:         property.PropertySlug,
			PropertyType: property.PropertyTypeCategory,
			Price:        property.USDPrice,
			ReviewScore:  property.ReviewScoreGeneral,
			StarRating:   property.StarRating,

			Counts: models.Counts{
				Bathroom:  property.BathroomCount,
				Bedroom:   property.BedroomCount,
				Reviews:   property.NumberOfReview,
				Occupancy: property.Occupancy,
			},
			Image: models.Image{
				Count:  len(property.Images),
				Images: property.Images,
			},
		},
	}

	return response, nil
}

// Healper Function for read json file and transfrom it to our Property struct

func loadProperties(path string) ([]models.Property, error) {

	data, err := os.ReadFile(path)

	if err != nil {
		return nil, fmt.Errorf("JSON file read failed : %w", err)
	}

	var properties []models.Property

	err = json.Unmarshal(data, &properties)

	if err != nil {

		return nil, fmt.Errorf("Converstion of Propeties data (json -> struct) failed : %w", err)
	}

	return properties, nil

}

// Struct for Holding properties data and all services

type AllPropertyServices struct {
	Properties []models.Property
}

// Property Service Declaration
var PropertyServices AllPropertyServices

// Constructor For Property Services
func NewPropertyServices(filePath string) (*AllPropertyServices, error) {

	filePath = strings.TrimSpace(filePath)

	if filePath == "" {
		return nil, errors.New("invalid empty file path")
	}

	properties, err := loadProperties(filePath)

	if err != nil {
		return nil, err
	}

	return &AllPropertyServices{Properties: properties}, nil

}

// Get Property by ID Service
func (s *AllPropertyServices) GetPropertyByID(id string) (*models.PropertyResponse, error) {

	for _, property := range s.Properties {
		if property.ID == id {

			transformProperty, err := transformProperty(property)

			if err != nil {

				return nil, fmt.Errorf("Can't transform data :%w", err)

			}

			return transformProperty, nil

		}
	}

	return nil, &utils.APIError{
		Status:  404,
		Message: "Property not found !",
	}
}

//Get properties by Query Service

func (s *AllPropertyServices) GetProperties(query *utils.PropertyFilter) (*models.PropertyListResponse, error) {

	var propertyList = make([]models.PropertyResponse, 0)

	for _, property := range s.Properties {

		if utils.FilterProperty(&property, query) {

			transformedProperty, err := transformProperty(property)

			if err != nil {
				return nil, fmt.Errorf("Tranfomation of data failed : %w", err)
			}

			propertyList = append(propertyList, *transformedProperty)

		}

	}
     
	//Handling the limit after proper filtering
	
	if query.Limit != nil && len(propertyList) > *query.Limit {
		propertyList = propertyList[:*query.Limit]
	}

	return &models.PropertyListResponse{
		Result: models.Result{
			Count: len(propertyList),
			Items: propertyList,
		},
	}, nil
}
