package stationfinder

import (
	"reflect"
	"sync"
	"testing"

	"github.com/Telenav/osrm-backend/integration/pkg/api/search/nearbychargestation"
)

var mockSearchResponse1 *nearbychargestation.Response = &nearbychargestation.Response{
	Results: []*nearbychargestation.Result{
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
}

var mockSearchResponse2 *nearbychargestation.Response = &nearbychargestation.Response{
	Results: []*nearbychargestation.Result{
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
			ID: "station5",
			Place: nearbychargestation.Place{
				Address: []*nearbychargestation.Address{
					&nearbychargestation.Address{
						GeoCoordinate: nearbychargestation.Coordinate{
							Latitude:  -12.333,
							Longitude: 122.333,
						},
						NavCoordinates: []*nearbychargestation.Coordinate{
							&nearbychargestation.Coordinate{
								Latitude:  -12.333,
								Longitude: 122.333,
							},
						},
					},
				},
			},
		},
	},
}

var mockSearchResponse3 *nearbychargestation.Response = &nearbychargestation.Response{
	Results: []*nearbychargestation.Result{
		&nearbychargestation.Result{
			ID: "station6",
			Place: nearbychargestation.Place{
				Address: []*nearbychargestation.Address{
					&nearbychargestation.Address{
						GeoCoordinate: nearbychargestation.Coordinate{
							Latitude:  30.333,
							Longitude: 122.333,
						},
						NavCoordinates: []*nearbychargestation.Coordinate{
							&nearbychargestation.Coordinate{
								Latitude:  30.333,
								Longitude: 122.333,
							},
						},
					},
				},
			},
		},
		&nearbychargestation.Result{
			ID: "station7",
			Place: nearbychargestation.Place{
				Address: []*nearbychargestation.Address{
					&nearbychargestation.Address{
						GeoCoordinate: nearbychargestation.Coordinate{
							Latitude:  -10.333,
							Longitude: 122.333,
						},
						NavCoordinates: []*nearbychargestation.Coordinate{
							&nearbychargestation.Coordinate{
								Latitude:  -10.333,
								Longitude: 122.333,
							},
						},
					},
				},
			},
		},
	},
}

var mockChargeStationInfo1 []ChargeStationInfo = []ChargeStationInfo{
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
		ID: "station4",
		Location: StationCoordinate{
			Lat: -32.333,
			Lon: 122.333,
		},
	},
}

var mockChargeStationInfo2 []ChargeStationInfo = []ChargeStationInfo{
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
		ID: "station5",
		Location: StationCoordinate{
			Lat: -12.333,
			Lon: 122.333,
		},
	},
}

var mockChargeStationInfo3 []ChargeStationInfo = []ChargeStationInfo{
	ChargeStationInfo{
		ID: "station6",
		Location: StationCoordinate{
			Lat: 30.333,
			Lon: 122.333,
		},
	},
	ChargeStationInfo{
		ID: "station7",
		Location: StationCoordinate{
			Lat: -10.333,
			Lon: 122.333,
		},
	},
}

func TestBasicFinderCorrectness(t *testing.T) {
	cases := []struct {
		input  []*nearbychargestation.Result
		expect []ChargeStationInfo
	}{
		{
			mockSearchResponse1.Results,
			mockChargeStationInfo1,
		},
	}

	for _, b := range cases {
		input := b.input
		expect := b.expect
		var bf basicFinder
		c := bf.iterateNearbyStations(input, nil)

		var wg sync.WaitGroup
		go func(wg *sync.WaitGroup) {
			wg.Add(1)
			defer wg.Done()

			var r []ChargeStationInfo
			for item := range c {
				r = append(r, item)
			}

			if !reflect.DeepEqual(r, expect) {
				t.Errorf("parse %v expect %v but got %v", b.input, b.expect, r)
			}
		}(&wg)
		wg.Wait()
	}
}

func TestBasicFinderAsync(t *testing.T) {
	cases := []struct {
		input     []*nearbychargestation.Result
		inputLock *sync.RWMutex
		expect    []ChargeStationInfo
	}{
		{
			mockSearchResponse1.Results,
			&sync.RWMutex{},
			mockChargeStationInfo1,
		},
	}

	for _, b := range cases {
		input := b.input
		expect := b.expect
		var bf basicFinder

		num := 20
		var wg sync.WaitGroup
		for i := 0; i < num; i++ {
			go func(wg *sync.WaitGroup) {
				wg.Add(1)

				c := bf.iterateNearbyStations(input, b.inputLock)
				go func(wg *sync.WaitGroup) {
					defer wg.Done()
					var r []ChargeStationInfo
					for item := range c {
						r = append(r, item)
					}
					if !reflect.DeepEqual(r, expect) {
						t.Errorf("parse %v expect %v but got %v", b.input, b.expect, r)
					}
				}(wg)
			}(&wg)
		}
		wg.Wait()

	}
}
