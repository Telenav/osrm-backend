package cloudfinder

import (
	"sync"

	"github.com/Telenav/osrm-backend/integration/pkg/api/search/nearbychargestation"
)

// createMockLowEnergyLocationStationFinder1 creates mock low energy location station finder with nearbychargestation.MockSearchResponse1
func createMockLowEnergyLocationStationFinder1() *lowEnergyLocationStationFinder {
	obj := &lowEnergyLocationStationFinder{
		nil,
		&basicFinder{
			tnSearchConnector: nil,

			searchResp:     nearbychargestation.MockSearchResponse1,
			searchRespLock: &sync.RWMutex{},
		},
	}
	return obj
}

// createMockLowEnergyLocationStationFinder2 creates mock low energy location station finder with nearbychargestation.MockSearchResponse2
func createMockLowEnergyLocationStationFinder2() *lowEnergyLocationStationFinder {
	obj := &lowEnergyLocationStationFinder{
		nil,
		&basicFinder{
			tnSearchConnector: nil,
			searchResp:        nearbychargestation.MockSearchResponse2,
			searchRespLock:    &sync.RWMutex{},
		},
	}
	return obj
}

// createMockLowEnergyLocationStationFinder3 creates mock low energy location station finder with nearbychargestation.MockSearchResponse3
func createMockLowEnergyLocationStationFinder3() *lowEnergyLocationStationFinder {
	obj := &lowEnergyLocationStationFinder{
		nil,
		&basicFinder{
			tnSearchConnector: nil,
			searchResp:        nearbychargestation.MockSearchResponse3,
			searchRespLock:    &sync.RWMutex{},
		},
	}
	return obj
}
