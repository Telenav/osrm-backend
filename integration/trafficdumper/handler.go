package trafficdumper

// Handler holds flags and methods for dumping traffic data.
type Handler struct {
	blockingOnly       bool
	writeToStdout      bool
	writeToFile        bool
	dumpFileNamePrefix string
	humanFriendlyCSV   bool
	dumpFileSplitIndex int
}

// NewHandler creates a new Handler with command-line flags.
func NewHandler() Handler {
	h := Handler{}
	h.blockingOnly = flags.blockingOnly
	h.writeToStdout = flags.stdout
	if len(flags.dumpFile) > 0 {
		h.writeToFile = true
		h.dumpFileNamePrefix = flags.dumpFile
	} else {
		h.writeToFile = false
	}
	h.humanFriendlyCSV = flags.humanFriendlyCSV
	return h
}
