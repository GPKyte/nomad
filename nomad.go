// Package nomad is a reflection of a reckless desire to be uproot and travel novel routes at deeply discounted rates
// To accomplish this goal, we utilize web scraping to collect airfares into a format convenient for analysis
// To read the full feature set and usage of NOMAD please review the README at the top of this directory.
package nomad

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/GPKyte/nomad/scrape"
	s "github.com/GPKyte/nomad/scrape"
	"github.com/jmoiron/sqlx"
)

// Here lay the options that may be supported one day. Eren solo un idea
const (
	Verbose = iota
	EditSettings
	WatchTrip
	WatchTime
	ScheduleRepeatingTask
	WeekAheadTask
	Automaton
)

// Environment settings, e.g. connection with Heroku PostgreSQL uses DATABASE_URL
const (
	databaseURL = "https://localhost"
	appPort     = 9716
)

const numNotRecorded = 997163173 /* An absurdly high and specific number to indicate something was not recorded */
const textNotRecorded = "N/A"

const createSchema = `
CREATE TABLE fare (
	airline text,
	cost	int,
	legs	int,
	begTime	time,
	begLoc	text,
	endTime	time,
	endLoc	text,
	srcTime time,
	srcLoc	text
);
`
const deleteSchema = `
DROP TABLE fare
`

// Scraper is the expected set of methods that any input source of data should provide
// implementations of to be relied on at this level
type Scraper interface {
	// It is expected that all data be sent back via the channel owned by the scraper on init
	AToAnywhereSoon(string)
	AToBNearDate(string, string, time.Time)
	AToBDuring(string, string, time.Time, time.Time)
}
type travelSettings struct {
	Destinations []string
	Origin       string
	Budget       int
}

func main() {
	var pipeline = make(chan scrape.Listing)
	var agent = skippy.New(pipeline) /* Skiplagged is one of the many travel websites which lists airfare that we collect */

	go record(pipeline)
	/* go watch(pipeline, criteria) Inspect all Listings through pipeline */

	/* Work with scheduler on initing defined tasks at specified intervals
	// Give several routines same(?) channel to send Listings back on

	*/
	// What are the searches I want to run? Pick three places tu honto ni querias viajar
	destinations = loadTravelSettings()
	for _, D := range destinations {
		doSearchAhead(originCity, D)
	}

	doSearchExact(originCity, destination, date)
	doSearchFares(originCity)

	checkForUpcomingDeals()
	checkForSpecificallyPlannedTrip()
	checkForHowEarlyToBook()

	close(pipeline)
}

/* TODO: Understand the formatting of data for insertion, debate whether to also leverage JSON as postgres supports in and embedded objects */
func record(data chan scrape.Listing) {
	var db = initDB()
	var tx = db.MustBegin() // Init first batch transaction

	// We have options for how often to write to the database
	// Bad soln: 	Write to DB (INSERT fare) with every. single. listing given. This would be bad. instead try to batch either by routine or...
	// Okay soln:	Every N Listings, and before closing: bulk insert what's left in channel
	dump := func() {
		tx.Commit()
	}
	insert := func(L scrape.Listing) {
		// "INSERT INTO person (first_name, last_name, email) VALUES (:first_name, :last_name, :email)
		err := tx.MustExec(`INSERT INTO fare
			(airline, cost, legs, begTime, begLoc, endTime, endLoc, srcTime, srcLoc)
			VALUES $1, $2, $3, $4, $5, $6, $7, $8, $9"`,
			textNotRecorded, L.Price, numNotRecorded, L.Depart.Time, L.Depart.Location, L.Arrive.Time, L.Arrive.Location, L.Scrape.Time, L.ScrapeLocation)

		if err != nil {
			log.Println(err.Error())
		}
	}

	for count := 0; true; count++ {
		fareListing, moreComing := <-data

		if !moreComing {
			/* Begin closing resources */
			dump()
			closeDB(persistentCollection)
			break
		}
		insert(fareListing)

		if count%100 == 99 {
			dump()
			tx = db.MustBegin()
		}
	}
}

