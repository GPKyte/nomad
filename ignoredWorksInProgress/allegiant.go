package scrape

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

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
func scrapeDealsPage() []Listing {
	likelyCap := 100
	listings := make([]Listing, likelyCap)
	col := newDefaultCollector()
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

	// Find deals page at `#mini-panel-allegiant3_bottom_menu > li.first leaf > a`.Text()
	url := "https://deals.allegiant.com/ats/url.aspx?cr=986&wu=11"
	col.Visit(url)

	for deal := range flightDeals {
		fmt.Println(deal.String())
		listings = append(listings, newPsuedoListingFromAllegiant(deal))
	}
	return listings
}

// Useful to make the time data from deal workabale, returning beg and end datetime
func makeUsefulTimeFrame(deal allegiantDealListing) (time.Time, time.Time) {
	//durationHours if I really want to patter match for this detail
	//durationMinutes same case as hours ^
	const nanoConvRate = 10 ^ 9 // Because Time lib uses nanoseconds

	var flightTime = 2 /* hours */ * 3600 /* sec/hr */ * nanoConvRate // Just place holder for simplification
	// "Rates sampled from {starttraveldateM_d_yyyy} through {endtraveldateM_d_yyyy} include all taxes and fees"
	datePattern := regexp.MustCompile(`(\d\d?/){2}\d{4}`)
	startDate := datePattern.FindAllString(deal.effectiveDateRange, 1)
	month, day, year := strings.Split(startDate, "/")
	hour, min, sec := 0 // Why bother with this for flight cards?

	return time.Date(year, month, day, hour, min, sec, 0, time.FixedZone("UTC-5")),
		time.Date(year, month, day, hour, min, sec, flightTime, time.FixedZone("UTC-5"))
}

func newPsuedoListingFromAllegiant(deal allegiantDealListing, srcURL string) Listing {
	var home string = strings.Split(deal.departLocation, " to")[:1] // Strip extra wording

	return Listing{
		Depart:     timeSpace{Time: before, Location: home},
		Arrive:     timeSpace{Time: after, Location: deal.arriveLocation},
		Scrape:     makeScrapeStamp(srcURL),
		Price:      int(deal.price),
		srcContent: deal.String(),
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

func loadAllegiantLocationData(index int) []location {
	var allegiantLocationData []byte = ioutil.ReadFile("resources/allegiant/locations.txt")
	var locations []location

	newlineMarkers := make(chan int)
	lines := make(chan string)
	done := make(chan bool)

	// Parse each location into data struct on init
	go func() {
		locations = make([]string, 100)

		for {
			locData, more := <-lines

			if more {
				commaIndex := strings.IndexByte(locData, byte(','))
				afterParenIndex := strings.IndexByte(locData, byte('(')) + 1 // Don't want paren in the code

				name := locData[:commaIndex]
				state := locData[commaIndex+len(", ") : commaIndex+len(", XY")] // Get 'XY'
				code := locData[afterParenIndex : afterParenIndex+len("ABC")]   // Get 'ABC'

				locations = append(locations, location{Name: name, State: state, Code: code})
			} else {
				done <- true
				return
			}
		}
	}()

	go func() {
		start := 0
		for {
			nextDelim, more := <-newlineMarkers

			if more {
				lines <- string(allegiantLocationData[start:nextDelim])
				start = nextDelim + 1

			} else {
				lines <- string(allegiantLocationData[start:])
				return
			}
		}
	}()

	// Find all newline characters deliminating lines
	for i := 0; i < len(allegiantLocationData); i++ {
		if allegiantLocationData[i] == '\n' {
			newlineMarkers <- i
		}
	}

	<-done
	return locations
}
