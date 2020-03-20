package historicalspeed

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseWaysMappingRecordFailure(t *testing.T) {
	cases := [][]string{
		strings.Split("LINK_PVID,TRAVEL_DIRECTION,U,M,T,W,R,F,S", ","),
		strings.Split("737019219,F,10788,14140,3561,14978,12324,2202,", ","),             // too less field
		strings.Split("737019219,F,10788,14140,3561,14978,12324,2202,2243,1", ","),       // too many fields
		strings.Split("a,F,10788,14140,3561,14978,12324,2202,2243", ","),                 // char instead of integer on wayID
		strings.Split("737019219,F,a,14140,3561,14978,12324,2202,2243", ","),             // char instead of integer on patternID
		strings.Split("-737019219,F,10788,14140,3561,14978,12324,2202,2243", ","),        // negative wayID
		strings.Split("737019219,F,-10788,14140,3561,14978,12324,2202,2243", ","),        // negative patternID
		strings.Split("737019219,F,10788000000,14140,3561,14978,12324,2202,2243", ","),   // patternID overflow
		strings.Split("737019219,F,10788,14140,3561,14978,12324,2202,2243,-130,0", ","),  // invalid timezone
		strings.Split("737019219,F,10788,14140,3561,14978,12324,2202,2243,-120,-1", ","), // invalid daylight saving
	}

	for _, c := range cases {
		_, _, err := parseWaysMappingRecord(c)
		if err == nil {
			t.Errorf("expect parse %v failed but got succeed", c)
		}
	}
}

func TestParseWaysMappingRecordSucceed(t *testing.T) {
	cases := []struct {
		record []string
		wayID  int64
		*mappingItem
	}{
		{
			strings.Split("737019219,F,10788,14140,3561,14978,12324,2202,2243", ","),
			-737019219,
			&mappingItem{
				[daysPerWeek]uint32{10788, 14140, 3561, 14978, 12324, 2202, 2243},
				0,
				0,
			},
		},
		{
			strings.Split("737019219,T,39,39,39,39,39,14949,39", ","),
			737019219,
			&mappingItem{
				[daysPerWeek]uint32{39, 39, 39, 39, 39, 14949, 39},
				0,
				0,
			},
		},
		{
			strings.Split("737019219,F,10788,14140,3561,14978,12324,2202,2243,-120,5", ","),
			-737019219,
			&mappingItem{
				[daysPerWeek]uint32{10788, 14140, 3561, 14978, 12324, 2202, 2243},
				-120,
				5,
			},
		},
		{
			strings.Split("737019219,F,10788,14140,3561,14978,12324,2202,2243,0,67", ","),
			-737019219,
			&mappingItem{
				[daysPerWeek]uint32{10788, 14140, 3561, 14978, 12324, 2202, 2243},
				0,
				67,
			},
		},
	}

	for _, c := range cases {
		wayID, mapping, err := parseWaysMappingRecord(c.record)
		if err != nil {
			t.Errorf("expect parse %v succeed but got err %v", c.record, err)
		}
		if wayID != c.wayID {
			t.Errorf("parse %v, expect wayID %d but got %d", c.record, c.wayID, wayID)
		}
		if !reflect.DeepEqual(*mapping, *c.mappingItem) {
			t.Errorf("parse %v, expect mapping %v but got %v", c.record, c.mappingItem, mapping)
		}
	}
}

func TestToWaysMappingRecord(t *testing.T) {
	cases := []struct {
		wayID int64
		mappingItem
		record []string
	}{
		{
			737019219,
			mappingItem{
				[daysPerWeek]uint32{39, 39, 39, 39, 39, 14949, 39}, 0, 0,
			},
			[]string{
				"737019219", "T", "39", "39", "39", "39", "39", "14949", "39", "000", "0",
			},
		},
		{
			737019219,
			mappingItem{
				[daysPerWeek]uint32{39, 39, 39, 39, 39, 14949, 39}, 80, 67,
			},
			[]string{
				"737019219", "T", "39", "39", "39", "39", "39", "14949", "39", "080", "67",
			},
		},
		{
			-737019219,
			mappingItem{
				[daysPerWeek]uint32{39, 39, 39, 39, 39, 14949, 39}, -70, 2,
			},
			[]string{
				"737019219", "F", "39", "39", "39", "39", "39", "14949", "39", "-070", "2",
			},
		},
	}

	for _, c := range cases {
		record := toWaysMappingRecord(c.wayID, &c.mappingItem)
		if !reflect.DeepEqual(record, c.record) {
			t.Errorf("wayID %d mappingItem %v to record, expect %v but got %v", c.wayID, c.mappingItem, c.record, record)
		}
	}
}
