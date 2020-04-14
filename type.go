package nomad

import (
	"io"
	"time"
)

/*
Listing contains the necessary data to describe an online travel listing
Which may be purchasable for a price, and may have layovers
*/
type Listing struct {
	Depart     timeSpace
	Arrive     timeSpace
	Scrape     timeSpace
	Price      int16
	srcContent string
}
/* Graph is a weighted and digraph impl with focus on short path Traversals
 * Nodes and edges are kept as slices, sorted for faster minCost selection */
 type graph struct {
	nodes []node
	edges []edge
}
type weightedGraph struct {
	graph
	edges []weightedEdge
}
type node struct {
	Stringer
	value interface{}
	edges []edges
	name string
}
type edge struct {
	start	node
	end		node
	isDirected bool
}
type weightedEdge struct {
	edge
	weight	float64
}
type timedEdge struct {
	edge
	time	time.Time
}
type trip struct {
	timedEdge
	weightedEdge
	Listing
}
type timeSpace struct {
	DateTime time.Time
	Location string /* standardized to lowercase, UTF8 default */
}
type readable interface {
	String() string
	Read() string
}
type location struct {
	node
	Stringer
	code	string
	name	string
}
func newGraphFrom(N []node, E []edge) *graph {
	G := new(graph)
	for n := range N {
		G.addNode(n)
	}
	for e := range E {
		G.addEdge(e)
	}
	return G
}
func (G *graph) addEdge(e edge) error {return nil}
func (G *graph) addNode(n node) error {return nil}
func (G *graph) delEdge(e edge) error {return nil}
func (G *graph) delNode(n node) error {return nil} /* Also removes all edges with node */
func (G *graph) checkEdgeExists(e Edge) bool {return false}
func (G *graph) isConnected() bool {return false}
func (G *graph) isTraversable() bool {return false}
func (G *graph) isHamiltonian() bool {return false}
func (G *graph) getHamiltonianCycle() cycle []edge {
	/* Traverse all Nodes at least, and preferably just once */
	if !(G.isConnected()) || !(G.isTraversable()) {
		cycle := []
	} else {
		cycle := G.doTraversal()
	}
	return cycle
}
func (G *graph)	doTraversal(start, finishAtStart) []edge {
	/* Look at starting Node
	 * Find incident nodes (neighbors)
	 * Keep +/- count of unique visited to verify "doneness"
	 * Choose edge's best neighbor
	 * Until dead-end or goal is reached choose edge's best neighbor and repeat
	 * If dead-end (all neighbors are visited) reached,
	 * go back to last branch and try next best neighbor, remember to count--
	 * If goal is reached (all nodes visited + finish node terms are met)
	 * Then return traversal path
	 * If count is reached but finish node terms are not met
	 * Then log path and save graph, inspect later and see if we can solve it
	 * Option to continue all paths remaining until terms are met
	 */
	return make([]edge, 0, 0)
}
func (G *graph) path(A node, Z node) []edge {
	/* Find nodes connected to Z
	 * Is A included?
	 * Yes – done
	 * No – find nodes with any of Z's neighbors
	 * Repeat searching edges for A with each neighbor
	 * considering parsing down edge set to save time on each subsequent pass
	 */
	return make([]edge,0,0)
}
func (G *graph) save(relativePathToSaveFile) error {
	/* What format would be proper for saving a graph?
	 * Would one file contain multiple graphs?
	 * Would I include both edges and Vertices?
	 * Can information be lossy or must all be retained?
	 * One line, N lines
	 * List of edges?
	 * How to represent and reinterpret Nodes?
	 * Will edge be NodeRepr -> NodeRepr?
	 * Make it easy to reload graph
	 */
	return nil
}
func (G *graph) load(relativePathToSaveFile) error {return nil} /* Modifies State of G */
func (G *graph) deepCopy() *graph {return new(graph)}
func String(G *graph) string {
	/* Assume Nodes and Edges are Stringers */
}

type json map[string]interface{}

func newInvalidListing() Listing {
	return Listing{}
}

type ValidDataObj interface {
	isNil()
	isValid()
}

func (T *timeSpace) isNil() bool {
	return false /* TODO: Check revamped timeSpace struct for nils or bad formatting */
}

func (T *timeSpace) isValid() bool {
	return false
}

func (L *Listing) isNil() (whetherAnyNilData bool) {
	whetherAnyNilData = (L.Price == 0 || L.Depart.isNil() || L.Arrive.isNil() || L.Scrape.isNil())

	return
}

/* TODO: change all isValid to isInvalid for consistent negative testing */
func (L *Listing) isValid() bool {
	return false
}

/* TODO: unexport default methods once testing confirms okay. fmt implies timeSpace should be exported if Record* is */

func makeScrapeStamp(srcURL string) {
	stamp := recordCurrentTimeSpace(srcURL)
	return stamp
}

func recordCurrentTimeSpace(location string) timeSpace {
	return timeSpace{
		DateTime: time.Now().UTC(),
		Location: location,
	}
}

// String will return JSON Repr of Listing Or Flat Repr
func (L *Listing) String() string {
	var gift string

	if weWantDefault := true; weWantDefault {
		gift = makeJSONString(L.json())
	} else {
		gift = L.csv()
	}
	return gift
}

func (L *Listing) csv() string {
	var csvRow = []string{
		L.Price,
		L.Depart.DateTime,
		L.Depart.Location,
		L.Arrive.DateTime,
		L.Arrive.Location,
		L.Scrape.DateTime,
		L.Scrape.Location,
	}

	return join(',', csvRow)
}

func (L *Listing) json() json {

	var jsonRepr = json{
		"Price":      L.Price,
		"DepartTime": L.Depart.DateTime,
		"DepartLoc":  L.Depart.Location,
		"ArriveTime": L.Arrive.DateTime,
		"ArriveLoc":  L.Arrive.Location,
		"ScrapeTime": L.Scrape.DateTime,
		"ScrapeURL":  L.Scrape.Location,
	}

	return jsonRepr
}

func makeJSONString(srcJSON json) {

}

func join(delim rune, csvRow []string) {

}

// NewListingsFromJSON is used to import Listing data that has been exported per standard
func NewListingsFromJSON(srcJSON io.Reader) []Listing {
	// TODO
}

// NewListingsFromCSV is used to import Listing data that has been exported per standard
func NewListingsFromCSV(srcCSV io.Reader) []Listing {
	// TODO
}

func NewListingRand() Listing {
	const nanoConvRate = 10 ^ 9
	var (
		now       = time.Now().UTC()
		before    = time.Now().UTC()
		twentyMin = time.Duration(nanoConvRate * (60 * 20)) // 60 sec/min * 20 min
		after     = time.Now().UTC().Add(twentyMin)
		departLoc = chooseLoc(5)
		arriveLoc = chooseLoc(9)
		url       = "https://random.local"
	)

	return Listing{
		Price:  50,
		Depart: timeSpace{before, departLoc},
		Arrive: timeSpace{after, arriveLoc},
		Scrape: timeSpace{now, url},
	}
}

// Value returns that which should be used in comparisons, this is the Listing price
func Value(L *Listing) int {
	return L.Price()
}

func values(collection ...interface{ Value() }) []int {
	var values = make([]int, len(collection))

	for v := range collection {
		values = append(values, Value(v))
	}
}

func max(collection []interface{ Value() }) int {
	top_value = Value(collection[0])

	for val := range values(collection) {

		if val > top_value {
			top_value = val
		}
	}

	return top_value
}
