# Rental Properties API (Beego)

A small REST API built with [Beego v2](https://github.com/beego/beego) that serves rental property data
(hotels, houses, apartments, villas, resorts and hostels) from a JSON file.
It supports fetching a single property by ID and listing properties with a rich set of filters.

The JSON file is read **once at startup** and kept in memory, so no database is needed.

---

## Features

- `GET` a single property by ID
- `GET` a list of properties with filters: feed, published, price range, star rating, review score,
  review count, property type, minimum bedrooms and amenities
- Raw data is transformed into a clean response (`geo_info`, `property`, `counts`, `image`, breadcrumbs)
- Validated query parameters with clear error messages
- Unit tests that run against the real dataset

---

## Prerequisites

| Tool | Version | Notes |
|---|---|---|
| [Go](https://go.dev/dl/) | 1.21 or newer | The code uses the standard `slices` package |
| [Bee tool](https://github.com/beego/bee) | latest | Used to run the server with `bee run` |
| curl | any | Only needed to try the sample requests |

Install the Bee tool once:

```bash
go install github.com/beego/bee/v2@latest
```

Make sure your Go bin folder (`$(go env GOPATH)/bin`) is on your `PATH`, then check:

```bash
bee version
```

---

## Project structure

```
Beego_API/
├── main.go                     # Loads the data file once, then starts the server
├── conf/
│   └── app.conf                # App settings, including the data_file path
├── data/
│   └── rental_properties.json  # The dataset (100 properties)
├── routers/                    # URL -> controller mapping
├── controllers/                # HTTP layer (reads query params, writes JSON)
├── services/
│   ├── property.go             # Business logic: load data, get by ID, list with filters, transform
│   └── property_test.go        # Tests for the service and filters
├── utils/
│   └── filter.go               # Query parsing/validation and the filter logic
└── models/
    ├── source.go               # Shape of the raw JSON data
    └── response.go             # Shape of the API responses
```

---

## Configuration

The path of the data file is read from `conf/app.conf`:

```ini
appname = rental-property-api
httpport = 8080
runmode = dev
copyrequestbody = true
autorender = false
data_file = data/rental_properties.json
```

- `data_file` is relative to the project root (the folder you run `bee run` from).
- `httpport` is the port the server listens on (Beego's default is `8080`).

---

## Run the project

```bash
# 1. Get the code and enter the project folder
cd Beego_API

# 2. Download dependencies
go mod tidy

# 3. Start the server (with auto-reload on code changes)
bee run
```

When it starts successfully you will see a line similar to:

```
http server Running on http://:8080
```

The API is now available at `http://localhost:8080`.

Prefer not to use Bee? This works too:

```bash
go run main.go
```

---

## Run the tests

```bash
go test ./... -v
```

The tests load the real `data/rental_properties.json` file once and check:

| Area | What is verified |
|---|---|
| Transform | All fields are mapped, `image.count` matches the images, `lat`/`lon` come from the right coordinate positions, breadcrumbs are parsed from the `categories` string |
| Filters (AND) | e.g. feed + published, price range, feed + published + property type |
| Amenities (OR) | A property matches if it has **any** of the listed amenities |
| Combined | Several filters plus amenities together |
| Empty result | No match returns an empty list (`[]`), never `null` |
| Get by ID | Found returns the property; missing returns a `404` error |

> You may see this line while testing: `init global config instance failed ... open conf/app.conf`.
> It is printed by Beego when the package is loaded outside the running app and can be ignored.

---

## API reference

Base URL: `http://localhost:8080`

### Get one property

```
GET /api/v1/properties/{id}
```

Returns the property, or a `404` error if the ID does not exist.

### List properties

```
GET /api/v1/properties
```

All filters are optional and are combined with **AND**, except `amenities`, which matches
properties that have **at least one** of the listed amenities (**OR**).
Results keep the order of the data file. If `limit` is left out, every matching property is returned.

| Query parameter | Type | Rules |
|---|---|---|
| `limit` | integer | 1–100, applied after filtering. Invalid values fall back to `20` |
| `feed` | integer | One of `11`, `12`, `22`, `24` |
| `published` | boolean | `true` or `false` |
| `min_price` | number | 0 or more (USD) |
| `max_price` | number | 0 or more, must not be lower than `min_price` |
| `min_star_rating` | integer | 0–5 |
| `min_review_score` | number | 0 or more |
| `min_reviews` | integer | 0 or more (minimum number of reviews) |
| `property_type` | string | `Hotel`, `House`, `Apartment`, `Villa`, `Resort` or `Hostel` (case-sensitive) |
| `min_bedroom` | integer | 1 or more |
| `amenities` | string | Comma-separated list, e.g. `Pool,Gym` (matches any) |

Invalid values are rejected with a descriptive error message.

### Sample response (single property)

```json
{
  "id": "BC-1000001",
  "feed": 11,
  "published": false,
  "geo_info": {
    "breadcrumbs": [
      { "LocationID": "89", "Name": "Japan", "Type": "country", "Slug": "japan", "Display": ["japan"] },
      { "LocationID": "6050001", "Name": "Tokyo", "Type": "state", "Slug": "japan/tokyo", "Display": ["japan", "tokyo"] },
      { "LocationID": "6100001", "Name": "Shinjuku", "Type": "city", "Slug": "japan/tokyo/shinjuku", "Display": ["japan", "tokyo", "shinjuku"] }
    ],
    "city": "Shinjuku",
    "country": "Japan",
    "country_code": "JP",
    "name": "Shinjuku, Japan",
    "location_id": "6100001",
    "lat": 35.6895,
    "lon": 139.6917,
    "state": "Tokyo",
    "state_abbr": ""
  },
  "property": {
    "amenities": ["Breakfast Included", "Child Friendly", "Pool"],
    "name": "Shinjuku Grand Resort",
    "slug": "shinjuku-grand-resort-0001",
    "property_type": "Resort",
    "price": 116.69,
    "review_score": 4.2,
    "star_rating": 5,
    "counts": { "bathroom": 2, "bedroom": 2, "reviews": 304, "occupancy": 3 },
    "image": {
      "count": 5,
      "images": ["image-1.jpg", "image-2.jpg", "image-3.jpg", "image-4.jpg", "image-5.jpg"]
    }
  }
}
```

### Sample response (list)

```json
{
  "result": {
    "count": 2,
    "items": [ { "id": "BC-1000022", "...": "..." }, { "id": "EP-4000019", "...": "..." } ]
  }
}
```

---

## Sample curl requests

Set the base URL once (adjust the port if you changed `httpport`):

```bash
BASE_URL=http://localhost:8080
```

Add `| jq` at the end of any command for pretty output if you have [jq](https://jqlang.github.io/jq/) installed.

**1. Get a single property by ID**

```bash
curl -s "$BASE_URL/api/v1/properties/HA-2000001"
```

**2. Property ID that does not exist (returns `404`)**

```bash
curl -i "$BASE_URL/api/v1/properties/XX-0000000"
```

**3. List all properties**

```bash
curl -s "$BASE_URL/api/v1/properties"
```

**4. Limit the number of results**

```bash
curl -s "$BASE_URL/api/v1/properties?limit=5"
```

**5. Feed + published** (feed 11 has 8 published properties)

```bash
curl -s "$BASE_URL/api/v1/properties?feed=11&published=true"
```

**6. Price range** (returns `BC-1000022` and `EP-4000019`)

```bash
curl -s "$BASE_URL/api/v1/properties?min_price=25&max_price=30"
```

**7. Property type** (published villas in feed 12: `HA-2000017`, `HA-2000018`, `HA-2000025`)

```bash
curl -s "$BASE_URL/api/v1/properties?feed=12&published=true&property_type=Villa"
```

**8. Star rating and review filters**

```bash
curl -s "$BASE_URL/api/v1/properties?min_star_rating=4&min_review_score=8&min_reviews=100"
```

**9. Minimum bedrooms** (feed 11 with 5 or more bedrooms: 4 results)

```bash
curl -s "$BASE_URL/api/v1/properties?feed=11&min_bedroom=5"
```

**10. Amenities, OR match** (published Japan properties with a Pool **or** a Gym: 5 results)

```bash
curl -s "$BASE_URL/api/v1/properties?feed=11&published=true&amenities=Pool,Gym"
```

**11. Amenity with a space in its name**

```bash
curl -s "$BASE_URL/api/v1/properties?amenities=Air%20Conditioner,Hot%20Tub"
```

**12. Everything combined** (returns only `BC-1000012`)

```bash
curl -s "$BASE_URL/api/v1/properties?feed=11&published=true&min_star_rating=3&max_price=250&amenities=Pool,Gym"
```

**13. No match: returns an empty list, not `null`**

```bash
curl -s "$BASE_URL/api/v1/properties?amenities=Helipad"
# {"result":{"count":0,"items":[]}}
```

**14. Invalid parameters are rejected**

```bash
curl -i "$BASE_URL/api/v1/properties?feed=99"
curl -i "$BASE_URL/api/v1/properties?min_price=abc"
curl -i "$BASE_URL/api/v1/properties?min_price=200&max_price=100"
curl -i "$BASE_URL/api/v1/properties?property_type=Castle"
```

---

## Troubleshooting

| Problem | Fix |
|---|---|
| `bee: command not found` | Run `go install github.com/beego/bee/v2@latest` and add `$(go env GOPATH)/bin` to your `PATH` |
| Server cannot find the data file | Check `data_file` in `conf/app.conf` and run `bee run` from the project root |
| `address already in use` | Another program uses the port. Stop it or change `httpport` in `conf/app.conf` |
| Tests cannot load the data | Make sure `data/rental_properties.json` exists (the tests read `../data/rental_properties.json` from the `services` folder) |