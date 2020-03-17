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

func recordCurrentTimeSpace(location string) timeSpace {
	return timeSpace{
		DateTime: time.Now().UTC(),
		Location: location,
	}
}

// String will return JSON Repr of Listing Or Flat Repr
func (L *Listing) String() string {
	return String(L.json())
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
	var (
		now       = time.Now().UTC()
		before    = time.Now().UTC()
		after     = time.Now().UTC()
		departLoc = chooseLoc(5)
		arriveLoc = chooseLoc(9)
		url       = "https://random.local"
	)

	return Listing{
		Price:  50,
		Depart: DateTime{before, departLoc},
		Arrive: DateTime{after, arriveLoc},
		Scrape: DateTime{now, url},
	}
}

func chooseLoc(index int) string {
	locations := []string{
		"A",
		"B",
		"C",
		"D",
		"E",
		"F",
		"G",
		"H",
		"I",
		"J",
		"K",
	}

	return locations[index]
}
