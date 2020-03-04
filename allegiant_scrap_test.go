package nomad

import (
	"fmt"
	"testing"
)

// var HTMLAllegiantListingContext string = readFile("resources/Allegiant1.html")

func TestRapidTest(t *testing.T) {
	justALinkForNow := "https://flight.deals.allegiant.com/ats/url.aspx?cr=986&wu=95&camp=20200301_SundaySavings&allow_search=1"
	flightDeals, err := ScrapeDealsPage(justALinkForNow)
	fmt.Println(flightDeals)

	if err != nil {
		t.Fatal()
	}
}
