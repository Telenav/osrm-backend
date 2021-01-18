package dotenw

import (
	"fmt"
	"io"

	"github.com/Telenav/osrm-backend/integration/util/builtinio"

	"github.com/Telenav/osrm-backend/integration/osrmfiles/meta"

	"github.com/Telenav/osrm-backend/integration/osrmfiles/fingerprint"
	"github.com/golang/glog"
)

// Contents represents `.osrm.enw` file structure.
type Contents struct {
	Fingerprint fingerprint.Fingerprint

	EdgeBasedNodeWeights       []int32
	EdgeBasedNodeWeightsMeta   meta.Num
	EdgeBasedNodeDurations     []int32
	EdgeBasedNodeDurationsMeta meta.Num
	EdgeBasedNodeDistances     []float32
	EdgeBasedNodeDistancesMeta meta.Num

	// for internal implementation
	writers  map[string]io.Writer
	filePath string
}

// New creates an empty Contents for `.osrm.enw`.
func New(file string) *Contents {
	c := Contents{}

	c.filePath = file

	// init writers
	c.writers = map[string]io.Writer{}
	c.writers["osrm_fingerprint.meta"] = &c.Fingerprint
	c.writers["/extractor/edge_based_node_weights"] = builtinio.BindWriterOnInt32Slice(&c.EdgeBasedNodeWeights)
	c.writers["/extractor/edge_based_node_weights.meta"] = &c.EdgeBasedNodeWeightsMeta
	c.writers["/extractor/edge_based_node_durations"] = builtinio.BindWriterOnInt32Slice(&c.EdgeBasedNodeDurations)
	c.writers["/extractor/edge_based_node_durations.meta"] = &c.EdgeBasedNodeDurationsMeta
	c.writers["/extractor/edge_based_node_distances"] = builtinio.BindWriterOnFloat32Slice(&c.EdgeBasedNodeDistances)
	c.writers["/extractor/edge_based_node_distances.meta"] = &c.EdgeBasedNodeDistancesMeta

	return &c
}

// PrintSummary prints summary and head lines of contents.
func (c *Contents) PrintSummary(head int) {
	glog.Infof("Loaded from %s\n", c.filePath)
	glog.Infof("  %s\n", &c.Fingerprint)

	glog.Infof("  edge_based_node_weights meta %d count %d\n", c.EdgeBasedNodeWeightsMeta, len(c.EdgeBasedNodeWeights))
	for i := 0; i < head && i < len(c.EdgeBasedNodeWeights); i++ {
		glog.Infof("    edge_based_node_weights[%d] %+v", i, c.EdgeBasedNodeWeights[i])
	}
	glog.Infof("  edge_based_node_durations meta %d count %d\n", c.EdgeBasedNodeDurationsMeta, len(c.EdgeBasedNodeDurations))
	for i := 0; i < head && i < len(c.EdgeBasedNodeDurations); i++ {
		glog.Infof("    edge_based_node_durations[%d] %+v", i, c.EdgeBasedNodeDurations[i])
	}
	glog.Infof("  edge_based_node_distances meta %d count %d\n", c.EdgeBasedNodeDistancesMeta, len(c.EdgeBasedNodeDistances))
	for i := 0; i < head && i < len(c.EdgeBasedNodeDistances); i++ {
		glog.Infof("    edge_based_node_distances[%d] %+v", i, c.EdgeBasedNodeDistances[i])
	}

}

// Validate checks whether the contents valid or not.
func (c *Contents) Validate() error {
	if !c.Fingerprint.IsValid() {
		return fmt.Errorf("invalid fingerprint %v", c.Fingerprint)
	}

	if uint64(c.EdgeBasedNodeWeightsMeta) != uint64(len(c.EdgeBasedNodeWeights)) {
		return fmt.Errorf("edge_based_node_weights meta not match, count in meta %d, but actual count %d", c.EdgeBasedNodeWeightsMeta, len(c.EdgeBasedNodeWeights))
	}
	if uint64(c.EdgeBasedNodeDurationsMeta) != uint64(len(c.EdgeBasedNodeDurations)) {
		return fmt.Errorf("edge_based_node_durations meta not match, count in meta %d, but actual count %d", c.EdgeBasedNodeDurationsMeta, len(c.EdgeBasedNodeDurations))
	}
	if uint64(c.EdgeBasedNodeDistancesMeta) != uint64(len(c.EdgeBasedNodeDistances)) {
		return fmt.Errorf("edge_based_node_distances meta not match, count in meta %d, but actual count %d", c.EdgeBasedNodeDistancesMeta, len(c.EdgeBasedNodeDistances))
	}

	return nil
}

// PostProcess post process the conents once contents loaded if necessary.
func (c *Contents) PostProcess() error {
	return nil
}

// FindWriter find io.Writer for the specified name.
func (c *Contents) FindWriter(name string) (io.Writer, bool) {
	w, b := c.writers[name]
	return w, b
}

// FilePath returns the file path that stores the contents.
func (c *Contents) FilePath() string {
	return c.filePath
}
