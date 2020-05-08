package haversine

import (
	"testing"

	"github.com/Telenav/osrm-backend/integration/util"
)

func TestGreatCircleDistance(t *testing.T) {
	// expect value got from http://www.onlineconversion.com/map_greatcircle_distance.htm
	expect := 111595.4865288326
	actual := GreatCircleDistance(32.333, 122.323, 31.333, 122.423)
	if !util.Float64Equal(expect, actual) {
		t.Errorf("Expected GreatCircleDistance returns %v, got %v", expect, actual)
	}
}
