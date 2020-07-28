package stationfindertype

import (
	"testing"

	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/entity"
)

func TestConvertIDFromStringToPlaceID(t *testing.T) {
	cases := []struct {
		idStr         string
		expectPlaceID entity.PlaceID
	}{
		{
			OrigLocationIDStr,
			OrigLocationID,
		},
		{
			DestLocationIDStr,
			DestLocationID,
		},
		{
			InvalidPlaceIDStr,
			InvalidPlaceID,
		},
		{
			"1234567",
			1234567,
		},
		{
			"b-1234567",
			1234567,
		},
		{
			"station1",
			1,
		},
	}

	for _, c := range cases {
		actualPlaceID := convertIDFromStringToPlaceID(c.idStr)
		if actualPlaceID != c.expectPlaceID {
			t.Errorf("For case %v, expect result is %v while actualPlace id is %v\n", c, c.expectPlaceID, actualPlaceID)
		}
	}
}
