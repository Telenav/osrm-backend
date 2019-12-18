package osrmtype

import (
	"math"
	"reflect"
	"testing"
)

func TestNodeBasedEdgesWrite(t *testing.T) {

	cases := []struct {
		p []byte
		NodeBasedEdges
	}{
		{
			[]byte{
				0x00, 0x00, 0x00, 0x00, 0x84, 0x70, 0x07, 0x00, 0x18, 0x00, 0x00, 0x00, 0x18, 0x00, 0x00, 0x00,
				0x73, 0x25, 0x5b, 0x42, 0xff, 0xff, 0xff, 0x7f, 0xce, 0xad, 0x00, 0x00, 0x21, 0x01, 0x02, 0x00,
			},
			NodeBasedEdges{
				NodeBasedEdge{
					Source:       0,
					Target:       487556,
					Weight:       24,
					Duration:     24,
					Distance:     54.7865715,
					GeometryID:   GeometryID{math.MaxUint32 >> 1, false},
					AnnotationID: 44494,
					Flags: NodeBasedEdgeClassification{
						Forward:    true,
						Backward:   false,
						IsSplit:    false,
						Roundabout: false,
						Circular:   false,
						Startpoint: true,
						Restricted: false,
						RoadClassification: RoadClassification{
							MotowayClass:      true,
							LinkClass:         false,
							MaybeIgnored:      false,
							RoadPriorityClass: RoadPriorityClassMotoway,
							NumberOfLanes:     2,
						},
						HighwayTurnClassification: 0,
						AccessTurnClassification:  0,
					},
				},
			},
		},
		{
			[]byte{
				0x00, 0x00, 0x00, 0x00, 0x84, 0x70, 0x07, 0x00, 0x18, 0x00, 0x00, 0x00, 0x18, 0x00, 0x00, 0x00,
				0x73, 0x25, 0x5b, 0x42, 0xff, 0xff, 0xff, 0x7f, 0xce, 0xad, 0x00, 0x00, 0x21, 0x01, 0x02, 0x00,
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, // redundant bytes
			},
			NodeBasedEdges{
				NodeBasedEdge{
					Source:       0,
					Target:       487556,
					Weight:       24,
					Duration:     24,
					Distance:     54.7865715,
					GeometryID:   GeometryID{math.MaxUint32 >> 1, false},
					AnnotationID: 44494,
					Flags: NodeBasedEdgeClassification{
						Forward:    true,
						Backward:   false,
						IsSplit:    false,
						Roundabout: false,
						Circular:   false,
						Startpoint: true,
						Restricted: false,
						RoadClassification: RoadClassification{
							MotowayClass:      true,
							LinkClass:         false,
							MaybeIgnored:      false,
							RoadPriorityClass: RoadPriorityClassMotoway,
							NumberOfLanes:     2,
						},
						HighwayTurnClassification: 0,
						AccessTurnClassification:  0,
					},
				},
			},
		},
	}

	for _, c := range cases {
		var n NodeBasedEdges
		writeLen, err := n.Write(c.p)
		if err != nil {
			t.Error(err)
		}
		if len(c.p)-writeLen != len(c.p)%nodeBasedEdgeBytes {
			t.Errorf("len(p) %d but write len %d", len(c.p), writeLen)
		}
		if !reflect.DeepEqual(n, c.NodeBasedEdges) {
			t.Errorf("parse %v expect %v but got %v", c.p, c.NodeBasedEdges, n)
		}
	}
}

func TestNodeBasedEdgeAnnotationsWrite(t *testing.T) {
	cases := []struct {
		p []byte
		NodeBasedEdgeAnnotations
	}{
		{
			[]byte{0x05, 0x00, 0x00, 0x00, 0xff, 0xff, 0x01, 0x01},
			NodeBasedEdgeAnnotations{
				NodeBasedEdgeAnnotation{
					NameID:            5,
					LaneDescriptionID: InvalidLaneDescriptionID,
					ClassData:         1,
					TravelMode:        TravelModeDriving,
					IsLeftHandDriving: false,
				},
			},
		},
		{
			[]byte{
				0x05, 0x00, 0x00, 0x00, 0xff, 0xff, 0x01, 0x01,
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, // redundant bytes
			},
			NodeBasedEdgeAnnotations{
				NodeBasedEdgeAnnotation{
					NameID:            5,
					LaneDescriptionID: InvalidLaneDescriptionID,
					ClassData:         1,
					TravelMode:        TravelModeDriving,
					IsLeftHandDriving: false,
				},
			},
		},
	}

	for _, c := range cases {
		var n NodeBasedEdgeAnnotations
		writeLen, err := n.Write(c.p)
		if err != nil {
			t.Error(err)
		}
		if len(c.p)-writeLen != len(c.p)%nodeBasedEdgeAnnotationBytes {
			t.Errorf("len(p) %d but write len %d", len(c.p), writeLen)
		}
		if !reflect.DeepEqual(n, c.NodeBasedEdgeAnnotations) {
			t.Errorf("parse %v expect %v but got %v", c.p, c.NodeBasedEdgeAnnotations, n)
		}
	}

}
