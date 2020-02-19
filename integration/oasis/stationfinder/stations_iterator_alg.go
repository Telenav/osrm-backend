package stationfinder

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

func buildChargeStationInfoDict(iter nearbyStationsIterator) map[string]bool {
	dict := make(map[string]bool)
	c := iter.iterateNearbyStations()
	for item := range c {
		dict[item.ID] = true
	}

	return dict
}
