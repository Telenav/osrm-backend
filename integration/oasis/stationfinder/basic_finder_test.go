package stationfinder

import (
	"reflect"
	"sync"
	"testing"

	"github.com/Telenav/osrm-backend/integration/pkg/api/search/nearbychargestation"
)

func TestBasicFinderCorrectness(t *testing.T) {
	cases := []struct {
		input  []*nearbychargestation.Result
		expect []ChargeStationInfo
	}{
		{
			[]*nearbychargestation.Result{
				&nearbychargestation.Result{
					ID: "station1",
					Place: nearbychargestation.Place{
						Address: []*nearbychargestation.Address{
							&nearbychargestation.Address{
								GeoCoordinate: nearbychargestation.Coordinate{
									Latitude:  32.333,
									Longitude: 122.333,
								},
								NavCoordinates: []*nearbychargestation.Coordinate{
									&nearbychargestation.Coordinate{
										Latitude:  32.333,
										Longitude: 122.333,
									},
								},
							},
						},
					},
				},
				&nearbychargestation.Result{
					ID: "station2",
					Place: nearbychargestation.Place{
						Address: []*nearbychargestation.Address{
							&nearbychargestation.Address{
								GeoCoordinate: nearbychargestation.Coordinate{
									Latitude:  -32.333,
									Longitude: -122.333,
								},
								NavCoordinates: []*nearbychargestation.Coordinate{
									&nearbychargestation.Coordinate{
										Latitude:  -32.333,
										Longitude: -122.333,
									},
								},
							},
						},
					},
				},
				&nearbychargestation.Result{
					ID: "station3",
					Place: nearbychargestation.Place{
						Address: []*nearbychargestation.Address{
							&nearbychargestation.Address{
								GeoCoordinate: nearbychargestation.Coordinate{
									Latitude:  32.333,
									Longitude: -122.333,
								},
								NavCoordinates: []*nearbychargestation.Coordinate{
									&nearbychargestation.Coordinate{
										Latitude:  32.333,
										Longitude: -122.333,
									},
								},
							},
						},
					},
				},
				&nearbychargestation.Result{
					ID: "station4",
					Place: nearbychargestation.Place{
						Address: []*nearbychargestation.Address{
							&nearbychargestation.Address{
								GeoCoordinate: nearbychargestation.Coordinate{
									Latitude:  -32.333,
									Longitude: 122.333,
								},
								NavCoordinates: []*nearbychargestation.Coordinate{
									&nearbychargestation.Coordinate{
										Latitude:  -32.333,
										Longitude: 122.333,
									},
								},
							},
						},
					},
				},
			},
			[]ChargeStationInfo{
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
				ChargeStationInfo{
					ID: "station3",
					Location: StationCoordinate{
						Lat: 32.333,
						Lon: -122.333,
					},
				},
				ChargeStationInfo{
					ID: "station2",
					Location: StationCoordinate{
						Lat: -32.333,
						Lon: 122.333,
					},
				},
			},
		},
	}

	for _, b := range cases {
		input := b.input
		expect := b.expect
		var bf basicFinder
		c := bf.iterateNearbyStations(input, nil)
		go func() {
			r := make([]ChargeStationInfo, 10)
			for item := range c {
				r = append(r, item)
			}
			if !reflect.DeepEqual(r, expect) {
				t.Errorf("parse %v expect %v but got %v", b.input, b.expect, r)
			}
		}()
	}
}

func TestBasicFinderAsync(t *testing.T) {
	cases := []struct {
		input     []*nearbychargestation.Result
		inputLock *sync.RWMutex
		expect    []ChargeStationInfo
	}{
		{
			[]*nearbychargestation.Result{
				&nearbychargestation.Result{
					ID: "station1",
					Place: nearbychargestation.Place{
						Address: []*nearbychargestation.Address{
							&nearbychargestation.Address{
								GeoCoordinate: nearbychargestation.Coordinate{
									Latitude:  32.333,
									Longitude: 122.333,
								},
								NavCoordinates: []*nearbychargestation.Coordinate{
									&nearbychargestation.Coordinate{
										Latitude:  32.333,
										Longitude: 122.333,
									},
								},
							},
						},
					},
				},
				&nearbychargestation.Result{
					ID: "station2",
					Place: nearbychargestation.Place{
						Address: []*nearbychargestation.Address{
							&nearbychargestation.Address{
								GeoCoordinate: nearbychargestation.Coordinate{
									Latitude:  -32.333,
									Longitude: -122.333,
								},
								NavCoordinates: []*nearbychargestation.Coordinate{
									&nearbychargestation.Coordinate{
										Latitude:  -32.333,
										Longitude: -122.333,
									},
								},
							},
						},
					},
				},
				&nearbychargestation.Result{
					ID: "station3",
					Place: nearbychargestation.Place{
						Address: []*nearbychargestation.Address{
							&nearbychargestation.Address{
								GeoCoordinate: nearbychargestation.Coordinate{
									Latitude:  32.333,
									Longitude: -122.333,
								},
								NavCoordinates: []*nearbychargestation.Coordinate{
									&nearbychargestation.Coordinate{
										Latitude:  32.333,
										Longitude: -122.333,
									},
								},
							},
						},
					},
				},
				&nearbychargestation.Result{
					ID: "station4",
					Place: nearbychargestation.Place{
						Address: []*nearbychargestation.Address{
							&nearbychargestation.Address{
								GeoCoordinate: nearbychargestation.Coordinate{
									Latitude:  -32.333,
									Longitude: 122.333,
								},
								NavCoordinates: []*nearbychargestation.Coordinate{
									&nearbychargestation.Coordinate{
										Latitude:  -32.333,
										Longitude: 122.333,
									},
								},
							},
						},
					},
				},
			},
			&sync.RWMutex{},
			[]ChargeStationInfo{
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
				ChargeStationInfo{
					ID: "station3",
					Location: StationCoordinate{
						Lat: 32.333,
						Lon: -122.333,
					},
				},
				ChargeStationInfo{
					ID: "station2",
					Location: StationCoordinate{
						Lat: -32.333,
						Lon: 122.333,
					},
				},
			},
		},
	}

	for _, b := range cases {
		input := b.input
		expect := b.expect
		var bf basicFinder

		for i := 0; i < 20; i++ {
			go func() {
				c := bf.iterateNearbyStations(input, b.inputLock)
				go func() {
					r := make([]ChargeStationInfo, 10)
					for item := range c {
						r = append(r, item)
					}
					if !reflect.DeepEqual(r, expect) {
						t.Errorf("parse %v expect %v but got %v", b.input, b.expect, r)
					}
				}()
			}()
		}

	}
}
