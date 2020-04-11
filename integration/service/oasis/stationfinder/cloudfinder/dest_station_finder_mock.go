package cloudfinder

import (
	"sync"

	"github.com/Telenav/osrm-backend/integration/pkg/api/search/nearbychargestation"
)

// CreateMockDestStationFinder1 creates mock dest station finder with nearbychargestation.MockSearchResponse1
func CreateMockDestStationFinder1() *destStationFinder {
	obj := &destStationFinder{
		oasisReq: nil,
		bf: &basicFinder{
			tnSearchConnector: nil,
			searchResp:        nearbychargestation.MockSearchResponse1,
			searchRespLock:    &sync.RWMutex{},
		},
	}
	return obj
}

// CreateMockDestStationFinder2 creates mock dest station finder with nearbychargestation.MockSearchResponse2
func createMockDestStationFinder2() *destStationFinder {
	obj := &destStationFinder{
		oasisReq: nil,
		bf: &basicFinder{
			tnSearchConnector: nil,
			searchResp:        nearbychargestation.MockSearchResponse2,
			searchRespLock:    &sync.RWMutex{},
		},
	}
	return obj
}

// CreateMockDestStationFinder3 creates mock dest station finder with nearbychargestation.MockSearchResponse3
func createMockDestStationFinder3() *destStationFinder {
	obj := &destStationFinder{
		oasisReq: nil,
		bf: &basicFinder{
			tnSearchConnector: nil,
			searchResp:        nearbychargestation.MockSearchResponse3,
			searchRespLock:    &sync.RWMutex{},
		},
	}
	return obj
}
