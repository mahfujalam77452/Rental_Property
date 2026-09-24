package services

// -----------------------------------------------------------------------------
// Tests for the property service (services/property.go) and the filter logic
// (utils/filter.go), run against the REAL properties JSON file.
//
// Run with:   go test ./... -v
//
// The dataset has 100 properties, 25 for each feed:
//
//	feed 11 -> IDs BC-1000001 .. BC-1000025   (Japan)
//	feed 12 -> IDs HA-2000001 .. HA-2000025   (United States)
//	feed 22 -> IDs HG-3000001 .. HG-3000025   (Canada)
//	feed 24 -> IDs EP-4000001 .. EP-4000025   (UK / Australia / Germany / ...)
//
// Filtered results always keep the order of the JSON file, so the expected ID
// lists below are written in file order.
// -----------------------------------------------------------------------------

import (
	"Beego_API/models"
	"Beego_API/utils"
	"encoding/json"
	"errors"
	"log"
	"os"
	"slices"
	"strings"
	"testing"
)


// Test setup and helpers


// testDataFilePath is the JSON file used by all tests. The path is relative to
// this /services folder, because "go test" runs inside the package folder.
const testDataFilePath = "../data/rental_properties.json"

// testPropertyService is loaded ONCE in TestMain and shared by every test.
// The tests only READ from it, so sharing it is safe. A test that needs to
// change the data must build its own copy with NewPropertyServices(testDataFilePath).
var testPropertyService *AllPropertyServices

// TestMain runs once before all the other tests in this package.
// It loads the data file one time (like main.go does at application start)
// and then runs every test. 
func TestMain(m *testing.M) {
	var err error
	testPropertyService, err = NewPropertyServices(testDataFilePath)
	if err != nil {
		// TestMain has no "t", so log.Fatalf is used instead of t.Fatalf.
		log.Fatalf("could not load property data from %q: %v", testDataFilePath, err)
	}

	os.Exit(m.Run()) // m.Run() runs all the tests and returns the exit code
}

// loadPropertyService returns the shared service that TestMain loaded.
// No file is read here, so calling it in every test is cheap.
func loadPropertyService(t *testing.T) *AllPropertyServices {
	t.Helper()
	return testPropertyService
}

// pointerTo returns a pointer to any value. PropertyFilter uses pointer fields
// ("nil" means "filter not set"), so this keeps the test cases short.
func pointerTo[T any](value T) *T {
	return &value
}

// findRawPropertyByID returns the original (untransformed) property from the
// loaded data so it can be passed straight to transformProperty.
func findRawPropertyByID(t *testing.T, propertyService *AllPropertyServices, propertyID string) models.Property {
	t.Helper()

	for _, rawProperty := range propertyService.Properties {
		if rawProperty.ID == propertyID {
			return rawProperty
		}
	}
	t.Fatalf("property %s does not exist in the dataset", propertyID)
	return models.Property{} // unreachable, t.Fatalf stops the test
}

// extractPropertyIDs collects the ID of every item in a response, in order.
func extractPropertyIDs(items []models.PropertyResponse) []string {
	propertyIDs := make([]string, 0, len(items))
	for _, item := range items {
		propertyIDs = append(propertyIDs, item.ID)
	}
	return propertyIDs
}

// assertPropertyIDs checks that the response contains exactly the expected IDs
// (same items, same order) and that Result.Count matches the number of items.
func assertPropertyIDs(t *testing.T, actualResponse *models.PropertyListResponse, expectedIDs []string) {
	t.Helper()

	if actualResponse == nil {
		t.Fatal("expected a response, got nil")
	}

	actualIDs := extractPropertyIDs(actualResponse.Result.Items)
	if !slices.Equal(actualIDs, expectedIDs) {
		t.Errorf("returned IDs = %v\n          want IDs = %v", actualIDs, expectedIDs)
	}
	if actualResponse.Result.Count != len(expectedIDs) {
		t.Errorf("Result.Count = %d, want %d", actualResponse.Result.Count, len(expectedIDs))
	}
}


