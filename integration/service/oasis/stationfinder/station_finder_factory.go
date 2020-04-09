package stationfinder

import (
	"fmt"

	"github.com/Telenav/osrm-backend/integration/service/oasis/searchconnector"
	"github.com/Telenav/osrm-backend/integration/service/oasis/stationfinder/tnsearchfinder"
	"github.com/golang/glog"
)

const (
	// TNSearchFinder is powered by Telenav Search's web services.
	TNSearchFinder = "TNSearchFinder"

	// LocalIndexerFinder is supported by pre-processed spatial index(such as google:s2) which is recorded on local.
	LocalIndexerFinder = "LocalIndexerFinder"
)

// CreateStationsFinder creates finder which implements StationFinder interface
func CreateStationsFinder(finderType, searchEndpoint, apiKey, apiSignature string) (StationFinder, error) {
	if err := checkInput(finderType, searchEndpoint, apiKey, apiSignature); err != nil {
		return nil, err
	}

	switch finderType {
	case TNSearchFinder:
		searchFinder := searchconnector.NewTNSearchConnector(searchEndpoint, apiKey, apiSignature)
		return tnsearchfinder.NewTnSearchStationFinder(searchFinder), nil
	}
	return nil, nil
}

// isValidStationFinderType returns false if finderType is unsupported, otherwise returns true
func isValidStationFinderType(finderType string) bool {
	return finderType == TNSearchFinder || finderType == LocalIndexerFinder
}

func checkInput(finderType, searchEndpoint, apiKey, apiSignature string) error {
	if !isValidStationFinderType(finderType) {
		glog.Error("Try to create finder not implemented yet, can only choose TNSearchFinder or LocalFinder for now.\n")
		err := fmt.Errorf("invalid station finder type")
		return err
	}

	if finderType == TNSearchFinder &&
		(len(searchEndpoint) == 0 ||
			len(apiKey) == 0 ||
			len(apiSignature) == 0) {
		err := fmt.Errorf("empty input for searchEndpoint/apiKey/apiSignature")
		return err
	}

	return nil
}
