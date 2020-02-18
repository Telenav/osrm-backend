package stationhandler

import (
	"fmt"

	"github.com/Telenav/osrm-backend/integration/oasis/osrmconnector"
	"github.com/Telenav/osrm-backend/integration/oasis/osrmhelper"
	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/coordinate"
	"github.com/golang/glog"
)

type cost struct {
	duration float64
	distance float64
}

type costBetweenChargeStations struct {
	fromID string
	toID   string
	cost
	err error
}

func calcCostBetweenChargeStationsPair(from nearbyStationsIterator, to nearbyStationsIterator, oc *osrmconnector.OSRMConnector) <-chan costBetweenChargeStations {
	c := make(chan costBetweenChargeStations, 1000)

	// collect (lat,lon)&ID for current location's nearby charge stations
	startPoints := make(coordinate.Coordinates, 50)
	startIDs := make([]string, 50)
	for v := range from.iterateNearbyStations() {
		startPoints = append(startPoints, v.location)
		startIDs = append(startIDs, v.id)
	}
	if len(startPoints) == 0 {
		closeChannelWithErrorInfo(c, fmt.Errorf("no station be found for current point during calcCostBetweenChargeStationsPair"))
		return c
	}

	// collect (lat,lon)&ID for target location's nearby charge stations
	targetPoints := make(coordinate.Coordinates, 50)
	targetIDs := make([]string, 50)
	for v := range to.iterateNearbyStations() {
		targetPoints = append(targetPoints, v.location)
		targetIDs = append(targetIDs, v.id)
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
			costPair := costBetweenChargeStations{
				fromID: startIDs[i],
				toID:   targetIDs[j],
				cost: cost{
					duration: *resp.Resp.Durations[i][j],
					distance: *resp.Resp.Distances[i][j],
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

func closeChannelWithErrorInfo(c chan costBetweenChargeStations, err error) {
	go func() {
		defer close(c)
		errobj := costBetweenChargeStations{
			err: err,
		}
		glog.Warningf("%v", errobj.err)
		c <- errobj
	}()
}
