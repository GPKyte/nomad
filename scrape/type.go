package scrape

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

/*
Listing contains the necessary data to describe an online travel listing
Which may be purchasable for a price, and may have layovers
*/
type Listing struct {
	Depart     timeSpace
	Arrive     timeSpace
	Scrape     timeSpace
	Price      money
	srcContent string
}

// Trip connects two TimeAndPlace for some cost
type Trip struct {
	depart TimeAndPlace
	arrive TimeAndPlace
	cost   money
}

type money int32

type timeSpace struct {
	DateTime time.Time
	Location string /* standardize to lowercase, UTF8 default */
}

// TimeAndPlace is used as a special node in a normal graph impl */
type TimeAndPlace struct { // Absolutely a duplicate of timeSpace
	T time.Time
	P string
}

func (T *TimeAndPlace) String() string {
	return fmt.Sprintf("%s @%s", T.P, T.T)
}

type Location struct {
	code string
	name string
}

func newInvalidListing() Listing {
	return Listing{}
}

func (T *timeSpace) isNil() bool {
	return false /* TODO: Check revamped timeSpace struct for nils or bad formatting */
}

func (T *timeSpace) isValid() bool {
	return false
}

func (L *Listing) isNil() (whetherAnyNilData bool) {
	whetherAnyNilData = (L.Price == 0 || L.Depart.isNil() || L.Arrive.isNil() || L.Scrape.isNil())

	return
}

/* TODO: change all isValid to isInvalid for consistent negative testing */
func (L *Listing) isValid() bool {
	return false
}

/* TODO: unexport default methods once testing confirms okay. fmt implies timeSpace should be exported if Record* is */

func makeScrapeStamp(srcURL string) timeSpace {
	stamp := recordCurrentTimeSpace(srcURL)
	return stamp
}

func recordCurrentTimeSpace(Location string) timeSpace {
	return timeSpace{
		DateTime: time.Now().UTC(),
		Location: Location,
	}
}

// String will return JSON Representation of listing
func (L *Listing) String() string {
	var gift string

	if weWantDefault := true; weWantDefault {
		gift = string(exportListingAsJSON(L))
	} else {
		gift = L.csv()
	}
	return gift
}

func (L *Listing) csv() string {
	var csvRow = []string{
		string(L.Price),
		L.Depart.DateTime.String(),
		L.Depart.Location,
		L.Arrive.DateTime.String(),
		L.Arrive.Location,
		L.Scrape.DateTime.String(),
		L.Scrape.Location,
	}

	return join(',', csvRow)
}

/* DEPRECATED in favor of json.Marshall() but keep for custom Marshaller */
func (L *Listing) json() map[string]interface{} {

	var jsonRepr = map[string]interface{}{
		"Price":      L.Price,
		"DepartTime": L.Depart.DateTime,
		"DepartLoc":  L.Depart.Location,
		"ArriveTime": L.Arrive.DateTime,
		"ArriveLoc":  L.Arrive.Location,
		"ScrapeTime": L.Scrape.DateTime,
		"ScrapeURL":  L.Scrape.Location,
	}

	return jsonRepr
}

/* Wrapper for the Marshalling function in the json library */
func exportListingAsJSON(L *Listing) []byte {
	b, err := json.MarshalIndent(L, "", "\t")
	if err != nil {
		panic(err)
	}
	return b
}

func importJSONAsListing(data []byte) Listing {
	var L = newInvalidListing()
	if err := json.Unmarshal(data, &L); err != nil {
		panic(err)
	}
	return L
}

func filter(input []byte, filterPredicate func(byte) bool) []byte {
	var out = make([]byte, len(input))

	for _, b := range input {
		if filterPredicate(b) {
			out = append(out, b)
		}
	}
	return out
}

func importJSONAsListings(data []byte) []Listing {
	var probableLineCount int = len(filter(data, func(b byte) bool {
		if b == byte('\n') {
			return true
		}
		return false
	}))
	var L = make([]Listing, 2*probableLineCount)

	/* Just assuming magic here for the time being */
	err := json.Unmarshal(data, &L)
	if err != nil {
		panic(err.Error())
	}

	return L
}

func join(delim rune, row []string) string {
	return strings.Join(row, string(delim))
}

// NewListingsFromJSON is used to import Listing data that has been exported per standard
func NewListingsFromJSON(srcJSON io.Reader) []Listing {
	// TODO
	return make([]Listing, 0, 0)
}

// NewListingsFromCSV is used to import Listing data that has been exported per standard
func NewListingsFromCSV(srcCSV io.Reader) []Listing {
	// TODO
	return make([]Listing, 0, 0)
}

func newListingRand() Listing {
	const nanoConvRate = 10 ^ 9
	var (
		now       = time.Now().UTC()
		before    = time.Now().UTC()
		twentyMin = time.Duration(nanoConvRate * (60 * 20)) // 60 sec/min * 20 min
		after     = time.Now().UTC().Add(twentyMin)
		departLoc = "A"
		arriveLoc = "B"
		url       = "https://random.local"
	)

	return Listing{
		Price:  50,
		Depart: timeSpace{before, departLoc},
		Arrive: timeSpace{after, arriveLoc},
		Scrape: timeSpace{now, url},
	}
}

// Value returns that which should be used in comparisons, this is the Listing price
func (L *Listing) Value() int {
	return int(L.Price)
}

func values(collection ...interface{ Value() int }) []int {
	var values = make([]int, len(collection))

	for _, v := range collection {
		next := v.Value()
		values = append(values, next)
	}

	return values
}

func max(collection []interface{ Value() int }) int {
	var max int = collection[0].Value()

	for _, val := range values(collection...) {

		if val > max {
			max = val
		}
	}

	return max
}

// Before is a wrapper around time package's Before method. Useful for sorting/comparison
func (T *TimeAndPlace) Before(otherTime TimeAndPlace) bool {
	return T.T.Before(otherTime.T)
}

func getLocationsByCode(codes ...string) []Location {
	var rez = make([]Location, len(codes))
	// Eventually will have mechanism for looking up known locations against a DB which has meta data on them
	for _, each := range codes {
		// For now, just give the caller what they want
		rez = append(rez, Location{name: "Unprovided", code: each})
	}

	return rez
}
