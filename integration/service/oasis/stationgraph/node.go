package stationgraph

import (
	"math"

	"github.com/Telenav/osrm-backend/integration/service/oasis/chargingstrategy"
	"github.com/golang/glog"
)

type chargeInfo struct {
	arrivalEnergy float64
	chargeTime    float64
	targetState   chargingstrategy.State
}

type locationInfo struct {
	lat float64
	lon float64
}

type node struct {
	id nodeID
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
			targetState: chargingstrategy.State{
				Energy: 0.0,
			},
		},
	}
}

// reachableByDistance is used to test whether target distance is reachable by current status
func (n *node) reachableByDistance(distance float64) bool {
	return n.targetState.Energy > distance
}

// calcChargeTime calculates time effort needed from enter level to exist level
// For example
// prevNode records energy level when left previous charge station, say 100
// the distance between previous charge station to current charge station is 50
// for current node, we want to charge to certain energy level, represented by targetState, say 200
// So, charge time will calculate how much time needed to charge from (100 - 50) to 200
func (n *node) calcChargeTime(prev *node, distance float64, strategy chargingstrategy.Strategy) float64 {
	arrivalEnergy := prev.targetState.Energy - distance
	if arrivalEnergy < 0 {
		glog.Fatalf("Before updateNode should check reachableByDistance() prev.arrivalEnergy=%#v distance=%#v", prev.arrivalEnergy, distance)
	}
	return strategy.EvaluateCost(arrivalEnergy, n.targetState).Duration
}

func (n *node) updateChargingTime(chargingTime float64) {
	// @todo: maybe node record chargestate, hide charging time information in charge status
	n.chargeTime = chargingTime
}

func (n *node) updateArrivalEnergy(prev *node, distance float64) {
	n.arrivalEnergy = calculateArrivalEnergy(prev, distance)
}

func calculateArrivalEnergy(prev *node, distance float64) float64 {
	energy := prev.targetState.Energy - distance
	if energy < 0 {
		glog.Fatalf("Before updateNode should check reachableByDistance() prev.arrivalEnergy=%#v distance=%#v", prev.arrivalEnergy, distance)
	}
	return energy
}
