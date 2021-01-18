package dotcellmetrics

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/Telenav/osrm-backend/integration/util/builtinio"

	"github.com/Telenav/osrm-backend/integration/osrmfiles/fingerprint"
	"github.com/Telenav/osrm-backend/integration/osrmfiles/meta"
	"github.com/golang/glog"
)

// Contents represents `.osrm.cell_metrics` file structure.
type Contents struct {
	Fingerprint fingerprint.Fingerprint

	// weight_name -> ExcludeCellMetrics mapping
	Metrics map[string]*ExcludeCellMetrics

	// for internal implementation
	writers  map[string]io.Writer
	filePath string
}

// New creates an empty Contents for `.osrm.cell_metrics`.
func New(file string) *Contents {
	c := Contents{Metrics: map[string]*ExcludeCellMetrics{}}

	c.filePath = file

	// init writers
	c.writers = map[string]io.Writer{}
	c.writers["osrm_fingerprint.meta"] = &c.Fingerprint
	// other writers will be created automatically in finding

	return &c
}

// PrintSummary prints summary and head lines of contents.
func (c *Contents) PrintSummary(head int) {
	glog.Infof("Loaded from %s\n", c.filePath)
	glog.Infof("  %s\n", &c.Fingerprint)

	for k, v := range c.Metrics {
		glog.Infof("  metrics for weight_name: %s", k)
		glog.Infof("    exclude metrics meta %d count %d\n", v.Meta, len(v.CellMetrics))
		for i := 0; i < len(v.CellMetrics); i++ {
			glog.Infof("      metrics[%d] weights meta %d count %d\n", i, v.CellMetrics[i].WeightsMeta, len(v.CellMetrics[i].Weights))
			for j := 0; j < head && j < len(v.CellMetrics[i].Weights); j++ {
				glog.Infof("        metrics[%d] weigths[%d] %+v", i, j, v.CellMetrics[i].Weights[j])
			}

			glog.Infof("      metrics[%d] durations meta %d count %d\n", i, v.CellMetrics[i].DurationsMeta, len(v.CellMetrics[i].Durations))
			for j := 0; j < head && j < len(v.CellMetrics[i].Durations); j++ {
				glog.Infof("        metrics[%d] durations[%d] %+v", i, j, v.CellMetrics[i].Durations[j])
			}

			glog.Infof("      metrics[%d] distances meta %d count %d\n", i, v.CellMetrics[i].DistancesMeta, len(v.CellMetrics[i].Distances))
			for j := 0; j < head && j < len(v.CellMetrics[i].Distances); j++ {
				glog.Infof("        metrics[%d] distances[%d] %+v", i, j, v.CellMetrics[i].Distances[j])
			}
		}
	}
}

// Validate checks whether the contents valid or not.
func (c *Contents) Validate() error {
	if !c.Fingerprint.IsValid() {
		return fmt.Errorf("invalid fingerprint %v", c.Fingerprint)
	}

	for k, v := range c.Metrics {

		if uint64(v.Meta) != uint64(len(v.CellMetrics)) {
			return fmt.Errorf("weight_name %s, metrics meta not match, count in meta %d, but actual metrics count %d", k, v.Meta, len(v.CellMetrics))
		}

		for i := 0; i < len(v.CellMetrics); i++ {
			if uint64(v.CellMetrics[i].WeightsMeta) != uint64(len(v.CellMetrics[i].Weights)) {
				return fmt.Errorf("weight_name %s, metrics[%d], weights meta not match, count in meta %d, but actual weights count %d", k, i, v.CellMetrics[i].WeightsMeta, len(v.CellMetrics[i].Weights))
			}

			if uint64(v.CellMetrics[i].DurationsMeta) != uint64(len(v.CellMetrics[i].Durations)) {
				return fmt.Errorf("weight_name %s, metrics[%d], durations meta not match, count in meta %d, but actual durations count %d", k, i, v.CellMetrics[i].DurationsMeta, len(v.CellMetrics[i].Durations))
			}

			if uint64(v.CellMetrics[i].DistancesMeta) != uint64(len(v.CellMetrics[i].Distances)) {
				return fmt.Errorf("weight_name %s, metrics[%d], distances meta not match, count in meta %d, but actual distances count %d", k, i, v.CellMetrics[i].DistancesMeta, len(v.CellMetrics[i].Distances))
			}
		}
	}

	return nil
}

