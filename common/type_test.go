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

func TestConvListingToGraphDetail(t *testing.T) {
	/* Generate listings by random */
	/* Inspect Mode on IF not expected output */
	/* Need what data? Listing.src,dest,time,duration,price */
	var L *Listing = new(Listing)
	L.Arrive.Location = chooseLocation()
	L.Depart.Location = chooseLocation()
	L.Price = 8
}

func TestGetDurationOfTrip(t *testing.T) {
	/* Build a trip with forced start and end times
	 * Then calc time difference in (Unit of time)
	 * Verify expectations */
	var L = NewListingRand()
}

func TestSaveAndLoadNewGraph(t *testing.T) {
	/* TODO: Decide whether preferred to overwrite files (simple) or deal with naming conventions */
	/* Start data Prep work */
	var G = new(Graph)
	var nodes []GraphNode
	var edges []GraphEdge

	anySmallNum := 4
	locations := getAirportLocations()
	locations = pick(anySmallNum, locations)
	intoPairs := 2
	trips := combine(locations, intoPairs)
	edges = makeEdgesFrom(trips)

	for n := range nodes {
		err = G.addNode(n)
		log(err)
	}
	for e := range edges {
		err = G.addEdge(e)
		log(err)
	}
	cacheDir := "test/cache/graph/"
	saveFile := cacheDir + "G0"
	G.save(saveFile)
	/* TODO: Test expected file exists */
	H := new(Graph)
	if G.equals(H) {
		t.FailNow()
	}

	H.load(saveFile)
	if G.notEqualTo(H) {
		t.FailNow()
	}
}
