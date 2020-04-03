package nomad

import (
	"fmt"
	"os"

	"github.com/PuerkitoBio/goquery"

	"github.com/gocolly/colly"
)

const (
	allegiantCacheDir = "cache/allegiant"
	rawHTMLSuffix     = "_raw.html"
)

// BasicHTMLError is a simple std for error reporting
type BasicHTMLError struct {
	context string
	want    string
	got     string
	html    *colly.HTMLElement
	DOM     *goquery.Selection
}

func (e *BasicHTMLError) Error() string {
	verbose := bool(false)
	var contextLine string

	if len(e.context) > 0 && verbose {
		contextLine = "Context: " + e.context + "\n"
	}
	return fmt.Sprintf("Probably Parsing related;\nWant: %s\nGot: %s\n%s", e.want, e.got, contextLine)
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

func newDefaultCollector() *colly.Collector {
	return colly.NewCollector(
		colly.AllowedDomains("www.allegiantair.com"),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:73.0), Gecko/20100101 Firefox/73.0"),
		colly.CacheDir(allegiantCacheDir),
		colly.Async(false), /* Maybe True... Allegiant uses AJAX for normal searches */
	)
}

// scrapeDealsPage will nevigate to a specific page, download the raw HTML and save that to a cache
func scrapeDealsPage(col *colly.Collector, url string) chan allegiantDealListing {
	flightDeals := make(chan allegiantDealListing)

	col.OnHTML("div.flight-deal-card", func(flightDealCard *colly.HTMLElement) {
		departLocation := flightDealCard.ChildText(".origin")
		arriveLocation := flightDealCard.ChildText(".destination")
		effectiveDateRange := flightDealCard.ChildText(".line3")
		price := flightDealCard.ChildText(".pc-price")
		duration := flightDealCard.ChildText("span.flight-time span")

		flightDeals <- allegiantDealListing{
			departLocation:     departLocation,
			arriveLocation:     arriveLocation,
			effectiveDateRange: effectiveDateRange,
			price:              price,
			duration:           duration,
		}
	})

	col.OnScraped(func(r *colly.Response) {
		cache(r.Request.URL.String()+rawHTMLSuffix, r.Body)
		close(flightDeals)
	})

	col.Visit(url)

	return flightDeals
}

func main() {
	// Instantiate default collector
	col := newDefaultCollector()
	// Find deals page at `#mini-panel-allegiant3_bottom_menu > li.first leaf > a`.Text()
	link2deals := "https://deals.allegiant.com/ats/url.aspx?cr=986&wu=11"
	flightDeals := scrapeDealsPage(col, link2deals)

	for deal := range flightDeals {
		fmt.Println(deal.String())
	}
}

func log(e error) {
	fmt.Println(e)
}

func cache(refname string, data []byte) {
	// Open a new file at refname, and push the data to it.
	// If the refname exists in the cache already, log the error
	// then generate a new refname
	cacheLoc := allegiantCacheDir + refname
	if _, err := os.Stat(cacheLoc); os.IsExist(err) {
		cache("9"+refname, data) // Think of a good cache naming system and avoid this problem
		// Maybe just report error, but don't I want to save all scrape data?
	}
	err := writeFile(cacheLoc, data)

	if err != nil {
		log(err)
	}
}

func parseFromCache(path string) []Listing {
	// TODO: read n files from cache location, should be agnostic of naming convention
	return make([]Listing, 0)
}
