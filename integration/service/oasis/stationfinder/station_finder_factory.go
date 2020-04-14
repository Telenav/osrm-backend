package stationfinder

import (
	"fmt"

	"github.com/Telenav/osrm-backend/integration/service/oasis/searchconnector"
	"github.com/Telenav/osrm-backend/integration/service/oasis/spatialindexer/s2indexer"
	"github.com/Telenav/osrm-backend/integration/service/oasis/stationfinder/cloudfinder"
	"github.com/golang/glog"
)

const (
	// CloudFinder is powered by Telenav Search's web services.
	CloudFinder = "CloudFinder"

	// LocalFinder is supported by pre-processed spatial index(such as google:s2) which is recorded on local.
	LocalFinder = "LocalFinder"
)

// CreateStationsFinder creates finder which implements StationFinder interface
func CreateStationsFinder(finderType, searchEndpoint, apiKey, apiSignature, dataFolderPath string) (StationFinder, error) {
	if err := checkInput(finderType, searchEndpoint, apiKey, apiSignature, dataFolderPath); err != nil {
		return nil, err
	}

	switch finderType {

	case CloudFinder:
		searchFinder := searchconnector.NewTNSearchConnector(searchEndpoint, apiKey, apiSignature)
		return cloudfinder.New(searchFinder), nil

	case LocalFinder:
		localIndex := s2indexer.NewS2Indexer().Load(dataFolderPath)
		if localIndex == nil {
			err := fmt.Errorf("failed to load s2Indexer")
			return nil, err
		}
	}

	return nil, nil
}

// isValidStationFinderType returns false if finderType is unsupported, otherwise returns true
func isValidStationFinderType(finderType string) bool {
	return finderType == CloudFinder || finderType == LocalFinder
}

func checkInput(finderType, searchEndpoint, apiKey, apiSignature, dataFolderPath string) error {
	if !isValidStationFinderType(finderType) {
		glog.Error("Try to create finder not implemented yet, can only choose CloudFinder or LocalFinder for now.\n")
		err := fmt.Errorf("invalid station finder type")
		return err
	}

	if finderType == CloudFinder &&
		(len(searchEndpoint) == 0 ||
			len(apiKey) == 0 ||
			len(apiSignature) == 0) {
		err := fmt.Errorf("empty input for searchEndpoint/apiKey/apiSignature")
		return err
	}

	if finderType == LocalFinder &&
		len(dataFolderPath) == 0 {
		err := fmt.Errorf("empty input for local index")
		return err
	}

	return nil
}
