package nomad

import (
	"fmt"
	"io"

	"github.com/PuerkitoBio/goquery"

	"github.com/gocolly/colly"
)

// BasicHTMLError is a simple std for error reporting
type BasicHTMLError struct {
	context string
	want    string
	got     string
	html    *colly.HTMLElement
	DOM     *goquery.Selection
}

func rapidTest(c *colly.Collector) error {
	c.Wait(800)

	return nil
}

func (e *BasicHTMLError) Error() string {
	verbose := bool(false)
	var contextLine string

	if len(e.context) > 0 && verbose {
		contextLine = "Context: " + e.context + "\n"
	}
	return fmt.Sprintf("Probably Parsing related;\nWant: %s\nGot: %s\n%s", e.want, e.got, contextLine)
}

func makeDefaultCollector() *colly.Collector {
	return colly.NewCollector(
		colly.AllowedDomains("allegiant.com"),
		colly.MaxDepth(2),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:73.0), Gecko/20100101 Firefox/73.0"),
		colly.CacheDir("resources/cache/allegiant"),
		colly.Async(false), /* Maybe True... Allegiant uses AJAX for normal searches */
	)
}

type allegiantDealListing struct {
	departLocation     string
	arriveLocation     string
	effectiveDateRange string
	price              string
	duration           string
}

func (card *allegiantDealListing) String() string {
	return fmt.Sprintf("From %v\nTo %v\nFor $%v", card.departLocation, card.arriveLocation, card.price)
}

func isNil(value string) bool {
	return value == ""
}

func NewListingsFrom(flightcard string) []Listing {

}

func extractListingFrom(flightDealCard *colly.HTMLElement) {
	departLocation := flightDealCard.ChildText(".origin")
	arriveLocation := flightDealCard.ChildText(".destination")
	effectiveDateRange := flightDealCard.ChildText(".line3")
	price := flightDealCard.ChildText(".pc-price")
	duration := flightDealCard.ChildText("span.flight-time span")

	if !isNil(price) {
		newCard := allegiantDealListing{
			departLocation:     departLocation,
			arriveLocation:     arriveLocation,
			effectiveDateRange: effectiveDateRange,
			price:              price,
			duration:           duration,
		}
		flightDeals = append(flightDeals, newCard)
	}

}

const allegiantCacheDir string = "cache/allegiant"

// ScrapeDealsPage will nevigate to a specific page, download the raw HTML and save that to a cache
func ScrapeDealsPage(link string) ([]allegiantDealListing, error) {
	c := makeDefaultCollector()

	if link != nil {
		c.Visit(link)
		Wait()

	}

	var flightDeals []allegiantDealListing
	c.OnHTML("div.flight-deal-card", extractListingFrom())

	return flightDeals, nil
}

func parseFromCache(path string) []Listing {
	cache := io.ReaderFrom(path)
	var doc goquery.Document = goquery.NewDocumentFromReader(cache)
	colly.Context
	doc

}

func NewListingFromAllegiantFlightCard(deal allegiantDealListing) []Listing {
	// Issue with the Listings from the deal page is that the average price is calculated
	// from a sampling of days which of the range, may not have full coverage. So what is shown
	// is bad data but a source of data as a test, thi means we should distinguish from the rest
	// with something like a marker such as the scrape source.

}

func main() {
	// Instantiate default collector
	justALinkForNow := "https://flight.deals.allegiant.com/ats/url.aspx?cr=986&wu=95&camp=20200301_SundaySavings&allow_search=1"
	flightDeals, err := ScrapeDealsPage(justALinkForNow)
	fmt.Println(flightDeals)

	if err != nil {
		return
	}
}
