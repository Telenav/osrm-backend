package appversion

import "flag"

var flags struct {
	version bool // print version
}

func init() {
	flag.BoolVar(&flags.version, "version", false, "Print version and exit.")
}

// VersionFlag indicate whether expect to print version and exit.
// Call it after `flag.Prase()` to make sure command line parameters available.
func VersionFlag() bool {
	return flags.version
}
