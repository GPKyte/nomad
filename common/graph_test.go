package common

import "testing"

func TestSaveAndLoadNewGraph(t *testing.T) {
	/* TODO: Decide whether preferred to overwrite files (simple) or deal with naming conventions */
	/* Start data Prep work */
	var G = new(graph)
	var nodes []node
	var edges []edge

	anySmallNum := 4
	locations := getAirportLocations()
	locations = pick(anySmallNum, locations)

	for counter, L := range locations {
		n := node{value: L, index: counter}
	}
	// Generate the simplized trips to fill the graph
	partA, partB := pairs(len(locations))
	for M := 0; M < len(partA); M++ {
		start := node{locations[partA[M]]}
		e := edge{
			start,
			locations[partB[M]],
		}
	}

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
