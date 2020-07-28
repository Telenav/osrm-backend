package stationfinderalg

import (
	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/entity"
	"github.com/Telenav/osrm-backend/integration/service/oasis/place"
	"github.com/Telenav/osrm-backend/integration/service/oasis/place/stationfinder/stationfindertype"
	"github.com/golang/glog"
)

type placeID2QueryResults map[entity.PlaceID][]*entity.RankedPlaceInfo
type placeID2Location map[entity.PlaceID]*nav.Location

type querier struct {
	id2QueryResults placeID2QueryResults
	id2Location     placeID2Location
}

// NewQuerierBasedOnWeightBetweenNeighborsChan creates place.TopoQuerier based on channel of charge station's WeightBetweenNeighbors
func NewQuerierBasedOnWeightBetweenNeighborsChan(c chan stationfindertype.WeightBetweenNeighbors) place.TopoQuerier {
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
				results := make([]*entity.RankedPlaceInfo, 0, 10)
				querier.id2QueryResults[neighborInfo.FromPlaceID()] = results
			}
			querier.id2QueryResults[neighborInfo.FromPlaceID()] = append(querier.id2QueryResults[neighborInfo.FromPlaceID()],
				&entity.RankedPlaceInfo{
					PlaceInfo: entity.PlaceInfo{
						ID: neighborInfo.ToPlaceID(),
						Location: &nav.Location{
							Lat: neighborInfo.ToLocation.Lat,
							Lon: neighborInfo.ToLocation.Lon,
						},
					},
					Weight: &entity.Weight{
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

// NearByStationQuery finds near by stations by given placeID and return them in recorded sequence
// Returns nil if given placeID is not found or no connectivity
func (q *querier) NearByStationQuery(placeID entity.PlaceID) []*entity.RankedPlaceInfo {
	if results, ok := q.id2QueryResults[placeID]; ok {
		return results
	}

	return nil
}

// GetLocation returns location of given station id
// Returns nil if given placeID is not found
func (q *querier) GetLocation(placeID entity.PlaceID) *nav.Location {
	if location, ok := q.id2Location[placeID]; ok {
		return location
	} else {
		return nil
	}
}
