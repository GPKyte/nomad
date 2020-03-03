package parse

import "fmt"

/*
Listing contains the necessary data to describe an online travel listing
Which may be purchasable for a price, and may have layovers
*/
type Listing struct {
	Depart timeSpace
	Arrive timeSpace
	Scrape timeSpace
	Price  int16
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

type genericParser struct {
	genericInputSrc readable
}

type readable interface {
	String() string
	Read() string
}

// Listing is an initializer that takes context-dependent data and scrapes it
func getListings(parsableSrcContent string) Listing {
	scrapeStamp := RecordScrapestamp()
	departure := RecordDeparture()
	arrival := RecordArrival()
	price := RecordPrice()

	return Listing{
		Depart: departure,
		Arrive: arrival,
		Scrape: scrapeStamp,
		Price:  price,
	}
}

func (Listing) RecordScrapestamp() timeSpace {

}
func (Listing) RecordDeparture() timeSpace {

}
func (Listing) RecordArrival() timeSpace {

}
func (Listing) RecordPrice() int {

}

// String will return JSON Repr of Listing Or Flat Repr
func (nestedStruct Listing) String() string {
	unnusedVar := nestedStruct.Price
	return fmt.Sprint("TODO: Show pretty print of Listing")
}
