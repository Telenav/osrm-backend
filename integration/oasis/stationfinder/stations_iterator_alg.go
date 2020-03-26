package stationfinder

import (
	"fmt"
	"sync"

	"github.com/Telenav/osrm-backend/integration/oasis/osrmconnector"
	"github.com/Telenav/osrm-backend/integration/oasis/osrmhelper"
	"github.com/Telenav/osrm-backend/integration/oasis/searchconnector"
	"github.com/Telenav/osrm-backend/integration/pkg/api/nav"
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
	Location nav.Location
	err      error
}

// // nav.Location represents location information
// type nav.Location nav.Location

// CalcWeightBetweenChargeStationsPair accepts two iterators and calculates weights between each pair of iterators
func CalcWeightBetweenChargeStationsPair(from nearbyStationsIterator, to nearbyStationsIterator, table osrmconnector.TableRequster) ([]NeighborInfo, error) {
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
		err := fmt.Errorf("empty iterator of from pass into CalcWeightBetweenChargeStationsPair")
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
	var result []NeighborInfo
	for i, startPoint := range startPoints {
		for j, targetPoint := range targetPoints {
			result = append(result, NeighborInfo{
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

// NeighborInfo represent cost information between two charge stations
type NeighborInfo struct {
	FromID       string
	FromLocation nav.Location
	ToID         string
	ToLocation   nav.Location
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

type WeightBetweenNeighbors struct {
	NeighborsInfo []NeighborInfo
	Err           error
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
// - The result of this function is channel of WeightBetweenNeighbors, the sequence of
//   WeightBetweenNeighbors is important for future logic: first result is start -> first
//   group of low energy charge stations, first group -> second group, ..., xxx group to
//   end
// - All iterators has been recorded in iterators array
//   @Todo: isIteratorReady could be removed later.  When iterator is not ready, should
//         pause inside iterator itself.  That need refactor the design of stationfinder.
func CalculateWeightBetweenNeighbors(locations []*nav.Location, oc *osrmconnector.OSRMConnector, sc *searchconnector.TNSearchConnector) chan WeightBetweenNeighbors {
	c := make(chan WeightBetweenNeighbors)

	if len(locations) > 2 {
		iterators := make([]nearbyStationsIterator, len(locations))
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
				iterators[index] = NewLowEnergyLocationStationFinder(sc, locations[index])
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

func putWeightBetweenChargeStationsIntoChannel(from nearbyStationsIterator, to nearbyStationsIterator, c chan WeightBetweenNeighbors, oc *osrmconnector.OSRMConnector) {
	r, err := CalcWeightBetweenChargeStationsPair(from, to, oc)
	if err != nil {
		glog.Errorf("CalculateWeightBetweenNeighbors failed with error %v", err)
	}
	result := WeightBetweenNeighbors{
		NeighborsInfo: r,
		Err:           err,
	}
	c <- result
}
