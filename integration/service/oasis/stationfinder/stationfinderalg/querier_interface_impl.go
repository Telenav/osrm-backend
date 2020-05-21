package stationfinderalg

import (
	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/connectivitymap"
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/common"
	"github.com/Telenav/osrm-backend/integration/service/oasis/stationfinder/stationfindertype"
	"github.com/golang/glog"
)

type stationID2QueryResults map[string][]*common.RankedPlaceInfo
type stationID2Location map[string]*nav.Location

type querier struct {
	id2QueryResults stationID2QueryResults
	id2Location     stationID2Location
}

// NewQuerierBasedOnWeightBetweenNeighborsChan creates connectivitymap.Querier based on channel of charge station's WeightBetweenNeighbors
func NewQuerierBasedOnWeightBetweenNeighborsChan(c chan stationfindertype.WeightBetweenNeighbors) connectivitymap.Querier {
	querier := &querier{
		id2QueryResults: make(stationID2QueryResults),
		id2Location:     make(stationID2Location),
	}

	for item := range c {
		if item.Err != nil {
			glog.Errorf("Met error during constructing stationgraph, error = %v", item.Err)
			return nil
		}

		for _, neighborInfo := range item.NeighborsInfo {

			if _, ok := querier.id2QueryResults[neighborInfo.FromID]; !ok {
				results := make([]*common.RankedPlaceInfo, 0, 10)
				querier.id2QueryResults[neighborInfo.FromID] = results
			}
			querier.id2QueryResults[neighborInfo.FromID] = append(querier.id2QueryResults[neighborInfo.FromID],
				&common.RankedPlaceInfo{
					PlaceInfo: common.PlaceInfo{
						ID: neighborInfo.ToPlaceID(),
						Location: &nav.Location{
							Lat: neighborInfo.ToLocation.Lat,
							Lon: neighborInfo.ToLocation.Lon,
						},
					},
					Weight: &common.Weight{
						Distance: neighborInfo.Distance,
						Duration: neighborInfo.Duration,
					},
				})

			if _, ok := querier.id2Location[neighborInfo.FromID]; !ok {
				querier.id2Location[neighborInfo.FromID] = &nav.Location{
					Lat: neighborInfo.FromLocation.Lat,
					Lon: neighborInfo.FromLocation.Lon,
				}
			}

			if _, ok := querier.id2Location[neighborInfo.ToID]; !ok {
				querier.id2Location[neighborInfo.ToID] = &nav.Location{
					Lat: neighborInfo.ToLocation.Lat,
					Lon: neighborInfo.ToLocation.Lon,
				}
			}
		}
	}

	return querier
}

func (q *querier) NearByStationQuery(stationID string) []*common.RankedPlaceInfo {
	if results, ok := q.id2QueryResults[stationID]; ok {
		return results
	} else {
		return nil
	}
}

func (q *querier) GetLocation(stationID string) *nav.Location {
	if location, ok := q.id2Location[stationID]; ok {
		return location
	} else {
		return nil
	}
}
