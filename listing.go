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
func pointToListingFromThis(parsableSrcContent string) *Listing {
	var L Listing = newListing()

	L.srcContent = parsableSrcContent
	L.RecordScrapestamp()
	L.RecordDeparture()
	L.RecordArrival()
	L.RecordPrice()

	return &L
}

type ValidDataObj interface {
	func isNil()
	/*
	func isValid()
	func is?
	*/
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

func (L *Listing) isValid() bool {
	return false
}

/* TODO: unexport default methods once testing confirms okay. fmt implies timeSpace should be exported if Record* is */

// RecordScrapestamp returns a default timeSpace because this is the default Listing class
func (Listing) RecordScrapestamp() timeSpace {
	return timeSpace{}
}

// RecordDeparture returns a default timeSpace because this is the default Listing class
func (Listing) RecordDeparture() timeSpace {
	return timeSpace{}
}

// RecordArrival returns a default timeSpace because this is the default Listing class
func (Listing) RecordArrival() timeSpace {
	return timeSpace{}
}

// RecordPrice returns a default timeSpace because this is the default Listing class
func (Listing) RecordPrice() int {
	return 0
}

// String will return JSON Repr of Listing Or Flat Repr
func (nestedStruct Listing) String() string {
	unnusedVar := nestedStruct.Price
	return fmt.Sprint("TODO: Show pretty print of Listing")
}
