package nomad

import (
	"io"
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

type timeSpace struct {
	DateTime time.Time
	Location string /* standardized to lowercase, UTF8 default */
}

type readable interface {
	String() string
	Read() string
}

type location struct {
	Name  string
	State string
	Code  string
}

type json map[string]interface{}

func newInvalidListing() Listing {
	return Listing{}
}

type ValidDataObj interface {
	isNil()
	isValid()
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

func makeScrapeStamp(srcURL string) {
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
		gift = makeJSONString(L.json())
	} else {
		gift = L.csv()
	}
	return gift
}

func (L *Listing) csv() string {
	var csvRow = []string{
		L.Price,
		L.Depart.DateTime,
		L.Depart.Location,
		L.Arrive.DateTime,
		L.Arrive.Location,
		L.Scrape.DateTime,
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

func makeJSONString(srcJSON json) {

}

func join(delim rune, csvRow []string) {

}

// NewListingsFromJSON is used to import Listing data that has been exported per standard
func NewListingsFromJSON(srcJSON io.Reader) []Listing {
	// TODO
}

// NewListingsFromCSV is used to import Listing data that has been exported per standard
func NewListingsFromCSV(srcCSV io.Reader) []Listing {
	// TODO
}

func NewListingRand() Listing {
	const nanoConvRate = 10 ^ 9
	var (
		now       = time.Now().UTC()
		before    = time.Now().UTC()
		twentyMin = time.Duration(nanoConvRate * (60 * 20)) // 60 sec/min * 20 min
		after     = time.Now().UTC().Add(twentyMin)
		departLoc = chooseLoc(5)
		arriveLoc = chooseLoc(9)
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
func Value(L *Listing) int {
	return L.Price()
}

func values(collection ...interface{ Value() }) []int {
	var values = make([]int, len(collection))

	for v := range collection {
		values = append(values, Value(v))
	}
}

func max(collection []interface{ Value() }) int {
	top_value = Value(collection[0])

	for val := range values(collection) {

		if val > top_value {
			top_value = val
		}
	}

	return top_value
}
