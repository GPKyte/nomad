package scrape

import (
	"fmt"
	"strings"
)

func chooseDate() string {/* want to use time.Time to format std dates like this */ return "2020-05-07" }
func chooseLocation() string { return "CVG" }
func concatURLArgs(kv map[string]string) string {
	var cat = make([]string, len(kv))

	for k, v := range kv {
		/* Could consider santizing arguments here, but leaving this chore for later */
		cat = append(cat, fmt.Sprintf("%s=%s"))
	}

	return strings.Join(cat, "&")
}

func main() {
	// Build Request Headers
	// format URL to visit
	// Decide some initial values, or create them dynamically for trips to find
	//      Consider leaving open for config
	//
	request_headers := 

	urlargs := map[string]string{
		"from":   chooseLocation(),
		"depart": chooseDate(), // YYYY-MM-DD
		"return": "",
		"format": "v2",
		"_":      "1587607414580", // Probably a timer of sorts
	}

}
