package stationgraph

import (
	"reflect"
	"testing"

	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/entity"
	"github.com/Telenav/osrm-backend/integration/service/oasis/place/stationfinder/stationfindertype"
)

func TestEdgeID2EdgeDataGetResultShouldEqualsToSetValue(t *testing.T) {
	cases := []struct {
		key   place2placeID
		value *entity.Weight
	}{
		// case 1
		{
			place2placeID{
				111,
				222,
			},
			&entity.Weight{
				Duration: 111.0,
				Distance: 222.0,
			},
		},
		// case 2
		{
			place2placeID{
				stationfindertype.OrigLocationID,
				stationfindertype.DestLocationID,
			},
			&entity.Weight{
				Duration: 333.0,
				Distance: 444.0,
			},
		},
	}

	edgeID2Data := newEdgeID2EdgeData()
	for _, c := range cases {
		edgeID2Data.add(c.key, c.value)
	}

	for _, c := range cases {
		actualValue := edgeID2Data.get(c.key)
		if !reflect.DeepEqual(c.value, actualValue) {
			t.Errorf("Expect to get value %#v for key %#v but got %#v\n", c.value, c.key, actualValue)
		}
	}
}

func TestEdgeID2EdgeDataWhenGetUnsettedKeyShouldGetNil(t *testing.T) {
	edgeID2Data := newEdgeID2EdgeData()
	actualValue := edgeID2Data.get(place2placeID{
		stationfindertype.OrigLocationID,
		stationfindertype.DestLocationID})
	if actualValue != nil {
		t.Errorf("Expect to get nil for unsetted key but got %#v\n", actualValue)
	}
}
