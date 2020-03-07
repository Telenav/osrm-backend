package stationgraph

import "github.com/golang/glog"

type queryHeapNodeInfo struct {
	prevNodeID nodeID
	pqElem     *pqElement
	minWeight  float64
	minDist    float64
	settled    bool
}

type queryHeap struct {
	pq *priorityQueue
	m  map[nodeID]*queryHeapNodeInfo
}

func newQueryHeap() *queryHeap {
	return &queryHeap{
		pq: newPriorityQueue(),
		m:  make(map[nodeID]*queryHeapNodeInfo),
	}
}

func (qh *queryHeap) add(currID, prevID nodeID, distance, duration float64) bool {
	newWeight := 0.0
	newDist := 0.0
	if prevID != invalidNodeID {
		newWeight = qh.m[prevID].minWeight + duration
		newDist = qh.m[prevID].minDist + distance
	} else {
		// Nothing to do. currID is start node.
	}

	if !qh.isVisited(currID) {
		e := qh.pq.push(currID, newWeight)
		glog.V(3).Infof("pq-push new element %#v with weight %#v\n", currID, newWeight)
		qh.m[currID] = &queryHeapNodeInfo{
			prevNodeID: prevID,
			pqElem:     e,
			minWeight:  newWeight,
			minDist:    newDist,
			settled:    false,
		}
		glog.V(3).Infof("query-heap insert new element %#v for %#v\n", qh.m[currID], currID)
		return true
	} else {
		if ok := qh.needUpdate(currID, newWeight); ok {
			if qh.isSettled(currID) {
				glog.Warning("Check your logic, settled node should not have smaller weight for dijkstra.")
			}
			qh.pq.decrease(qh.m[currID].pqElem, newWeight)
			qh.update(currID, prevID, newWeight, newDist)
			glog.V(3).Infof("pq-update element %#v with weight %#v\n", currID, newWeight)
			return true
		}
	}
	return false
}

func (qh *queryHeap) next() nodeID {
	if qh.pq.empty() {
		return invalidNodeID
	}

	n := qh.pq.pop()
	glog.V(3).Infof("pq-pop element %#v\n", n)
	qh.settle(n)
	return n
}

// node id list: invalidNodeID -> start -> mid1 -> mid2 -> end
// will return {mid1, mid2}
func (qh *queryHeap) retrieve(endNodeID nodeID) []nodeID {
	var r []nodeID
	if !qh.isSettled(endNodeID) {
		return r
	}

	currID := endNodeID
	for {
		currV, _ := qh.m[currID]
		if currV.prevNodeID == invalidNodeID {
			return r
		}
		prevV, _ := qh.m[currV.prevNodeID]
		if prevV.prevNodeID == invalidNodeID {
			for i := len(r)/2 - 1; i >= 0; i-- {
				oppsiteIndex := len(r) - 1 - i
				r[i], r[oppsiteIndex] = r[oppsiteIndex], r[i]
			}
			return r
		}
		r = append(r, currV.prevNodeID)
		currID = currV.prevNodeID
	}
}

func (qh *queryHeap) isVisited(id nodeID) bool {
	_, ok := qh.m[id]
	return ok
}

func (qh *queryHeap) isSettled(id nodeID) bool {
	v, ok := qh.m[id]
	if !ok {
		glog.Fatal("Check your logic, isSettled() should be called when isVisited() returns true")
		return false
	} else {
		return v.settled
	}
}

func (qh *queryHeap) needUpdate(id nodeID, weight float64) bool {
	v, ok := qh.m[id]
	if !ok {
		glog.Fatal("Check your logic, needUpdate() should be called when isVisited() returns true")
		return true
	} else {
		return v.minWeight > weight
	}
}

func (qh *queryHeap) update(id, prevNodeID nodeID, weight, dist float64) {
	v, ok := qh.m[id]
	if !ok {
		glog.Fatal("Check your logic, update() should be called when isVisited() returns true")
		return
	} else {
		v.minWeight = weight
		v.minDist = dist
		v.prevNodeID = prevNodeID
	}
}

func (qh *queryHeap) settle(id nodeID) {
	if v, ok := qh.m[id]; ok {
		if v.settled {
			glog.Warningf("Check your logic, settle() should be called with unsettled node(%b)", id)
		}
		v.settled = true
	} else {
		glog.Fatalf("Check your logic, settle() should be called with visited node(%b)", id)
	}
}
