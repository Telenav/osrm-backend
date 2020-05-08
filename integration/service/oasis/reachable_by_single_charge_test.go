package oasis

import (
	"testing"

	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/mock"
)

func TestRankingSingleChargeStation(t *testing.T) {
	index, err := rankingSingleChargeStation(&mock.Mock1To4TableResponse1, &mock.Mock4To1TableResponse1)
	if err != nil || index != 1 {
		t.Errorf("expect %v but got %v", 1, index)
	}
}
