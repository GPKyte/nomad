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

type apiResponse struct {
	Trips []*trip `json:"trips"`
}

type skippyResponse apiResponse

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

func parseFromAPIv3Search(results []byte) {
	hmm := new(topNest)
	if err := json.Unmarshal(results, hmm); err != nil {
		panic(err.Error())
	}

	for _, fare := range hmm.Plans["outbound"] {

		var budget int64 = 100 /*cents*/ * 2000 /*USDollars ($)*/
		if fare.OneWayCost > budget {
			continue // skip this price, but be warned that this generates incomplete data
		}
		var id string = fare.Flight      /* Site-generated ID for flight Itinerary makes convenient lookup across nested structures */
		var cost int64 = fare.OneWayCost /* Round Trip Costs are not being considered at this time */
		var flight nestedFlight = hmm.Flights[id]
		var format string = FullDateTimeFormat

		var departTime, arriveTime time.Time
		var departLoc, arriveLoc string /* TODO: Enforce Location-type name lookup by code and fill-in data */
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

			/* TODO: Make real Listing and send it back on a channel */
			fmt.Printf("\nThis is a Listing:\n%s:%s -> %s:%s\nPrice:%v\n", departLoc, departTime, arriveLoc, arriveTime, cost)
		}
	}
}

type trip struct {
	City       string `json:"city"`
	Cost       int    `json:"cost"`
	HiddenCity bool   `json:"hiddenCity"`
}

type logger struct{}

// Determine the impact of Booking ahead of Departure date by N days
// Find patterns of "best" for N
// This only collects data and tags it for this purpose
func checkWhenTheEarlyBirdRises() {
	/* Define a handful of static locations to use as reference points */
	/* It is important to also record airline information when present */
	/* Look as soon as next day and up to 6mo */
	/* This means # Requests = X routes * (30*6 dates) */
	/* That's about 500 Requests, this should happen infrequently, like each month */
	from := getLocationsByCode("CLE", "CVG", "PIT")
	to := getLocationsByCode("DEN", "PIE", "LAX")
	response := make(chan string)

	for w := range from {
		/* As we learn more, be more selective with dates checked */
		for _, date := range getDatesForNextN(60 /* days */) {
			url := formatURL(from[w], to[w], date)
			response <- fakeVisit(url)

			waitVariableTime()
		}
	}

	/* Collect and interpret results as they arrive */
	processFareData := func() {
		var data string = <-response
		fmt.Printf(data)
	}

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
		"format": "v2",
		"_":      currentTime,
	}
	var endpoint string

	if to.Code != "" {
		urlargs["to"] = to.Code
		endpoint = "search.php"
	} else {
		endpoint = "skipsy.php"
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

/* Overwrite (should be infrequent) the known cache of airport locations
 * Format is NOT Validated, be careful to follow conventions or be surprised
 * Intended for use with the data from skiplagged API in the future */
func updateCacheOfAirports(withNewJSON []byte) error {
	var mode = os.FileMode(int(0444))
	err := ioutil.WriteFile(pathToLocationCache, withNewJSON, mode)
	return err
}

func (L *logger) Write(p []byte) (n int, err error) {
	fmt.Print(string(p))
	return len(p), nil
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

func getDatesForNextN(days int) []string {
	now := time.Now()
	var years, monthsAnd, daysAhead int
	daysAhead = days
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
func visit(url string) {
	/* TODO
	 * Cache either parsed or original data using a timestamp filename
	 * URL holds some meta data we might want to use, but avoid this
	 * In every case we want what?
	 * To receive Listings as they become available?
	 * To download the requested page to memory&disk?
	 * To return Trips?
	 * ...
	 */
	log.Println("Visiting: ", url)
	resp, err := http.Get(url)

	if err != nil {
		panic(string(err.Error()))
	}

	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var responseAsJSON apiResponse
	json.Unmarshal(b, &responseAsJSON)
}

func waitVariableTime() {
	const minimumWait = 2 /*seconds*/
	seconds := time.Duration(minimumWait+rand.Intn(20)) * time.Second
	time.Sleep(seconds)
}
