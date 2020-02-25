package stationgraph

import (
	"math"

	"github.com/golang/glog"
)

type chargeInfo struct {
	arrivalEnergy float64
	chargeTime    float64
	chargeEnergy  float64
}

type node struct {
	id    nodeID
	edges []*edge
	chargeInfo
}

type nodeID int32

const invalidNodeID = math.MinInt32

func newNode() *node {
	return &node{
		id: invalidNodeID,
		chargeInfo: chargeInfo{
			arrivalEnergy: 0.0,
			chargeTime:    0.0,
			chargeEnergy:  0.0,
		},
	}
}

func (n *node) isLocationReachable(distance float64) bool {
	return (n.arrivalEnergy + n.chargeEnergy) > distance
}

func (n *node) updateArrivalEnergy(prev *node, distance float64) {
	n.arrivalEnergy = prev.arrivalEnergy + prev.chargeEnergy - distance
	if n.arrivalEnergy < 0 {
		glog.Fatal("Before updateNode should check isLocationReachable()")
	}
}
