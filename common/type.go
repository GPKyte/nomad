package common

import (
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
	Price      int16
	srcContent string
}

type trip struct {
	startArea location
	endArea   location
	startTime time.Time
	endTime   time.Time
	price     int
}
type timeSpace struct {
	DateTime time.Time
	Location string /* standardize to lowercase, UTF8 default */
}
type readable interface {
	String() string
	Read() string
}
type location struct {
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

func recordCurrentTimeSpace(location string) timeSpace {
	return timeSpace{
		DateTime: time.Now().UTC(),
		Location: location,
	}
}

// String will return JSON Repr of Listing Or Flat Repr
func (L *Listing) String() string {
	var gift string

	if weWantDefault := true; weWantDefault {
		gift = string(makeJSON(L.json()))
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

func (L *Listing) json() json {

	var jsonRepr = json{
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

func makeJSON(srcJSONstruct json) []byte {
	return make([]byte, 0, 0)
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
