package scrape

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
	"testing"
	"time"
)

const (
	Verbose = false
)

func TestURLArgsHandling(t *testing.T) {
	var args = map[string]string{
		"1":  "are",
		"2":  "you",
		"A":  "ready",
		"b":  "to",
		"Ca": "Rum",
		"_":  "mmmm",
		"Pa": "ble!!!?!?",
	}

	want := "1=are&2=you&A=ready&b=to&Ca=Rum&_=mmmm&Pa=ble!!!?!?"
	got := concatURLArgs(args)

	/* Unfortunately, maps do not preserve order
	 * Fortunately, we don't actually care, but it does change our simple test
	 * Either one can Unmarshall the "got" string and confirm it
	 * Or one can simply validate a few substrings and a count of the delimiter as a smoke test */
	var goodCount bool = (len(args)-1 == strings.Count(got, "&"))
	var foundEnoughMatches bool = (strings.Contains(got, "1=are") && strings.Contains(got, "Ca=Rum") && strings.Contains(got, "_=mmm"))
	var sameLength bool = (len(want) == len(got))

	if !foundEnoughMatches || !goodCount || !sameLength {
		// fmt.Printf("goodCount: %v, enoughSubstringsMatch: %v, sameLength: %v", goodCount, foundEnoughMatches, sameLength)
		t.Fatal(explain(want, got))
	}
}

func bad(e error) bool {
	return e != nil
}

func TestUnmarshalFromCache(t *testing.T) {
	responseAsJSON := new(apiResponse)
	b, err := ioutil.ReadFile("cache/cvg/any/2020.5.7.json")

	if bad(err) {
		t.Fatal(err.Error())
	}

	err = json.Unmarshal(b, responseAsJSON)
	trips := responseAsJSON.Trips

	if bad(err) || len(trips) <= 1 {
		t.Fatal(err.Error())
	}
	any := trips[rand.Intn(len(trips))].City
	if Verbose {
		fmt.Printf("Just to check...Found a trip from CVG to %s\n", any)

		for _, t := range trips {
			fmt.Println(t.City, " ", t.Cost)
		}
	}
}

/* Testing for these traits
* That Dates generated match the expected #, and start/end
* Format must be accurate, if wrong, inspect struct */
func TestDateGenerationForURLArgs(t *testing.T) {
	var locale, _ = time.LoadLocation("UTC")
	var onceUponATime = time.Date(2050 /*yr*/, 5 /*mo*/, 20 /*d*/, 15 /*hr*/, 0, 0, 0, locale)

	/* Only checking Formatting Once, but all should comply once the Format string const is accurate */
	if got := onceUponATime.Format(DateFormat); got != "2050-05-20" {
		panic(got)
	}

	type testCondition struct {
		now, then time.Time
		length    int
	}
	testBoundaries := []testCondition{
		// While other specifics may come up, always test 0, 1, and N
		{onceUponATime, onceUponATime.AddDate(1, 0, 0), 365}, // Careful about leap year, so determinate start date used
		{onceUponATime, onceUponATime.AddDate(0, 0, 1), 1},
		{onceUponATime, onceUponATime, 0},
	}
	/* Test expected length of resultant slices */
	for _, test := range testBoundaries {
		before, after, want := test.now, test.then, test.length
		got := getDatesBetween(before, after)

		if len(got) != want {
			t.Fatalf("\nWanted length:\t%v,\nGot:\t\t%s", want, got)
		}
	}
	/* Random Testing additionally used on similar abstracted method */
	for i := 0; i < 100; i++ {
		days := rand.Intn(400)

		/* Starting from Tomorrow or Today, TODO: decide later when it gets used more often */
		if list := getDatesForNext(days); len(list) != days {
			t.Fatal("Not enough days generated from getDatesForNext(N(days)")
		}
	}
}

func TestCheckWhenEarlyBirdRises(t *testing.T) {
	checkWhenTheEarlyBirdRises()
}

func emptyStringSlice(this []string) bool {
	return len(this) == 0
}

func TestLoadCacheOfAirports(t *testing.T) {
	airports := loadCacheOfAirports()

	if len(airports) == 0 {
		t.Fatal("Empty results from loading cache of airports")
	}

	if len(airports[0].Code) != len("CVG") {
		t.Fatalf("Wrong format:\t%s\n\tXYZ three digit code preferred for Airports", airports[0])
	}
}

func TestInspectRecentTestResults(t *testing.T) {
	t.SkipNow()
	fmt.Println(loadCacheOfAirports())
}

func TestScrapeOne2Any(t *testing.T) {

}
func TestScrapeOne2Another(t *testing.T) {

}

func failedExpectations(want string, got string) bool {
	return (want != got)
}

func explain(want string, got string) string {
	return fmt.Sprintf("\nWanted: \t%s\nGot Instead: \t%s", want, got)
}
