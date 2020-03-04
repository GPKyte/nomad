package nomad

import (
	"fmt"
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

func newListing() Listing {
	return Listing{}
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

func recordCurrentTimeSpace(location string) timeSpace {
	return timeSpace{
		DateTime: time.Now().UTF8(),
		Location: location,
	}
}

// String will return JSON Repr of Listing Or Flat Repr
func (L *Listing) String() string {
	return fmt.Sprintf("From: %s\nTo: %s\nCost: %v\nStamp: %s\n.............................", L.Depart, L.Arrive, L.Price, L.Scrape)
}
