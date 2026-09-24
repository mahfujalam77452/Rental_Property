package utils

import (
	"Beego_API/models"
	"fmt"
	"slices"
	"strconv"
	"strings"

	beego "github.com/beego/beego/v2/server/web"
)

// This the struct for storing the upcoming query from the url

type PropertyFilter struct {
	Limit          *int
	MinPrice       *float64
	MaxPrice       *float64
	MinStarRating  *int
	MinReviewScore *float64
	MinReviews     *int
	Published      *bool
	PropertyType   *string
	Feed           *int
	MinBedroom     *int
	Amenities      *[]string
}

// For parsing query and validate them
func ParsePropertyQueryWithValidation(c *beego.Controller) (*PropertyFilter, error) {

	filter := &PropertyFilter{}

	//limit
	if value := strings.TrimSpace(c.GetString("limit")); value != "" {

		limit, err := strconv.Atoi(value)
		//just handling unexpected value in limit
		if err != nil || limit <= 0 || limit > 100 {
			limit = 20
		}

		filter.Limit = &limit
	}

	// min_price

	if value := strings.TrimSpace(c.GetString("min_price")); value != "" {

		minPrice, err := strconv.ParseFloat(value, 64)

		if err != nil || minPrice < 0 {
			return nil, fmt.Errorf("invalid min_price: must be a valid positive number")
		}

		filter.MinPrice = &minPrice
	}

	// max_price

	if value := strings.TrimSpace(c.GetString("max_price")); value != "" {

		maxPrice, err := strconv.ParseFloat(value, 64)

		if err != nil || maxPrice < 0 {
			return nil, fmt.Errorf("invalid max_price: must be a valid positive number")
		}

		filter.MaxPrice = &maxPrice
	}

	//min-price & max-price validation

	if filter.MaxPrice != nil && filter.MinPrice != nil && *filter.MaxPrice < *filter.MinPrice {
		return nil, fmt.Errorf("Maximum Price should greater then the Minimum Price ")
	}

	// min_star_rating

	if value := strings.TrimSpace(c.GetString("min_star_rating")); value != "" {

		rating, err := strconv.Atoi(value)

		if err != nil || rating < 0 || rating > 5 {
			return nil, fmt.Errorf("invalid min_star_rating: must be a positive integer (0-5)")
		}

		filter.MinStarRating = &rating
	}

	// min_review_score

	if value := strings.TrimSpace(c.GetString("min_review_score")); value != "" {

		score, err := strconv.ParseFloat(value, 64)

		if err != nil || score < 0 {
			return nil, fmt.Errorf("invalid min_review_score: must be a positive number")
		}

		filter.MinReviewScore = &score
	}

	// min_reviews

	if value := strings.TrimSpace(c.GetString("min_reviews")); value != "" {

		reviews, err := strconv.Atoi(value)

		if err != nil || reviews < 0 {
			return nil, fmt.Errorf("invalid min_reviews: must be a positive integer")
		}

		filter.MinReviews = &reviews
	}

	// published

	if value := strings.TrimSpace(c.GetString("published")); value != "" {

		published, err := strconv.ParseBool(value)

		if err != nil {
			return nil, fmt.Errorf("invalid published: must be true or false")
		}

		filter.Published = &published
	}

	// property_type

	if value := strings.TrimSpace(c.GetString("property_type")); value != "" {

		validTypes := map[string]bool{
			"Hotel":     true,
			"House":     true,
			"Apartment": true,
			"Villa":     true,
			"Resort":    true,
			"Hostel":    true,
		}

		if !validTypes[value] {
			return nil, fmt.Errorf(
				"invalid property_type: must be Hotel, House, Apartment, Villa, Resort, or Hostel",
			)
		}

		filter.PropertyType = &value
	}

	// feed
	if value := strings.TrimSpace(c.GetString("feed")); value != "" {

		feed, err := strconv.Atoi(value)

		if err != nil {
			return nil, fmt.Errorf("invalid feed: must be an integer")
		}

		validFeeds := map[int]bool{
			11: true,
			12: true,
			22: true,
			24: true,
		}

		if !validFeeds[feed] {
			return nil, fmt.Errorf(
				"invalid feed: must be 11, 12, 22, or 24",
			)
		}

		filter.Feed = &feed
	}

	// min_bedroom

	if value := strings.TrimSpace(c.GetString("min_bedroom")); value != "" {

		bedroom, err := strconv.Atoi(value)

		if err != nil || bedroom < 1 {
			return nil, fmt.Errorf("invalid min_bedroom: must be a positive integer (>0)")
		}

		filter.MinBedroom = &bedroom
	}

	// amenities

	if value := strings.TrimSpace(c.GetString("amenities")); value != "" {

		amenities := strings.Split(value, ",")

		filter.Amenities = &amenities
	}

	return filter, nil

}

// For apply filters on a property.

func FilterProperty(property *models.Property, query *PropertyFilter) bool {

	//Filtering minprice

	if query.MinPrice != nil && property.USDPrice < *query.MinPrice {
		return false
	}

	//Filtering Maxprice

	if query.MaxPrice != nil && property.USDPrice > *query.MaxPrice {
		return false
	}

	//Filtering Minimum Rating

	if query.MinStarRating != nil && property.StarRating < *query.MinStarRating {
		return false
	}

	//Filtering Minimum Review Score

	if query.MinReviewScore != nil && property.ReviewScoreGeneral < *query.MinReviewScore {
		return false
	}

	//Filtering Minimum review counts

	if query.MinReviews != nil && property.NumberOfReview < *query.MinReviews {
		return false
	}

	//Filtering Published or not

	if query.Published != nil && property.Published != *query.Published {
		return false
	}

	//Filtering Property Types

	if query.PropertyType != nil && property.PropertyTypeCategory != *query.PropertyType {
		return false
	}

	//Filtering Property Feed

	if query.Feed != nil && property.Feed != *query.Feed {
		return false
	}

	//Filtering Minimum Bed Room

	if query.MinBedroom != nil && property.BedroomCount < *query.MinBedroom {
		return false
	}

	// Filtering Amenities

	if query.Amenities != nil {

		for _, amenity := range *query.Amenities {

			if slices.Contains(property.AmenityCategories, strings.TrimSpace(amenity)) {
				return true
			}
		}

		return false
	}

	return true
}
