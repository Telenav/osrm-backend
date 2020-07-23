package resourcemanager

import (
	"fmt"
	"time"

	"github.com/Telenav/osrm-backend/integration/service/oasis/connectivitymap"
	"github.com/Telenav/osrm-backend/integration/service/oasis/osrmconnector"
	"github.com/Telenav/osrm-backend/integration/service/oasis/spatialindexer"
	"github.com/Telenav/osrm-backend/integration/service/oasis/spatialindexer/s2indexer"
	"github.com/Telenav/osrm-backend/integration/service/oasis/stationfinder"
	"github.com/golang/glog"
)

// ResourceMgr defines strategy be used by charge station selection
type ResourceMgr struct {
	osrmConnector          *osrmconnector.OSRMConnector        // osrmConnector represents communication with OSRM backend
	stationFinder          stationfinder.StationFinder         // stationFinder generates nearby stations based cloud search or local spatial index
	spatialIndexerFinder   spatialindexer.Finder               // spatialIndexerFinder answers spatial query based on pre-generated spatial data
	connectivityMap        *connectivitymap.ConnectivityMap    // connectivityMap contains connectivity information for stations
	stationLocationQuerier spatialindexer.PlaceLocationQuerier // stationLocationQuerier answers location information for specific station
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

	stationFinder, err := stationfinder.CreateStationsFinder(finderType, searchEndpoint, apiKey, apiSignature, s2indexer)
	if err != nil {
		glog.Errorf("Failed to call stationfinder.CreateStationsFinder, met error = %+v\n", err)
		return nil, err
	}

	connectivityMap := connectivitymap.New(0.0).Load(dataFolderPath)
	if connectivityMap == nil {
		err := fmt.Errorf("failed to load ConnectivityMap")
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

// StationFinder returns interface of StationFinder used for finding nearby stations
// based cloud search or local spatial index
func (r *ResourceMgr) StationFinder() stationfinder.StationFinder {
	return r.stationFinder
}

// SpatialIndexerFinder answers spatial query based on pre-generated spatial data
func (r *ResourceMgr) SpatialIndexerFinder() spatialindexer.Finder {
	return r.spatialIndexerFinder
}

// ConnectivityMap returns connectivity information for stations
func (r *ResourceMgr) ConnectivityMap() *connectivitymap.ConnectivityMap {
	return r.connectivityMap
}

// StationLocationQuerier answers location information for specific station
func (r *ResourceMgr) StationLocationQuerier() spatialindexer.PlaceLocationQuerier {
	return r.stationLocationQuerier
}
