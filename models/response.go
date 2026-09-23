package models

type Breadcrumb struct {
	LocationID string
	Name       string
	Type       string
	Slug       string
	Display    []string
}

type GeoInfo struct {
	Breadcrumbs []Breadcrumb `json:"breadcrumbs"`
	City        string       `json:"city"`
	Country     string       `json:"country"`
	CountryCode string       `json:"country_code"`
	Name        string       `json:"name"`
	LocationID  string       `json:"location_id"`
	Lat         float64      `json:"lat"`
	Lon         float64      `json:"lon"`
	State       string       `json:"state"`
	StateAbbr   string       `json:"state_abbr"`
}

type Counts struct {
	Bathroom  int `json:"bathroom"`
	Bedroom   int `json:"bedroom"`
	Reviews   int `json:"reviews"`
	Occupancy int `json:"occupancy"`
}

type Image struct {
	Count  int      `json:"count"`
	Images []string `json:"images"`
}

type PropertyInfo struct {
	Amenities    []string `json:"amenities"`
	Name         string   `json:"name"`
	Slug         string   `json:"slug"`
	PropertyType string   `json:"property_type"`
	Price        float64  `json:"price"`
	ReviewScore  float64  `json:"review_score"`
	StarRating   int      `json:"star_rating"`
	Counts       Counts   `json:"counts"`
	Image        Image    `json:"image"`
}

type PropertyResponse struct {
	ID        string       `json:"id"`
	Feed      int          `json:"feed"`
	Published bool         `json:"published"`
	GeoInfo   GeoInfo      `json:"geo_info"`
	Property  PropertyInfo `json:"property"`
}

type Result struct {
	Count int                `json:"count"`
	Items []PropertyResponse `json:"items"`
}

type PropertyListResponse struct {
	Result Result `json:"result"`
}
