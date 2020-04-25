package scrape

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
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

func head(n int, items ...interface{}) []interface{} {
	return items[:n]
}

func TestUnmarshalFromCache(t *testing.T) {
	responseAsJSON := new(apiResponse)
	b, err := ioutil.ReadFile("cache/cvg/any/2020.5.7")

	if bad(err) {
		t.Fatal(err.Error())
	}

	err = json.Unmarshal(b, responseAsJSON)
	if bad(err) || len(responseAsJSON.Trips) <= 1 {
		t.Fatal(err.Error())
	}

	for _, t := range head(100, responseAsJSON.Trips) {
		t, ok := t.(trip)
		if !ok {
			panic("How the heck could this have happened?")
		}
		fmt.Printf("%s: $%v\n", t.City, t.Cost/100)
	}
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
