package connectivitymap

import (
	"reflect"
	"runtime"
	"testing"

	"github.com/Telenav/osrm-backend/integration/service/oasis/spatialindexer"
	"github.com/Telenav/osrm-backend/integration/service/oasis/spatialindexer/ranker"
)

// Mock env:
// Iterator returns 100 of fixed location point(IDs are 1000, 1001, 1002, ...)
// Finder returns fixed array result(10 points, IDs are 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 )
// Ranker use great circle distance to calculate
// Expect result:
// build() generate same result as pre-calculated map
// map[1000] = {results of 10 points}
// map[1001] = {results of 10 points}
// ...
// map[1099] = {results of 10 points}
func TestBuilderWithMockIteratorAndFinder(t *testing.T) {
	builder := newConnectivityMapBuilder(&spatialindexer.MockOneHundredPointsIterator{},
		&spatialindexer.MockFinder{},
		ranker.CreateRanker(ranker.SimpleRanker, nil),
		100,
		runtime.NumCPU())

	actual := builder.build()

	// construct expect map
	expect := make(ID2NearByIDsMap)

	var idAndDistanceArray = []IDAndDistance{
		IDAndDistance{
			ID:       3,
			Distance: 345.220003472554,
		},
		IDAndDistance{
			ID:       2,
			Distance: 402.8536530341791,
		},
		IDAndDistance{
			ID:       4,
			Distance: 1627.1858848458571,
		},
		IDAndDistance{
			ID:       5,
			Distance: 4615.586636153461,
		},
		IDAndDistance{
			ID:       1,
			Distance: 5257.70008125706,
		},
		IDAndDistance{
			ID:       6,
			Distance: 6888.7486674247,
		},
		IDAndDistance{
			ID:       7,
			Distance: 7041.893747628621,
		},
		IDAndDistance{
			ID:       10,
			Distance: 8622.213424347745,
		},
		IDAndDistance{
			ID:       9,
			Distance: 9438.804320070916,
		},
		IDAndDistance{
			ID:       8,
			Distance: 9897.44482638937,
		},
	}

	for i := 0; i < 100; i++ {
		index := i + 1000
		expect[(spatialindexer.PointID(index))] = idAndDistanceArray
	}

	if !reflect.DeepEqual(actual, expect) {
		t.Errorf("Failed to pass TestBuilder with mock data, \nexpect \n%+v\n but got \n%v\n", expect, actual)
	}
}

func TestBuilderWithSingleWorker(t *testing.T) {
	builder := newConnectivityMapBuilder(&spatialindexer.MockOneHundredPointsIterator{},
		&spatialindexer.MockFinder{},
		ranker.CreateRanker(ranker.SimpleRanker, nil),
		100,
		runtime.NumCPU())

	actual := builder.build()

	expect := builder.buildInSerial()

	if !reflect.DeepEqual(actual, expect) {
		t.Errorf("Failed to pass TestBuilder with mock data, \nexpect \n%+v\n but got \n%v\n", expect, actual)
	}
}
