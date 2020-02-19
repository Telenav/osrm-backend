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
	overlap := make([]ChargeStationInfo, 10)
	dict := buildChargeStationInfoDict(iterF)
	c := iterS.iterateNearbyStations()
	for item := range c {
		if _, has := dict[item.ID]; has {
			overlap = append(overlap, item)
		}
	}

	return overlap
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
	Err error
}

func CalcCostBetweenChargeStationsPair(from nearbyStationsIterator, to nearbyStationsIterator, oc *osrmconnector.OSRMConnector) <-chan CostBetweenChargeStations {
	c := make(chan CostBetweenChargeStations, 1000)

	// collect (lat,lon)&ID for current location's nearby charge stations
	startPoints := make(coordinate.Coordinates, 50)
	startIDs := make([]string, 50)
	for v := range from.iterateNearbyStations() {
		startPoints = append(startPoints, coordinate.Coordinate{
			Lat: v.Location.Lat,
			Lon: v.Location.Lon,
		})
		startIDs = append(startIDs, v.ID)
	}
	if len(startPoints) == 0 {
		closeChannelWithErrorInfo(c, fmt.Errorf("no station be found for current point during calcCostBetweenChargeStationsPair"))
		return c
	}

	// collect (lat,lon)&ID for target location's nearby charge stations
	targetPoints := make(coordinate.Coordinates, 50)
	targetIDs := make([]string, 50)
	for v := range to.iterateNearbyStations() {
		targetPoints = append(targetPoints, coordinate.Coordinate{
			Lat: v.Location.Lat,
			Lon: v.Location.Lon,
		})
		targetIDs = append(targetIDs, v.ID)
	}
	if len(targetPoints) == 0 {
		closeChannelWithErrorInfo(c, fmt.Errorf("no station be found for target point during calcCostBetweenChargeStationsPair"))
		return c
	}

	// generate table request
	req, err := osrmhelper.GenerateTableReq4Points(startPoints, targetPoints)
	if err != nil {
		closeChannelWithErrorInfo(c, err)
		return c
	}

	// request for table
	respC := oc.Request4Table(req)
	resp := <-respC
	if resp.Err != nil {
		closeChannelWithErrorInfo(c, resp.Err)
		return c
	}

	// iterate table response result
	if len(resp.Resp.Sources) != len(startPoints) || len(resp.Resp.Destinations) != len(targetPoints) {
		closeChannelWithErrorInfo(c, fmt.Errorf("Incorrect osrm table response for url: %s", req.RequestURI()))
		return c
	}
	for i := range startPoints {
		for j := range targetPoints {
			costPair := CostBetweenChargeStations{
				FromID: startIDs[i],
				ToID:   targetIDs[j],
				Cost: Cost{
					Duration: *resp.Resp.Durations[i][j],
					Distance: *resp.Resp.Distances[i][j],
				},
			}
			go func() {
				c <- costPair
			}()
		}
	}

	close(c)
	return c
}

func closeChannelWithErrorInfo(c chan CostBetweenChargeStations, err error) {
	go func() {
		defer close(c)
		errobj := CostBetweenChargeStations{
			Err: err,
		}
		glog.Warningf("%v", errobj.Err)
		c <- errobj
	}()
}

func buildChargeStationInfoDict(iter nearbyStationsIterator) map[string]bool {
	dict := make(map[string]bool)
	c := iter.iterateNearbyStations()
	for item := range c {
		dict[item.ID] = true
	}

	return dict
}
