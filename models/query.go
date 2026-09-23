type PropertyFilter struct {
	Limit *int 
	MinPrice *float64
	MaxPrice *float64
	MinStarRating *int
	MinReviewScore *float64
	MinReviews *int
	Published *bool
	PropertyType *string
	Feed *int
	MinBedroom *int
	Amenities *[]string
}