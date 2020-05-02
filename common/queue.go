package common

/* TODO Research way of generalizing */

type queue interface {
	enQ()
	deQ()
}

type qnode interface {
	value() string
	next()
}

/* Will break with high volume */
type basicQ struct {
	container []qnode
	head      int
	tail      int
}

func newBasicQueue() *basicQ {
	var Q = new(basicQ)
	Q.container = make([]qnode, 1<<10)
	Q.head = 0 /* If head < tail, there exist enQ'd qnodes */
	Q.tail = 0
	return Q
}

func (Q *basicQ) enQ(n *qnode) {
	Q.tail++
	index := Q.tail % len(Q.container)
	Q.container[index] = *n
}

func (Q *basicQ) deQ() *qnode {
	var result *qnode

	if Q.head >= Q.tail {
		result = nil
	} else {
		Q.head++
		index := Q.head % len(Q.container)
		result = &(Q.container[index])
	}

	return result
}
