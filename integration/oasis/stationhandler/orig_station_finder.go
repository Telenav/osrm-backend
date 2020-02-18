package stationhandler

import (
	"sync"

	"github.com/Telenav/osrm-backend/integration/oasis/osrmconnector"
	"github.com/Telenav/osrm-backend/integration/oasis/searchconnector"
	"github.com/Telenav/osrm-backend/integration/oasis/searchhelper"
	"github.com/Telenav/osrm-backend/integration/pkg/api/oasis"
	searchcoordinate "github.com/Telenav/osrm-backend/integration/pkg/api/search/coordinate"
	"github.com/Telenav/osrm-backend/integration/pkg/api/search/nearbychargestation"
	"github.com/golang/glog"
)

//@todo: This number need to be adjusted based on charge station profile
const origMaxSearchCandidateNumber int = 999

type origStationFinder struct {
	osrmConnector     *osrmconnector.OSRMConnector
	tnSearchConnector *searchconnector.TNSearchConnector
	oasisReq          *oasis.Request
	searchResp        *nearbychargestation.Response
	searchRespLock    *sync.RWMutex
}

func NewOrigStationFinder(oc *osrmconnector.OSRMConnector, sc *searchconnector.TNSearchConnector, oasisReq *oasis.Request) *origStationFinder {
	obj := &origStationFinder{
		osrmConnector:     oc,
		tnSearchConnector: sc,
		oasisReq:          oasisReq,
		searchResp:        nil,
		searchRespLock:    &sync.RWMutex{},
	}
	obj.prepare()
	return obj
}

func (sf *origStationFinder) prepare() {
	req, _ := searchhelper.GenerateSearchRequest(
		searchcoordinate.Coordinate{
			Lat: sf.oasisReq.Coordinates[0].Lat,
			Lon: sf.oasisReq.Coordinates[0].Lon},
		origMaxSearchCandidateNumber,
		sf.oasisReq.CurrRange)

	respC := sf.tnSearchConnector.ChargeStationSearch(req)
	resp := <-respC
	if resp.Err != nil {
		glog.Warningf("Search failed during prepare orig search for url: %s", req.RequestURI())
		return
	}

	sf.searchRespLock.Lock()
	sf.searchResp = resp.Resp
	sf.searchRespLock.Unlock()
	return
}

func (sf *origStationFinder) iterateNearbyStations() <-chan chargeStationInfo {
	return iterateNearbyStations(sf.searchResp.Results, sf.searchRespLock)
}

// func (sf *origStationFinder) iterateNearbyStations() <-chan chargeStationInfo {
// 	if sf.searchResp == nil || len(sf.searchResp.Results) == 0 {
// 		c := make(chan chargeStationInfo)
// 		go func() {
// 			defer close(c)
// 		}()
// 		return c
// 	}

// 	c := make(chan chargeStationInfo, len(sf.searchResp.Results))
// 	results := make([]*nearbysearch.Result, len(sf.searchResp.Results))
// 	go func() {
// 		defer close(c)
// 		for result := range sf.searchResp.Results {
// 			if len(result.Place.Address) == 0 {
// 				continue
// 			}
// 			station := chargeStationInfo{
// 				id: result.ID,
// 				location: coordinate.Coordinate{
// 					Lat: result.Place.Address[0].GeoCoordinate.Lat,
// 					Lon: result.Place.Address[0].GeoCoordinate.Lon},
// 			}
// 			c <- station
// 		}
// 	}()
// }

// func (sf *origStationFinder) calcCostBetweenChargeStationsPair(another nearbyStationFinder) <- chan costBetweenChargeStations {
// 	c := make(chan costBetweenChargeStations, 1000)

// 	// collect (lat,lon)&ID for current location's nearby charge stations
// 	startPoints := make(Coordinate.Coordinates, 50)
// 	startIDs := make([]string, 50)
// 	for v := range sf.iterateNearbyStations() {
// 		startPoints = append(startPoints, v.location)
// 		startIDs = append(startIDs, v.ID)
// 	}
// 	if len(startPoints) == 0 {
// 		closeChannelWithErrorInfo(c, fmt.Error("No station be found for current point during calcCostBetweenChargeStationsPair.")
// 		return c
// 	}

// 	// collect (lat,lon)&ID for target location's nearby charge stations
// 	targetPoints := make(Coordinate.Coordinates, 50)
// 	targetIDs := make([]string, 50)
// 	for v := range another.iterateNearbyStations() {
// 		targetPoints = append(targetPoints, v.location)
// 		targetIDs = append(targetIDs, v.ID)
// 	}
// 	if len(targetPoints) == 0 {
// 		closeChannelWithErrorInfo(c, fmt.Error("No station be found for target point during calcCostBetweenChargeStationsPair."))
// 		return c
// 	}

// 	// generate table request
// 	req, err := generateTableReq4Points(startPoints, targetPoints)
// 	if err != nil {
// 		closeChannelWithErrorInfo(c, err)
// 		return c
// 	}

// 	// request for table
// 	respC := osrmConnector.Request4Table(req)
// 	resp := <- respC
// 	if resp.Err != nil {
// 		closeChannelWithErrorInfo(c, resp.Err)
// 		return c
// 	}

// 	// iterate table response result
// 	if len(resp.Resp.Sources) != len(startPoints) || len(resp.Resp.Destinations) != len(targetPoints) {
// 		closeChannelWithErrorInfo(c, fmt.Errorf("Incorrect osrm table response for url: %s", req.RequestURI()))
// 		return c
// 	}
// 	for i := range len(startPoints) {
// 		for j := range len(targetPoints) {
// 			costPair := costBetweenChargeStations {
// 				from_id : startIDs[i],
// 				to_id : targetIDs[j],
// 				duration : resp.Resp.Durations[i][j]
// 				distance : resp.Resp.Distance[i][j]
// 			}
// 			go func() {
// 				c <- costPair
// 			}()
// 		}
// 	}

// 	close(c)
// 	return c
// }

// func closeChannelWithErrorInfo(c chan costBetweenChargeStations, err error) {
// 	go func() {
// 		defer close(c)
// 		errobj := costBetweenChargeStations {
// 			err : err,
// 		}
// 		glog.Warningf("%v", errobj.err)
// 		c <- errobj
// 	}
// }