// PostProcess post process the conents once contents loaded if necessary.
func (c *Contents) PostProcess() error {
	return nil
}

// FindWriter find io.Writer for the specified name.
func (c *Contents) FindWriter(name string) (io.Writer, bool) {
	if w, ok := c.writers[name]; ok {
		return w, ok
	}

	// typical format
	// /mld/metrics/routability/exclude.meta
	// /mld/metrics/routability/exclude/0/weights.meta
	// /mld/metrics/routability/exclude/0/weights
	// /mld/metrics/routability/exclude/0/durations.meta
	// /mld/metrics/routability/exclude/0/durations
	// /mld/metrics/routability/exclude/0/distances.meta
	// /mld/metrics/routability/exclude/0/distances
	// /mld/metrics/routability/exclude/1/weights.meta
	// /mld/metrics/routability/exclude/1/weights
	// /mld/metrics/routability/exclude/1/durations.meta
	// /mld/metrics/routability/exclude/1/durations
	// /mld/metrics/routability/exclude/1/distances.meta
	// /mld/metrics/routability/exclude/1/distances

	s := strings.Split(name, "/")
	if len(s) < 5 { // insufficient to parse weight_name
		return nil, false
	}
	weightName := s[3]
	if _, ok := c.Metrics[weightName]; !ok { // first touch weight name, create the structs in mapping
		c.Metrics[weightName] = &ExcludeCellMetrics{
			CellMetrics: []CellMetric{},
		}
	}

	// typically: '/mld/metrics/routability/exclude.meta'
	if strings.HasSuffix(name, "exclude.meta") {
		return &c.Metrics[weightName].Meta, true
	}

	if len(s) < 7 { // insufficient to parse id
		return nil, false
	}

	id, err := strconv.Atoi(s[5])
	if err != nil {
		glog.Warningf("parse exclude id from %s, err: %v", name, err)
		return nil, false
	}
	for len(c.Metrics[weightName].CellMetrics) <= id {
		c.Metrics[weightName].CellMetrics = append(c.Metrics[weightName].CellMetrics, CellMetric{
			Weights:   []int32{},
			Durations: []int32{},
			Distances: []float32{},
		})
	}
	switch s[len(s)-1] { // last field
	case "weights.meta":
		return &c.Metrics[weightName].CellMetrics[id].WeightsMeta, true
	case "weights":
		return builtinio.BindWriterOnInt32Slice(&c.Metrics[weightName].CellMetrics[id].Weights), true
	case "durations.meta":
		return &c.Metrics[weightName].CellMetrics[id].DurationsMeta, true
	case "durations":
		return builtinio.BindWriterOnInt32Slice(&c.Metrics[weightName].CellMetrics[id].Durations), true
	case "distances.meta":
		return &c.Metrics[weightName].CellMetrics[id].DistancesMeta, true
	case "distances":
		return builtinio.BindWriterOnFloat32Slice(&c.Metrics[weightName].CellMetrics[id].Distances), true
	}

	return nil, false
}

// FilePath returns the file path that stores the contents.
func (c *Contents) FilePath() string {
	return c.filePath
}

// CellMetric represents weights,durations and distances for cell.
// C++ Implementation: https://github.com/Telenav/osrm-backend/blob/40015847054011efbd61c2912e7ff4c135b6a570/include/customizer/cell_metric.hpp#L17
type CellMetric struct {
	WeightsMeta   meta.Num
	Weights       []int32
	DurationsMeta meta.Num
	Durations     []int32
	DistancesMeta meta.Num
	Distances     []float32
}

// ExcludeCellMetrics represents all exclude flags related metrics for one cell.
// C++ Implementation: https://github.com/Telenav/osrm-backend/blob/40015847054011efbd61c2912e7ff4c135b6a570/src/customize/customizer.cpp#L262
type ExcludeCellMetrics struct {
	Meta        meta.Num
	CellMetrics []CellMetric
}
