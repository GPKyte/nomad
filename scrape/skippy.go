package scrape

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// DateFormat for URL param with skiplagged
const (
	DateFormat         = "2006-01-02"
	FullDateTimeFormat = "2006-01-02T15:04:05-07:00"
)
const (
	/* Dir must end in trailing slash and files must be typed */
	pathToLocationCache = "cache/locations.json"
	pathToRawCache      = "cache/raw/"
	pathToTripCacheDir  = "cache/trips/"
)

// SkippyScraper is the entry point to this program and stores the meta-data needed for externally-facing methods
// Optionally Grow out meta data as needed
type SkippyScraper struct {
	storage chan Listing
	name    string
}

// NewSkippy will provide a consistent destination across scrape jobs, supply the channel to send Listings through
func NewSkippy(storage chan Listing) *SkippyScraper {
	S := new(SkippyScraper)
	S.storage = storage
	S.name = "Skiplagged.com v3"

	return S
}

// AToBNearDate uses its own discretion to find airfare from A to B in the days after this date
func (S *SkippyScraper) AToBNearDate(A, B string, date time.Time) {
	S.AToBDuring(A, B, date, date.AddDate(0, 0, 10 /*days ahead*/))
}

// AToBDuring searches every day in the given time frame
func (S *SkippyScraper) AToBDuring(A, B string, start, end time.Time) {
	var dates []string = getDatesBetween(start, end)

	for _, D := range dates {
		S.searchAndScrape(A, B, D)
	}
}

// searchAndScrape is the essence of each search and requires preformatting the parameters
// "from" and "to" are three-letter airport codes
// date follow the DateFormat constant defined in this file
func (S *SkippyScraper) searchAndScrape(A, B, date string) {
	from, to := SanitizeLocations(A, B)

	var url string = formatURL(from, to, date)
	var responseBody []byte = visit(url)

	S.scrape(responseBody)
	waitVariableTime()
}

// AToAnywhereSoon may check standard and non-standard listings for deals and this needs to be accounted for if using the data
// This is because Skiplagged.com offers two main APIs and this one will trigger a call to Skippy,
// which aggregates data by mincost to any Location, instead of per Flight to one Location
func (S *SkippyScraper) AToAnywhereSoon(A string) {
	// Need to specifically handle v2 Skippy API response here instead of reusing abstraction :(
	S.AToBNearDate(A, "", time.Now())
}

/* In order to parse version 3 api/search? response
 * we define several nested types for convenient Unmarshalling */
type topNest struct {
	Flights map[string]nestedFlight `json:"flights"`
	Plans   map[string][]nestedFare `json:"itineraries"`
}
type nestedFlight struct {
	Segments []nestedSegment `json:"segments"`
	Duration int64           `json:"duration"`
	Count    int             `json:"count"`
}
type nestedSegment struct {
	Airline      string          `json:"airline"`
	FlightNumber int             `json:"flight_number"`
	Departure    nestedTimeSpace `json:"departure"`
	Arrival      nestedTimeSpace `json:"arrival"`
	Duration     int             `json:"duration"`
}
type nestedTimeSpace struct {
	Time    string `json:"time"`
	Airport string `json:"airport"`
}
type nestedPlan struct {
	Outbound []nestedFare `json:"outbound"`
	Inbound  []nestedFare `json:"inbound"`
}
type nestedFare struct {
	Flight        string `json:"flight"`
	OneWayCost    int64  `json:"one_way_price"`
	RoundTripCost int64  `json:"min_round_trip_price"`
}

func (S *SkippyScraper) scrape(results []byte) {
	var apiSearchV3Response = new(topNest)

	if err := json.Unmarshal(results, apiSearchV3Response); err != nil {
		log.Panicln(err.Error())
	}
	for _, fare := range apiSearchV3Response.Plans["outbound"] {

		var budget int64 = 100 /*cents*/ * 2000 /*USDollars ($)*/
		if fare.OneWayCost > budget {
			continue // skip this price, but be warned that this generates incomplete data
		}
		var id string = fare.Flight      /* Site-generated ID for flight Itinerary makes convenient lookup across nested structures */
		var cost int64 = fare.OneWayCost /* Round Trip Costs are not being considered at this time */
		var flight nestedFlight = apiSearchV3Response.Flights[id]
		var format string = FullDateTimeFormat

		var departTime, arriveTime time.Time
		var departLoc, arriveLoc string
		var err error

		/* TODO: Add legs to a single Listing, instead of making individual Listings */
		for _, leg := range flight.Segments {
			departTime, err = time.Parse(format, leg.Departure.Time)
			departLoc = leg.Departure.Airport
			if err != nil {
				panic(err.Error())
			}
			arriveTime, err = time.Parse(format, leg.Arrival.Time)
			arriveLoc = leg.Arrival.Airport
			if err != nil {
				panic(err.Error())
			}

			S.storage <- newSkipLaggedListing(departLoc, departTime, arriveLoc, arriveTime, cost)
		}
	}
	recover()
}

type trip struct {
	City       string `json:"city"`
	Cost       int    `json:"cost"`
	HiddenCity bool   `json:"hiddenCity"`
}

func cacheRaw(response []byte, name string) error {
	var mode = os.FileMode(int(0444))
	var fullpath = fmt.Sprintf("%s%s", pathToRawCache, name)

	err := ioutil.WriteFile(fullpath, response, mode)
	return err
}

