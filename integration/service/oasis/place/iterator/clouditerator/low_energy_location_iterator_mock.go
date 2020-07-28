package clouditerator

import (
	"sync"

	"github.com/Telenav/osrm-backend/integration/api/search/nearbychargestation"
)

// createMockLowEnergyLocationIterator1 creates mock low energy location station iterator with nearbychargestation.MockSearchResponse1
func createMockLowEnergyLocationIterator1() *lowEnergyLocationIterator {
	obj := &lowEnergyLocationIterator{
		nil,
		&basicIterator{
			tnSearchConnector: nil,

			searchResp:     nearbychargestation.MockSearchResponse1,
			searchRespLock: &sync.RWMutex{},
		},
	}
	return obj
}

// createMockLowEnergyLocationIterator2 creates mock low energy location station iterator with nearbychargestation.MockSearchResponse2
func createMockLowEnergyLocationIterator2() *lowEnergyLocationIterator {
	obj := &lowEnergyLocationIterator{
		nil,
		&basicIterator{
			tnSearchConnector: nil,
			searchResp:        nearbychargestation.MockSearchResponse2,
			searchRespLock:    &sync.RWMutex{},
		},
	}
	return obj
}

// createMockLowEnergyLocationIterator3 creates mock low energy location station iterator with nearbychargestation.MockSearchResponse3
func createMockLowEnergyLocationIterator3() *lowEnergyLocationIterator {
	obj := &lowEnergyLocationIterator{
		nil,
		&basicIterator{
			tnSearchConnector: nil,
			searchResp:        nearbychargestation.MockSearchResponse3,
			searchRespLock:    &sync.RWMutex{},
		},
	}
	return obj
}
