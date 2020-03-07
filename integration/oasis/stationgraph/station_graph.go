package stationgraph

import (
	"github.com/Telenav/osrm-backend/integration/oasis/chargingstrategy"
	"github.com/Telenav/osrm-backend/integration/oasis/stationfinder"
	"github.com/golang/glog"
)

type stationGraph struct {
	g               *graph
	stationID2Nodes map[string][]*node
	num2StationID   map[uint32]string // from number to original stationID
	// stationID is converted to numbers(0, 1, 2 ...) based on visit sequence

	stationsCount uint32
	strategy      chargingstrategy.ChargingStrategyCreator
}

// NewStationGraph creates station graph from channel
func NewStationGraph(c chan stationfinder.WeightBetweenNeighbors, currEnergyLevel, maxEnergyLevel float64, strategy chargingstrategy.ChargingStrategyCreator) *stationGraph {
	sg := &stationGraph{
		g: &graph{
			startNodeID: invalidNodeID,
			endNodeID:   invalidNodeID,
		},
		stationID2Nodes: make(map[string][]*node),
		num2StationID:   make(map[uint32]string),
		stationsCount:   0,
		strategy:        strategy,
	}

	for item := range c {
		if item.Err != nil {
			glog.Errorf("Met error during constructing stationgraph, error = %v", item.Err)
			return nil
		}

		for _, neighborInfo := range item.NeighborsInfo {
			sg.buildNeighborInfoBetweenNodes(neighborInfo, currEnergyLevel, maxEnergyLevel)
		}
	}

	return sg.constructGraph()
}

func (sg *stationGraph) buildNeighborInfoBetweenNodes(neighborInfo stationfinder.NeighborInfo, currEnergyLevel, maxEnergyLevel float64) {
	for _, fromNode := range sg.getChargeStationsNodes(neighborInfo.FromID, currEnergyLevel, maxEnergyLevel) {
		for _, toNode := range sg.getChargeStationsNodes(neighborInfo.ToID, currEnergyLevel, maxEnergyLevel) {
			fromNode.neighbors = append(fromNode.neighbors, &neighbor{
				targetNodeID: toNode.id,
				distance:     neighborInfo.Distance,
				duration:     neighborInfo.Duration,
			})
		}
	}
}

func (sg *stationGraph) getChargeStationsNodes(id string, currEnergyLevel, maxEnergyLevel float64) []*node {
	if _, ok := sg.stationID2Nodes[id]; !ok {
		if sg.isStart(id) {
			sg.constructStartNode(id, currEnergyLevel)
		} else if sg.isEnd(id) {
			sg.constructEndNode(id)
		} else {
			var nodes []*node
			for _, strategy := range sg.strategy.CreateChargingStrategies() {
				n := &node{
					id: nodeID(sg.stationsCount),
					chargeInfo: chargeInfo{
						chargeEnergy: strategy.ChargingEnergy,
					},
				}
				nodes = append(nodes, n)
				sg.num2StationID[sg.stationsCount] = id
				sg.stationsCount += 1
			}

			sg.stationID2Nodes[id] = nodes
		}
	}
	return sg.stationID2Nodes[id]
}

func (sg *stationGraph) isStart(id string) bool {
	return id == stationfinder.OrigLocationID
}

func (sg *stationGraph) isEnd(id string) bool {
	return id == stationfinder.DestLocationID
}

func (sg *stationGraph) getStationID(id nodeID) string {
	return sg.num2StationID[uint32(id)]
}

func (sg *stationGraph) constructStartNode(id string, currEnergyLevel float64) {

	n := &node{
		id:         nodeID(sg.stationsCount),
		chargeInfo: chargeInfo{arrivalEnergy: currEnergyLevel},
	}
	sg.stationID2Nodes[id] = []*node{n}
	sg.num2StationID[sg.stationsCount] = id
	sg.stationsCount += 1
}

func (sg *stationGraph) constructEndNode(id string) {

	n := &node{
		id: nodeID(sg.stationsCount),
	}
	sg.stationID2Nodes[id] = []*node{n}
	sg.num2StationID[sg.stationsCount] = id
	sg.stationsCount += 1
}

func (sg *stationGraph) constructGraph() *stationGraph {
	for k, v := range sg.stationID2Nodes {
		if sg.isStart(k) {
			sg.g.startNodeID = v[0].id
		}

		if sg.isEnd(k) {
			sg.g.endNodeID = v[0].id
		}

		for _, n := range v {
			sg.g.nodes = append(sg.g.nodes, n)
		}
	}

	if sg.g.startNodeID == invalidNodeID {
		glog.Error("Invalid nodeid generated for start node.\n")
		return nil
	} else if sg.g.endNodeID == invalidNodeID {
		glog.Error("Invalid nodeid generated for start node.\n")
		return nil
	} else if len(sg.g.nodes) != int(sg.stationsCount) {
		glog.Errorf("Invalid nodes generated, len(sg.g.nodes) is %d while sg.stationsCount is %d.\n", len(sg.g.nodes), sg.stationsCount)
		return nil
	}

	return sg
}
