package connectivitymap

import (
	"reflect"
	"runtime"
	"testing"

	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/mock"
	"github.com/Telenav/osrm-backend/integration/service/oasis/spatialindexer"
	"github.com/Telenav/osrm-backend/integration/service/oasis/spatialindexer/ranker"
)

// Mock env:
// Iterator returns 100 of fixed location places(IDs are 1000, 1001, 1002, ...)
// Finder returns fixed array result(10 places, IDs are 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 )
// Ranker use great circle distance to calculate
// Expect result:
// build() generate same result as pre-calculated map
// map[1000] = {results of 10 places}
// map[1001] = {results of 10 places}
// ...
// map[1099] = {results of 10 places}
func TestBuilderWithMockIteratorAndFinder(t *testing.T) {
	builder := newConnectivityMapBuilder(&mock.MockOneHundredPlacesIterator{},
		&mock.MockFinder{},
		ranker.CreateRanker(ranker.SimpleRanker, nil),
		100,
		runtime.NumCPU())

	actual := builder.build()

	// construct expect map
	expect := make(ID2NearByIDsMap)

	var idAndWeightArray = []IDAndWeight{
		{
			3,
			Weight{
				345.220003472554,
				15.550450606871804,
			},
		},
		{
			2,
			Weight{
				402.8536530341791,
				18.146560947485547,
			},
		},
		{
			4,
			Weight{
				1627.1858848458571,
				73.29666147954312,
			},
		},
		{
			5,
			Weight{
				4615.586636153461,
				207.9093079348406,
			},
		},
		{
			1,
			Weight{
				5257.70008125706,
				236.8333369935613,
			},
		},
		{
			6,
			Weight{
				6888.7486674247,
				310.30399402813964,
			},
		},
		{
			7,
			Weight{
				7041.893747628621,
				317.2024210643523,
			},
		},
		{
			10,
			Weight{
				8622.213424347745,
				388.3879920877363,
			},
		},
		{
			9,
			Weight{
				9438.804320070916,
				425.1713657689602,
			},
		},
		{
			8,
			Weight{
				9897.44482638937,
				445.8308480355572,
			},
		},
	}

	for i := 0; i < 100; i++ {
		index := i + 1000
		expect[(spatialindexer.PlaceID(index))] = idAndWeightArray
	}

	if !reflect.DeepEqual(actual, expect) {
		t.Errorf("Failed to pass TestBuilder with mock data, \nexpect \n%+v\n but got \n%v\n", expect, actual)
	}
}

func TestBuilderWithSingleWorker(t *testing.T) {
	builder := newConnectivityMapBuilder(&mock.MockOneHundredPlacesIterator{},
		&mock.MockFinder{},
		ranker.CreateRanker(ranker.SimpleRanker, nil),
		100,
		runtime.NumCPU())

	actual := builder.build()

	expect := builder.buildInSerial()

	if !reflect.DeepEqual(actual, expect) {
		t.Errorf("Failed to pass TestBuilder with mock data, \nexpect \n%+v\n but got \n%v\n", expect, actual)
	}
}
