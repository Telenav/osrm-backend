package iterator

import (
	"fmt"

	"github.com/Telenav/osrm-backend/integration/service/oasis/place"
	"github.com/Telenav/osrm-backend/integration/service/oasis/place/iterator/clouditerator"
	"github.com/Telenav/osrm-backend/integration/service/oasis/place/iterator/localiterator"
	"github.com/Telenav/osrm-backend/integration/util/searchconnector"
	"github.com/golang/glog"
)

const (
	// CloudFinder is powered by Telenav Search's web services.
	CloudFinder = "CloudFinder"

	// LocalFinder is supported by pre-processed spatial index(such as google:s2) which is recorded on local.
	LocalFinder = "LocalFinder"
)

// CreateIteratorGenerator creates finder which implements IteratorGenerator interface
func CreateIteratorGenerator(finderType, searchEndpoint, apiKey, apiSignature string, finder place.Finder) (place.IteratorGenerator, error) {
	if err := checkInput(finderType, searchEndpoint, apiKey, apiSignature, finder); err != nil {
		return nil, err
	}

	switch finderType {

	case CloudFinder:
		searchFinder := searchconnector.NewTNSearchConnector(searchEndpoint, apiKey, apiSignature)
		return clouditerator.New(searchFinder), nil

	case LocalFinder:
		return localiterator.New(finder), nil
	}

	return nil, nil
}

// isValidStationFinderType returns false if finderType is unsupported, otherwise returns true
func isValidStationFinderType(finderType string) bool {
	return finderType == CloudFinder || finderType == LocalFinder
}

func checkInput(finderType, searchEndpoint, apiKey, apiSignature string, finder place.Finder) error {
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
		finder == nil {
		err := fmt.Errorf("empty input for local index")
		return err
	}

	return nil
}
