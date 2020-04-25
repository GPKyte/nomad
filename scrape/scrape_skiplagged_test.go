package scrape

import (
	"fmt"
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

func failedExpectations(want string, got string) bool {
	return (want != got)
}

func explain(want string, got string) string {
	return fmt.Sprintf("\nWanted: \t%s\nGot Instead: \t%s", want, got)
}
