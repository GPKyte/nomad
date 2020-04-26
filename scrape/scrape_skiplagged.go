package scrape

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

// DateFormat for URL param with skiplagged
const DateFormat = "2006-01-02"

const (
	/* Dir must end in trailing slash and files must be typed */
	pathToLocationCache = "cache/locations.json"
	pathToTripCacheDir  = "cache/trips/"
)

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

func formatURL(from, depart string) string {
	urlargs := map[string]string{
		"from":   from,
		"depart": depart,
		"return": "", /* No Roundtrip searches */
		"format": "v2",
	}
	return fmt.Sprintf("http://skiplagged.com/api/skipsy.php?%s", concatURLArgs(urlargs))
}

/* Helps to return top N airports which will be a focus of focused outbound trips */
func getYourMostFrequentLayoverAirports() []string {
	top := make([]string, 0, 10)
	/* Once statistical methods are prevalent as utils, revisit this method */
	cache := loadCacheOfAirports()

	for i := range top {
		top[i] = cache[rand.Intn(len(cache))]
	}

	return top
}

/* There are roughly 80-120 Airports depending on scope of site */
func loadCacheOfAirports() (airports []string) {
	type airport struct {
		Code string `json:"code"`
	}
	var loaderStruct = make([]*airport, 0, 200)

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
		airports = append(airports, string(loc.Code))
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

/* Idea is to look for deals out of popular pit stops then later find deals to those pitstops and beyond
 * checkOutBoundFromMajorAirports will check fare from selected airports for the next N=5 days from date provided */
func checkOutboundFromMajorAirports(fromThisDay time.Time) {
	for _, airport := range getYourMostFrequentLayoverAirports() {
		for _, date := range getDatesBetween(fromThisDay, fromThisDay.AddDate(0, 0, 5 /*days*/)) {
			visit(formatURL(airport, date))
		}
	}
}

type logger struct{}

func (L *logger) Write(p []byte) (n int, err error) {
	fmt.Print(string(p))
	return len(p), nil
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
	log.Println(url)
	resp, err := http.Get(url)

	if err != nil {
		panic(string(err.Error()))
	}

	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var responseAsJSON apiResponse
	json.Unmarshal(b, &responseAsJSON)

	for _, t := range responseAsJSON.Trips {
		fmt.Printf("%s: $%v", t.City, t.Cost)
	}
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

type apiResponse struct {
	Trips []*trip `json:"trips"`
}

type trip struct {
	City       string `json:"city"`
	Cost       int    `json:"cost"`
	HiddenCity bool   `json:"hiddenCity"`
}
