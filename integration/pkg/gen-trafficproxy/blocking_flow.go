package proxy

const blockingSpeedThreshold = 1 // Think it's blocking if flow speed smaller than this threshold.

// IsBlocking tests whether the Flow is blocking or not.
//   This function extends protoc-gen-go generated code on testing whether is blocking for Flow.
func (f *Flow) IsBlocking() bool {

	return f.TrafficLevel == TrafficLevel_CLOSED || f.Speed < blockingSpeedThreshold
}
