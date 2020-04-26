package util

/* TODO Research way of generalizing */

type queue interface {
	enQ()
	deQ()
}

type node interface {
	value() string
	next()
}

/* Will break with high volume */
type basicQ struct {
	container []node
	head      int
	tail      int
}

func newBasicQueue() *basicQ {
	var Q = new(basicQ)
	Q.container = make([]node, 1<<10)
	Q.head = 0 /* If head < tail, there exist enQ'd nodes */
	Q.tail = 0
	return Q
}

func (Q *basicQ) enQ(n *node) {
	Q.tail++
	index := Q.tail % len(Q.container)
	Q.container[index] = *n
}

func (Q *basicQ) deQ() *node {
	var result *node

	if Q.head >= Q.tail {
		result = nil
	} else {
		Q.head++
		index := Q.head % len(Q.container)
		result = &(Q.container[index])
	}

	return result
}
