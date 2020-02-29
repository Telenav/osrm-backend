package stationfinder

import (
	"fmt"

	"github.com/Telenav/osrm-backend/integration/oasis/osrmconnector"
	"github.com/Telenav/osrm-backend/integration/oasis/osrmhelper"
	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/coordinate"
	"github.com/golang/glog"
)

// FindOverlapBetweenStations finds overlap charge stations based on two iterator
func FindOverlapBetweenStations(iterF nearbyStationsIterator, iterS nearbyStationsIterator) []ChargeStationInfo {
	var overlap []ChargeStationInfo
	dict := buildChargeStationInfoDict(iterF)
	c := iterS.iterateNearbyStations()
	for item := range c {
		if _, has := dict[item.ID]; has {
			overlap = append(overlap, item)
		}
	}

	return overlap
}

// ChargeStationInfo defines charge station information
type ChargeStationInfo struct {
	ID       string
	Location StationCoordinate
	err      error
}

// StationCoordinate represents location information
type StationCoordinate struct {
	Lat float64
	Lon float64
}

func CalcCostBetweenChargeStationsPair(from nearbyStationsIterator, to nearbyStationsIterator, table osrmconnector.TableRequster) ([]CostBetweenChargeStations, error) {
	// collect (lat,lon)&ID for current location's nearby charge stations
	var startPoints coordinate.Coordinates
	var startIDs []string
	for v := range from.iterateNearbyStations() {
		startPoints = append(startPoints, coordinate.Coordinate{
			Lat: v.Location.Lat,
			Lon: v.Location.Lon,
		})
		startIDs = append(startIDs, v.ID)
	}
	if len(startPoints) == 0 {
		err := fmt.Errorf("Empty iterator of from pass into calcCostBetweenChargeStationsPair")
		glog.Warningf("%v", err)
		return nil, err
	}

	// collect (lat,lon)&ID for target location's nearby charge stations
	var targetPoints coordinate.Coordinates
	var targetIDs []string
	for v := range to.iterateNearbyStations() {
		targetPoints = append(targetPoints, coordinate.Coordinate{
			Lat: v.Location.Lat,
			Lon: v.Location.Lon,
		})
		targetIDs = append(targetIDs, v.ID)
	}
	if len(targetPoints) == 0 {
		err := fmt.Errorf("Empty iterator of to pass into calcCostBetweenChargeStationsPair")
		glog.Warningf("%v", err)
		return nil, err
	}

	// generate table request
	req, err := osrmhelper.GenerateTableReq4Points(startPoints, targetPoints)
	if err != nil {
		glog.Warningf("%v", err)
		return nil, err
	}

	// request for table
	respC := table.Request4Table(req)
	resp := <-respC
	if resp.Err != nil {
		glog.Warningf("%v", resp.Err)
		return nil, resp.Err
	}

	// iterate table response result
	if len(resp.Resp.Sources) != len(startPoints) || len(resp.Resp.Destinations) != len(targetPoints) {
		err := fmt.Errorf("Incorrect osrm table response for url: %s", req.RequestURI())
		return nil, err
	}

	var result []CostBetweenChargeStations
	for i := range startPoints {
		for j := range targetPoints {
			result = append(result, CostBetweenChargeStations{
				FromID: startIDs[i],
				ToID:   targetIDs[j],
				Cost: Cost{
					Duration: *resp.Resp.Durations[i][j],
					Distance: *resp.Resp.Distances[i][j],
				},
			})
		}
	}

	return result, nil
}

// Cost represent cost information
type Cost struct {
	Duration float64
	Distance float64
}

// CostBetweenChargeStations represent cost information between two charge stations
type CostBetweenChargeStations struct {
	FromID string
	ToID   string
	Cost
}

func buildChargeStationInfoDict(iter nearbyStationsIterator) map[string]bool {
	dict := make(map[string]bool)
	c := iter.iterateNearbyStations()
	for item := range c {
		dict[item.ID] = true
	}

	return dict
}

func CalculateWeightBetweenNeighbors()
