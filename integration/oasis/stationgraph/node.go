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

func (n *node) isLocationReachable(distance float64) bool {
	return n.chargeEnergy > distance
	// return (n.arrivalEnergy + n.chargeEnergy) > distance
}

func (n *node) calcChargeTime(prev *node, distance float64, strategy chargingstrategy.ChargingStrategyCreator) float64 {
	arrivalEnergy := prev.arrivalEnergy - distance
	if arrivalEnergy < 0 {
		glog.Fatal("Before updateNode should check isLocationReachable()")
	}
	return strategy.EvaluateCost(arrivalEnergy, chargingstrategy.ChargingStrategy{ChargingEnergy: n.chargeEnergy}).Duration
}

func (n *node) updateChargingTime(chargingTime float64) {
	// @todo: maybe node record chargestate
	n.chargeTime = chargingTime
}

func (n *node) updateArrivalEnergy(prev *node, distance float64) {
	n.arrivalEnergy = prev.arrivalEnergy - distance
	if n.arrivalEnergy < 0 {
		glog.Fatal("Before updateNode should check isLocationReachable()")
	}
}
