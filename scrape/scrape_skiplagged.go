package scrape

import (
	"fmt"
	"strings"
)

func chooseDate() string     { /* want to use time.Time to format std dates like this */ return "2020-05-07" }
func chooseLocation() string { return "CVG" }
func concatURLArgs(kv map[string]string) string {
	var cat []string

	for k, v := range kv {
		/* Could consider santizing arguments here, but leaving this chore for later */
		cat = append(cat, fmt.Sprintf("%s=%s", k, v))
	}

	return strings.Join(cat, "&")
}

func main() {
	// Build Request Headers
	// format URL to visit
	// Decide some initial values, or create them dynamically for trips to find
	//      Consider leaving open for config
	//

	urlargs := map[string]string{
		"from":   chooseLocation(),
		"depart": chooseDate(), // YYYY-MM-DD
		"return": "",
		"format": "v2",
		"_":      "1587607414580", // Probably a timer of sorts
	}

	reqHeaders := map[string]string{
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
		"Accept-Encoding":           "gzip, deflate, br",
		"Accept-Language":           "en-US,en;q=0.5",
		"Connection":                "keep-alive",
		"DNT":                       "1",
		"Host":                      "skiplagged.com",
		"TE":                        "Trailers",
		"Upgrade-Insecure-Requests": "1",
		"User-Agent":                "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:75.0) Gecko/20100101 Firefox/75.0",
	}

	url := fmt.Sprintf("https://skiplagged.com/api/skipsy.php?%s", concatURLArgs(urlargs))
	// Visit site
	visit(url, reqHeaders)
	// Check response for errors, if any wait and try again up to so many times

	// Open response in json and parse to trips
	// Use combo of custom Unmarshalling and interpretation to create Listings. Consider how to be 'lax on timing aspect. Choose middle of day?

}

func visit(url string, reqHeaders map[string]string) {
	fmt.Println(url)
}
