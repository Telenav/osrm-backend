package stationfinderalg

import (
	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/connectivitymap"
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/common"
	"github.com/Telenav/osrm-backend/integration/service/oasis/stationfinder/stationfindertype"
	"github.com/golang/glog"
)

type placeID2QueryResults map[common.PlaceID][]*common.RankedPlaceInfo
type placeID2Location map[common.PlaceID]*nav.Location

type querier struct {
	id2QueryResults placeID2QueryResults
	id2Location     placeID2Location
}

// NewQuerierBasedOnWeightBetweenNeighborsChan creates connectivitymap.Querier based on channel of charge station's WeightBetweenNeighbors
func NewQuerierBasedOnWeightBetweenNeighborsChan(c chan stationfindertype.WeightBetweenNeighbors) connectivitymap.Querier {
	querier := &querier{
		id2QueryResults: make(placeID2QueryResults),
		id2Location:     make(placeID2Location),
	}

	for item := range c {
		if item.Err != nil {
			glog.Errorf("Met error during constructing stationgraph, error = %v", item.Err)
			return nil
		}

		for _, neighborInfo := range item.NeighborsInfo {

			if _, ok := querier.id2QueryResults[neighborInfo.FromPlaceID()]; !ok {
				results := make([]*common.RankedPlaceInfo, 0, 10)
				querier.id2QueryResults[neighborInfo.FromPlaceID()] = results
			}
			querier.id2QueryResults[neighborInfo.FromPlaceID()] = append(querier.id2QueryResults[neighborInfo.FromPlaceID()],
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

			if _, ok := querier.id2Location[neighborInfo.FromPlaceID()]; !ok {
				querier.id2Location[neighborInfo.FromPlaceID()] = &nav.Location{
					Lat: neighborInfo.FromLocation.Lat,
					Lon: neighborInfo.FromLocation.Lon,
				}
			}

			if _, ok := querier.id2Location[neighborInfo.ToPlaceID()]; !ok {
				querier.id2Location[neighborInfo.ToPlaceID()] = &nav.Location{
					Lat: neighborInfo.ToLocation.Lat,
					Lon: neighborInfo.ToLocation.Lon,
				}
			}
		}
	}

	return querier
}

func (q *querier) NearByStationQuery(placeID common.PlaceID) []*common.RankedPlaceInfo {
	if results, ok := q.id2QueryResults[placeID]; ok {
		return results
	} else {
		return nil
	}
}

func (q *querier) GetLocation(placeID common.PlaceID) *nav.Location {
	if location, ok := q.id2Location[placeID]; ok {
		return location
	} else {
		return nil
	}
}
