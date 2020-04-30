package stationfinderalg

import (
	"fmt"
	"sync"

	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/api/osrm/coordinate"
	"github.com/Telenav/osrm-backend/integration/service/oasis/osrmconnector"
	"github.com/Telenav/osrm-backend/integration/service/oasis/osrmhelper"
	"github.com/Telenav/osrm-backend/integration/service/oasis/stationfinder"
	"github.com/Telenav/osrm-backend/integration/service/oasis/stationfinder/stationfindertype"
	"github.com/golang/glog"
)

// FindOverlapBetweenStations finds overlap charge stations based on two iterator
func FindOverlapBetweenStations(iterF stationfindertype.NearbyStationsIterator, iterS stationfindertype.NearbyStationsIterator) []*stationfindertype.ChargeStationInfo {
	var overlap []*stationfindertype.ChargeStationInfo
	dict := buildChargeStationInfoDict(iterF)
	c := iterS.IterateNearbyStations()
	for item := range c {
		if _, has := dict[item.ID]; has {
			overlap = append(overlap, item)
		}
	}

	return overlap
}

// // nav.Location represents location information
// type nav.Location nav.Location

// CalcWeightBetweenChargeStationsPair accepts two iterators and calculates weights between each pair of iterators
func CalcWeightBetweenChargeStationsPair(from stationfindertype.NearbyStationsIterator, to stationfindertype.NearbyStationsIterator, table osrmconnector.TableRequster) ([]stationfindertype.NeighborInfo, error) {
	// collect (Lat,Lon)&ID for current location's nearby charge stations
	var startPoints coordinate.Coordinates
	var startIDs []string
	for v := range from.IterateNearbyStations() {
		startPoints = append(startPoints, coordinate.Coordinate{
			Lat: v.Location.Lat,
			Lon: v.Location.Lon,
		})
		startIDs = append(startIDs, v.ID)
	}
	if len(startPoints) == 0 {
		err := fmt.Errorf("empty iterator of from pass into CalcWeightBetweenChargeStationsPair")
		glog.Warningf("%v", err)
		return nil, err
	}

	// collect (Lat,Lon)&ID for target location's nearby charge stations
	var targetPoints coordinate.Coordinates
	var targetIDs []string
	for v := range to.IterateNearbyStations() {
		targetPoints = append(targetPoints, coordinate.Coordinate{
			Lat: v.Location.Lat,
			Lon: v.Location.Lon,
		})
		targetIDs = append(targetIDs, v.ID)
	}
	if len(targetPoints) == 0 {
		err := fmt.Errorf("empty iterator of to pass into CalcWeightBetweenChargeStationsPair")
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

	if len(resp.Resp.Sources) != len(startPoints) || len(resp.Resp.Destinations) != len(targetPoints) {
		err := fmt.Errorf("incorrect osrm table response for url: %s", req.RequestURI())
		return nil, err
	}

	// iterate table response result
	var result []stationfindertype.NeighborInfo
	for i, startPoint := range startPoints {
		for j, targetPoint := range targetPoints {
			result = append(result, stationfindertype.NeighborInfo{
				FromID: startIDs[i],
				FromLocation: nav.Location{
					Lat: startPoint.Lat,
					Lon: startPoint.Lon,
				},
				ToID: targetIDs[j],
				ToLocation: nav.Location{
					Lat: targetPoint.Lat,
					Lon: targetPoint.Lon,
				},
				Weight: stationfindertype.Weight{
					Duration: resp.Resp.Durations[i][j],
					Distance: resp.Resp.Distances[i][j],
				},
			})
		}
	}

	return result, nil
}

func buildChargeStationInfoDict(iter stationfindertype.NearbyStationsIterator) map[string]bool {
	dict := make(map[string]bool)
	c := iter.IterateNearbyStations()
	for item := range c {
		dict[item.ID] = true
	}

	return dict
}

// CalculateWeightBetweenNeighbors accepts locations array, which will search for nearby
// charge stations and then calculate weight between stations, the result is used to
// construct graph.
// - The input of locations contains: orig location -> first place to search for charge ->
//   second location to search for charge -> ... -> dest location
// - Both search nearby charge stations and calculate weight between stations are heavy
//   operations, so put them into go-routine and use waitgroup to guarantee result channel
//   is closed after everything is done.
// - CalcWeightBetweenChargeStationsPair needs two iterators, one for nearbystationiterator
//   represents from location and one for next location.  An array of channel is created
//   to represent whether specific iterator is ready or not.
// - The result of this function is channel of stationfindertype.WeightBetweenNeighbors, the sequence of
//   stationfindertype.WeightBetweenNeighbors is important for future logic: first result is start -> first
//   group of low energy charge stations, first group -> second group, ..., xxx group to
//   end
// - All iterators has been recorded in iterators array
//   @Todo: isIteratorReady could be removed later.  When iterator is not ready, should
//         pause inside iterator itself.  That need refactor the design of stationfinder.
func CalculateWeightBetweenNeighbors(locations []*nav.Location, oc *osrmconnector.OSRMConnector, finder stationfinder.StationFinder) chan stationfindertype.WeightBetweenNeighbors {
	c := make(chan stationfindertype.WeightBetweenNeighbors)

	if len(locations) > 2 {
		iterators := make([]stationfindertype.NearbyStationsIterator, len(locations))
		isIteratorReady := make([]chan bool, len(locations))
		for i := range isIteratorReady {
			isIteratorReady[i] = make(chan bool)
		}
		var wg sync.WaitGroup

		for i := 0; i < len(locations); i++ {
			if i == 0 {
				wg.Add(1)
				go func(first int) {
					iterators[first] = NewOrigIter(locations[first])
					isIteratorReady[first] <- true
					wg.Done()
					glog.Info("Finish generating NewOrigIter")
				}(i)
				continue
			}

			if i == len(locations)-1 {
				wg.Add(1)
				go func(last int) {
					iterators[last] = NewDestIter(locations[last])
					glog.Info("Finish generating NewDestIter")
					<-isIteratorReady[last-1]
					putWeightBetweenChargeStationsIntoChannel(iterators[last-1], iterators[last], c, oc)
					glog.Infof("Finish generating putWeightBetweenChargeStationsIntoChannel for %d", last)
					wg.Done()
				}(i)

				break
			}

			wg.Add(1)
			go func(index int) {
				iterators[index] = finder.NewLowEnergyLocationStationFinder(locations[index])
				glog.Infof("Finish generating NewLowEnergyLocationStationFinder for %d", index)
				<-isIteratorReady[index-1]
				isIteratorReady[index] <- true
				putWeightBetweenChargeStationsIntoChannel(iterators[index-1], iterators[index], c, oc)
				glog.Infof("Finish generating putWeightBetweenChargeStationsIntoChannel for %d", index)
				wg.Done()
			}(i)
		}

		go func(wg *sync.WaitGroup) {
			wg.Wait()
			glog.Info("Finish all tasks in CalculateWeightBetweenNeighbors")
			close(c)
			for _, cI := range isIteratorReady {
				close(cI)
			}
		}(&wg)
	}

	return c
}

func putWeightBetweenChargeStationsIntoChannel(from stationfindertype.NearbyStationsIterator, to stationfindertype.NearbyStationsIterator, c chan stationfindertype.WeightBetweenNeighbors, oc *osrmconnector.OSRMConnector) {
	r, err := CalcWeightBetweenChargeStationsPair(from, to, oc)
	if err != nil {
		glog.Errorf("CalculateWeightBetweenNeighbors failed with error %v", err)
	}
	result := stationfindertype.WeightBetweenNeighbors{
		NeighborsInfo: r,
		Err:           err,
	}
	c <- result
}
