package parse

import "fmt"

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
	Date      string
	TimeOfDay int    /* Value between {0 ... (24 * 60)} */
	Location  string /* standardized to lowercase, UTF8 default */
}

// ListingParser will form the template for site-specific parsers
type ListingParser interface {
	RecordScrapestamp() timeSpace
	RecordDeparture() timeSpace
	RecordArrival() timeSpace
	RecordPrice() int
}

type readable interface {
	String() string
	Read() string
}

func newListing() Listing {
	return Listing{}
}

// Listing is an initializer that takes context-dependent data and scrapes it
func (LP *ListingParser) pointToListingFromThis(parsableSrcContent string) *Listing {
	var L Listing = newListing()

	L.srcContent = parsableSrcContent
	LP.RecordScrapestamp(L)
	LP.RecordDeparture(L)
	LP.RecordArrival(L)
	LP.RecordPrice(L)

	return &L
}

type ValidDataObj interface {
	isNil()
	isValid()
}

func WhetherNilOrNot(obj ValidDataObj) bool {
	/* TODO: Consider diff between dot notation and method as is. Recall readability versus convenience */
	return obj.isNil()
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

// RecordScrapestamp returns a default timeSpace because this is the default Listing class
func (LP *ListingParser) RecordScrapestamp(L *Listing) err {
	L.ScrapeStamp = timeSpace{}
	return nil
}

// RecordDeparture returns a default timeSpace because this is the default Listing class
func (LP *ListingParser) RecordDeparture(L *Listing) err {
	L.Depart = timeSpace{}
	return nil
}

// RecordArrival returns a default timeSpace because this is the default Listing class
func (LP *ListingParser) RecordArrival(L *Listing) err {
	L.Arrive = timeSpace{}
	return nil
}

// RecordPrice returns a default timeSpace because this is the default Listing class
func (LP *ListingParser) RecordPrice(L *Listing) err {
	L.Price = 0
	return nil
}

// String will return JSON Repr of Listing Or Flat Repr
func (nestedStruct Listing) String() string {
	return fmt.Sprint("TODO: Show pretty print of Listing")
}
