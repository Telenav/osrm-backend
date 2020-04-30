package ranker

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/Telenav/osrm-backend/integration/api/osrm/coordinate"
	"github.com/Telenav/osrm-backend/integration/api/osrm/genericoptions"
	"github.com/Telenav/osrm-backend/integration/api/osrm/table"
	"github.com/Telenav/osrm-backend/integration/service/oasis/osrmconnector"
	"github.com/Telenav/osrm-backend/integration/service/oasis/spatialindexer"
)

func TestGenerateTableRequest(t *testing.T) {
	cases := []struct {
		center     spatialindexer.Location
		targets    []*spatialindexer.PointInfo
		startIndex int
		endIndex   int
		expect     *table.Request
	}{
		// case 1: test center -> {1, 2, 3, 4, 5}
		{
			center: spatialindexer.Location{
				Lat: 0,
				Lon: 0,
			},
			targets: []*spatialindexer.PointInfo{
				&spatialindexer.PointInfo{
					ID: 1,
					Location: spatialindexer.Location{
						Lat: 1.1,
						Lon: 1.1,
					},
				},
				&spatialindexer.PointInfo{
					ID: 2,
					Location: spatialindexer.Location{
						Lat: 2.2,
						Lon: 2.2,
					},
				},
				&spatialindexer.PointInfo{
					ID: 3,
					Location: spatialindexer.Location{
						Lat: 3.3,
						Lon: 3.3,
					},
				},
				&spatialindexer.PointInfo{
					ID: 4,
					Location: spatialindexer.Location{
						Lat: 4.4,
						Lon: 4.4,
					},
				},
				&spatialindexer.PointInfo{
					ID: 5,
					Location: spatialindexer.Location{
						Lat: 5.5,
						Lon: 5.5,
					},
				},
			},
			startIndex: 0,
			endIndex:   4,
			expect: &table.Request{
				Service: "table",
				Version: "v1",
				Profile: "driving",
				Coordinates: coordinate.Coordinates{
					coordinate.Coordinate{
						Lat: 0,
						Lon: 0,
					},
					coordinate.Coordinate{
						Lat: 1.1,
						Lon: 1.1,
					},
					coordinate.Coordinate{
						Lat: 2.2,
						Lon: 2.2,
					},
					coordinate.Coordinate{
						Lat: 3.3,
						Lon: 3.3,
					},
					coordinate.Coordinate{
						Lat: 4.4,
						Lon: 4.4,
					},
					coordinate.Coordinate{
						Lat: 5.5,
						Lon: 5.5,
					},
				},
				Sources: genericoptions.Elements{
					"0",
				},
				Destinations: genericoptions.Elements{
					"1",
					"2",
					"3",
					"4",
					"5",
				},
				Annotations: "distance,duration",
			},
		},
		// case 2: test center -> {1, 2, 3}
		{
			center: spatialindexer.Location{
				Lat: 0,
				Lon: 0,
			},
			targets: []*spatialindexer.PointInfo{
				&spatialindexer.PointInfo{
					ID: 1,
					Location: spatialindexer.Location{
						Lat: 1.1,
						Lon: 1.1,
					},
				},
				&spatialindexer.PointInfo{
					ID: 2,
					Location: spatialindexer.Location{
						Lat: 2.2,
						Lon: 2.2,
					},
				},
				&spatialindexer.PointInfo{
					ID: 3,
					Location: spatialindexer.Location{
						Lat: 3.3,
						Lon: 3.3,
					},
				},
				&spatialindexer.PointInfo{
					ID: 4,
					Location: spatialindexer.Location{
						Lat: 4.4,
						Lon: 4.4,
					},
				},
				&spatialindexer.PointInfo{
					ID: 5,
					Location: spatialindexer.Location{
						Lat: 5.5,
						Lon: 5.5,
					},
				},
			},
			startIndex: 1,
			endIndex:   3,
			expect: &table.Request{
				Service: "table",
				Version: "v1",
				Profile: "driving",
				Coordinates: coordinate.Coordinates{
					coordinate.Coordinate{
						Lat: 0,
						Lon: 0,
					},
					coordinate.Coordinate{
						Lat: 2.2,
						Lon: 2.2,
					},
					coordinate.Coordinate{
						Lat: 3.3,
						Lon: 3.3,
					},
					coordinate.Coordinate{
						Lat: 4.4,
						Lon: 4.4,
					},
				},
				Sources: genericoptions.Elements{
					"0",
				},
				Destinations: genericoptions.Elements{
					"1",
					"2",
					"3",
				},
				Annotations: "distance,duration",
			},
		},
	}

	for _, c := range cases {
		actual := generateTableRequest(c.center, c.targets, c.startIndex, c.endIndex)
		if !reflect.DeepEqual(actual, c.expect) {
			t.Errorf("During TestGenerateTableRequest, expect table request is \n%+v\n but actual is \n%+v\n", c.expect, actual)
		}
	}

}

