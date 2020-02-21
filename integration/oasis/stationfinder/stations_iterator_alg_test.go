package stationfinder

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/Telenav/osrm-backend/integration/oasis/osrmconnector"
	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/table"
)

var mockDict1 map[string]bool = map[string]bool{
	"station1": true,
	"station2": true,
	"station3": true,
	"station4": true,
}

func TestBuildChargeStationInfoDict1(t *testing.T) {
	sf := createMockOrigStationFinder1()
	m := buildChargeStationInfoDict(sf)
	if !reflect.DeepEqual(m, mockDict1) {
		t.Errorf("expect %v but got %v", mockDict1, m)
	}
}

var overlapChargeStationInfo1 []ChargeStationInfo = []ChargeStationInfo{
	ChargeStationInfo{
		ID: "station1",
		Location: StationCoordinate{
			Lat: 32.333,
			Lon: 122.333,
		},
	},
	ChargeStationInfo{
		ID: "station2",
		Location: StationCoordinate{
			Lat: -32.333,
			Lon: -122.333,
		},
	},
}

func TestFindOverlapBetweenStations1(t *testing.T) {
	sf1 := createMockOrigStationFinder2()
	sf2 := createMockDestStationFinder1()
	r := FindOverlapBetweenStations(sf1, sf2)

	if !reflect.DeepEqual(r, overlapChargeStationInfo1) {
		t.Errorf("expect %v but got %v", overlapChargeStationInfo1, r)
	}
}

type fakeTableResponse struct {
}

func (ft *fakeTableResponse) Request4Table(r *table.Request) <-chan osrmconnector.TableResponse {
	fmt.Printf("fuck")
	c := make(chan osrmconnector.TableResponse)
	go func() {
		defer close(c)
		c <- osrmconnector.TableResponse{
			Resp: &table.Mock4To2TableResponse1,
			Err:  nil,
		}
	}()
	return c
}

func TestCalcCostBetweenChargeStationsPair(t *testing.T) {
	// from: station1, station2, station3, station4
	sf1 := createMockOrigStationFinder1()
	// to: station6, station7
	sf2 := createMockOrigStationFinder3()

	table := &fakeTableResponse{}
	r, err := CalcCostBetweenChargeStationsPair(sf1, sf2, table)

	if err != nil {
		t.Errorf("expect no error but generate error of %v", err)
	}
	expect := []CostBetweenChargeStations{
		CostBetweenChargeStations{
			FromID: "station1",
			ToID:   "station6",
			Cost: Cost{
				Duration: 2,
				Distance: 2,
			},
		},
		CostBetweenChargeStations{
			FromID: "station1",
			ToID:   "station7",
			Cost: Cost{
				Duration: 3,
				Distance: 3,
			},
		},
		CostBetweenChargeStations{
			FromID: "station2",
			ToID:   "station6",
			Cost: Cost{
				Duration: 4,
				Distance: 4,
			},
		},
		CostBetweenChargeStations{
			FromID: "station2",
			ToID:   "station7",
			Cost: Cost{
				Duration: 5,
				Distance: 5,
			},
		},
		CostBetweenChargeStations{
			FromID: "station3",
			ToID:   "station6",
			Cost: Cost{
				Duration: 6,
				Distance: 6,
			},
		},
		CostBetweenChargeStations{
			FromID: "station3",
			ToID:   "station7",
			Cost: Cost{
				Duration: 7,
				Distance: 7,
			},
		},
		CostBetweenChargeStations{
			FromID: "station4",
			ToID:   "station6",
			Cost: Cost{
				Duration: 8,
				Distance: 8,
			},
		},
		CostBetweenChargeStations{
			FromID: "station4",
			ToID:   "station7",
			Cost: Cost{
				Duration: 9,
				Distance: 9,
			},
		},
	}
	if !reflect.DeepEqual(r, expect) {
		t.Errorf("expect %v but got %v", expect, r)
	}
}
