// Package appversion unifies common versions for all binaries.
package appversion

import (
	"fmt"
	"io"
	"os"
)

// these information will be collected when build, by `-ldflags "-X appversion.appVersion=0.1"`
// They have to be global variables in package to accept the data by `-ldflags`.
var (
	appVersion string
	buildTime  string
	gitCommit  string
)

// VersionInfo represents version information structure.
type VersionInfo struct {
	AppVersion string `json:"app version"`
	BuildTime  string `json:"build time"`
	GitCommit  string `json:"git commit"`
}

// Version returns version information.
func Version() VersionInfo {
	return VersionInfo{
		appVersion, buildTime, gitCommit,
	}
}

// Fprint prints version information to writer.
func Fprint(w io.Writer) {
	fmt.Fprintf(w, "Version:    %s\n", appVersion)
	fmt.Fprintf(w, "Build Time: %s\n", buildTime)
	fmt.Fprintf(w, "Git Commit: %s\n", gitCommit)
}

// Print prints version information to stdout.
func Print() {
	Fprint(os.Stdout)
}

// PrintExit prints version to stdout and os.Exit(2) if have `-version` flag.
// Call it after `flag.Parse()`.
func PrintExit() {
	if VersionFlag() {
		Print()
		os.Exit(2)
	}
}