// TestRankPointsByOSRMShortestPathWithDifferentPointThreshold tries with different threshold and assert for same response
// In this case, target array contains 6 points {1, 2, 3, 4, 5, 6}
// when pointThreshold = pointsThresholdPerRequest(MaxInt32), all of them will be put in the same request
// when pointThreshold = 3, will divide points into two request, {1, 2, 3}, {4, 5, 6}
//                          but the result will be merged and sorted, so is the same with single request
// when pointThreshold = 4, will divide points into two request, {1, 2, 3, 4}, {5, 6}
//                          but the result will be merged and sorted, so is the same with single request
func TestRankPointsByOSRMShortestPathWithDifferentPointThreshold(t *testing.T) {
	cases := []struct {
		center  spatialindexer.Location
		targets []*spatialindexer.PointInfo
		expect  []*spatialindexer.RankedPointInfo
	}{
		{
			center: spatialindexer.Location{
				Lat: 0,
				Lon: 0,
			},
			targets: []*spatialindexer.PointInfo{
				&spatialindexer.PointInfo{
					ID: 1,
					Location: spatialindexer.Location{
						Lat: 1.1,
						Lon: 1.1,
					},
				},
				&spatialindexer.PointInfo{
					ID: 2,
					Location: spatialindexer.Location{
						Lat: 2.2,
						Lon: 2.2,
					},
				},
				&spatialindexer.PointInfo{
					ID: 3,
					Location: spatialindexer.Location{
						Lat: 3.3,
						Lon: 3.3,
					},
				},
				&spatialindexer.PointInfo{
					ID: 4,
					Location: spatialindexer.Location{
						Lat: 4.4,
						Lon: 4.4,
					},
				},
				&spatialindexer.PointInfo{
					ID: 5,
					Location: spatialindexer.Location{
						Lat: 5.5,
						Lon: 5.5,
					},
				},
				&spatialindexer.PointInfo{
					ID: 6,
					Location: spatialindexer.Location{
						Lat: 6.6,
						Lon: 6.6,
					},
				},
			},
			expect: []*spatialindexer.RankedPointInfo{
				&spatialindexer.RankedPointInfo{
					PointInfo: spatialindexer.PointInfo{
						ID: 1,
						Location: spatialindexer.Location{
							Lat: 1.1,
							Lon: 1.1,
						},
					},
					Distance: 1.1,
				},
				&spatialindexer.RankedPointInfo{
					PointInfo: spatialindexer.PointInfo{
						ID: 2,
						Location: spatialindexer.Location{
							Lat: 2.2,
							Lon: 2.2,
						},
					},
					Distance: 2.2,
				},
				&spatialindexer.RankedPointInfo{
					PointInfo: spatialindexer.PointInfo{
						ID: 3,
						Location: spatialindexer.Location{
							Lat: 3.3,
							Lon: 3.3,
						},
					},
					Distance: 3.3,
				},
				&spatialindexer.RankedPointInfo{
					PointInfo: spatialindexer.PointInfo{
						ID: 4,
						Location: spatialindexer.Location{
							Lat: 4.4,
							Lon: 4.4,
						},
					},
					Distance: 4.4,
				},
				&spatialindexer.RankedPointInfo{
					PointInfo: spatialindexer.PointInfo{
						ID: 5,
						Location: spatialindexer.Location{
							Lat: 5.5,
							Lon: 5.5,
						},
					},
					Distance: 5.5,
				},
				&spatialindexer.RankedPointInfo{
					PointInfo: spatialindexer.PointInfo{
						ID: 6,
						Location: spatialindexer.Location{
							Lat: 6.6,
							Lon: 6.6,
						},
					},
					Distance: 6.6,
				},
			},
		},
	}

	// fake OSRM server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		if r.Method != "GET" {
			t.Errorf("Expected 'GET' request, got '%s'", r.Method)
		}

		if strings.HasPrefix(r.URL.EscapedPath(), "/table/v1/driving/") {
			req, _ := table.ParseRequestURL(r.URL)

			s := len(req.Sources)
			d := len(req.Destinations)
			if s == 1 && d == 6 {
				var tableResponseBytes, _ = json.Marshal(mock1To6TableResponse)
				w.Write(tableResponseBytes)
			} else if s == 1 && d == 3 &&
				reflect.DeepEqual(req.Coordinates, coordinate.Coordinates{
					coordinate.Coordinate{Lat: 0, Lon: 0},
					coordinate.Coordinate{Lat: 1.1, Lon: 1.1},
					coordinate.Coordinate{Lat: 2.2, Lon: 2.2},
					coordinate.Coordinate{Lat: 3.3, Lon: 3.3}}) {
				var tableResponseBytes, _ = json.Marshal(mock1To3TableResponsePart1)
				w.Write(tableResponseBytes)
			} else if s == 1 && d == 3 &&
				reflect.DeepEqual(req.Coordinates, coordinate.Coordinates{
					coordinate.Coordinate{Lat: 0, Lon: 0},
					coordinate.Coordinate{Lat: 4.4, Lon: 4.4},
					coordinate.Coordinate{Lat: 5.5, Lon: 5.5},
					coordinate.Coordinate{Lat: 6.6, Lon: 6.6}}) {
				var tableResponseBytes, _ = json.Marshal(mock1To3TableResponsePart2)
				w.Write(tableResponseBytes)
			} else if s == 1 && d == 4 &&
				reflect.DeepEqual(req.Coordinates, coordinate.Coordinates{
					coordinate.Coordinate{Lat: 0, Lon: 0},
					coordinate.Coordinate{Lat: 1.1, Lon: 1.1},
					coordinate.Coordinate{Lat: 2.2, Lon: 2.2},
					coordinate.Coordinate{Lat: 3.3, Lon: 3.3},
					coordinate.Coordinate{Lat: 4.4, Lon: 4.4}}) {
				var tableResponseBytes, _ = json.Marshal(mock1To4TableResponsePart1)
				w.Write(tableResponseBytes)
			} else if s == 1 && d == 2 &&
				reflect.DeepEqual(req.Coordinates, coordinate.Coordinates{
					coordinate.Coordinate{Lat: 0, Lon: 0},
					coordinate.Coordinate{Lat: 5.5, Lon: 5.5},
					coordinate.Coordinate{Lat: 6.6, Lon: 6.6}}) {
				var tableResponseBytes, _ = json.Marshal(mock1To4TableResponsePart2)
				w.Write(tableResponseBytes)
			}
			return
		}

	}))
	defer ts.Close()

	oc := osrmconnector.NewOSRMConnector(ts.URL)

	for _, c := range cases {
		actual := rankPointsByOSRMShortestPath(c.center, c.targets, oc, pointsThresholdPerRequest)
		if !reflect.DeepEqual(actual, c.expect) {
			t.Errorf("During TestRankerInterfaceViaOSRMRanker, expect \n%s \nwhile actual result is \n%s\n",
				printRankedPointInfoArray(c.expect),
				printRankedPointInfoArray(actual))
		}

		actual = rankPointsByOSRMShortestPath(c.center, c.targets, oc, 3)
		if !reflect.DeepEqual(actual, c.expect) {
			t.Errorf("During TestRankerInterfaceViaOSRMRanker, expect \n%s \nwhile actual result is \n%s\n",
				printRankedPointInfoArray(c.expect),
				printRankedPointInfoArray(actual))
		}

		actual = rankPointsByOSRMShortestPath(c.center, c.targets, oc, 4)
		if !reflect.DeepEqual(actual, c.expect) {
			t.Errorf("During TestRankerInterfaceViaOSRMRanker, expect \n%s \nwhile actual result is \n%s\n",
				printRankedPointInfoArray(c.expect),
				printRankedPointInfoArray(actual))
		}
	}
}
