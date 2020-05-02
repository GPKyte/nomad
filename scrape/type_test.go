package scrape

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"testing"
)

func TestRandomListingGen(t *testing.T) {
	L := newListingRand()

	if L.Arrive.DateTime.Before(L.Depart.DateTime) {
		t.Fail()
	}
	if L.Arrive.Location == L.Depart.Location {
		t.Fail()
	}
}

func chooseLocation() string {
	var locs = []string{"A", "B", "C", "D", "E"}
	return locs[rand.Intn(5)]
}

func getAirportLocations() []location {
	var result = make([]location, 200)
	var locationsRaw []struct {
		name string
		code string
	}

	locationCache, err := ioutil.ReadFile("../resources/test/cache/locations.json")
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(locationCache, locationsRaw); err != nil {
		panic(err)
	}

	for _, L := range locationsRaw {
		result = append(result, location{L.name, L.code})
	}
	return result
}

func pick(howMany int, ofThese []location) []location {
	var seed = 5090716181  // Any number, skip seeding override when determinism wanted
	rand.Seed(int64(seed)) // Override while considering a rand.Int() soln or when determinism wanted

	if howMany > len(ofThese) { // Sanity check
		panic(howMany)
	}

	// Build up a map to select unique elements and ignore repeats
	var uniquePicks = make(map[location]int, howMany)
	for len(uniquePicks) < howMany {
		p := rand.Int() % len(ofThese)
		thisPick := ofThese[p]
		uniquePicks[thisPick]++
	}

	// Retrieve keys into a simple slice and return address of it
	var chosen = make([]location, howMany)
	for k := range uniquePicks {
		chosen = append(chosen, k)
	}

	return chosen
}
