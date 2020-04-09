package way2nodescsv

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseRecord(t *testing.T) {
	cases := []struct {
		record  []string
		wayID   int64
		nodeIDs []int64
	}{
		{strings.Split("24418325,84760891102,19496208102", ","), 24418325, []int64{84760891102, 19496208102}},
		{strings.Split("24418325,84760891102,19496208102,", ","), 24418325, []int64{84760891102, 19496208102}},
		{strings.Split("24418325,84760891102,19496208102,,,,,", ","), 24418325, []int64{84760891102, 19496208102}},
	}

	for _, c := range cases {
		wayID, nodeIDs, err := ParseRecord(c.record)
		if err != nil {
			t.Error(err)
		}
		if wayID != c.wayID || !reflect.DeepEqual(nodeIDs, c.nodeIDs) {
			t.Errorf("ParseRecord %v, expect %d,%v, but got %d,%v", c.record, c.wayID, c.nodeIDs, wayID, nodeIDs)
		}
	}

	expectFailCases := []struct {
		record []string
	}{
		{strings.Split("", ",")},
		{strings.Split("24418325,84760891102", ",")},
	}
	for _, c := range expectFailCases {
		wayID, nodeIDs, err := ParseRecord(c.record)
		if err == nil {
			t.Errorf("ParseRecord %v, expect fail but got %d,%v", c.record, wayID, nodeIDs)
		}
	}

}

func TestParseLine(t *testing.T) {

	cases := []struct {
		line    string
		wayID   int64
		nodeIDs []int64
	}{
		{"24418325,84760891102,19496208102", 24418325, []int64{84760891102, 19496208102}},
		{"24418325,84760891102,19496208102,", 24418325, []int64{84760891102, 19496208102}},
		{"24418325,84760891102,19496208102,,,,,", 24418325, []int64{84760891102, 19496208102}},
	}

	for _, c := range cases {
		wayID, nodeIDs, err := ParseLine(c.line)
		if err != nil {
			t.Error(err)
		}
		if wayID != c.wayID || !reflect.DeepEqual(nodeIDs, c.nodeIDs) {
			t.Errorf("ParseLine %s, expect %d,%v, but got %d,%v", c.line, c.wayID, c.nodeIDs, wayID, nodeIDs)
		}
	}

	expectFailCases := []struct {
		line string
	}{
		{""},
		{"24418325,84760891102"},
	}
	for _, c := range expectFailCases {
		wayID, nodeIDs, err := ParseLine(c.line)
		if err == nil {
			t.Errorf("ParseLine %s, expect fail but got %d,%v", c.line, wayID, nodeIDs)
		}
	}
}

func TestFormat(t *testing.T) {
	cases := []struct {
		wayID   int64
		nodeIDs []int64
		line    string
	}{
		{24418325, []int64{84760891102, 19496208102}, "24418325,84760891102,19496208102"},
		{24418325, []int64{84760891102, 19496208102, 12345}, "24418325,84760891102,19496208102,12345"},
	}

	for _, c := range cases {
		r := FormatToRecord(c.wayID, c.nodeIDs)
		expectRecord := strings.Split(c.line, ",")
		if !reflect.DeepEqual(r, expectRecord) {
			t.Errorf("FormatToRecord %d,%v expect %v but got %v", c.wayID, c.nodeIDs, expectRecord, r)
		}

		line := FormatToString(c.wayID, c.nodeIDs)
		if line != c.line {
			t.Errorf("FormatToString %d,%v expect %s but got %s", c.wayID, c.nodeIDs, c.line, line)
		}
	}
}