/* Idea is to look for deals out of popular pit stops then later find deals to those pitstops and beyond
 * checkOutBoundFromMajorAirports will check fare from selected airports for the next N=5 days from date provided */
func checkOutboundFromMajorAirports(fromThisDay time.Time) {
	for _, airport := range getYourMostFrequentLayoverAirports() {
		for _, date := range getDatesBetween(fromThisDay, fromThisDay.AddDate(0, 0, 5 /*days*/)) {
			visit(formatURL(airport, Location{}, date))
		}
	}
}
func chooseDate() string     { /* want to use time.Time to format std dates like this */ return "2020-05-07" }
func chooseLocation() string { return "CVG" }
func concatURLArgs(kv map[string]string) string {
	var cat []string

	for k, v := range kv {
		/* Could consider santizing arguments here, but leaving this chore for later */
		cat = append(cat, fmt.Sprintf("%s=%s", k, v))
	}

	return strings.Join(cat, "&")
}

/* TODO: Deprecate because this function may be growing without bound and is basic enough to duplicate code */
func formatURL(from, to Location, prettyDate string) string {
	// Example: https://skiplagged.com/api/search.php?from=CLE&to=SVQ&depart=2020-05-16&return=&poll=true&format=v3&_=1588452120703
	currentTime := strconv.FormatInt(time.Now().Unix()*1000, 10)

	if len("1588452120703") != len(currentTime) {
		panic(fmt.Sprintf("Wrong time format! Should be in milliseconds (Unix), %v", currentTime))
	}

	urlargs := map[string]string{
		"from":   from.Code,
		"depart": prettyDate,
		"return": "", /* No Roundtrip searches */
		"_":      currentTime,
	}
	var endpoint string

	if to.Code != "" {
		urlargs["to"] = to.Code
		urlargs["format"] = "v3"
		endpoint = "search.php"
	} else {
		endpoint = "skipsy.php"
		urlargs["format"] = "v2"
	}
	return fmt.Sprintf("http://skiplagged.com/api/%s?%s", endpoint, concatURLArgs(urlargs))
}

/* There are roughly 80-120 Airports depending on scope of site */
func loadCacheOfAirports() (airports []Location) {
	var loaderStruct = make([]Location, 0, 200)

	raw, err := ioutil.ReadFile(pathToLocationCache)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = json.Unmarshal(raw, &loaderStruct)
	if err != nil {
		log.Fatal(err.Error())
	}

	/* Reduce data to just the airport code before sending back []string */
	for _, loc := range loaderStruct {
		airports = append(airports, loc)
	}

	return airports
}

func newSkipLaggedListing(departLoc string, departTime time.Time, arriveLoc string, arriveTime time.Time, cost int64) Listing {
	return Listing{
		Depart: timeSpace{Location: departLoc, DateTime: departTime},
		Arrive: timeSpace{Location: arriveLoc, DateTime: arriveTime},
		Scrape: makeScrapeStamp("SkipLagged.com"),
		Price:  money(cost),
	}
}

/* Overwrite (should be infrequent) the known cache of airport locations
 * Format is NOT Validated, be careful to follow conventions or be surprised
 * Intended for use with the data from skiplagged API in the future */
func updateCacheOfAirports(withNewJSON []byte) error {
	var mode = os.FileMode(int(0444))
	err := ioutil.WriteFile(pathToLocationCache, withNewJSON, mode)
	return err
}

/* Helps to return top N airports which will be a focus of focused outbound trips */
func getYourMostFrequentLayoverAirports() []Location {
	top := make([]Location, 0, 10)
	/* Once statistical methods are prevalent as utils, revisit this method */
	cache := loadCacheOfAirports()

	for i := range top {
		top[i] = cache[rand.Intn(len(cache))]
	}

	return top
}

func getDatesForNext(nDays int) []string {
	now := time.Now()
	var years, monthsAnd, daysAhead int
	daysAhead = nDays
	andThen := now.AddDate(years, monthsAnd, daysAhead)

	return getDatesBetween(now, andThen)
}

func getDatesBetween(then, andNow time.Time) []string {
	const oneDayAtATime = 1

	datesThat := make([]string, 0, 50) /* Setting capacity low because it is expected to be low and can expand peacefully thanks to std Go */
	shouldBeBefore, this := then, andNow

	allTheTimeThat := func(a, b time.Time) bool { return a.Before(b) }

	for allTheTimeThat(shouldBeBefore, this) {
		datesThat = append(datesThat, shouldBeBefore.Format(DateFormat))
		shouldBeBefore = shouldBeBefore.AddDate(0, 0, oneDayAtATime)
	}
	return datesThat
}

func fakeVisit(url string) string {
	return fmt.Sprintf(url)
}
func visit(url string) []byte {
	var resp *http.Response
	var body []byte
	var err error

	fmt.Println("Visiting: ", url)

	if resp, err = http.Get(url); err != nil {
		panic(err.Error())
	}

	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		panic(err.Error())
	}
	/* TODO: Cache response for future usage? */
	defer resp.Body.Close()

	if len(body) == 0 {
		panic("No Response received from visiting " + url)
	}
	return body
}

// According to robots.txt in 2020.4.1 the request per second rate limited to 1 second between bot calls
// We can respect that and we're in no rush so add more
func waitVariableTime() {
	const minimumWait = 2 /*seconds*/
	seconds := time.Duration(minimumWait+rand.Intn(20)) * time.Second
	time.Sleep(seconds)
}
