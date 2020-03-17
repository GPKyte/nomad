package nomad

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

// var HTMLAllegiantListingContext string = readFile("resources/Allegiant1.html")
func check(err error) {
	if err != nil {
		panic(err)
	}
}
func testRapidTest(t *testing.T) {
	justALinkForNow := "https://flight.deals.allegiant.com/ats/url.aspx?cr=986&wu=95&camp=20200301_SundaySavings&allow_search=1"
	flightDeals, err := ScrapeDealsPage(justALinkForNow)
	fmt.Println(flightDeals)

	check(err)
}

func TestSelectFlightCards(t *testing.T) {
	f, err := os.Open("resources/Allegiant.html")
	defer f.Close()

	doc, err := goquery.NewDocumentFromReader(io.Reader(f))
	t.Log(err)
	check(err)
	fmt.Print(doc)
}
