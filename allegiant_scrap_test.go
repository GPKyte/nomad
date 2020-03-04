package nomad

import (
	"fmt"
	"testing"
)

var HTMLAllegiantListingContext string = readFile("resources/Allegiant1.html")

func TestRapidTest(t *testing.T) {
	justALinkForNow := "https://flight.deals.allegiant.com/ats/url.aspx?cr=986&wu=95&camp=20200301_SundaySavings&allow_search=1"
	flightDeals, err := scrapeDealsPage(justALinkForNow)
	fmt.Println(flightDeals...)
}

func TestGetAllDailyDealListings(t *testing.T) {

}

func TestGetADailyDealListing(t *testing.T) {
	siteContent := HTMLAllegiantListingContext
	var LP LPAllegiant
	LP.URL = ""

	L * Listing, err = LP.pointToListingFromThis(siteContent)
	if err != nil {
		t.Fail(err)
	}
}

func TestStringMethod(t *testing.T) {
	L := Listing{}
	goal := ""
	got := string(L)
}

func TestComparison(t *testing.T) {
	t.Fail
}
