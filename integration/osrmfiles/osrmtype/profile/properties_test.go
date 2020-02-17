package profile

import (
	"io/ioutil"
	"reflect"
	"testing"
)

func TestWriteProperties(t *testing.T) {

	propertiesBytes, err := ioutil.ReadFile("testdata/properties")
	if err != nil {
		t.Error(err)
	}

	cases := []struct {
		p  []byte
		pr Properties
	}{
		{
			propertiesBytes,
			Properties{
				TrafficSignalPenalty:   0,
				UTurnPenalty:           200,
				MaxSpeedForMapMatching: 50,

				ContinueStraightAtWaypoint: true,
				UseTurnRestrictions:        true,
				LeftHandDriving:            false,
				FallbackToDuration:         false,

				WeightName:        "routability",
				ClassNames:        []string{"motorway", "restricted", "tunnel", "toll", "ferry", "", "", ""},
				ExcludableClasses: []uint8{0, 8, 1, 16, 255, 255, 255, 255},

				WeightPrecision:         1,
				ForceSplitEdges:         false,
				CallTaglessNodeFunction: true,
			},
		},
	}

	for _, c := range cases {
		var pr Properties
		n, err := pr.Write(c.p)
		if err != nil {
			t.Error(err)
		}
		if n != len(c.p) {
			t.Errorf("expect to write %d bytes but actually wrote %d bytes", len(c.p), n)
		}
		if !reflect.DeepEqual(pr, c.pr) {
			t.Errorf("parse %v, expect %#v but got %#v", c.p, c.pr, pr)
		}
	}
}
