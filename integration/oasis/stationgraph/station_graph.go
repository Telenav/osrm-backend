package stationgraph

import (
	"github.com/Telenav/osrm-backend/integration/oasis/chargingstrategy"
	"github.com/Telenav/osrm-backend/integration/oasis/stationfinder"
	"github.com/golang/glog"
)

type stationGraph struct {
	g                 *graph
	stationName2Nodes map[string][]*node
	stationID2Name    map[uint32]string
	nodesCount        uint32
	strategy          chargingstrategy.ChargingStrategyCreator
}

// NewStationGraph creates station graph from channel
func NewStationGraph(c chan stationfinder.WeightBetweenNeighbors, currEnergyLevel, maxEnergyLevel float64, strategy chargingstrategy.ChargingStrategyCreator) *stationGraph {
	sg := &stationGraph{
		g: &graph{
			startNodeID: invalidNodeID,
			endNodeID:   invalidNodeID,
		},
		stationName2Nodes: make(map[string][]*node),
		stationID2Name:    make(map[uint32]string),
		nodesCount:        0,
		strategy:          strategy,
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
	for _, fn := range sg.getChargeStationsNodes(neighborInfo.FromName, currEnergyLevel, maxEnergyLevel) {
		for _, tn := range sg.getChargeStationsNodes(neighborInfo.ToName, currEnergyLevel, maxEnergyLevel) {
			fn.neighbors = append(fn.neighbors, &neighbor{
				targetNodeID: tn.id,
				distance:     neighborInfo.Distance,
				duration:     neighborInfo.Duration,
			})
		}
	}
}

func (sg *stationGraph) getChargeStationsNodes(name string, currEnergyLevel, maxEnergyLevel float64) []*node {
	if _, ok := sg.stationName2Nodes[name]; !ok {
		if sg.isStart(name) {
			sg.constructStartNode(name, currEnergyLevel)
		} else if sg.isEnd(name) {
			sg.constructEndNode(name)
		} else {
			var nodes []*node
			for _, strategy := range sg.strategy.CreateChargingStrategies() {
				n := &node{
					id: nodeID(sg.nodesCount),
					chargeInfo: chargeInfo{
						chargeTime:   strategy.ChargingTime,
						chargeEnergy: strategy.ChargingEnergy,
					},
				}
				nodes = append(nodes, n)
				sg.stationID2Name[sg.nodesCount] = name
				sg.nodesCount += 1
			}

			sg.stationName2Nodes[name] = nodes
		}
	}
	return sg.stationName2Nodes[name]
}

func (sg *stationGraph) isStart(name string) bool {
	return name == stationfinder.OrigLocationName
}

func (sg *stationGraph) isEnd(name string) bool {
	return name == stationfinder.DestLocationName
}

func (sg *stationGraph) getName(id nodeID) string {
	return sg.stationID2Name[uint32(id)]
}

func (sg *stationGraph) constructStartNode(name string, currEnergyLevel float64) {

	n := &node{
		id:         nodeID(sg.nodesCount),
		chargeInfo: chargeInfo{arrivalEnergy: currEnergyLevel},
	}
	sg.stationName2Nodes[name] = []*node{n}
	sg.stationID2Name[sg.nodesCount] = name
	sg.nodesCount += 1
}

func (sg *stationGraph) constructEndNode(name string) {

	n := &node{
		id: nodeID(sg.nodesCount),
	}
	sg.stationName2Nodes[name] = []*node{n}
	sg.stationID2Name[sg.nodesCount] = name
	sg.nodesCount += 1
}

func (sg *stationGraph) constructGraph() *stationGraph {
	for k, v := range sg.stationName2Nodes {
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
	} else if len(sg.g.nodes)-1 != int(sg.nodesCount) {
		glog.Error("Invalid nodes generated, len(sg.g.nodes) is %d while sg.nodesCount is %d.\n", len(sg.g.nodes), sg.nodesCount)
		return nil
	}

	return sg
}
