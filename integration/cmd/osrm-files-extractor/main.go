package main

import (
	"flag"
	"strings"

	"github.com/Telenav/osrm-backend/integration/osrmfiles"
	"github.com/Telenav/osrm-backend/integration/osrmfiles/dotnbgnodes"
	"github.com/Telenav/osrm-backend/integration/osrmfiles/dotosrm"
	"github.com/Telenav/osrm-backend/integration/osrmfiles/dotproperties"
	"github.com/Telenav/osrm-backend/integration/osrmfiles/dottimestamp"

	"github.com/golang/glog"
)

const (
	dotOSRMSuffix = ".osrm"

	dotTimestampSuffix        = ".timestamp"
	dotOSRMDotTimestampSuffix = dotOSRMSuffix + dotTimestampSuffix

	dotNBGNodesSuffix        = ".nbg_nodes"
	dotOSRMDotNBGNodesSuffix = dotOSRMSuffix + dotNBGNodesSuffix

	dotPropertiesSuffix        = ".properties"
	dotOSRMDotPropertiesSuffix = dotOSRMSuffix + dotPropertiesSuffix
)

// osrmBasefilePath should be 'xxx.osrm'
func createEmptyOSRMFilesContents(osrmBasefilePath string) map[string]osrmfiles.ContentsOperator {

	m := map[string]osrmfiles.ContentsOperator{}
	m[dotOSRMSuffix] = dotosrm.New(osrmBasefilePath)
	m[dotOSRMDotTimestampSuffix] = dottimestamp.New(osrmBasefilePath + dotTimestampSuffix)
	m[dotOSRMDotNBGNodesSuffix] = dotnbgnodes.New(osrmBasefilePath+dotNBGNodesSuffix, flags.packBits)
	m[dotOSRMDotPropertiesSuffix] = dotproperties.New(osrmBasefilePath + dotPropertiesSuffix)

	return m
}

func main() {
	flag.Parse()
	defer glog.Flush()

	suffixIndex := strings.LastIndex(flags.filePath, dotOSRMSuffix)
	if suffixIndex < 0 {
		glog.Errorf("file path %s should end by .osrm[.xxx]\n", flags.filePath)
		return
	}
	suffix := flags.filePath[suffixIndex:]                       // should be '.osrm' or '.osrm.xxx'
	baseFilePath := flags.filePath[:suffixIndex] + dotOSRMSuffix // should be xxx.osrm

	// create empty files and contents mapping for loading later
	osrmContents := createEmptyOSRMFilesContents(baseFilePath)

	if suffix != dotOSRMSuffix || (suffix == dotOSRMSuffix && flags.singleFile) {
		// only keep the specified contents if want to load single file
		for k := range osrmContents {
			if k != suffix {
				delete(osrmContents, k)
			}
		}
	}

	if len(osrmContents) == 0 {
		glog.Warningf("nothing need to load for %s", flags.filePath)
		return
	}

	// load contents and print summary
	for k, c := range osrmContents {
		if c == nil {
			glog.Errorf("nil Contents to load %s", k)
			continue
		}
		if err := osrmfiles.Load(c); err != nil {
			glog.Error(err)
			continue
		}

		if flags.printSummary >= 0 {
			c.PrintSummary(flags.printSummary)
		}
	}

}
