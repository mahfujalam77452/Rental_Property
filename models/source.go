package models

type lonlat struct {
	Coordinates []string `json:"coordinates`
}

type Property struct{
    ID string `json:"id"`
	Feed int `json:"feed"`
	Country string `json:"country"`
	Country_code string `json:"country_code"`
	State string `json:"state"`
	State_abbr string `json:"state_abbr`
	City string `json:"city"`
	Display string `json:"display"`
	Location_id string `json:"loacation_id"`
    
	Property_name string `json:"property_name"`
	Property_slug string `json:"property_slug"`
	Property_type_category string `json:"property_type_category"`
	
	Usd_price float64 `json:"usd_price"`
	Occupancy int `json:"occupancy"`
	Bedroom_count int `json:"bedroom_count"`
	Bathroom_count int `json:"bathroom_count"`
	Number_of_review int `json:"number_of_review"`
	Review_score_general float64 `json:"review_score_general"`
	Star_rating int `json:"star_rating"`
	
	Amenity_categories []string `json:"amenity_categories"`

	Lonlat lonlat `json:"lonlat"`

	Categories string `json:"categories"`

	Published bool `json:"published"`

	Images []string `json:"images"`




}