// Loading the data


func TestNewPropertyServices_LoadsAllProperties(t *testing.T) {
	// Call the constructor directly (not the shared service), because loading
	// the file is exactly what this test checks.
	propertyService, err := NewPropertyServices(testDataFilePath)
	if err != nil {
		t.Fatalf("could not load property data: %v", err)
	}

	const expectedPropertyCount = 100
	if len(propertyService.Properties) != expectedPropertyCount {
		t.Errorf("loaded %d properties, want %d", len(propertyService.Properties), expectedPropertyCount)
	}
}

func TestNewPropertyServices_ReturnsErrorForBadPath(t *testing.T) {
	badPaths := map[string]string{
		"empty path":      "",
		"whitespace path": "   ",
		"missing file":    "does/not/exist.json",
	}

	for caseName, badPath := range badPaths {
		t.Run(caseName, func(t *testing.T) {
			propertyService, err := NewPropertyServices(badPath)

			if err == nil {
				t.Error("expected an error, got nil")
			}
			if propertyService != nil {
				t.Errorf("expected a nil service, got %v", propertyService)
			}
		})
	}
}

// With no filters every property is returned. This also proves that all 100
// records have valid "categories" JSON and valid coordinates, because
// GetProperties fails if any single record cannot be transformed.
func TestGetProperties_WithoutFilters_ReturnsEveryProperty(t *testing.T) {
	propertyService := loadPropertyService(t)

	actualResponse, err := propertyService.GetProperties(&utils.PropertyFilter{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	const expectedPropertyCount = 100
	if actualResponse.Result.Count != expectedPropertyCount || len(actualResponse.Result.Items) != expectedPropertyCount {
		t.Errorf("Count/len(Items) = %d/%d, want %d/%d",
			actualResponse.Result.Count, len(actualResponse.Result.Items),
			expectedPropertyCount, expectedPropertyCount)
	}
}

// Transform  (transformProperty: raw Property -> PropertyResponse)


// Checks every mapped field using the first record, "Shinjuku Grand Resort".
func TestTransformProperty_MapsAllFields(t *testing.T) {
	propertyService := loadPropertyService(t)
	rawProperty := findRawPropertyByID(t, propertyService, "BC-1000001")

	transformed, err := transformProperty(rawProperty)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// --- top-level fields ---
	if transformed.ID != "BC-1000001" || transformed.Feed != 11 || transformed.Published != false {
		t.Errorf("top-level fields wrong: id=%s feed=%d published=%v",
			transformed.ID, transformed.Feed, transformed.Published)
	}

	// --- latitude / longitude ---
	// The JSON stores coordinates as [longitude, latitude] = [139.6917, 35.6895],
	// so Lon must come from index 0 and Lat from index 1.
	if transformed.GeoInfo.Lat != 35.6895 {
		t.Errorf("Lat = %v, want 35.6895", transformed.GeoInfo.Lat)
	}
	if transformed.GeoInfo.Lon != 139.6917 {
		t.Errorf("Lon = %v, want 139.6917", transformed.GeoInfo.Lon)
	}

	// --- images: Count must equal the number of image files ---
	expectedImages := []string{"image-1.jpg", "image-2.jpg", "image-3.jpg", "image-4.jpg", "image-5.jpg"}
	if transformed.Property.Image.Count != 5 {
		t.Errorf("Image.Count = %d, want 5", transformed.Property.Image.Count)
	}
	if !slices.Equal(transformed.Property.Image.Images, expectedImages) {
		t.Errorf("Image.Images = %v, want %v", transformed.Property.Image.Images, expectedImages)
	}

	// --- breadcrumbs: parsed from the "categories" JSON string ---
	expectedBreadcrumbs := []models.Breadcrumb{
		{LocationID: "89", Name: "Japan", Type: "country", Slug: "japan", Display: []string{"japan"}},
		{LocationID: "6050001", Name: "Tokyo", Type: "state", Slug: "japan/tokyo", Display: []string{"japan", "tokyo"}},
		{LocationID: "6100001", Name: "Shinjuku", Type: "city", Slug: "japan/tokyo/shinjuku", Display: []string{"japan", "tokyo", "shinjuku"}},
	}
	actualBreadcrumbs := transformed.GeoInfo.Breadcrumbs
	if len(actualBreadcrumbs) != len(expectedBreadcrumbs) {
		t.Fatalf("len(Breadcrumbs) = %d, want %d", len(actualBreadcrumbs), len(expectedBreadcrumbs))
	}
	for index, expected := range expectedBreadcrumbs {
		actual := actualBreadcrumbs[index]
		if actual.LocationID != expected.LocationID || actual.Name != expected.Name ||
			actual.Type != expected.Type || actual.Slug != expected.Slug ||
			!slices.Equal(actual.Display, expected.Display) {
			t.Errorf("Breadcrumbs[%d] = %+v, want %+v", index, actual, expected)
		}
	}

	// --- remaining geo fields (state_abbr is null in the JSON, so it becomes "") ---
	geoInfo := transformed.GeoInfo
	if geoInfo.City != "Shinjuku" || geoInfo.Country != "Japan" || geoInfo.CountryCode != "JP" ||
		geoInfo.State != "Tokyo" || geoInfo.StateAbbr != "" || geoInfo.Name != "Shinjuku, Japan" {
		t.Errorf("GeoInfo wrong: %+v", geoInfo)
	}

	// --- property details ---
	propertyInfo := transformed.Property
	if propertyInfo.Name != "Shinjuku Grand Resort" || propertyInfo.Slug != "shinjuku-grand-resort-0001" ||
		propertyInfo.PropertyType != "Resort" || propertyInfo.Price != 116.69 ||
		propertyInfo.ReviewScore != 4.2 || propertyInfo.StarRating != 5 {
		t.Errorf("PropertyInfo wrong: %+v", propertyInfo)
	}
	expectedAmenities := []string{"Breakfast Included", "Child Friendly", "Pool"}
	if !slices.Equal(propertyInfo.Amenities, expectedAmenities) {
		t.Errorf("Amenities = %v, want %v", propertyInfo.Amenities, expectedAmenities)
	}

	// --- counts ---
	expectedCounts := models.Counts{Bathroom: 2, Bedroom: 2, Reviews: 304, Occupancy: 3}
	if propertyInfo.Counts != expectedCounts {
		t.Errorf("Counts = %+v, want %+v", propertyInfo.Counts, expectedCounts)
	}
}

// The JSON key is "location_id". This test protects against a wrong struct tag
// (models.Property must use `json:"location_id"`), which would leave it empty.
func TestTransformProperty_MapsLocationID(t *testing.T) {
	propertyService := loadPropertyService(t)
	rawProperty := findRawPropertyByID(t, propertyService, "BC-1000001")

	transformed, err := transformProperty(rawProperty)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if transformed.GeoInfo.LocationID != "6100001" {
		t.Errorf("GeoInfo.LocationID = %q, want %q", transformed.GeoInfo.LocationID, "6100001")
	}
}

// Negative longitude (western hemisphere) and a property with exactly one image.
func TestTransformProperty_NegativeLongitudeAndSingleImage(t *testing.T) {
	propertyService := loadPropertyService(t)

	// New York coordinates are [-74.006, 40.7128]: longitude negative, latitude positive.
	newYork, err := transformProperty(findRawPropertyByID(t, propertyService, "HA-2000001"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if newYork.GeoInfo.Lat != 40.7128 || newYork.GeoInfo.Lon != -74.006 {
		t.Errorf("Lat/Lon = %v/%v, want 40.7128/-74.006", newYork.GeoInfo.Lat, newYork.GeoInfo.Lon)
	}
	if newYork.GeoInfo.StateAbbr != "NY" || newYork.Property.Image.Count != 3 {
		t.Errorf("StateAbbr=%q Image.Count=%d, want NY/3", newYork.GeoInfo.StateAbbr, newYork.Property.Image.Count)
	}

	// "Umeda Garden House" has a single image.
	umeda, err := transformProperty(findRawPropertyByID(t, propertyService, "BC-1000004"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if umeda.Property.Image.Count != 1 || len(umeda.Property.Image.Images) != 1 {
		t.Errorf("Image.Count/len(Images) = %d/%d, want 1/1",
			umeda.Property.Image.Count, len(umeda.Property.Image.Images))
	}
}

// The real data is always valid, so the two error branches of transformProperty
// are tested with small hand-made records instead.
func TestTransformProperty_ReturnsErrorForInvalidData(t *testing.T) {
	// newValidProperty returns a minimal record that transforms without error.
	newValidProperty := func() models.Property {
		property := models.Property{
			Categories: `[{"LocationID":"1","Name":"X","Type":"country","Slug":"x","Display":["x"]}]`,
		}
		property.LonLat.Coordinates = []float64{1.0, 2.0}
		return property
	}

	t.Run("categories is not valid JSON", func(t *testing.T) {
		property := newValidProperty()
		property.Categories = "not-json"

		transformed, err := transformProperty(property)
		if err == nil || transformed != nil {
			t.Fatalf("expected error and nil response, got response=%v err=%v", transformed, err)
		}
	})

	t.Run("fewer than two coordinates", func(t *testing.T) {
		property := newValidProperty()
		property.LonLat.Coordinates = []float64{1.0} // longitude only, latitude missing

		transformed, err := transformProperty(property)
		if err == nil || transformed != nil {
			t.Fatalf("expected error and nil response, got response=%v err=%v", transformed, err)
		}
	})
}


// Filters combined with AND  (every filter that is set must match)


// Feed 11 has 25 properties, but only 8 of them are published.
func TestGetProperties_FilterByFeedAndPublished(t *testing.T) {
	propertyService := loadPropertyService(t)

	actualResponse, err := propertyService.GetProperties(&utils.PropertyFilter{
		Feed:      pointerTo(11),
		Published: pointerTo(true),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedIDs := []string{
		"BC-1000004", "BC-1000007", "BC-1000009", "BC-1000012",
		"BC-1000020", "BC-1000021", "BC-1000022", "BC-1000023",
	}
	assertPropertyIDs(t, actualResponse, expectedIDs)
}

// Feed 12 has five villas (07, 14, 17, 18, 25); 07 and 14 are not published.
func TestGetProperties_FilterByFeedPublishedAndPropertyType(t *testing.T) {
	propertyService := loadPropertyService(t)

	actualResponse, err := propertyService.GetProperties(&utils.PropertyFilter{
		Feed:         pointerTo(12),
		Published:    pointerTo(true),
		PropertyType: pointerTo("Villa"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertPropertyIDs(t, actualResponse, []string{"HA-2000017", "HA-2000018", "HA-2000025"})
}

func TestGetProperties_FilterByPriceRange(t *testing.T) {
	testCases := []struct {
		name        string
		filter      utils.PropertyFilter
		expectedIDs []string
	}{
		{
			// Across all 100 properties only two prices are between 25 and 30:
			// 25.67 (Japan) and 26.39 (Germany).
			name:        "min 25 and max 30 across every feed",
			filter:      utils.PropertyFilter{MinPrice: pointerTo(25.0), MaxPrice: pointerTo(30.0)},
			expectedIDs: []string{"BC-1000022", "EP-4000019"},
		},
		{
			// Both bounds are inclusive, so min == max still matches an exact price.
			name:        "min and max both equal to 25.67 (inclusive bounds)",
			filter:      utils.PropertyFilter{MinPrice: pointerTo(25.67), MaxPrice: pointerTo(25.67)},
			expectedIDs: []string{"BC-1000022"},
		},
		{
			// Only two Canadian properties cost between 40 and 50: 43.08 and 48.04.
			name:        "feed 22 combined with a price range",
			filter:      utils.PropertyFilter{Feed: pointerTo(22), MinPrice: pointerTo(40.0), MaxPrice: pointerTo(50.0)},
			expectedIDs: []string{"HG-3000003", "HG-3000008"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			propertyService := loadPropertyService(t)
			filter := testCase.filter // copy so each sub-test owns its filter

			actualResponse, err := propertyService.GetProperties(&filter)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			assertPropertyIDs(t, actualResponse, testCase.expectedIDs)
		})
	}
}


// Amenities filter  (OR: a property needs at least ONE of the listed amenities)


// Base set for every case below: feed 11 + published, i.e.
// BC-1000004, 07, 09, 12, 20, 21, 22, 23. Within that set:
//
//	has "Pool" -> 04, 21
//	has "Gym"  -> 07, 09, 12, 21
//
// so "Pool OR Gym" is the union of both lists (5 properties). If the filter
// wrongly required BOTH amenities, only BC-1000021 would be returned.
func TestGetProperties_FilterByAmenities_MatchesAnyOfTheListedAmenities(t *testing.T) {
	testCases := []struct {
		name        string
		amenities   []string
		expectedIDs []string
	}{
		{
			name:        "only Pool",
			amenities:   []string{"Pool"},
			expectedIDs: []string{"BC-1000004", "BC-1000021"},
		},
		{
			name:        "only Gym",
			amenities:   []string{"Gym"},
			expectedIDs: []string{"BC-1000007", "BC-1000009", "BC-1000012", "BC-1000021"},
		},
		{
			name:        "Pool OR Gym (each property matches at least one)",
			amenities:   []string{"Pool", "Gym"},
			expectedIDs: []string{"BC-1000004", "BC-1000007", "BC-1000009", "BC-1000012", "BC-1000021"},
		},
		{
			name:        "one unknown amenity and one known amenity",
			amenities:   []string{"Helipad", "Pool"},
			expectedIDs: []string{"BC-1000004", "BC-1000021"},
		},
		{
			name:        "spaces around an amenity are ignored",
			amenities:   []string{"Pool", " Gym"},
			expectedIDs: []string{"BC-1000004", "BC-1000007", "BC-1000009", "BC-1000012", "BC-1000021"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			propertyService := loadPropertyService(t)

			actualResponse, err := propertyService.GetProperties(&utils.PropertyFilter{
				Feed:      pointerTo(11),
				Published: pointerTo(true),
				Amenities: pointerTo(testCase.amenities),
			})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			assertPropertyIDs(t, actualResponse, testCase.expectedIDs)
		})
	}
}


// Combined filters  (several AND filters together with the amenities OR filter)


// Narrowing step by step:
//
//	feed 11 + published                 -> 04, 07, 09, 12, 20, 21, 22, 23
//	amenities "Pool" OR "Gym"           -> 04, 07, 09, 12, 21
//	min star rating 3                   -> 07 (5 stars), 09 (4), 12 (3)     [04 and 21 are 1 star]
//	max price 250                       -> 12 ($98.59)                      [07 is $340.53, 09 is $336.24]
//
// Every filter removes at least one property, so this fails if any one of them is ignored.
func TestGetProperties_CombinedFiltersWithAmenities(t *testing.T) {
	propertyService := loadPropertyService(t)

	actualResponse, err := propertyService.GetProperties(&utils.PropertyFilter{
		Feed:          pointerTo(11),
		Published:     pointerTo(true),
		MinStarRating: pointerTo(3),
		MaxPrice:      pointerTo(250.0),
		Amenities:     pointerTo([]string{"Pool", "Gym"}),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertPropertyIDs(t, actualResponse, []string{"BC-1000012"})
}

// The limit is applied AFTER filtering, so the first 3 feed-11 properties come back.
func TestGetProperties_LimitIsAppliedAfterFiltering(t *testing.T) {
	propertyService := loadPropertyService(t)

	actualResponse, err := propertyService.GetProperties(&utils.PropertyFilter{
		Feed:  pointerTo(11),
		Limit: pointerTo(3),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertPropertyIDs(t, actualResponse, []string{"BC-1000001", "BC-1000002", "BC-1000003"})
}

// "min_bedroom" must compare BedroomCount.
// Japanese properties with 5 or more bedrooms: BC-1000003, 13, 14, 19.
func TestGetProperties_FilterByMinimumBedrooms(t *testing.T) {
	propertyService := loadPropertyService(t)

	actualResponse, err := propertyService.GetProperties(&utils.PropertyFilter{
		Feed:       pointerTo(11),
		MinBedroom: pointerTo(5),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertPropertyIDs(t, actualResponse, []string{"BC-1000003", "BC-1000013", "BC-1000014", "BC-1000019"})
}

// Empty result  (no match must give an empty slice, never nil)


func TestGetProperties_NoMatch_ReturnsEmptySliceNotNil(t *testing.T) {
	testCases := []struct {
		name   string
		filter utils.PropertyFilter
	}{
		{
			// The most expensive feed-24 property costs 328.07.
			name:   "feed 24 with a minimum price above every feed-24 property",
			filter: utils.PropertyFilter{Feed: pointerTo(24), MinPrice: pointerTo(330.0)},
		},
		{
			name:   "an amenity that no property has",
			filter: utils.PropertyFilter{Amenities: pointerTo([]string{"Helipad"})},
		},
		{
			// The highest review score in the data is 9.8.
			name:   "a minimum review score above every property",
			filter: utils.PropertyFilter{MinReviewScore: pointerTo(9.9)},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			propertyService := loadPropertyService(t)
			filter := testCase.filter

			actualResponse, err := propertyService.GetProperties(&filter)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// The items must be an empty (non-nil) slice with a count of 0.
			if actualResponse.Result.Items == nil {
				t.Fatal("Result.Items is nil, want an empty non-nil slice")
			}
			if len(actualResponse.Result.Items) != 0 || actualResponse.Result.Count != 0 {
				t.Errorf("len(Items)/Count = %d/%d, want 0/0",
					len(actualResponse.Result.Items), actualResponse.Result.Count)
			}

			// A nil slice would be sent to API clients as "items": null instead of [].
			responseJSON, err := json.Marshal(actualResponse)
			if err != nil {
				t.Fatalf("json.Marshal failed: %v", err)
			}
			if !strings.Contains(string(responseJSON), `"items":[]`) {
				t.Errorf("response JSON = %s, want it to contain \"items\":[]", responseJSON)
			}
		})
	}
}


// Get property by ID


func TestGetPropertyByID_ReturnsPropertyWhenFound(t *testing.T) {
	propertyService := loadPropertyService(t)

	property, err := propertyService.GetPropertyByID("HA-2000001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if property == nil {
		t.Fatal("expected a property, got nil")
	}

	if property.ID != "HA-2000001" || property.Feed != 12 || property.Published != false {
		t.Errorf("top-level fields wrong: id=%s feed=%d published=%v",
			property.ID, property.Feed, property.Published)
	}
	if property.Property.Name != "New York Grand Hotel" || property.Property.Price != 264.49 {
		t.Errorf("Property details wrong: %+v", property.Property)
	}
	if property.GeoInfo.City != "New York" || len(property.GeoInfo.Breadcrumbs) != 3 {
		t.Errorf("GeoInfo wrong: %+v", property.GeoInfo)
	}
}

func TestGetPropertyByID_ReturnsNotFoundErrorWhenMissing(t *testing.T) {
	propertyService := loadPropertyService(t)

	property, err := propertyService.GetPropertyByID("XX-0000000")

	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if property != nil {
		t.Errorf("expected a nil property, got %+v", property)
	}

	// The error must be an APIError carrying HTTP status 404.
	var apiError *utils.APIError
	if !errors.As(err, &apiError) {
		t.Fatalf("error type = %T, want *utils.APIError", err)
	}
	if apiError.Status != 404 {
		t.Errorf("APIError.Status = %d, want 404", apiError.Status)
	}
}