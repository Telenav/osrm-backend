package cloudfinder

import (
	"sync"

	"github.com/Telenav/osrm-backend/integration/api/search/nearbychargestation"
)

// CreateMockOrigStationFinder1 creates mock orig station finder with nearbychargestation.MockSearchResponse1
func CreateMockOrigStationFinder1() *origStationFinder {
	obj := &origStationFinder{
		nil,
		&basicFinder{
			tnSearchConnector: nil,
			searchResp:        nearbychargestation.MockSearchResponse1,
			searchRespLock:    &sync.RWMutex{},
		},
	}
	return obj
}

// CreateMockOrigStationFinder2 creates mock orig station finder with nearbychargestation.MockSearchResponse2
func CreateMockOrigStationFinder2() *origStationFinder {
	obj := &origStationFinder{
		nil,
		&basicFinder{
			tnSearchConnector: nil,
			searchResp:        nearbychargestation.MockSearchResponse2,
			searchRespLock:    &sync.RWMutex{},
		},
	}
	return obj
}

// CreateMockOrigStationFinder3 creates mock orig station finder with nearbychargestation.MockSearchResponse3
func CreateMockOrigStationFinder3() *origStationFinder {
	obj := &origStationFinder{
		nil,
		&basicFinder{
			tnSearchConnector: nil,
			searchResp:        nearbychargestation.MockSearchResponse3,
			searchRespLock:    &sync.RWMutex{},
		},
	}
	return obj
}
