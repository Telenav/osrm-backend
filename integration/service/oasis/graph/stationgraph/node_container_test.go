package stationgraph

import (
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/graph/chargingstrategy"
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/entity"
)

func TestAddAndGetFunctionsForNodeContainer(t *testing.T) {
	input := []struct {
		placeID     entity.PlaceID
		chargeState chargingstrategy.State
		location    nav.Location
	}{
		{
			1,
			chargingstrategy.State{
				Energy: 10.0,
			},
			nav.Location{
				Lat: 1.1,
				Lon: 1.1,
			},
		},
		{
			1,
			chargingstrategy.State{
				Energy: 20.0,
			},
			nav.Location{
				Lat: 1.1,
				Lon: 1.1,
			},
		},
		{
			1,
			chargingstrategy.State{
				Energy: 30.0,
			},
			nav.Location{
				Lat: 1.1,
				Lon: 1.1,
			},
		},
		{
			2,
			chargingstrategy.State{
				Energy: 10.0,
			},
			nav.Location{
				Lat: 2.2,
				Lon: 2.2,
			},
		},
		{
			2,
			chargingstrategy.State{
				Energy: 20.0,
			},
			nav.Location{
				Lat: 2.2,
				Lon: 2.2,
			},
		},
		{
			3,
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
		n       *node
		placeID entity.PlaceID
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
			1,
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
			1,
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
			1,
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
			2,
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
			2,
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
			3,
		},
	}

	if len(input) != len(expect) {
		t.Errorf("Incorrect test case, array of input and array of expect should 1-to-1 match.\n")
	}

	nc := newNodeContainer()

	for i := 0; i < len(input); i++ {
		tmpNode := nc.addNode(input[i].placeID, input[i].chargeState /*, input[i].location*/)

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

		tmpStationID := nc.nodeID2PlaceID((nodeID)(i))
		if !reflect.DeepEqual(tmpStationID, expect[i].placeID) {
			t.Errorf("Calling nodeContainer's placeID() generate incorrect result, expect %#v but got %#v.\n", expect[i].placeID, tmpStationID)
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
		placeID     entity.PlaceID
		chargeState chargingstrategy.State
		location    nav.Location
	}{
		{
			1,
			chargingstrategy.State{
				Energy: 10.0,
			},
			nav.Location{
				Lat: 1.1,
				Lon: 1.1,
			},
		},
		{
			1,
			chargingstrategy.State{
				Energy: 10.0,
			},
			nav.Location{
				Lat: 1.1,
				Lon: 1.1,
			},
		},
		{
			1,
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
		n       *node
		placeID entity.PlaceID
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
		1,
	}

	nc := newNodeContainer()

	for i := 0; i < len(input); i++ {
		tmpNode := nc.addNode(input[i].placeID, input[i].chargeState /*, input[i].location*/)

		if !reflect.DeepEqual(tmpNode, expect.n) {
			t.Errorf("Calling nodeContainer's addNode() generate incorrect result, expect %#v but got %#v.\n", expect.n, tmpNode)
		}
	}

}
