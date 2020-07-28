package topograph

import (
	"reflect"
	"runtime"
	"testing"

	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/entity"
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/mock"
	"github.com/Telenav/osrm-backend/integration/service/oasis/place/spatialindexer/ranker"
)

// Mock env:
// IteratorGenerator returns 100 of fixed location places(IDs are 1000, 1001, 1002, ...)
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

	var idAndWeightArray = []*entity.TransferInfo{
		{
			PlaceWithLocation: entity.PlaceWithLocation{
				ID: 3,
				Location: &nav.Location{
					Lat: 37.401948,
					Lon: -121.977384,
				},
			},
			Weight: &entity.Weight{
				Distance: 345.220003472554,
				Duration: 15.550450606871804,
			},
		},
		{
			PlaceWithLocation: entity.PlaceWithLocation{
				ID: 2,
				Location: &nav.Location{
					Lat: 37.399331,
					Lon: -121.981193,
				},
			},
			Weight: &entity.Weight{
				Distance: 402.8536530341791,
				Duration: 18.146560947485547,
			},
		},
		{
			PlaceWithLocation: entity.PlaceWithLocation{
				ID: 4,
				Location: &nav.Location{
					Lat: 37.407082,
					Lon: -121.991937,
				},
			},
			Weight: &entity.Weight{
				Distance: 1627.1858848458571,
				Duration: 73.29666147954312,
			},
		},
		{
			PlaceWithLocation: entity.PlaceWithLocation{
				ID: 5,
				Location: &nav.Location{
					Lat: 37.407277,
					Lon: -121.925482,
				},
			},
			Weight: &entity.Weight{
				Distance: 4615.586636153461,
				Duration: 207.9093079348406,
			},
		},
		{
			PlaceWithLocation: entity.PlaceWithLocation{
				ID: 1,
				Location: &nav.Location{
					Lat: 37.355204,
					Lon: -121.953901,
				},
			},
			Weight: &entity.Weight{
				Distance: 5257.70008125706,
				Duration: 236.8333369935613,
			},
		},
		{
			PlaceWithLocation: entity.PlaceWithLocation{
				ID: 6,
				Location: &nav.Location{
					Lat: 37.375024,
					Lon: -121.904706,
				},
			},
			Weight: &entity.Weight{
				Distance: 6888.7486674247,
				Duration: 310.30399402813964,
			},
		},
		{
			PlaceWithLocation: entity.PlaceWithLocation{
				ID: 7,
				Location: &nav.Location{
					Lat: 37.359592,
					Lon: -121.914164,
				},
			},
			Weight: &entity.Weight{
				Distance: 7041.893747628621,
				Duration: 317.2024210643523,
			},
		},
		{
			PlaceWithLocation: entity.PlaceWithLocation{
				ID: 10,
				Location: &nav.Location{
					Lat: 37.373546,
					Lon: -122.068904,
				},
			},
			Weight: &entity.Weight{
				Distance: 8622.213424347745,
				Duration: 388.3879920877363,
			},
		},
		{
			PlaceWithLocation: entity.PlaceWithLocation{
				ID: 9,
				Location: &nav.Location{
					Lat: 37.368453,
					Lon: -122.0764,
				},
			},
			Weight: &entity.Weight{
				Distance: 9438.804320070916,
				Duration: 425.1713657689602,
			},
		},
		{
			PlaceWithLocation: entity.PlaceWithLocation{
				ID: 8,
				Location: &nav.Location{
					Lat: 37.366023,
					Lon: -122.080777,
				},
			},
			Weight: &entity.Weight{
				Distance: 9897.44482638937,
				Duration: 445.8308480355572,
			},
		},
	}

	for i := 0; i < 100; i++ {
		index := i + 1000
		expect[(entity.PlaceID(index))] = idAndWeightArray
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
