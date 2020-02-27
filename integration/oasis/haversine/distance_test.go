package haversine

import "testing"

func TestGreatCircleDistance(t *testing.T) {
	// expect value got from http://www.onlineconversion.com/map_greatcircle_distance.htm
	expect := 111595.4865288326
	actual := GreatCircleDistance(32.333, 122.323, 31.333, 122.423)
	if !floatEquals(expect, actual) {
		t.Errorf("Expected GreatCircleDistance returns %v, got %v", expect, actual)
	}
}

var epsilon float64 = 0.00000001

func floatEquals(a, b float64) bool {
	if (a-b) < epsilon && (b-a) < epsilon {
		return true
	}
	return false
}
