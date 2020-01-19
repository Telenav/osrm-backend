package trafficproxy

import "fmt"

// IsBlocking tests whether the Flow is blocking or not.
//   This function extends protoc-gen-go generated code on testing whether is blocking for Flow.
func (f *Flow) IsBlocking() bool {

	return f.TrafficLevel == TrafficLevel_CLOSED
}

// CSVString represents Flow as defined CSV format.
// I.e. 'wayID,Speed,TrafficLevel,Timestamp'
func (f *Flow) CSVString() string {
	return fmt.Sprintf("%d,%f,%d,%d", f.WayID, f.Speed, f.TrafficLevel, f.Timestamp)
}

// HumanFriendlyCSVString represents Flow as defined CSV format, but prefer human friendly string instead of integer.
// I.e. 'wayID,Speed,TrafficLevel,Timestamp'
func (f *Flow) HumanFriendlyCSVString() string {
	return fmt.Sprintf("%d,%f,%s,%d", f.WayID, f.Speed, f.TrafficLevel, f.Timestamp)
}
