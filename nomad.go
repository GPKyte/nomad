// Package main is the entry point to project NOMAD which is a reflection of a reckless desire to be uproot and travel novel routes at deeply discounted rates
// To accomplish this goal, we utilize web scraping to collect airfares into a format convenient for analysis
// To read the full feature set and usage of NOMAD please review the README at the top of this directory.
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/GPKyte/nomad/scrape"
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

// Maintains collection of scrapers and may be used to link collected Listings to Graph structure
type nomad struct {
	bots []Scraper
}
type travelSettings struct {
	Destinations []string
	Origin       string
	Budget       int
	PlannedTrips []time.Time
}

func main() {
	var pipeline = make(chan scrape.Listing) // Give several routines same(?) channel to send Listings back on */
	var nomad = new(nomad)
	nomad.bots = append(nomad.bots, scrape.NewSkippy(pipeline))
	// As more bot types are created, append them to nomad

	go record(pipeline, true)
	/* go watch(pipeline, criteria) Inspect all Listings through pipeline */

	var settings = loadTravelSettings()
	var originCity = settings.Origin
	for _, D := range settings.Destinations {
		nomad.doSearchAhead(originCity, D)
	}

	checkForHowEarlyToBook(nomad)
	checkForSpecificallyPlannedTrip(nomad, settings)

	close(pipeline)
}

/* TODO: Understand the formatting of data for insertion, debate whether to also leverage JSON as postgres supports in and embedded objects */
func record(data chan scrape.Listing, verbose bool) {
	var db = initDB()
	var tx = db.MustBegin() // Init first batch transaction

	// We have options for how often to write to the database
	// Bad soln: 	Write to DB (INSERT fare) with every. single. listing given. This would be bad. instead try to batch either by routine or...
	// Okay soln:	Every N Listings, and before closing: bulk insert what's left in channel
	dump := func() {
		tx.Commit()
	}

	for count := 0; true; count++ {
		fareListing, moreComing := <-data

		if !moreComing /* Listenig to a closed channel, wrap up */ {
			/* Begin closing resources */
			dump()
			db.Close()
			break
		}
		query := `INSERT INTO fare
		(airline, cost, legs, begTime, begLoc, endTime, endLoc, srcTime, srcLoc)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"`
		L := fareListing // alias

		tx.MustExec(query, textNotRecorded, L.Price, numNotRecorded, L.Depart.DateTime, L.Depart.Location, L.Arrive.DateTime, L.Arrive.Location, L.Scrape.DateTime, L.Scrape.Location)

		if verbose {
			fmt.Println("Insert into fare table: ", fareListing)
		}
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
func (N *nomad) doSearchExact(from, to string, date time.Time) {
	for _, search := range N.bots {
		search.AToBDuring(from, to, date, date.AddDate(0, 0, 5))
	}
}

// doSearchAhead will Looks at upcoming fares as part of the attempt to gather consistent lookahead data throughout the year
func (N *nomad) doSearchAhead(from, to string) {
	for _, search := range N.bots {
		search.AToBNearDate(from, to, time.Now())
	}
}

// doSearchFares acknowledges that we may know where to go but for the right price it can be many more places than we originally thought
// while this can't provide the deaggregated data useful for the envisioned table of detailed listings
// it is useful in quickly populating a relaxed-contraint graph for route finding with low-cost edges
// especially if given origins at the most strongly connected layover destinations
func (N *nomad) doSearchFares(origins ...string) {
	for _, search := range N.bots {
		for _, o := range origins {
			search.AToAnywhereSoon(o)
		}
	}
}

// watch incoming listings in the background and attempt route-finding magic
func watch(this chan scrape.Listing, criteria travelSettings) {
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

// Given that any are defined, look for trips to any desired destinations on indicated dates of travel
func checkForSpecificallyPlannedTrip(n *nomad, settings travelSettings) {
	var F = settings.Origin
	for _, T := range settings.Destinations {
		for _, date := range settings.PlannedTrips {
			n.doSearchExact(F, T, date)
		}
	}
}

// Determine the impact on price from days until departure from date of purchase
// by collecting a variety of dates; denser in the immediate future and sparser farther away to reduce requested data
func checkForHowEarlyToBook(n *nomad) {
	/* Define a handful of static locations to use as reference points */
	/* That's about 500 Requests, this should happen infrequently, like each month */
	from := []string{"CLE", "CVG", "PIT"}
	to := []string{"DEN", "PIE", "LAX"}
	base := time.Now()

	scrape := n.bots[0]

	for _, F := range from {
		for _, T := range to {
			scrape.AToBDuring(F, T, base, base.AddDate(0, 1, 0))                      // 0-30d
			scrape.AToBDuring(F, T, base.AddDate(0, 0, 6*7), base.AddDate(0, 0, 7*7)) // Buddy said 6 Weeks in the best, so let's check that
			scrape.AToBNearDate(F, T, base.AddDate(0, 2, 0))                          // 60-65d
			scrape.AToBNearDate(F, T, base.AddDate(0, 6, 0))                          // 180-185d
		}
	}
	/* TODO: Consider reducing data here into a selection of the answer to the question named by the method header */
}

func loadTravelSettings() travelSettings {
	tzone, _ := time.LoadLocation("America/New_York")
	// This method stubs loading user profiles until expanding to multi-users
	// Can easily write up this data in a config file or as a user document stored in DB when it becomes important to do so
	return travelSettings{
		Destinations: []string{"SVQ", "LAX", "IDA", "AKJ", "PIT", "SGU"}, /* These are airport codes */
		Origin:       "CLE",                                              /* May expand to multiple origins, but ideally one per user */
		Budget:       250,                                                /* Used in filtering */
		PlannedTrips: []time.Time{time.Date(2020, 8, 1, 0, 0, 0, 0, tzone)},
	}
}
