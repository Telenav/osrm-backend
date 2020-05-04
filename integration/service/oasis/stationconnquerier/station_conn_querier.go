package stationconnquerier

import (
	"strconv"

	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/connectivitymap"
	"github.com/Telenav/osrm-backend/integration/service/oasis/spatialindexer"
	"github.com/Telenav/osrm-backend/integration/service/oasis/stationfinder/stationfindertype"
	"github.com/golang/glog"
)

type StationConnectivityQuerier struct {
	stationLocationQuerier   spatialindexer.PlaceLocationQuerier
	stationConnectivity      *connectivitymap.ConnectivityMap
	reachableStationsByStart []*connectivitymap.QueryResult
	reachableStationToEnd    map[string]*connectivitymap.QueryResult
	startLocation            *nav.Location
	endLocation              *nav.Location
}

func New(stationFinder spatialindexer.Finder, stationRanker spatialindexer.Ranker,
	stationLocationQuerier spatialindexer.PlaceLocationQuerier,
	stationConnectivity *connectivitymap.ConnectivityMap,
	start, end *nav.Location,
	currEnergyLevel, maxEnergyLevel float64) connectivitymap.Querier {

	querier := &StationConnectivityQuerier{
		stationLocationQuerier: stationLocationQuerier,
		stationConnectivity:    stationConnectivity,
		startLocation:          start,
		endLocation:            end,
	}
	querier.connectStartIntoStationGraph(stationFinder, stationRanker, start, currEnergyLevel)
	querier.connectEndIntoStationGraph(stationFinder, stationRanker, end, maxEnergyLevel)

	return querier
}

func (querier *StationConnectivityQuerier) connectStartIntoStationGraph(stationFinder spatialindexer.Finder, stationRanker spatialindexer.Ranker,
	start *nav.Location, currEnergyLevel float64) {
	center := spatialindexer.Location{Lat: start.Lat, Lon: start.Lon}
	nearByPoints := stationFinder.FindNearByPointIDs(center, currEnergyLevel, spatialindexer.UnlimitedCount)
	rankedPoints := stationRanker.RankPointIDsByShortestDistance(center, nearByPoints)

	reachableStationsByStart := make([]*connectivitymap.QueryResult, 0, len(rankedPoints))
	for _, rankedPointInfo := range rankedPoints {
		tmp := &connectivitymap.QueryResult{
			StationID:       rankedPointInfo.ID.String(),
			StationLocation: &nav.Location{Lat: rankedPointInfo.Location.Lat, Lon: rankedPointInfo.Location.Lon},
			Distance:        rankedPointInfo.Distance,
			// TODO codebear801 Replace with pre-calculate duration https://github.com/Telenav/osrm-backend/issues/321
			Duration: rankedPointInfo.Distance,
		}
		reachableStationsByStart = append(reachableStationsByStart, tmp)
	}

	querier.reachableStationsByStart = reachableStationsByStart
}

func (querier *StationConnectivityQuerier) connectEndIntoStationGraph(stationFinder spatialindexer.Finder, stationRanker spatialindexer.Ranker,
	end *nav.Location, maxEnergyLevel float64) {
	center := spatialindexer.Location{Lat: end.Lat, Lon: end.Lon}
	nearByPoints := stationFinder.FindNearByPointIDs(center, maxEnergyLevel, spatialindexer.UnlimitedCount)
	rankedPoints := stationRanker.RankPointIDsByShortestDistance(center, nearByPoints)

	reachableStationToEnd := make(map[string]*connectivitymap.QueryResult)
	for _, rankedPointInfo := range rankedPoints {
		reachableStationToEnd[rankedPointInfo.ID.String()] = &connectivitymap.QueryResult{
			StationID:       stationfindertype.DestLocationID,
			StationLocation: end,
			Distance:        rankedPointInfo.Distance,
			//TODO codebear801 https://github.com/Telenav/osrm-backend/issues/321
			Duration: rankedPointInfo.Distance,
		}
	}

	querier.reachableStationToEnd = reachableStationToEnd
}

// NearByStationQuery finds near by stations by given stationID and return them in recorded sequence
// Returns nil if given stationID is not found or no connectivity
func (querier *StationConnectivityQuerier) NearByStationQuery(stationID string) []*connectivitymap.QueryResult {

	if stationID == stationfindertype.OrigLocationID {
		return querier.reachableStationsByStart
	}

	if stationID == stationfindertype.DestLocationID {
		return nil
	}

	placeID, err := strconv.Atoi(stationID)
	if err != nil {
		glog.Errorf("Incorrect station ID passed to NearByStationQuery %+v, got error %#v", stationID, err)
		return nil
	}
	if connectivityResults, ok := querier.stationConnectivity.QueryConnectivity((spatialindexer.PointID)(placeID)); ok {
		size := len(connectivityResults)
		if querier.isStationConnectsToEnd(stationID) {
			size += 1
		}

		results := make([]*connectivitymap.QueryResult, 0, size)
		for _, idAndWeight := range connectivityResults {
			tmp := &connectivitymap.QueryResult{
				StationID:       idAndWeight.ID.String(),
				StationLocation: querier.GetLocation(idAndWeight.ID.String()),
				Distance:        idAndWeight.Distance,
				//TODO codebear801 https://github.com/Telenav/osrm-backend/issues/321
				Duration: idAndWeight.Distance,
			}
			results = append(results, tmp)
		}

		return querier.connectEndIntoGraph(stationID, results)
	} else {
		if querier.isStationConnectsToEnd(stationID) {
			results := make([]*connectivitymap.QueryResult, 0, 1)
			return querier.connectEndIntoGraph(stationID, results)
		}
	}

	return nil
}

// GetLocation returns location of given station id
// Returns nil if given stationID is not found
func (querier *StationConnectivityQuerier) GetLocation(stationID string) *nav.Location {
	switch stationID {
	case stationfindertype.OrigLocationID:
		return querier.startLocation
	case stationfindertype.DestLocationID:
		return querier.endLocation
	default:
		return querier.stationLocationQuerier.GetLocation(stationID)
	}
}

func (querier *StationConnectivityQuerier) isStationConnectsToEnd(stationID string) bool {
	_, ok := querier.reachableStationToEnd[stationID]
	return ok
}

func (querier *StationConnectivityQuerier) connectEndIntoGraph(stationID string, results []*connectivitymap.QueryResult) []*connectivitymap.QueryResult {
	if queryResult4End, ok := querier.reachableStationToEnd[stationID]; ok {
		results = append(results, queryResult4End)
	}

	return results
}
