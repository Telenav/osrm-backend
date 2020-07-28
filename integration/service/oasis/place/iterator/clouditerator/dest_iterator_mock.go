package clouditerator

import (
	"sync"

	"github.com/Telenav/osrm-backend/integration/api/search/nearbychargestation"
)

// CreateMockDestIterator1 creates mock dest station iterator with nearbychargestation.MockSearchResponse1
func CreateMockDestIterator1() *destIterator {
	obj := &destIterator{
		nil,
		&basicIterator{
			tnSearchConnector: nil,
			searchResp:        nearbychargestation.MockSearchResponse1,
			searchRespLock:    &sync.RWMutex{},
		},
	}
	return obj
}

// CreateMockDestIterator2 creates mock dest station iterator with nearbychargestation.MockSearchResponse2
func CreateMockDestIterator2() *destIterator {
	obj := &destIterator{
		nil,
		&basicIterator{
			tnSearchConnector: nil,
			searchResp:        nearbychargestation.MockSearchResponse2,
			searchRespLock:    &sync.RWMutex{},
		},
	}
	return obj
}

// CreateMockDestIterator3 creates mock dest station iterator with nearbychargestation.MockSearchResponse3
func CreateMockDestIterator3() *destIterator {
	obj := &destIterator{
		nil,
		&basicIterator{
			tnSearchConnector: nil,
			searchResp:        nearbychargestation.MockSearchResponse3,
			searchRespLock:    &sync.RWMutex{},
		},
	}
	return obj
}