/* Return a pool of connections to Travel fare database */
func initDB() *sqlx.DB {
	connParam := fmt.Sprint(os.Getenv("DATABASE_URL"), "sslmode=enable") // May need additional param como user, y password, pero no intereso en ese ahora
	conn, err := sqlx.Connect("postgres", connParam)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return conn
}

// The following doSearch helpers are helpful abstractions
// doSearchExact can be used to track one particular planned vaction trip
// giving primarily data on the (time of departure - time of scraping) impact on price
// But also useful for watching for a surprise deal matching budget constraints
func (bot *s.Scraper) doSearchExact(from, to string, date time.Time) {
	bot.AToBDuring(from, to, date, date.AddDate(0, 0, 5))
}

// doSearchAhead will Looks at upcoming fares as part of the attempt to gather consistent lookahead data throughout the year
func (bot *s.Scraper) doSearchAhead(from, to string) {
	bot.AToBNearDate(from, to, time.Now())
}

// doSearchFares acknowledges that we may know where to go but for the right price it can be many more places than we originally thought
// while this can't provide the deaggregated data useful for the envisioned table of detailed listings
// it is useful in quickly populating a relaxed-contraint graph for route finding with low-cost edges
// especially if given origins at the most strongly connected layover destinations
func (bot *s.Scraper) doSearchFares(origins ...string) {
	for O := range origins {
		bot.AToAnywhereSoon(O)
	}
}

// watch incoming listings in the background and attempt route-finding magic
func watch(this chan Listing, criteria travelSettings) {
	// Make the compass as the interface between graph model and Fare data
	// This should be in memory and maybe should be on a different channel than the main pipeline

	// We look for good deals individually to watched locations
	// And we are checking for unicorn routes (MultiDay Layover over between multiple great deals)
	for {
		// ...
	}
	// But we should also be managing a graph and building the routes for it and
	// When unicorn sighted
	// Log sighting
	// Rescrape the path given between all Locations involved
	// Check for roundtrip too??
}

func checkForSpecificallyPlannedTrip(s *Scraper, settings travelSettings) {
	for trip := range settings.plannedTrips {
		trip.start
		trip.end
		trip.date
	}
}

// Determine the impact on price from days until departure from date of purchase
// by collecting a variety of dates; denser in the immediate future and sparser farther away to reduce requested data
func checkForHowEarlyToBook(s *Scraper) {
	/* Define a handful of static locations to use as reference points */
	/* That's about 500 Requests, this should happen infrequently, like each month */
	from := getLocationsByCode("CLE", "CVG", "PIT")
	to := getLocationsByCode("DEN", "PIE", "LAX")
	base := time.Now()

	for _, F := range from {
		for _, T := range to {
			s.AToBDuring(F, T, base, base.AddDate(0, 1, 0))                      // 0-30d
			s.AToBDuring(F, T, base.AddDate(0, 0, 6*7), base.AddDate(0, 0, 7*7)) // Buddy said 6 Weeks in the best, so let's check that
			s.AToBNearDate(F, T, base.AddDate(0, 2, 0))                          // 60-65d
			s.AToBNearDate(F, T, base.AddDate(0, 6, 0))                          // 180-185d
		}
	}
	/* TODO: Consider reducing data here into a selection of the answer to the question named by the method header */
}

func loadTravelSettings() travelSettings {
	// This method stubs loading user profiles until expanding to multi-users
	// Can easily write up this data in a config file or as a user document stored in DB when it becomes important to do so
	return travelSettings{
		Destinations: []string{"SVQ", "LAX", "IDA", "AKJ", "PIT", "SGU"}, /* These are airport codes */
		Origin:       "CLE",                                              /* May expand to multiple origins, but ideally one per user */
		Budget:       250,                                                /* Used in filtering */
		PlannedTrips: trip
	}
}
