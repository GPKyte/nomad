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
	DateFormat = "2006-01-02"
	UnixMilli  = "05000"
)
const (
	/* Dir must end in trailing slash and files must be typed */
	pathToLocationCache = "cache/locations.json"
	pathToTripCacheDir  = "cache/trips/"
)

type apiResponse struct {
	Trips []*trip `json:"trips"`
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
	var mode = os.FileMode(int(0777))
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
