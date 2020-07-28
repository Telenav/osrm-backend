package resourcemanager

import (
	"fmt"
	"time"

	"github.com/Telenav/osrm-backend/integration/service/oasis/place"
	"github.com/Telenav/osrm-backend/integration/service/oasis/place/iterator"
	"github.com/Telenav/osrm-backend/integration/service/oasis/place/spatialindexer/s2indexer"
	"github.com/Telenav/osrm-backend/integration/service/oasis/place/topograph"
	"github.com/Telenav/osrm-backend/integration/util/osrmconnector"
	"github.com/golang/glog"
)

// ResourceMgr defines strategy be used by charge station selection
type ResourceMgr struct {
	osrmConnector          *osrmconnector.OSRMConnector // osrmConnector represents communication with OSRM backend
	stationFinder          place.IteratorGenerator      // stationFinder generates nearby stations based cloud search or local spatial index
	spatialIndexerFinder   place.Finder                 // spatialIndexerFinder answers spatial query based on pre-generated spatial data
	connectivityMap        *topograph.MemoryTopoGraph   // connectivityMap contains connectivity information for stations
	stationLocationQuerier place.LocationQuerier        // stationLocationQuerier answers location information for specific station
}

// NewResourceMgr creates ResourceMgr object
func NewResourceMgr(osrmBackend, finderType, searchEndpoint, apiKey, apiSignature, dataFolderPath string) (*ResourceMgr, error) {
	// @todo: need make sure connectivity is on and continues available
	//        simple request to guarantee server is alive after init
	startTime := time.Now()

	if len(osrmBackend) == 0 {
		err := fmt.Errorf("empty osrmBackend end point")
		return nil, err
	}

	s2indexer := s2indexer.NewS2Indexer().Load(dataFolderPath)
	if s2indexer == nil {
		err := fmt.Errorf("failed to load s2Indexer")
		return nil, err
	}

	stationFinder, err := iterator.CreateIteratorGenerator(finderType, searchEndpoint, apiKey, apiSignature, s2indexer)
	if err != nil {
		glog.Errorf("Failed to call iterator.CreateIteratorGenerator, met error = %+v\n", err)
		return nil, err
	}

	connectivityMap := topograph.New(0.0).Load(dataFolderPath)
	if connectivityMap == nil {
		err := fmt.Errorf("failed to load MemoryTopoGraph")
		return nil, err
	}

	glog.Infof("Initialize OASIS resource manager takes %f seconds.", time.Since(startTime).Seconds())

	return &ResourceMgr{
		osrmConnector:          osrmconnector.NewOSRMConnector(osrmBackend),
		stationFinder:          stationFinder,
		spatialIndexerFinder:   s2indexer,
		connectivityMap:        connectivityMap,
		stationLocationQuerier: s2indexer,
	}, nil
}

// OSRMConnector gets osrmConnector used for communicating with OSRM backend
func (r *ResourceMgr) OSRMConnector() *osrmconnector.OSRMConnector {
	return r.osrmConnector
}

// IteratorGenerator returns interface of IteratorGenerator used for finding nearby stations
// based cloud search or local spatial index
func (r *ResourceMgr) IteratorGenerator() place.IteratorGenerator {
	return r.stationFinder
}

// SpatialIndexerFinder answers spatial query based on pre-generated spatial data
func (r *ResourceMgr) SpatialIndexerFinder() place.Finder {
	return r.spatialIndexerFinder
}

// MemoryTopoGraph returns connectivity information for stations
func (r *ResourceMgr) MemoryTopoGraph() *topograph.MemoryTopoGraph {
	return r.connectivityMap
}

// StationLocationQuerier answers location information for specific station
func (r *ResourceMgr) StationLocationQuerier() place.LocationQuerier {
	return r.stationLocationQuerier
}
