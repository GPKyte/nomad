package common

import "testing"

func TestRandomListingGen(t *testing.T) {
	L := NewListingRand()

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

func TestGetDurationOfTrip(t *testing.T) {
	/* Build a trip with forced start and end times
	 * Then calc time difference in (Unit of time)
	 * Verify expectations */
	var L = NewListingRand()
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
	}
	}

	}
}
