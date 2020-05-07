package ranker

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/Telenav/osrm-backend/integration/api"
	"github.com/Telenav/osrm-backend/integration/api/osrm"
	"github.com/Telenav/osrm-backend/integration/api/osrm/table"
	"github.com/Telenav/osrm-backend/integration/api/osrm/table/options"
	"github.com/Telenav/osrm-backend/integration/service/oasis/osrmconnector"
	"github.com/Telenav/osrm-backend/integration/service/oasis/spatialindexer"
	"github.com/golang/glog"
)

// pointsThresholdPerRequest limits max point count for each single table request.
// During pre-processing, its possible to calculate distance between thousands of points, which will
// cause too big request and might reach potential limitation of different parts.
// The scenario here is 1-to-N table request, use pointsLimit4SingleTableRequest to limit N
const pointsThresholdPerRequest = 500

func rankPointsByOSRMShortestPath(center spatialindexer.Location, targets []*spatialindexer.PointInfo,
	oc *osrmconnector.OSRMConnector, pointsThreshold int) []*spatialindexer.RankedPointInfo {
	if len(targets) == 0 {
		glog.Warning("When try to rankPointsByGreatCircleDistanceToCenter, input array is empty\n")
		return nil
	}

	var wg sync.WaitGroup
	pointWithDistanceC := make(chan *spatialindexer.RankedPointInfo, len(targets))
	startIndex := 0 // startIndex is a valid index of targets
	endIndex := 0   // endIndex is valid index of targets
	for {
		if startIndex >= len(targets) {
			break
		}

		endIndex = startIndex + pointsThreshold - 1
		if endIndex >= len(targets) {
			endIndex = len(targets) - 1
		}

		wg.Add(1)
		go func(wg *sync.WaitGroup, startIndex, endIndex int) {
			rankedPoints, err := calcCenter2TargetsDistanceViaShortestPath(center, targets, oc, startIndex, endIndex)

			if err != nil {
				glog.Errorf("Failed to calculate shortest path for range [%d, %d] for center = %+v, len(targets) = %+v, err = %+v\n",
					startIndex, endIndex, center, len(targets), err)
				// @todo: add retry logic when failed or may be put retry logic in connector
			} else {
				for _, item := range rankedPoints {
					pointWithDistanceC <- item
				}
			}

			wg.Done()
		}(&wg, startIndex, endIndex)

		startIndex = endIndex + 1
	}

	wg.Wait()
	close(pointWithDistanceC)

	rankAgent := newRankAgent(len(targets))
	return rankAgent.RankByDistance(pointWithDistanceC)

}

func calcCenter2TargetsDistanceViaShortestPath(center spatialindexer.Location, targets []*spatialindexer.PointInfo, oc *osrmconnector.OSRMConnector, startIndex, endIndex int) ([]*spatialindexer.RankedPointInfo, error) {
	req := generateTableRequest(center, targets, startIndex, endIndex)
	respC := oc.Request4Table(req)
	resp := <-respC

	if resp.Err != nil {
		glog.Errorf("Failed to generate table response for \n %s with \n err =%v \n", req.RequestURI(), resp.Err)
		return nil, resp.Err
	}

	if resp.Resp.Code != osrm.CodeOK {
		err := fmt.Errorf("failed to get correct table response for Request = %+v \n Resp = %+v", req.RequestURI(), resp.Resp)
		return nil, err
	}

	if len(resp.Resp.Distances[0]) != endIndex-startIndex+1 {
		err := fmt.Errorf("incorrect table response for request %+v, expect %d items but returns %d", req.RequestURI(), endIndex-startIndex+1, len(resp.Resp.Distances[0]))
		return nil, err
	}

	glog.V(3).Infof("Inside ranker, get table response for request %+v\n", resp.Resp)
	glog.V(3).Infof("In the response,  len(resp.Resp.Distances[0]) = %+v\n", len(resp.Resp.Distances[0]))

	result := make([]*spatialindexer.RankedPointInfo, 0, endIndex-startIndex+1)
	for i := 0; i < endIndex-startIndex+1; i++ {
		result = append(result, &spatialindexer.RankedPointInfo{
			PointInfo: spatialindexer.PointInfo{
				ID:       targets[startIndex+i].ID,
				Location: targets[startIndex+i].Location,
			},
			Distance: resp.Resp.Distances[0][i],
		})
	}
	return result, nil
}

// generateTableRequest generates table requests from center to [startIndex, endIndex] of targets
func generateTableRequest(center spatialindexer.Location, targets []*spatialindexer.PointInfo, startIndex, endIndex int) *table.Request {
	if startIndex < 0 || startIndex > endIndex || endIndex >= len(targets) {
		glog.Fatalf("startIndex should be smaller equal to endIndex, and both of them should in the range of len(targets), while (startIndex, endIndex, len(targets)) = (%d, %d, %d)",
			startIndex, endIndex, len(targets))
	}

	req := table.NewRequest()
	req.Coordinates = append(convertLocation2Coordinates(center),
		convertPointInfos2Coordinates(targets, startIndex, endIndex)...)

	req.Sources = append(req.Sources, strconv.Itoa(0))
	pointsCount4Sources := 1
	for i := 0; i <= endIndex-startIndex; i++ {
		str := strconv.Itoa(i + pointsCount4Sources)
		req.Destinations = append(req.Destinations, str)
	}

	req.Annotations = options.AnnotationsValueDistance + api.Comma + options.AnnotationsValueDuration

	return req
}

func convertLocation2Coordinates(location spatialindexer.Location) osrm.Coordinates {
	result := make(osrm.Coordinates, 0, 1)
	result = append(result, osrm.Coordinate{
		Lat: location.Lat,
		Lon: location.Lon,
	})
	return result
}

func convertPointInfos2Coordinates(targets []*spatialindexer.PointInfo, startIndex, endIndex int) osrm.Coordinates {
	result := make(osrm.Coordinates, 0, endIndex-startIndex+1)
	for i := startIndex; i <= endIndex; i++ {
		result = append(result, osrm.Coordinate{
			Lat: targets[i].Location.Lat,
			Lon: targets[i].Location.Lon,
		})
	}
	return result
}
