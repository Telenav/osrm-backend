package clouditerator

import (
	"sync"

	"github.com/Telenav/osrm-backend/integration/api/search/nearbychargestation"
)

// CreateMockOrigIterator1 creates mock orig station finder with nearbychargestation.MockSearchResponse1
func CreateMockOrigIterator1() *origIterator {
	obj := &origIterator{
		nil,
		&basicIterator{
			tnSearchConnector: nil,
			searchResp:        nearbychargestation.MockSearchResponse1,
			searchRespLock:    &sync.RWMutex{},
		},
	}
	return obj
}

// CreateMockOrigIterator2 creates mock orig station finder with nearbychargestation.MockSearchResponse2
func CreateMockOrigIterator2() *origIterator {
	obj := &origIterator{
		nil,
		&basicIterator{
			tnSearchConnector: nil,
			searchResp:        nearbychargestation.MockSearchResponse2,
			searchRespLock:    &sync.RWMutex{},
		},
	}
	return obj
}

// CreateMockOrigIterator3 creates mock orig station finder with nearbychargestation.MockSearchResponse3
func CreateMockOrigIterator3() *origIterator {
	obj := &origIterator{
		nil,
		&basicIterator{
			tnSearchConnector: nil,
			searchResp:        nearbychargestation.MockSearchResponse3,
			searchRespLock:    &sync.RWMutex{},
		},
	}
	return obj
}
