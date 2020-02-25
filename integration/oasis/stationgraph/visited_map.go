package stationgraph

import "github.com/golang/glog"

type visitedNodeInfo struct {
	prevNodeID nodeID
	pqElem     *pqElement
	minCost    float64
	minDist    float64
	settled    bool
}

type visitedMap struct {
	pq *priorityQueue
	m  map[nodeID]*visitedNodeInfo
}

func newVisitedMap() *visitedMap {
	return &visitedMap{
		pq: newPriorityQueue(),
		m:  make(map[nodeID]*visitedNodeInfo),
	}
}

func (vm *visitedMap) add(currID, prevID nodeID, distance, duration float64) bool {
	newCost := 0.0
	newDist := 0.0
	if prevID != invalidNodeID {
		newCost = vm.m[prevID].minCost + duration
		newDist = vm.m[prevID].minDist + distance
	}

	if !vm.isVisited(currID) {
		e := vm.pq.push(currID, newCost)
		vm.m[currID] = &visitedNodeInfo{
			prevNodeID: prevID,
			pqElem:     e,
			minCost:    newCost,
			minDist:    newDist,
			settled:    false,
		}
		return true
	} else {
		if ok := vm.needUpdate(currID, newCost); ok {
			if vm.isSettled(currID) {
				glog.Warning("Check your logic, settled node should not have smaller cost for dijkstra.")
			}
			vm.pq.decrease(vm.m[currID].pqElem, newCost)
			vm.update(currID, prevID, newCost, newDist)
			return true
		}
	}
	return false
}

func (vm *visitedMap) next() nodeID {
	if vm.pq.empty() {
		return invalidNodeID
	}

	n := vm.pq.pop()
	vm.settle(n)
	return n
}

// node id list: invalidNodeID -> start -> mid1 -> mid2 -> end
// will return {mid1, mid2}
func (vm *visitedMap) retrieve(endNodeID nodeID) []nodeID {
	var r []nodeID
	if !vm.isSettled(endNodeID) {
		return r
	}

	currID := endNodeID
	for {
		currV, _ := vm.m[currID]
		if currV.prevNodeID == invalidNodeID {
			return r
		}
		prevV, _ := vm.m[currV.prevNodeID]
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

func (vm *visitedMap) isVisited(id nodeID) bool {
	_, ok := vm.m[id]
	return ok
}

func (vm *visitedMap) isSettled(id nodeID) bool {
	v, ok := vm.m[id]
	if !ok {
		glog.Fatal("Check your logic, isSettled() should be called when isVisited() returns true")
		return false
	} else {
		return v.settled
	}
}

func (vm *visitedMap) needUpdate(id nodeID, cost float64) bool {
	v, ok := vm.m[id]
	if !ok {
		glog.Fatal("Check your logic, needUpdate() should be called when isVisited() returns true")
		return true
	} else {
		return v.minCost > cost
	}
}

func (vm *visitedMap) update(id, prevNodeID nodeID, cost, dist float64) {
	v, ok := vm.m[id]
	if !ok {
		glog.Fatal("Check your logic, update() should be called when isVisited() returns true")
		return
	} else {
		v.minCost = cost
		v.minDist = dist
		v.prevNodeID = prevNodeID
	}
}

func (vm *visitedMap) settle(id nodeID) {
	if v, ok := vm.m[id]; ok {
		if v.settled {
			glog.Warningf("Check your logic, settle() should be called with unsettled node(%b)", id)
		}
		v.settled = true
	} else {
		glog.Fatalf("Check your logic, settle() should be called with visited node(%b)", id)
	}
}
