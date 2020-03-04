package nomad

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"

	"github.com/gocolly/colly"
)

type BasicHTMLError struct {
	context string
	want    string
	got     string
	html    *colly.HTMLElement
	DOM     *goquery.Selection
}

func rapidTest(c *colly.Collector) error {
	c.Visit("https://flight.deals.allegiant.com/ats/url.aspx?cr=986&wu=95&camp=20200223_SundaySavings&allow_search=1")
	err := c.OnHTML("div.flight-deal-card", func(e *colly.HTMLElement) {
		fmt.Print(e.ChildTexts("div.origin"))
	})
	return err
}

func (e *BasicHTMLError) Error() string {
	verbose := bool(false)
	var contextLine string

	if len(e.context) > 0 && verbose {
		contextLine = "Context: " + e.context + "\n"
	}
	return fmt.Sprintf("Probably Parsing related;\nWant: %s\nGot: %s\n%s", e.want, e.got, contextLine)
}

func makeDefaultCollector() colly.Collector {
	return colly.NewCollector(
		colly.AllowedDomains("allegiant.com"),
		colly.MaxDepth(2),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:73.0), Gecko/20100101 Firefox/73.0"),
		colly.ParseHTTPErrorResponse(global.DEBUG), /* TODO: make this global var work */
		colly.MaxBodySize(maxint),
		colly.CacheDir(""),
		colly.Async(false),
	)
}

type LPAllegiant ListingParser

// RecordScrapestamp returns a default timeSpace because this is the default Listing class
func (LP *ListingParser) RecordScrapestamp(L *Listing) err {
	location := LP.url
	datetime := Now()
	L.ScrapeStamp = timeSpace{Location: location, DateTime: datetime}
	return nil
}

// RecordDeparture returns a default timeSpace because this is the default Listing class
func (LP *ListingParser) RecordDeparture(L *Listing) err {

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

type flightDealCard struct {
	departLocation     string
	arriveLocation     string
	effectiveDateRange string
	price              int
	duration           string
}

func (card *flightDealCard) String() string {
	return fmt.Sprintf("From %v\nTo %v\nFor $%v", card.departLocation, card.arriveLocation, card.price)
}

func (c *colly.Collector) scrapeDealsPage(link string) ([]flightDealCard, error) {
	c.Visit(link)

	Wait(504)
	var flightDeals []flightDealCard
	c.OnHTML("div.flight-deal-card", func(flightDealCard *colly.HTMLElement) {
		departLocation := flightDealCard.ChildAttr("div.origin")
		arriveLocation := flightDealCard.ChildAttr("div.destination")
		effectiveDateRange := flightDealCard.ChildAttr("div.line3")
		price := flightDealCard.ChildAttr("div.pc-price")
		duration := flightDealCard.ChildAttr("span.flight-time span")

		if price != nil {
			append(flightDeals, fmt.Println(flightDealCard{
				departLocation:     departLocation,
				arriveLocation:     arriveLocation,
				effectiveDateRange: effectiveDateRange,
				price:              price,
				duration:           duration,
			}))
		}
	})

	return flightDeals, nil
}

func main() {
	// Instantiate default collector
	c = makeCollector()

	justALinkForNow := "https://flight.deals.allegiant.com/ats/url.aspx?cr=986&wu=95&camp=20200301_SundaySavings&allow_search=1"
	flightDeals, err := scrapeDealsPage(justALinkForNow)
	fmt.Println(flightDeals...)
}
