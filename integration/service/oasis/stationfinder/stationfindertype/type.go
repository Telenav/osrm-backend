package stationfindertype

import (
	"regexp"
	"strconv"

	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/common"
	"github.com/golang/glog"
)

// NearbyStationsIterator provide interator for near by stations
type NearbyStationsIterator interface {
	// IterateNearbyStations returns a channel which contains near by charge station under certain conditions
	IterateNearbyStations() <-chan *ChargeStationInfo
}

// ChargeStationInfo defines charge station information
type ChargeStationInfo struct {
	ID       string
	Location nav.Location
	Err      error
}

// Weight represent weight information
type Weight struct {
	Duration float64
	Distance float64
}

// NeighborInfo represent cost information between two charge stations
type NeighborInfo struct {
	FromID       string
	FromLocation nav.Location
	ToID         string
	ToLocation   nav.Location
	Weight
}

// WeightBetweenNeighbors contains a group of neighbors information
type WeightBetweenNeighbors struct {
	NeighborsInfo []NeighborInfo
	Err           error
}

// FromPlaceID converts from String to PlaceID for From
func (neighbor NeighborInfo) FromPlaceID() common.PlaceID {
	return convertIDFromStringToPlaceID(neighbor.FromID)
}

// ToPlaceID converts from String to PlaceID for To
func (neighbor NeighborInfo) ToPlaceID() common.PlaceID {
	return convertIDFromStringToPlaceID(neighbor.ToID)
}

// Telenav web response format b-12345678
// While pre-processed result is value itself
func convertIDFromStringToPlaceID(s string) common.PlaceID {
	switch s {
	case OrigLocationIDStr:
		return OrigLocationID
	case DestLocationIDStr:
		return DestLocationID
	case InvalidPlaceIDStr:
		return InvalidPlaceID
	default:
		num, err := strconv.ParseInt(s, 10, 64)
		if err == nil { // pre-processed result
			placeID := (common.PlaceID)(num)
			if !isPredefinedValueTakenPlaceIDValue(placeID) {
				return placeID
			} // else: Fatal assert inside isPredefinedValueTakenPlaceIDValue
		} else { // from Web Service
			re := regexp.MustCompile("[0-9]+")
			idStrArray := re.FindAllString(s, -1)
			if len(idStrArray) != 1 {
				glog.Fatalf("Assumption of ID from search response is wrong, expect pattern is b- with single group of number, but got %v.\n", s)
			}
			if num, err := strconv.ParseInt(idStrArray[0], 10, 64); err == nil {
				placeID := (common.PlaceID)(num)
				if !isPredefinedValueTakenPlaceIDValue(placeID) {
					return placeID
				} // else: Fatal assert inside isPredefinedValueTakenPlaceIDValue
			} else {
				glog.Fatalf("Assumption of ID from search response is wrong, expect format is b- with number, but when converting %v got error %#v.\n", s, err)
			}
		}
	}
	glog.Fatalf("PlaceID %v could not be decodinged\n", s)
	return InvalidPlaceID
}

func isPredefinedValueTakenPlaceIDValue(id common.PlaceID) bool {
	if id == OrigLocationID || id == DestLocationID || id == InvalidPlaceID {
		glog.Fatal("Predefined ID use the same value as value from PlaceID, please adjust either one of them.\n")
		return true
	}
	return false
}
