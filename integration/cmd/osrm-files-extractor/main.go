package main

import (
	"flag"
	"strings"

	"github.com/Telenav/osrm-backend/integration/util/appversion"

	"github.com/Telenav/osrm-backend/integration/osrmfiles"
	"github.com/Telenav/osrm-backend/integration/osrmfiles/dotcellmetrics"
	"github.com/Telenav/osrm-backend/integration/osrmfiles/dotcnbg"
	"github.com/Telenav/osrm-backend/integration/osrmfiles/dotcnbgtoebg"
	"github.com/Telenav/osrm-backend/integration/osrmfiles/dotebg"
	"github.com/Telenav/osrm-backend/integration/osrmfiles/dotebgnodes"
	"github.com/Telenav/osrm-backend/integration/osrmfiles/dotenw"
	"github.com/Telenav/osrm-backend/integration/osrmfiles/dotgeometry"
	"github.com/Telenav/osrm-backend/integration/osrmfiles/dotnames"
	"github.com/Telenav/osrm-backend/integration/osrmfiles/dotnbgnodes"
	"github.com/Telenav/osrm-backend/integration/osrmfiles/dotosrm"
	"github.com/Telenav/osrm-backend/integration/osrmfiles/dotproperties"
	"github.com/Telenav/osrm-backend/integration/osrmfiles/dotrestrictions"
	"github.com/Telenav/osrm-backend/integration/osrmfiles/dottimestamp"
	"github.com/Telenav/osrm-backend/integration/osrmfiles/dotturndurationpenalties"
	"github.com/Telenav/osrm-backend/integration/osrmfiles/dotturnpenaltiesindex"
	"github.com/Telenav/osrm-backend/integration/osrmfiles/dotturnweightpenalties"

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

	dotNamesSuffix        = ".names"
	dotOSRMDotNamesSuffix = dotOSRMSuffix + dotNamesSuffix

	dotCNBGSuffix        = ".cnbg"
	dotOSRMDotCNBGSuffix = dotOSRMSuffix + dotCNBGSuffix

	dotCNBGToEBGSuffix        = ".cnbg_to_ebg"
	dotOSRMDotCNBGToEBGSuffix = dotOSRMSuffix + dotCNBGToEBGSuffix

	dotRestrictionsSuffix        = ".restrictions"
	dotOSRMDotRestrictionsSuffix = dotOSRMSuffix + dotRestrictionsSuffix

	dotGeometrySuffix        = ".geometry"
	dotOSRMDotGeometrySuffix = dotOSRMSuffix + dotGeometrySuffix

	dotENWSuffix        = ".enw"
	dotOSRMDotENWSuffix = dotOSRMSuffix + dotENWSuffix

	dotTurnPenaltiesIndexSuffix        = ".turn_penalties_index"
	dotOSRMDotTurnPenaltiesIndexSuffix = dotOSRMSuffix + dotTurnPenaltiesIndexSuffix

	dotTurnWeightPenaltiesSuffix        = ".turn_weight_penalties"
	dotOSRMDotTurnWeightPenaltiesSuffix = dotOSRMSuffix + dotTurnWeightPenaltiesSuffix

	dotTurnDurationPenaltiesSuffix        = ".turn_duration_penalties"
	dotOSRMDotTurnDurationPenaltiesSuffix = dotOSRMSuffix + dotTurnDurationPenaltiesSuffix

	dotEBGNodesSuffix        = ".ebg_nodes"
	dotOSRMDotEBGNodesSuffix = dotOSRMSuffix + dotEBGNodesSuffix

	dotCellMetricsSuffix        = ".cell_metrics"
	dotOSRMDotCellMetricsSuffix = dotOSRMSuffix + dotCellMetricsSuffix

	dotEBGSuffix        = ".ebg"
	dotOSRMDotEBGSuffix = dotOSRMSuffix + dotEBGSuffix
)

// osrmBasefilePath should be 'xxx.osrm'
func createEmptyOSRMFilesContents(osrmBasefilePath string) map[string]osrmfiles.ContentsOperator {

	m := map[string]osrmfiles.ContentsOperator{}
	m[dotOSRMSuffix] = dotosrm.New(osrmBasefilePath)
	m[dotOSRMDotTimestampSuffix] = dottimestamp.New(osrmBasefilePath + dotTimestampSuffix)
	m[dotOSRMDotNBGNodesSuffix] = dotnbgnodes.New(osrmBasefilePath+dotNBGNodesSuffix, flags.packBits)
	m[dotOSRMDotPropertiesSuffix] = dotproperties.New(osrmBasefilePath + dotPropertiesSuffix)
	m[dotOSRMDotNamesSuffix] = dotnames.New(osrmBasefilePath + dotNamesSuffix)
	m[dotOSRMDotCNBGSuffix] = dotcnbg.New(osrmBasefilePath + dotCNBGSuffix)
	m[dotOSRMDotCNBGToEBGSuffix] = dotcnbgtoebg.New(osrmBasefilePath + dotCNBGToEBGSuffix)
	m[dotOSRMDotRestrictionsSuffix] = dotrestrictions.New(osrmBasefilePath + dotRestrictionsSuffix)
	m[dotOSRMDotGeometrySuffix] = dotgeometry.New(osrmBasefilePath + dotGeometrySuffix)
	m[dotOSRMDotENWSuffix] = dotenw.New(osrmBasefilePath + dotENWSuffix)
	m[dotOSRMDotTurnPenaltiesIndexSuffix] = dotturnpenaltiesindex.New(osrmBasefilePath + dotTurnPenaltiesIndexSuffix)
	m[dotOSRMDotTurnWeightPenaltiesSuffix] = dotturnweightpenalties.New(osrmBasefilePath + dotTurnWeightPenaltiesSuffix)
	m[dotOSRMDotTurnDurationPenaltiesSuffix] = dotturndurationpenalties.New(osrmBasefilePath + dotTurnDurationPenaltiesSuffix)
	m[dotOSRMDotEBGNodesSuffix] = dotebgnodes.New(osrmBasefilePath + dotEBGNodesSuffix)
	m[dotOSRMDotCellMetricsSuffix] = dotcellmetrics.New(osrmBasefilePath + dotCellMetricsSuffix)
	m[dotOSRMDotEBGSuffix] = dotebg.New(osrmBasefilePath + dotEBGSuffix)

	return m
}

func main() {
	flag.Parse()
	appversion.PrintExit()
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
