package ranker

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/Telenav/osrm-backend/integration/api/osrm/code"
	"github.com/Telenav/osrm-backend/integration/api/osrm/table"
	"github.com/Telenav/osrm-backend/integration/service/oasis/osrmconnector"
	"github.com/Telenav/osrm-backend/integration/service/oasis/spatialindexer"
)

func TestRankerInterfaceViaOSRMRanker(t *testing.T) {
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
					Duration: 1.1,
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
					Duration: 2.2,
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
					Duration: 3.3,
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
					Duration: 4.4,
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
					Duration: 5.5,
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
					Duration: 6.6,
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
				var tableResponseBytesOrig2Location1, _ = json.Marshal(mock1To6TableResponse)
				w.Write(tableResponseBytesOrig2Location1)
			}
			return
		}

	}))
	defer ts.Close()

	oc := osrmconnector.NewOSRMConnector(ts.URL)
	ranker := CreateRanker(OSRMBasedRanker, oc)

	for _, c := range cases {
		actual := ranker.RankPointIDsByShortestDistance(c.center, c.targets)
		if !reflect.DeepEqual(actual, c.expect) {
			t.Errorf("During TestRankerInterfaceViaOSRMRanker, expect \n%s \nwhile actual result is \n%s\n",
				printRankedPointInfoArray(c.expect),
				printRankedPointInfoArray(actual))
		}
	}
}

var mockFloatArray1To6 []float64 = []float64{1.1, 2.2, 3.3, 4.4, 5.5, 6.6}

var mock1To6TableResponse table.Response = table.Response{
	Code: code.OK,
	Durations: [][]float64{
		[]float64{
			mockFloatArray1To6[0],
			mockFloatArray1To6[1],
			mockFloatArray1To6[2],
			mockFloatArray1To6[3],
			mockFloatArray1To6[4],
			mockFloatArray1To6[5],
		},
	},
	Distances: [][]float64{
		[]float64{
			mockFloatArray1To6[0],
			mockFloatArray1To6[1],
			mockFloatArray1To6[2],
			mockFloatArray1To6[3],
			mockFloatArray1To6[4],
			mockFloatArray1To6[5],
		},
	},
}

var mock1To3TableResponsePart1 table.Response = table.Response{
	Code: code.OK,
	Durations: [][]float64{
		[]float64{
			mockFloatArray1To6[0],
			mockFloatArray1To6[1],
			mockFloatArray1To6[2],
		},
	},
	Distances: [][]float64{
		[]float64{
			mockFloatArray1To6[0],
			mockFloatArray1To6[1],
			mockFloatArray1To6[2],
		},
	},
}

var mock1To3TableResponsePart2 table.Response = table.Response{
	Code: code.OK,
	Durations: [][]float64{
		[]float64{
			mockFloatArray1To6[3],
			mockFloatArray1To6[4],
			mockFloatArray1To6[5],
		},
	},
	Distances: [][]float64{
		[]float64{
			mockFloatArray1To6[3],
			mockFloatArray1To6[4],
			mockFloatArray1To6[5],
		},
	},
}

var mock1To4TableResponsePart1 table.Response = table.Response{
	Code: code.OK,
	Durations: [][]float64{
		[]float64{
			mockFloatArray1To6[0],
			mockFloatArray1To6[1],
			mockFloatArray1To6[2],
			mockFloatArray1To6[3],
		},
	},
	Distances: [][]float64{
		[]float64{
			mockFloatArray1To6[0],
			mockFloatArray1To6[1],
			mockFloatArray1To6[2],
			mockFloatArray1To6[3],
		},
	},
}

var mock1To4TableResponsePart2 table.Response = table.Response{
	Code: code.OK,
	Durations: [][]float64{
		[]float64{
			mockFloatArray1To6[4],
			mockFloatArray1To6[5],
		},
	},
	Distances: [][]float64{
		[]float64{
			mockFloatArray1To6[4],
			mockFloatArray1To6[5],
		},
	},
}
