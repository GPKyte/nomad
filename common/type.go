package nomad

import (
	"fmt"
	"io"
	"strings"
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
	value interface{}
	edges []edge
	name  string
}
type edge struct {
	start      node
	end        node
	isDirected bool
}
type weightedEdge struct {
	edge
	weight float64
}
type timedEdge struct {
	edge
	time time.Time
}
type trip struct {
	timedEdge
	weightedEdge
	Listing
}
type timeSpace struct {
	DateTime time.Time
	Location string /* standardize to lowercase, UTF8 default */
}
type readable interface {
	String() string
	Read() string
}
type location struct {
	code string
	name string
}

func newGraphFrom(N []node, E []edge) *graph {
	G := new(graph)

	for _, ne := range N {
		G.addNode(ne)
	}
	for _, ew := range E {
		G.addEdge(ew)
	}
	return G
}
func (G *graph) addEdge(e edge) error        { return nil }
func (G *graph) addNode(n node) error        { return nil }
func (G *graph) anyNode() node               { return G.nodes[0] }
func (G *graph) delEdge(e edge) error        { return nil }
func (G *graph) delNode(n node) error        { return nil } /* Also removes all edges with node */
func (G *graph) checkEdgeExists(e edge) bool { return false }
func (G *graph) isConnected() bool           { return false }
func (G *graph) isTraversable() bool         { return false }
func (G *graph) isHamiltonian() bool         { return false }
func (G *graph) getHamiltonianCycle() []edge {
	var cycle []edge
	/* Traverse all Nodes at least, and preferably just once */
	if !(G.isConnected()) || !(G.isTraversable()) {
		cycle = make([]edge, 0)
	} else {
		cycle = G.doTraversal(G.anyNode(), true)
	}
	return cycle
}
func (G *graph) doTraversal(start node, finishAtStart bool) []edge {
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
	return make([]edge, 0, 0)
}
func (G *graph) save(relativePathToSaveFile string) error {
	/* What format would be proper for saving a graph? JSON all the way down...
	 * Would one file contain multiple graphs? No
	 * Would I include both edges and Vertices? Yes
	 * Can information be lossy or must all be retained? Lossless
	 * One line, N lines? One line preferred, can expand in tool later to pprint
	 * List of edges? Instead of node.[]neighbors
	 * How to represent and reinterpret Nodes?
	 * Will edge be NodeRepr -> NodeRepr?
	 * Make it easy to reload graph
	 */
	return nil
}
func (G *graph) load(relativePathToSaveFile string) error { return nil } /* Modifies State of G */
func (G *graph) deepCopy() *graph                         { return new(graph) }
func (G *graph) equals(H *graph) bool {
	/* Mostly care about edge set, but perhaps also node set */
	/* Cannot directly compare either as both are slices and that op/n is not supported */
	return false
}
func (G *graph) notEqualTo(H *graph) bool { return !G.equals(H) }
func (G *graph) String() string {
	return fmt.Sprintf("N:%s")
}
func (N *node) String() string {
	return fmt.Sprint(N.name)
}
func (E *edge) String() string {
	var direction string
	if E.isDirected {
		direction = "->"
	} else {
		direction = "<->"
	}
	return fmt.Sprintf("%s%s%s", E.start, direction, E.end)
}

type json map[string]interface{}

func newInvalidListing() Listing {
	return Listing{}
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

func makeScrapeStamp(srcURL string) timeSpace {
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
		gift = string(makeJSON(L.json()))
	} else {
		gift = L.csv()
	}
	return gift
}

func (L *Listing) csv() string {
	var csvRow = []string{
		string(L.Price),
		L.Depart.DateTime.String(),
		L.Depart.Location,
		L.Arrive.DateTime.String(),
		L.Arrive.Location,
		L.Scrape.DateTime.String(),
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

func makeJSON(srcJSONstruct json) []byte {
	return make([]byte, 0, 0)
}

func join(delim rune, row []string) string {
	return strings.Join(row, string(delim))
}

// NewListingsFromJSON is used to import Listing data that has been exported per standard
func NewListingsFromJSON(srcJSON io.Reader) []Listing {
	// TODO
	return make([]Listing, 0, 0)
}

// NewListingsFromCSV is used to import Listing data that has been exported per standard
func NewListingsFromCSV(srcCSV io.Reader) []Listing {
	// TODO
	return make([]Listing, 0, 0)
}

func newListingRand() Listing {
	const nanoConvRate = 10 ^ 9
	var (
		now       = time.Now().UTC()
		before    = time.Now().UTC()
		twentyMin = time.Duration(nanoConvRate * (60 * 20)) // 60 sec/min * 20 min
		after     = time.Now().UTC().Add(twentyMin)
		departLoc = "A"
		arriveLoc = "B"
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
func (L *Listing) Value() int {
	return int(L.Price)
}

func values(collection ...interface{ Value() int }) []int {
	var values = make([]int, len(collection))

	for _, v := range collection {
		next := v.Value()
		values = append(values, next)
	}

	return values
}

func max(collection []interface{ Value() int }) int {
	var max int = collection[0].Value()

	for _, val := range values(collection...) {

		if val > max {
			max = val
		}
	}

	return max
}
