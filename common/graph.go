package common

import (
	"fmt"
	"sort"
)

/* Design Questions
 * Is the graph weighted? Yes
 * Is there custom traversal logic? Not in this impl
 * How will the s=structure be used? s will find paths of low cost, and provide them for review
 * What algorithms support Pathfinding? A*
 * Will my graph be sparse or dense? Sparse
 * We will use a adjacency list, i.e. edges go in a (sorted) list, AND into
 *
 * What algorithm can identify which airports are most heavily traversed and in-a-way useful?
 * How should I tag my data to know if it was crafted to answer specific questions?
 * 		Just that, a tag. Easy to add a field to a datastruct and Marshall that to JSON or the DB we don't have
 * How should I deal with non-standard names? */

// Graph is a weighted and digraph impl with focus on short path Traversals
// Nodes and edges are kept as slices, sorted for faster minCost selection
type graph struct {
	nodes []nodeList
	edges []edge
	shortcut map[TimeAndPlace]node
}

func (G *graph) getNode(fromReferenceToThis TimeAndPlace) (node, error) {
	if G.has(fromReferenceToThis) {
		return G.shortcut[fromReferenceToThis]
	}
	return fmt.Errorf("Could not find Node: %s", fromReferenceToThis)
}

func ()

func newGraphFrom(rawNodes []TimeAndPlace, rawEdges []trip) *graph {
	G := new(graph)

	for i, ne := range rawNodes {
		G.addNode(node{value: ne, index: i})
	}
	for _, ew := range rawEdges {
		G.addEdge(G.get(ew.depart), G.get(ew.arrive), ew.price)
	}
	return G
}
func (G *graph) addEdge(e edge) error {
	if !(G.hasNode(e.start) || !G.hasNode(e.end)) {
		panic("Missing Node to Edge")
	}

	G.edges = append(G.edges, e)
	return nil
}
func (G *graph) addNode(n node) error {
	n.index = len(G.nodes)
	G.nodes = append(G.nodes, n)

	if G.nodes[n.index] != n {
		panic("Location of node in graph is errant")
	}

	return nil
}
func (G *graph) getAnyNode() node { return G.nodes[0] }
func (G *graph) getByValue(want interface{}) int {
	for _, node := range G.nodes {
		if node.value == want {
			return node.index
		}
	}
	return -1
}

func (G *graph) hasEdge(start, end int) bool {
	for _, any := range G.edges {
		if any.start == start && any.end == end {
			return true
		}
	}
	return false
}

func (G *graph) has(some interface{}) bool {
	switch some.(type) {
	case TimeAndPlace:
		return G.shortcut[some.(TimeAndPlace)]
	case node:
		return G.nodes[some.(node).index]
	case edge:
		e := some.(edge)
		return hasEdge(e.start.index, e.end.index)
	default:
		return false
	}
}

func (G *graph) isConnected() bool { return false }
func (G *graph) getHamiltonianCycle() []edge {
	var cycle []edge
	/* Traverse all Nodes at least, and preferably just once */
	if !(G.isConnected()) || !(G.isTraversable()) {
		cycle = make([]edge, 0)
	} else {
		cycle = G.doTraversal(G.getAnyNode(), true)
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
	Q = new(Queue)

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
	return fmt.Sprint(N.index)
}

type node struct {
	value TimeAndPlace
	index int
}

type nodeList []node

func (slice nodeList) Len() int {
	return len(slice)
}
func (slice nodeList) Less(i, j int) bool {
	return slice[i] < slice[j]
}
func (slice nodeList) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

type edge struct {
	start  int
	end    int
	weight int
}

func (E *edge) String() string {
	var direction string = "->"
	return fmt.Sprintf("%s%s%s", E.start, direction, E.end)
}
