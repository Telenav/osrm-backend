package stationfinder

import (
	"reflect"
	"testing"
)

var mockDict1 map[string]bool = map[string]bool{
	"station1": true,
	"station2": true,
	"station3": true,
	"station4": true,
}

func TestbuildChargeStationInfoDict1(t *testing.T) {
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
