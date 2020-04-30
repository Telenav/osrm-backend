package oasis

import (
	"testing"

	"github.com/Telenav/osrm-backend/integration/api/osrm/table"
)

func TestRankingSingleChargeStation(t *testing.T) {
	index, err := rankingSingleChargeStation(&table.Mock1ToNTableResponse1, &table.MockNTo1TableResponse1)
	if err != nil || index != 1 {
		t.Errorf("expect %v but got %v", 1, index)
	}
}
