package stationgraph

import (
	"math"

	"github.com/Telenav/osrm-backend/integration/oasis/chargingstrategy"
	"github.com/golang/glog"
)

type chargeInfo struct {
	arrivalEnergy float64
	chargeTime    float64
	chargeEnergy  float64
}

type locationInfo struct {
	lat float64
	lon float64
}

type node struct {
	id        nodeID
	neighbors []*neighbor
	chargeInfo
	locationInfo
}

type nodeID uint32

const invalidNodeID = math.MaxUint32

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

// Function isLocationReachable is used to test whether target node is reachable
func (n *node) isLocationReachable(distance float64) bool {
	return n.chargeEnergy > distance
}

// calcChargeTime calculates time effort needed from previous final status to current
func (n *node) calcChargeTime(prev *node, distance float64, strategy chargingstrategy.ChargingStrategyCreator) float64 {
	arrivalEnergy := prev.chargeEnergy - distance
	if arrivalEnergy < 0 {
		glog.Fatalf("Before updateNode should check isLocationReachable() prev.arrivalEnergy=%#v distance=%#v", prev.arrivalEnergy, distance)
	}
	return strategy.EvaluateCost(arrivalEnergy, chargingstrategy.State{ChargingEnergy: n.chargeEnergy}).Duration
}

func (n *node) updateChargingTime(chargingTime float64) {
	// @todo: maybe node record chargestate
	n.chargeTime = chargingTime
}

func (n *node) updateArrivalEnergy(prev *node, distance float64) {
	n.arrivalEnergy = prev.chargeEnergy - distance
	if n.arrivalEnergy < 0 {
		glog.Fatalf("Before updateNode should check isLocationReachable() prev.arrivalEnergy=%#v distance=%#v", prev.arrivalEnergy, distance)
	}
}
