package stationgraph

import (
	"container/heap"
	"math"

	"github.com/golang/glog"
)

type priorityQueue struct {
	pqImpl *priorityQueueImpl
}

func newPriorityQueue() *priorityQueue {
	return &priorityQueue{
		pqImpl: &priorityQueueImpl{},
	}
}

func (pq *priorityQueue) push(id nodeID, weight float64) *pqElement {
	e := &pqElement{
		nodeID: id,
		weight: weight,
	}
	heap.Push(pq.pqImpl, e)
	return e
}

func (pq *priorityQueue) pop() nodeID {
	e := heap.Pop(pq.pqImpl)
	return e.(*pqElement).nodeID
}

func (pq *priorityQueue) decrease(e *pqElement, newWeight float64) {
	e.weight = newWeight
	heap.Fix(pq.pqImpl, e.index)
}

func (pq *priorityQueue) empty() bool {
	return len(*(pq.pqImpl)) == 0
}

type pqElement struct {
	nodeID nodeID
	weight float64
	index  int
}

const invalidIndex = math.MaxUint32

type priorityQueueImpl []*pqElement

func (pq priorityQueueImpl) Len() int           { return len(pq) }
func (pq priorityQueueImpl) Less(i, j int) bool { return pq[i].weight < pq[j].weight }
func (pq priorityQueueImpl) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueueImpl) Push(x interface{}) {
	n := len(*pq)
	item := x.(*pqElement)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *priorityQueueImpl) Pop() interface{} {
	n := len(*pq)
	if n == 0 {
		return nil
	}

	old := *pq
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

func (pq *priorityQueueImpl) get(i int) *pqElement {
	if i < 0 || i > len(*pq)-1 {
		glog.Fatalf("Query PQ with invalid index %v while size of pq is %v", i, len(*pq))
	}
	return (*pq)[i]
}
