package stationgraph

import (
	"reflect"
	"testing"
)

func TestDijkstraAlgorithm(t *testing.T) {
	cases := []struct {
		graph          Graph
		start          nodeID
		end            nodeID
		chargeStations []nodeID
	}{
		// case 1
		// - input graph:
		//    node_0 -> node_1, duration = 30, distance = 30
		//    node_0 -> node_2, duration = 20, distance = 20
		//    node_1 -> node_3, duration = 10, distance = 10
		//    node_2 -> node_4, duration = 50, distance = 50
		//    node_2 -> node_3, duration = 50, distance = 50
		//    node_3 -> node_4, duration = 10, distance = 10
		// - start = 0, end = 3
		// - set charge information to fixed status to ignore situation of lack of energy
		//
		// expect:  0 -> 1 -> 3
		{
			newMockGraph1(),
			0,
			3,
			[]nodeID{1},
		},
		// case 2
		// - input graph:
		//    same graph as case 1
		// - start = 0, end = 4
		// - set charge information to fixed status to ignore situation of lack of energy
		//
		// expect:  0 -> 1 -> 3 -> 4
		{
			newMockGraph1(),
			0,
			4,
			[]nodeID{1, 3},
		},
		// case 3
		// - input graph:
		//    node_0 -> node_1, duration = 30, distance = 30
		//    node_0 -> node_2, duration = 20, distance = 20
		//    node_1 -> node_3, duration = 20, distance = 20
		//    node_1 -> node_4, duration = 15, distance = 15
		//    node_2 -> node_3, duration = 30, distance = 30
		//    node_2 -> node_4, duration = 20, distance = 20
		//    node_3 -> node_5, duration = 10, distance = 10
		//    node_3 -> node_6, duration = 10, distance = 10
		//    node_3 -> node_7, duration = 10, distance = 10
		//    node_4 -> node_5, duration = 15, distance = 15
		//    node_4 -> node_6, duration = 15, distance = 15
		//    node_4 -> node_7, duration = 15, distance = 15
		//    node_5 -> node_8, duration = 10, distance = 10
		//    node_6 -> node_8, duration = 20, distance = 20
		//    node_7 -> node_8, duration = 30, distance = 30
		// - start = 0, end = 8
		// - set charge information to fixed status to ignore situation of lack of energy
		//
		// expect:  0 -> 2 -> 4 -> 5 -> 8
		{
			newMockGraph2(),
			0,
			8,
			[]nodeID{2, 4, 5},
		},
		// case 4
		// - input graph:
		//    node_0 -> node_1, duration = 15, distance = 15
		//    node_0 -> node_2, duration = 20, distance = 20
		//    node_1 -> node_3, duration = 20, distance = 20
		//    node_1 -> node_4, duration = 15, distance = 15
		//    node_2 -> node_3, duration = 30, distance = 30
		//    node_2 -> node_4, duration = 20, distance = 20
		//    node_3 -> node_5, duration = 10, distance = 10
		//    node_3 -> node_6, duration = 10, distance = 10
		//    node_3 -> node_7, duration = 10, distance = 10
		//    node_4 -> node_5, duration = 15, distance = 15
		//    node_4 -> node_6, duration = 15, distance = 15
		//    node_4 -> node_7, duration = 15, distance = 15
		//    node_5 -> node_8, duration = 10, distance = 10
		//    node_6 -> node_8, duration = 20, distance = 20
		//    node_7 -> node_8, duration = 30, distance = 30
		// - start = 0, end = 8
		// - set charge information to fixed status to ignore situation of lack of energy
		//
		// expect:  0 -> 1 -> 4 -> 5 -> 8
		{
			newMockGraph3(),
			0,
			8,
			[]nodeID{1, 4, 5},
		},
		// case 5
		// - input graph:
		//    node_0 -> node_1, duration = 15, distance = 15
		//    node_0 -> node_2, duration = 20, distance = 20
		//    node_1 -> node_3, duration = 20, distance = 20
		//    node_1 -> node_4, duration = 15, distance = 15
		//    node_2 -> node_3, duration = 30, distance = 30
		//    node_2 -> node_4, duration = 5, distance = 5
		//    node_3 -> node_5, duration = 10, distance = 10
		//    node_3 -> node_6, duration = 10, distance = 10
		//    node_3 -> node_7, duration = 10, distance = 10
		//    node_4 -> node_5, duration = 15, distance = 15
		//    node_4 -> node_6, duration = 15, distance = 15
		//    node_4 -> node_7, duration = 15, distance = 15
		//    node_5 -> node_8, duration = 10, distance = 10
		//    node_6 -> node_8, duration = 20, distance = 20
		//    node_7 -> node_8, duration = 30, distance = 30
		// - start = 0, end = 8
		// - Each station only charges 16 unit of energy
		//
		// expect: 0 -> 1 -> 4 -> 5 -> 8, without consider charging shortest path is 0 -> 2 -> 4 -> 5 -> 8
		// But 16 could not make from node_0 to node_2
		{
			newMockGraph4(),
			0,
			8,
			[]nodeID{1, 4, 5},
		},
	}

	for i, c := range cases {
		s := dijkstra(c.graph, c.start, c.end)
		if !reflect.DeepEqual(c.chargeStations, s) {
			t.Errorf("for test case %d expect %v but got %v", i, c.chargeStations, s)
		}
	}
}
