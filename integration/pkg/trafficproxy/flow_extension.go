package trafficproxy

import "fmt"

const blockingSpeedThreshold = 1 // Think it's blocking if flow speed smaller than this threshold.

// IsBlocking tests whether the Flow is blocking or not.
//   This function extends protoc-gen-go generated code on testing whether is blocking for Flow.
func (f *Flow) IsBlocking() bool {

	return f.TrafficLevel == TrafficLevel_CLOSED || f.Speed < blockingSpeedThreshold
}

// CSVString represents Flow as defined CSV format.
// I.e. 'wayID,Speed,TrafficLevel'
func (f *Flow) CSVString() string {
	return fmt.Sprintf("%d,%f,%d", f.WayID, f.Speed, f.TrafficLevel)
}

// HumanFriendlyCSVString represents Flow as defined CSV format, but prefer human friendly string instead of integer.
// I.e. 'wayID,Speed,TrafficLevel'
func (f *Flow) HumanFriendlyCSVString() string {
	return fmt.Sprintf("%d,%f,%s", f.WayID, f.Speed, f.TrafficLevel)
}
