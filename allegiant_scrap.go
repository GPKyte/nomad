package nomad

import (
	"fmt"

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

func scrapeDealsPage(link string) ([]allegiantDealListing, error) {
	c := makeDefaultCollector()
	c.Visit(link)

	Wait(504)
	var flightDeals []allegiantDealListing
	c.OnHTML("div.flight-deal-card", func(flightDealCard *colly.HTMLElement) {
		departLocation := flightDealCard.ChildText("div.origin")
		arriveLocation := flightDealCard.ChildText("div.destination")
		effectiveDateRange := flightDealCard.ChildText("div.line3")
		price := flightDealCard.ChildText("div.pc-price")
		duration := flightDealCard.ChildText("span.flight-time span")

		if isNil(price) {
			newCard := allegiantDealListing{
				departLocation:     departLocation,
				arriveLocation:     arriveLocation,
				effectiveDateRange: effectiveDateRange,
				price:              price,
				duration:           duration,
			}
			flightDeals = append(flightDeals, newCard)
		}
	})

	return flightDeals, nil
}

func main() {
	// Instantiate default collector
	justALinkForNow := "https://flight.deals.allegiant.com/ats/url.aspx?cr=986&wu=95&camp=20200301_SundaySavings&allow_search=1"
	flightDeals, err := scrapeDealsPage(justALinkForNow)
	fmt.Println(flightDeals)
}
