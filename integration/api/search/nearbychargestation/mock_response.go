package nearbychargestation

// MockSearchResponse1 provides 4 charge stations in the result: station1,
// station2, station3, station4
var MockSearchResponse1 *Response = &Response{
	Results: []*Result{
		&Result{
			ID: "station1",
			Place: Place{
				Address: []*Address{
					&Address{
						GeoCoordinate: Coordinate{
							Latitude:  32.333,
							Longitude: 122.333,
						},
						NavCoordinates: []*Coordinate{
							&Coordinate{
								Latitude:  32.333,
								Longitude: 122.333,
							},
						},
					},
				},
			},
		},
		&Result{
			ID: "station2",
			Place: Place{
				Address: []*Address{
					&Address{
						GeoCoordinate: Coordinate{
							Latitude:  -32.333,
							Longitude: -122.333,
						},
						NavCoordinates: []*Coordinate{
							&Coordinate{
								Latitude:  -32.333,
								Longitude: -122.333,
							},
						},
					},
				},
			},
		},
		&Result{
			ID: "station3",
			Place: Place{
				Address: []*Address{
					&Address{
						GeoCoordinate: Coordinate{
							Latitude:  32.333,
							Longitude: -122.333,
						},
						NavCoordinates: []*Coordinate{
							&Coordinate{
								Latitude:  32.333,
								Longitude: -122.333,
							},
						},
					},
				},
			},
		},
		&Result{
			ID: "station4",
			Place: Place{
				Address: []*Address{
					&Address{
						GeoCoordinate: Coordinate{
							Latitude:  -32.333,
							Longitude: 122.333,
						},
						NavCoordinates: []*Coordinate{
							&Coordinate{
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

// MockSearchResponse2 provides 3 charge stations in the result: station1,
// station2, station5
var MockSearchResponse2 *Response = &Response{
	Results: []*Result{
		&Result{
			ID: "station1",
			Place: Place{
				Address: []*Address{
					&Address{
						GeoCoordinate: Coordinate{
							Latitude:  32.333,
							Longitude: 122.333,
						},
						NavCoordinates: []*Coordinate{
							&Coordinate{
								Latitude:  32.333,
								Longitude: 122.333,
							},
						},
					},
				},
			},
		},
		&Result{
			ID: "station2",
			Place: Place{
				Address: []*Address{
					&Address{
						GeoCoordinate: Coordinate{
							Latitude:  -32.333,
							Longitude: -122.333,
						},
						NavCoordinates: []*Coordinate{
							&Coordinate{
								Latitude:  -32.333,
								Longitude: -122.333,
							},
						},
					},
				},
			},
		},
		&Result{
			ID: "station5",
			Place: Place{
				Address: []*Address{
					&Address{
						GeoCoordinate: Coordinate{
							Latitude:  -12.333,
							Longitude: 122.333,
						},
						NavCoordinates: []*Coordinate{
							&Coordinate{
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

// MockSearchResponse3 provides 2 charge stations in the result: station6,
// station7
var MockSearchResponse3 *Response = &Response{
	Results: []*Result{
		&Result{
			ID: "station6",
			Place: Place{
				Address: []*Address{
					&Address{
						GeoCoordinate: Coordinate{
							Latitude:  30.333,
							Longitude: 122.333,
						},
						NavCoordinates: []*Coordinate{
							&Coordinate{
								Latitude:  30.333,
								Longitude: 122.333,
							},
						},
					},
				},
			},
		},
		&Result{
			ID: "station7",
			Place: Place{
				Address: []*Address{
					&Address{
						GeoCoordinate: Coordinate{
							Latitude:  -10.333,
							Longitude: 122.333,
						},
						NavCoordinates: []*Coordinate{
							&Coordinate{
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
