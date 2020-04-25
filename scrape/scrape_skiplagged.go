package scrape

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
		"Accept-Encoding":           "gzip, deflate",
		"Accept-Language":           "en-US,en;q=0.5",
		"Cache-control":             "max-age=0",
		"Connection":                "keep-alive",
		"DNT":                       "1",
		"Cookie":                    "session=a6d88f9a76de15278402f46a047b0724; currencyRate=1; currencyFormat=%24%7Bamount%7D; currencyCode=USD; when=2020-05-07; whenBack=; G_ENABLED_IDPS=google",
		"Upgrade-Insecure-Requests": "1",
		"User-Agent":                "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:75.0) Gecko/20100101 Firefox/75.0",
	}

	url := fmt.Sprintf("http://skiplagged.com/api/skipsy.php?%s", concatURLArgs(urlargs))
	// Visit site
	visit(url, reqHeaders)
	// Check response for errors, if any wait and try again up to so many times

	// Open response in json and parse to trips
	// Use combo of custom Unmarshalling and interpretation to create Listings. Consider how to be 'lax on timing aspect. Choose middle of day?

}

type logger struct{}

func (L *logger) Write(p []byte) (n int, err error) {
	fmt.Print(string(p))
	return len(p), nil
}

func visit(url string, reqHeaders map[string]string) {
	fmt.Println(url)
	resp, err := http.Get(url)

	if err != nil {
		panic(string(err.Error()))
	}

	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var responseAsJSON apiResponse
	json.Unmarshal(b, &responseAsJSON)

	for _, t := range responseAsJSON.Trips {
		fmt.Printf("%s: $%v", t.City, t.Cost)
	}
}

type apiResponse struct {
	Trips []*trip `json:"trips"`
}

type trip struct {
	City       string `json:"city"`
	Cost       int    `json:"cost"`
	HiddenCity bool   `json:"hiddenCity"`
}
