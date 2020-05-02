package common

import (
	"fmt"
)

// placeholder is a placeholder for the time being
type placeholder interface{}

type node struct {
	value placeholder
	label int
}
type nodeList []node
type edge struct {
	start  int
	end    int
	weight int
}

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
	nodes    []node
	edges    map[int][]edge
	shortcut map[placeholder]node
}

func (G *graph) getNode(fromReferenceToThis placeholder) (node, error) {
	if G.has(fromReferenceToThis) {
		return G.shortcut[fromReferenceToThis], nil
	}
	return node{}, fmt.Errorf("Could not find Node: %s", fromReferenceToThis)
}

func (G *graph) get(fromThisValue placeholder) (label int) {
	return G.shortcut[fromThisValue].label
}

func newGraphFrom(rawNodes []placeholder, rawEdges []Trip) *graph {
	G := new(graph)
	G.nodes = make([]node, len(rawNodes))
	G.edges = make(map[int][]edge)
	G.shortcut = make(map[placeholder]node)

	for i, ne := range rawNodes {
		G.addNode(node{ne, i})
	}
	for _, ew := range rawEdges {
		G.addEdge(G.get(ew.depart), G.get(ew.arrive), int(ew.cost))
	}
	return G
}

// addEdge between ALREADY EXISTING nodes referenced by their labels
func (G *graph) addEdge(to, from, weight int) error {
	if !G.has(to) || !G.has(from) {
		/* TODO: Save and display context of Graph at this moment and log it */
		return fmt.Errorf("Missing at least one node to add new edge")
	}

	/* Init empty list */
	if G.edges[to] == nil || len(G.edges[to]) == 0 {
		G.edges[to] = make([]edge, 10) /* Choose a low capacity based on current expectations, can always grow out as needed */
	}

	G.edges[to] = append(G.edges[to], edge{from, to, weight})

	return nil
}
func (G *graph) addNode(n node) error {
	G.nodes[n.label] = n
	G.shortcut[n.value] = n

	return nil
}
func (G *graph) getAnyNode() node { return G.nodes[0] }
func (G *graph) getByValue(want interface{}) int {
	for _, node := range G.nodes {
		if node.value == want {
			return node.label
		}
	}
	return -1
}

func (G *graph) hasEdge(start, end int) bool {
	for _, any := range G.edges[start] {
		if any.end == end {
			return true
		}
	}
	return false
}

func (G *graph) has(some interface{}) bool {
	switch some.(type) {
	case placeholder:
		return G.shortcut[some.(placeholder)].label >= 0
	case node:
		return G.nodes[some.(node).label].label >= 0
	case edge:
		e := some.(edge)
		return G.hasEdge(e.start, e.end)
	default:
		return false
	}
}

func (G *graph) isConnected() bool { return false }

//DEPRECATED
func (G *graph) getHamiltonianCycle() []edge {
	var cycle []edge
	/* Traverse all Nodes at least, and preferably just once */
	if !G.isConnected() {
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
	return fmt.Sprint(N.label)
}

func (E *edge) String() string {
	var direction string = "->"
	return fmt.Sprintf("%s%s%s", E.start, direction, E.end)
}
