package stationgraph

import (
	"reflect"
	"testing"

	"github.com/Telenav/osrm-backend/integration/oasis/chargingstrategy"
)

func TestGraphGeneral(t *testing.T) {
	cases := []struct {
		g              *graph
		chargeStations []nodeID
	}{
		{
			// case 1: 0 -> 1 -> 3, avoid information of charging
			// node_0 -> node_1, duration = 30, distance = 30
			// node_0 -> node_2, duration = 20, distance = 20
			// node_1 -> node_3, duration = 10, distance = 10
			// node_2 -> node_4, duration = 50, distance = 50
			// node_2 -> node_3, duration = 50, distance = 50
			// node_3 -> node_4, duration = 10, distance = 10
			&graph{
				strategy: chargingstrategy.NewNullChargeStrategy(),
				nodes: []*node{
					&node{
						id: 0,
						neighbors: []*neighbor{
							&neighbor{
								targetNodeID: 1,
								distance:     30,
								duration:     30,
							},
							&neighbor{
								targetNodeID: 2,
								distance:     20,
								duration:     20,
							},
						},
						chargeInfo: chargeInfo{
							arrivalEnergy: 999,
							chargeEnergy:  999,
						},
					},
					&node{
						id: 1,
						neighbors: []*neighbor{
							&neighbor{
								targetNodeID: 3,
								distance:     10,
								duration:     10,
							},
						},
						chargeInfo: chargeInfo{
							arrivalEnergy: 999,
							chargeEnergy:  999,
						},
					},
					&node{
						id: 2,
						neighbors: []*neighbor{
							&neighbor{
								targetNodeID: 4,
								distance:     50,
								duration:     50,
							},
							&neighbor{
								targetNodeID: 3,
								distance:     50,
								duration:     50,
							},
						},
						chargeInfo: chargeInfo{
							arrivalEnergy: 999,
							chargeEnergy:  999,
						},
					},
					&node{
						id: 3,
						neighbors: []*neighbor{
							&neighbor{
								targetNodeID: 4,
								distance:     10,
								duration:     10,
							},
						},
						chargeInfo: chargeInfo{
							arrivalEnergy: 999,
							chargeEnergy:  999,
						},
					},
					&node{
						id:        4,
						neighbors: []*neighbor{},
						chargeInfo: chargeInfo{
							arrivalEnergy: 999,
							chargeEnergy:  999,
						},
					},
				},
				startNodeID: 0,
				endNodeID:   3,
			},
			[]nodeID{1},
		},
		{
			// case 2: 0 -> 1 -> 3 -> 4, avoid information of charging
			// node_0 -> node_1, duration = 30, distance = 30
			// node_0 -> node_2, duration = 20, distance = 20
			// node_1 -> node_3, duration = 10, distance = 10
			// node_2 -> node_4, duration = 50, distance = 50
			// node_2 -> node_3, duration = 50, distance = 50
			// node_3 -> node_4, duration = 10, distance = 10
			&graph{
				strategy: chargingstrategy.NewNullChargeStrategy(),
				nodes: []*node{
					&node{
						id: 0,
						neighbors: []*neighbor{
							&neighbor{
								targetNodeID: 1,
								distance:     30,
								duration:     30,
							},
							&neighbor{
								targetNodeID: 2,
								distance:     20,
								duration:     20,
							},
						},
						chargeInfo: chargeInfo{
							arrivalEnergy: 999,
							chargeEnergy:  999,
						},
					},
					&node{
						id: 1,
						neighbors: []*neighbor{
							&neighbor{
								targetNodeID: 3,
								distance:     10,
								duration:     10,
							},
						},
						chargeInfo: chargeInfo{
							arrivalEnergy: 999,
							chargeEnergy:  999,
						},
					},
					&node{
						id: 2,
						neighbors: []*neighbor{
							&neighbor{
								targetNodeID: 4,
								distance:     50,
								duration:     50,
							},
							&neighbor{
								targetNodeID: 3,
								distance:     50,
								duration:     50,
							},
						},
						chargeInfo: chargeInfo{
							arrivalEnergy: 999,
							chargeEnergy:  999,
						},
					},
					&node{
						id: 3,
						neighbors: []*neighbor{
							&neighbor{
								targetNodeID: 4,
								distance:     10,
								duration:     10,
							},
						},
						chargeInfo: chargeInfo{
							arrivalEnergy: 999,
							chargeEnergy:  999,
						},
					},
					&node{
						id:        4,
						neighbors: []*neighbor{},
						chargeInfo: chargeInfo{
							arrivalEnergy: 999,
							chargeEnergy:  999,
						},
					},
				},
				startNodeID: 0,
				endNodeID:   4,
			},
			[]nodeID{1, 3},
		},
		{
			// case 2: 0 -> 2 -> 4, avoid information of charging
			// node_0 -> node_1, duration = 30, distance = 30
			// node_0 -> node_2, duration = 20, distance = 20
			// node_1 -> node_3, duration = 10, distance = 10
			// node_2 -> node_4, duration = 20, distance = 20
			// node_2 -> node_3, duration = 50, distance = 50
			// node_3 -> node_4, duration = 10, distance = 10
			&graph{
				strategy: chargingstrategy.NewNullChargeStrategy(),
				nodes: []*node{
					&node{
						id: 0,
						neighbors: []*neighbor{
							&neighbor{
								targetNodeID: 1,
								distance:     30,
								duration:     30,
							},
							&neighbor{
								targetNodeID: 2,
								distance:     20,
								duration:     20,
							},
						},
						chargeInfo: chargeInfo{
							arrivalEnergy: 999,
							chargeEnergy:  999,
						},
					},
					&node{
						id: 1,
						neighbors: []*neighbor{
							&neighbor{
								targetNodeID: 3,
								distance:     10,
								duration:     10,
							},
						},
						chargeInfo: chargeInfo{
							arrivalEnergy: 999,
							chargeEnergy:  999,
						},
					},
					&node{
						id: 2,
						neighbors: []*neighbor{
							&neighbor{
								targetNodeID: 4,
								distance:     20,
								duration:     20,
							},
							&neighbor{
								targetNodeID: 3,
								distance:     50,
								duration:     50,
							},
						},
						chargeInfo: chargeInfo{
							arrivalEnergy: 999,
							chargeEnergy:  999,
						},
					},
					&node{
						id: 3,
						neighbors: []*neighbor{
							&neighbor{
								targetNodeID: 4,
								distance:     10,
								duration:     10,
							},
						},
						chargeInfo: chargeInfo{
							arrivalEnergy: 999,
							chargeEnergy:  999,
						},
					},
					&node{
						id:        4,
						neighbors: []*neighbor{},
						chargeInfo: chargeInfo{
							arrivalEnergy: 999,
							chargeEnergy:  999,
						},
					},
				},
				startNodeID: 0,
				endNodeID:   4,
			},
			[]nodeID{2},
		},
		{
			// case 3:  0 -> 2 -> 4 -> 5 -> 8, avoid information of charging
			// node_0 -> node_1, duration = 30, distance = 30
			// node_0 -> node_2, duration = 20, distance = 20
			// node_1 -> node_3, duration = 20, distance = 20
			// node_1 -> node_4, duration = 15, distance = 15
			// node_2 -> node_3, duration = 30, distance = 30
			// node_2 -> node_4, duration = 20, distance = 20
			// node_3 -> node_5, duration = 10, distance = 10
			// node_3 -> node_6, duration = 10, distance = 10
			// node_3 -> node_7, duration = 10, distance = 10
			// node_4 -> node_5, duration = 15, distance = 15
			// node_4 -> node_6, duration = 15, distance = 15
			// node_4 -> node_7, duration = 15, distance = 15
			// node_5 -> node_8, duration = 10, distance = 10
			// node_6 -> node_8, duration = 20, distance = 20
			// node_7 -> node_8, duration = 30, distance = 30
			&graph{
				strategy: chargingstrategy.NewNullChargeStrategy(),
				nodes: []*node{
					&node{
						id: 0,
						neighbors: []*neighbor{
							&neighbor{
								targetNodeID: 1,
								distance:     30,
								duration:     30,
							},
							&neighbor{
								targetNodeID: 2,
								distance:     20,
								duration:     20,
							},
						},
						chargeInfo: chargeInfo{
							arrivalEnergy: 999,
							chargeEnergy:  999,
						},
					},
					&node{
						id: 1,
						neighbors: []*neighbor{
							&neighbor{
								targetNodeID: 3,
								distance:     20,
								duration:     20,
							},
							&neighbor{
								targetNodeID: 4,
								distance:     15,
								duration:     15,
							},
						},
						chargeInfo: chargeInfo{
							arrivalEnergy: 999,
							chargeEnergy:  999,
						},
					},
					&node{
						id: 2,
						neighbors: []*neighbor{
							&neighbor{
								targetNodeID: 4,
								distance:     20,
								duration:     20,
							},
							&neighbor{
								targetNodeID: 3,
								distance:     30,
								duration:     30,
							},
						},
						chargeInfo: chargeInfo{
							arrivalEnergy: 999,
							chargeEnergy:  999,
						},
					},
					&node{
						id: 3,
						neighbors: []*neighbor{
							&neighbor{
								targetNodeID: 5,
								distance:     10,
								duration:     10,
							},
							&neighbor{
								targetNodeID: 6,
								distance:     10,
								duration:     10,
							},
							&neighbor{
								targetNodeID: 7,
								distance:     10,
								duration:     10,
							},
						},
						chargeInfo: chargeInfo{
							arrivalEnergy: 999,
							chargeEnergy:  999,
						},
					},
					&node{
						id: 4,
						neighbors: []*neighbor{
							&neighbor{
								targetNodeID: 5,
								distance:     15,
								duration:     15,
							},
							&neighbor{
								targetNodeID: 6,
								distance:     15,
								duration:     15,
							},
							&neighbor{
								targetNodeID: 7,
								distance:     15,
								duration:     15,
							},
						},
						chargeInfo: chargeInfo{
							arrivalEnergy: 999,
							chargeEnergy:  999,
						},
					},
					&node{
						id: 5,
						neighbors: []*neighbor{
							&neighbor{
								targetNodeID: 8,
								distance:     10,
								duration:     10,
							},
						},
						chargeInfo: chargeInfo{
							arrivalEnergy: 999,
							chargeEnergy:  999,
						},
					},
					&node{
						id: 6,
						neighbors: []*neighbor{
							&neighbor{
								targetNodeID: 8,
								distance:     20,
								duration:     20,
							},
						},
						chargeInfo: chargeInfo{
							arrivalEnergy: 999,
							chargeEnergy:  999,
						},
					},
					&node{
						id: 7,
						neighbors: []*neighbor{
							&neighbor{
								targetNodeID: 8,
								distance:     30,
								duration:     30,
							},
						},
						chargeInfo: chargeInfo{
							arrivalEnergy: 999,
							chargeEnergy:  999,
						},
					},
					&node{
						id: 8,
						neighbors: []*neighbor{
							&neighbor{},
						},
						chargeInfo: chargeInfo{
							arrivalEnergy: 999,
							chargeEnergy:  999,
						},
					},
				},
				startNodeID: 0,
				endNodeID:   8,
			},
			[]nodeID{2, 4, 5},
		},
		{
			// case 4:  0 -> 1 -> 4 -> 5 -> 8, avoid information of charging
			// node_0 -> node_1, duration = 15, distance = 15
			// node_0 -> node_2, duration = 20, distance = 20
			// node_1 -> node_3, duration = 20, distance = 20
			// node_1 -> node_4, duration = 15, distance = 15
			// node_2 -> node_3, duration = 30, distance = 30
			// node_2 -> node_4, duration = 20, distance = 20
			// node_3 -> node_5, duration = 10, distance = 10
			// node_3 -> node_6, duration = 10, distance = 10
			// node_3 -> node_7, duration = 10, distance = 10
			// node_4 -> node_5, duration = 15, distance = 15
			// node_4 -> node_6, duration = 15, distance = 15
			// node_4 -> node_7, duration = 15, distance = 15
			// node_5 -> node_8, duration = 10, distance = 10
			// node_6 -> node_8, duration = 20, distance = 20
			// node_7 -> node_8, duration = 30, distance = 30
			&graph{
				strategy: chargingstrategy.NewNullChargeStrategy(),
				nodes: []*node{
					&node{
						id: 0,
						neighbors: []*neighbor{
							&neighbor{
								targetNodeID: 1,
								distance:     15,
								duration:     15,
							},
							&neighbor{
								targetNodeID: 2,
								distance:     20,
								duration:     20,
							},
						},
						chargeInfo: chargeInfo{
							arrivalEnergy: 999,
							chargeEnergy:  999,
						},
					},
					&node{
						id: 1,
						neighbors: []*neighbor{
							&neighbor{
								targetNodeID: 3,
								distance:     20,
								duration:     20,
							},
							&neighbor{
								targetNodeID: 4,
								distance:     15,
								duration:     15,
							},
						},
						chargeInfo: chargeInfo{
							arrivalEnergy: 999,
							chargeEnergy:  999,
						},
					},
					&node{
						id: 2,
						neighbors: []*neighbor{
							&neighbor{
								targetNodeID: 4,
								distance:     20,
								duration:     20,
							},
							&neighbor{
								targetNodeID: 3,
								distance:     30,
								duration:     30,
							},
						},
						chargeInfo: chargeInfo{
							arrivalEnergy: 999,
							chargeEnergy:  999,
						},
					},
					&node{
						id: 3,
						neighbors: []*neighbor{
							&neighbor{
								targetNodeID: 5,
								distance:     10,
								duration:     10,
							},
							&neighbor{
								targetNodeID: 6,
								distance:     10,
								duration:     10,
							},
							&neighbor{
								targetNodeID: 7,
								distance:     10,
								duration:     10,
							},
						},
						chargeInfo: chargeInfo{
							arrivalEnergy: 999,
							chargeEnergy:  999,
						},
					},
					&node{
						id: 4,
						neighbors: []*neighbor{
							&neighbor{
								targetNodeID: 5,
								distance:     15,
								duration:     15,
							},
							&neighbor{
								targetNodeID: 6,
								distance:     15,
								duration:     15,
							},
							&neighbor{
								targetNodeID: 7,
								distance:     15,
								duration:     15,
							},
						},
						chargeInfo: chargeInfo{
							arrivalEnergy: 999,
							chargeEnergy:  999,
						},
					},
					&node{
						id: 5,
						neighbors: []*neighbor{
							&neighbor{
								targetNodeID: 8,
								distance:     10,
								duration:     10,
							},
						},
						chargeInfo: chargeInfo{
							arrivalEnergy: 999,
							chargeEnergy:  999,
						},
					},
					&node{
						id: 6,
						neighbors: []*neighbor{
							&neighbor{
								targetNodeID: 8,
								distance:     20,
								duration:     20,
							},
						},
						chargeInfo: chargeInfo{
							arrivalEnergy: 999,
							chargeEnergy:  999,
						},
					},
					&node{
						id: 7,
						neighbors: []*neighbor{
							&neighbor{
								targetNodeID: 8,
								distance:     30,
								duration:     30,
							},
						},
						chargeInfo: chargeInfo{
							arrivalEnergy: 999,
							chargeEnergy:  999,
						},
					},
					&node{
						id: 8,
						neighbors: []*neighbor{
							&neighbor{},
						},
						chargeInfo: chargeInfo{
							arrivalEnergy: 999,
							chargeEnergy:  999,
						},
					},
				},
				startNodeID: 0,
				endNodeID:   8,
			},
			[]nodeID{1, 4, 5},
		},
		{
			// case 5: each station only charges 16
			// without considering charing, shortest path is 0 -> 2 -> 4 -> 5 -> 8
			// considering charing, shortest path is 0 -> 1 -> 4 -> 5 -> 8
			// 0 -> 1 -> 4 -> 5 -> 8, avoid information of charging
			// node_0 -> node_1, duration = 15, distance = 15
			// node_0 -> node_2, duration = 20, distance = 20
			// node_1 -> node_3, duration = 20, distance = 20
			// node_1 -> node_4, duration = 15, distance = 15
			// node_2 -> node_3, duration = 30, distance = 30
			// node_2 -> node_4, duration = 5, distance = 5
			// node_3 -> node_5, duration = 10, distance = 10
			// node_3 -> node_6, duration = 10, distance = 10
			// node_3 -> node_7, duration = 10, distance = 10
			// node_4 -> node_5, duration = 15, distance = 15
			// node_4 -> node_6, duration = 15, distance = 15
			// node_4 -> node_7, duration = 15, distance = 15
			// node_5 -> node_8, duration = 10, distance = 10
			// node_6 -> node_8, duration = 20, distance = 20
			// node_7 -> node_8, duration = 30, distance = 30
			&graph{
				strategy: chargingstrategy.NewNullChargeStrategy(),
				nodes: []*node{
					&node{
						id: 0,
						neighbors: []*neighbor{
							&neighbor{
								targetNodeID: 1,
								distance:     15,
								duration:     15,
							},
							&neighbor{
								targetNodeID: 2,
								distance:     20,
								duration:     20,
							},
						},
						chargeInfo: chargeInfo{
							arrivalEnergy: 16,
							chargeEnergy:  0,
						},
					},
					&node{
						id: 1,
						neighbors: []*neighbor{
							&neighbor{
								targetNodeID: 3,
								distance:     20,
								duration:     20,
							},
							&neighbor{
								targetNodeID: 4,
								distance:     15,
								duration:     15,
							},
						},
						chargeInfo: chargeInfo{
							arrivalEnergy: 0,
							chargeEnergy:  16,
						},
					},
					&node{
						id: 2,
						neighbors: []*neighbor{
							&neighbor{
								targetNodeID: 4,
								distance:     5,
								duration:     5,
							},
							&neighbor{
								targetNodeID: 3,
								distance:     30,
								duration:     30,
							},
						},
						chargeInfo: chargeInfo{
							arrivalEnergy: 0,
							chargeEnergy:  16,
						},
					},
					&node{
						id: 3,
						neighbors: []*neighbor{
							&neighbor{
								targetNodeID: 5,
								distance:     10,
								duration:     10,
							},
							&neighbor{
								targetNodeID: 6,
								distance:     10,
								duration:     10,
							},
							&neighbor{
								targetNodeID: 7,
								distance:     10,
								duration:     10,
							},
						},
						chargeInfo: chargeInfo{
							arrivalEnergy: 0,
							chargeEnergy:  16,
						},
					},
					&node{
						id: 4,
						neighbors: []*neighbor{
							&neighbor{
								targetNodeID: 5,
								distance:     15,
								duration:     15,
							},
							&neighbor{
								targetNodeID: 6,
								distance:     15,
								duration:     15,
							},
							&neighbor{
								targetNodeID: 7,
								distance:     15,
								duration:     15,
							},
						},
						chargeInfo: chargeInfo{
							arrivalEnergy: 0,
							chargeEnergy:  16,
						},
					},
					&node{
						id: 5,
						neighbors: []*neighbor{
							&neighbor{
								targetNodeID: 8,
								distance:     10,
								duration:     10,
							},
						},
						chargeInfo: chargeInfo{
							arrivalEnergy: 999,
							chargeEnergy:  999,
						},
					},
					&node{
						id: 6,
						neighbors: []*neighbor{
							&neighbor{
								targetNodeID: 8,
								distance:     20,
								duration:     20,
							},
						},
						chargeInfo: chargeInfo{
							arrivalEnergy: 0,
							chargeEnergy:  16,
						},
					},
					&node{
						id: 7,
						neighbors: []*neighbor{
							&neighbor{
								targetNodeID: 8,
								distance:     30,
								duration:     30,
							},
						},
						chargeInfo: chargeInfo{
							arrivalEnergy: 0,
							chargeEnergy:  16,
						},
					},
					&node{
						id: 8,
						neighbors: []*neighbor{
							&neighbor{},
						},
						chargeInfo: chargeInfo{
							arrivalEnergy: 0,
							chargeEnergy:  0,
						},
					},
				},
				startNodeID: 0,
				endNodeID:   8,
			},
			[]nodeID{1, 4, 5},
		},
	}

	for i, c := range cases {
		s := c.g.dijkstra()
		if !reflect.DeepEqual(c.chargeStations, s) {
			t.Errorf("for test case %d expect %v but got %v", i, c.chargeStations, s)
		}
	}
}
