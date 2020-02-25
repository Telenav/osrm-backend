package stationgraph

import (
	"reflect"
	"testing"
)

func popElementsInPQAndCompareResult(t *testing.T, pq *priorityQueue, expect []nodeID) {
	var actual []nodeID
	for {
		if pq.empty() {
			break
		}
		v := pq.pop()
		actual = append(actual, v)
	}
	if !reflect.DeepEqual(expect, actual) {
		t.Errorf("expect %v but got %v", expect, actual)
	}
}

func TestPQGeneral(t *testing.T) {
	pq := newPriorityQueue()
	pq.push(1, 10)
	pq.push(2, 20)
	pq.push(3, 30)
	pq.push(4, 15)
	pq.push(5, 25)

	expect := []nodeID{1, 4, 2, 5, 3}
	popElementsInPQAndCompareResult(t, pq, expect)
}

func TestPQUpdate(t *testing.T) {
	id2index := make(map[nodeID]*pqElement)

	pq := newPriorityQueue()
	id2index[1] = pq.push(1, 10)
	id2index[2] = pq.push(2, 20)
	id2index[3] = pq.push(3, 30)
	id2index[4] = pq.push(4, 15)
	id2index[5] = pq.push(5, 25)

	pq.decrease(id2index[3], 1)
	expect := []nodeID{3, 1, 4, 2, 5}
	popElementsInPQAndCompareResult(t, pq, expect)

}

func TestPQUpdate2(t *testing.T) {
	id2index := make(map[nodeID]*pqElement)

	pq := newPriorityQueue()
	id2index[1] = pq.push(1, 10)
	id2index[2] = pq.push(2, 20)
	id2index[3] = pq.push(3, 30)
	id2index[4] = pq.push(4, 15)
	id2index[5] = pq.push(5, 25)

	pq.decrease(id2index[3], 1)
	pq.decrease(id2index[5], 2)
	expect := []nodeID{3, 5, 1, 4, 2}
	popElementsInPQAndCompareResult(t, pq, expect)
}
