package common

import "fmt"

// Graph is a weighted and digraph impl with focus on short path Traversals
// Nodes and edges are kept as slices, sorted for faster minCost selection
type Graph struct {
	nodes []node
	edges []edge
}
type node struct {
	value interface{}
	index int
}
type edge struct {
	start  int
	end    int
	weight int
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
func (G *graph) addEdge(e edge) error {
	if !(G.hasNode(e.start) && G.hasNode(e.end)) {
		panic("Missing Node")
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
func (G *graph) delEdge(e edge) error { return nil }
func (G *graph) delNode(n node) error { return nil } /* Also removes all edges with node */
func (G *graph) hasEdge(start, end int) bool {
	for _, any := range G.edges {
		if any.start == start && any.end == end {
			return true
		}
	}
	return false
}
func (G *graph) hasNode(index int) bool {
	return index < len(G.nodes)
}
func (G *graph) isConnected() bool   { return false }
func (G *graph) isTraversable() bool { return false }
func (G *graph) isHamiltonian() bool { return false }
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
	return fmt.Sprint(n.index)
}
func (E *edge) String() string {
	var direction string = ":"
	return fmt.Sprintf("%s%s%s", E.start, direction, E.end)
}
