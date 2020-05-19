package stationgraph

import (
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/chargingstrategy"
)

func TestAddAndGetFunctionsForNodeContainer(t *testing.T) {
	input := []struct {
		stationID   string
		chargeState chargingstrategy.State
		location    nav.Location
	}{
		{
			"station1",
			chargingstrategy.State{
				Energy: 10.0,
			},
			nav.Location{
				Lat: 1.1,
				Lon: 1.1,
			},
		},
		{
			"station1",
			chargingstrategy.State{
				Energy: 20.0,
			},
			nav.Location{
				Lat: 1.1,
				Lon: 1.1,
			},
		},
		{
			"station1",
			chargingstrategy.State{
				Energy: 30.0,
			},
			nav.Location{
				Lat: 1.1,
				Lon: 1.1,
			},
		},
		{
			"station2",
			chargingstrategy.State{
				Energy: 10.0,
			},
			nav.Location{
				Lat: 2.2,
				Lon: 2.2,
			},
		},
		{
			"station2",
			chargingstrategy.State{
				Energy: 20.0,
			},
			nav.Location{
				Lat: 2.2,
				Lon: 2.2,
			},
		},
		{
			"station3",
			chargingstrategy.State{
				Energy: 15.0,
			},
			nav.Location{
				Lat: 3.3,
				Lon: 3.3,
			},
		},
	}

	expect := []struct {
		n         *node
		stationID string
	}{
		{
			&node{
				0,
				chargeInfo{
					targetState: chargingstrategy.State{
						Energy: 10.0,
					},
				},
				// nav.Location{
				// 	Lat: 1.1,
				// 	Lon: 1.1,
				// },
			},
			"station1",
		},
		{
			&node{
				1,
				chargeInfo{
					targetState: chargingstrategy.State{
						Energy: 20.0,
					},
				},
				// nav.Location{
				// 	Lat: 1.1,
				// 	Lon: 1.1,
				// },
			},
			"station1",
		},
		{
			&node{
				2,
				chargeInfo{
					targetState: chargingstrategy.State{
						Energy: 30.0,
					},
				},
				// nav.Location{
				// 	Lat: 1.1,
				// 	Lon: 1.1,
				// },
			},
			"station1",
		},
		{
			&node{
				3,
				chargeInfo{
					targetState: chargingstrategy.State{
						Energy: 10.0,
					},
				},
				// nav.Location{
				// 	Lat: 2.2,
				// 	Lon: 2.2,
				// },
			},
			"station2",
		},
		{
			&node{
				4,
				chargeInfo{
					targetState: chargingstrategy.State{
						Energy: 20.0,
					},
				},
				// nav.Location{
				// 	Lat: 2.2,
				// 	Lon: 2.2,
				// },
			},
			"station2",
		},
		{
			&node{
				5,
				chargeInfo{
					targetState: chargingstrategy.State{
						Energy: 15.0,
					},
				},
				// nav.Location{
				// 	Lat: 3.3,
				// 	Lon: 3.3,
				// },
			},
			"station3",
		},
	}

	if len(input) != len(expect) {
		t.Errorf("Incorrect test case, array of input and array of expect should 1-to-1 match.\n")
	}

	nc := newNodeContainer()

	for i := 0; i < len(input); i++ {
		tmpNode := nc.addNode(input[i].stationID, input[i].chargeState /*, input[i].location*/)

		if !reflect.DeepEqual(tmpNode, expect[i].n) {
			t.Errorf("Calling nodeContainer's addNode() generate incorrect result, expect %#v but got %#v.\n", expect[i].n, tmpNode)
		}
	}

	for i := 0; i < len(input); i++ {
		isVisited := nc.isNodeVisited((nodeID)(i))
		if !isVisited {
			t.Errorf("Calling nodeContainer's isNodeVisited() generate incorrect result, expect node %v is visited but returns no.\n", i)
		}

		tmpNode := nc.getNode((nodeID)(i))
		if !reflect.DeepEqual(tmpNode, expect[i].n) {
			t.Errorf("Calling nodeContainer's getNode() generate incorrect result, expect %#v but got %#v.\n", expect[i].n, tmpNode)
		}

		tmpStationID := nc.stationID((nodeID)(i))
		if !reflect.DeepEqual(tmpStationID, expect[i].stationID) {
			t.Errorf("Calling nodeContainer's stationID() generate incorrect result, expect %#v but got %#v.\n", expect[i].stationID, tmpStationID)
		}
	}

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 10; i++ {
		randomGen := func(min, max int) int {
			return rand.Intn(max-min) + min
		}
		unVisitedNodeID := randomGen(7, 100)
		isVisited := nc.isNodeVisited((nodeID)(unVisitedNodeID))
		if isVisited {
			t.Errorf("Calling nodeContainer's isNodeVisited() generate incorrect result, expect node %v is unvisited but returns yes.\n", unVisitedNodeID)
		}
	}

}

func TestAddDuplicateNodeForNodeContainer(t *testing.T) {
	input := []struct {
		stationID   string
		chargeState chargingstrategy.State
		location    nav.Location
	}{
		{
			"station1",
			chargingstrategy.State{
				Energy: 10.0,
			},
			nav.Location{
				Lat: 1.1,
				Lon: 1.1,
			},
		},
		{
			"station1",
			chargingstrategy.State{
				Energy: 10.0,
			},
			nav.Location{
				Lat: 1.1,
				Lon: 1.1,
			},
		},
		{
			"station1",
			chargingstrategy.State{
				Energy: 10.0,
			},
			nav.Location{
				Lat: 1.1,
				Lon: 1.1,
			},
		},
	}

	expect := struct {
		n         *node
		stationID string
	}{
		&node{
			0,
			chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 10.0,
				},
			},
			// nav.Location{
			// 	Lat: 1.1,
			// 	Lon: 1.1,
			// },
		},
		"station1",
	}

	nc := newNodeContainer()

	for i := 0; i < len(input); i++ {
		tmpNode := nc.addNode(input[i].stationID, input[i].chargeState /*, input[i].location*/)

		if !reflect.DeepEqual(tmpNode, expect.n) {
			t.Errorf("Calling nodeContainer's addNode() generate incorrect result, expect %#v but got %#v.\n", expect.n, tmpNode)
		}
	}

}
