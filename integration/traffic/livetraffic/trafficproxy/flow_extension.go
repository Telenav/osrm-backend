package trafficproxy

import (
	"fmt"
	"sort"
	"strings"
)

// IsBlocking tests whether the Flow is blocking or not.
//   This function extends protoc-gen-go generated code on testing whether is blocking for Flow.
func (f *Flow) IsBlocking() bool {

	return f.TrafficLevel == TrafficLevel_CLOSED
}

// CSVString represents Flow as defined CSV format.
// I.e. 'wayID,Speed,TrafficLevel,Timestamp,speed,level,begin,end'
func (f *Flow) CSVString() string {
	var result string

	flowSegmentsString := getSegmentsString(f.SegmentedFlow, false)

	if len(flowSegmentsString) > 0 {
		result = fmt.Sprintf("%d,%f,%d,%d,%s", f.WayID, f.Speed, f.TrafficLevel, f.Timestamp, flowSegmentsString)
	} else {
		result = fmt.Sprintf("%d,%f,%d,%d", f.WayID, f.Speed, f.TrafficLevel, f.Timestamp)
	}

	return result
}

// HumanFriendlyCSVString represents Flow as defined CSV format, but prefer human friendly string instead of integer.
// I.e. 'wayID,Speed,TrafficLevel,Timestamp'
func (f *Flow) HumanFriendlyCSVString() string {
	var result string

	flowSegmentsString := getSegmentsString(f.SegmentedFlow, true)

	if len(flowSegmentsString) > 0 {
		result = fmt.Sprintf("%d,%f,%s,%d,%s", f.WayID, f.Speed, f.TrafficLevel, f.Timestamp, flowSegmentsString)
	} else {
		result = fmt.Sprintf("%d,%f,%s,%d", f.WayID, f.Speed, f.TrafficLevel, f.Timestamp)

	}

	return result
}

func getSegmentsString(segments []*SegmentedFlow, humanFmt bool) string {
	var format string

	if humanFmt {
		format = "%f,%s,%d,%d"
	} else {
		format = "%f,%d,%d,%d"
	}

	segmentsString := []string{}

	sort.Slice(segments, func(i, j int) bool { return segments[i].Begin < segments[j].Begin })

	for _, segment := range segments {
		segmentsString = append(segmentsString, fmt.Sprintf(format, segment.Speed, segment.TrafficLevel, segment.Begin, segment.End))
	}

	return strings.Join(segmentsString, ",")
}
