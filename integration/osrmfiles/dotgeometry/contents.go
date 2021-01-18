package dotgeometry

import (
	"fmt"
	"io"

	"github.com/Telenav/osrm-backend/integration/util/builtinio"

	"github.com/Telenav/osrm-backend/integration/osrmfiles/fingerprint"
	"github.com/Telenav/osrm-backend/integration/osrmfiles/meta"
	"github.com/Telenav/osrm-backend/integration/osrmfiles/osrmtype"
	"github.com/Telenav/osrm-backend/integration/osrmfiles/osrmtype/packed"
	"github.com/golang/glog"
)

// Contents represents `.osrm.geometry` file structure.
type Contents struct {
	Fingerprint            fingerprint.Fingerprint
	IndexMeta              meta.Num
	Index                  []uint32 // https://github.com/Telenav/osrm-backend/blob/6900e30070a4ed3f1ca59004d57010a344cc7c9b/include/extractor/segment_data_container.hpp#L215
	NodesMeta              meta.Num
	Nodes                  osrmtype.NodeIDs
	ForwardDataSourcesMeta meta.Num
	ForwardDataSources     []uint8
	ReverseDataSourcesMeta meta.Num
	ReverseDataSources     []uint8

	ForwardWeights   packed.Uint64Vector
	ReverseWeights   packed.Uint64Vector
	ForwardDurations packed.Uint64Vector
	ReverseDurations packed.Uint64Vector

	// for internal implementation
	writers  map[string]io.Writer
	filePath string
}

// New creates an empty Contents for `.osrm.nbg_nodes`.
func New(file string) *Contents {
	c := Contents{
		ForwardWeights:   packed.NewUint64Vector(22), // https://github.com/Telenav/osrm-backend/blob/6900e30070a4ed3f1ca59004d57010a344cc7c9b/include/util/typedefs.hpp#L109
		ReverseWeights:   packed.NewUint64Vector(22),
		ForwardDurations: packed.NewUint64Vector(22), // https://github.com/Telenav/osrm-backend/blob/6900e30070a4ed3f1ca59004d57010a344cc7c9b/include/util/typedefs.hpp#L110
		ReverseDurations: packed.NewUint64Vector(22),
	}

	c.filePath = file

	// init writers
	c.writers = map[string]io.Writer{}
	c.writers["osrm_fingerprint.meta"] = &c.Fingerprint
	c.writers["/common/segment_data/index.meta"] = &c.IndexMeta
	c.writers["/common/segment_data/index"] = builtinio.BindWriterOnUint32Slice(&c.Index)
	c.writers["/common/segment_data/nodes.meta"] = &c.NodesMeta
	c.writers["/common/segment_data/nodes"] = &c.Nodes
	c.writers["/common/segment_data/forward_data_sources.meta"] = &c.ForwardDataSourcesMeta
	c.writers["/common/segment_data/forward_data_sources"] = builtinio.BindWriterOnUint8Slice(&c.ForwardDataSources)
	c.writers["/common/segment_data/reverse_data_sources.meta"] = &c.ReverseDataSourcesMeta
	c.writers["/common/segment_data/reverse_data_sources"] = builtinio.BindWriterOnUint8Slice(&c.ReverseDataSources)
	c.writers["/common/segment_data/forward_weights/number_of_elements.meta"] = &c.ForwardWeights.NumOfElements
	c.writers["/common/segment_data/forward_weights/packed.meta"] = &c.ForwardWeights.PackedMeta
	c.writers["/common/segment_data/forward_weights/packed"] = &c.ForwardWeights
	c.writers["/common/segment_data/reverse_weights/number_of_elements.meta"] = &c.ReverseWeights.NumOfElements
	c.writers["/common/segment_data/reverse_weights/packed.meta"] = &c.ReverseWeights.PackedMeta
	c.writers["/common/segment_data/reverse_weights/packed"] = &c.ReverseWeights
	c.writers["/common/segment_data/forward_durations/number_of_elements.meta"] = &c.ForwardDurations.NumOfElements
	c.writers["/common/segment_data/forward_durations/packed.meta"] = &c.ForwardDurations.PackedMeta
	c.writers["/common/segment_data/forward_durations/packed"] = &c.ForwardDurations
	c.writers["/common/segment_data/reverse_durations/number_of_elements.meta"] = &c.ReverseDurations.NumOfElements
	c.writers["/common/segment_data/reverse_durations/packed.meta"] = &c.ReverseDurations.PackedMeta
	c.writers["/common/segment_data/reverse_durations/packed"] = &c.ReverseDurations

	return &c
}

// PrintSummary prints summary and head lines of contents.
func (c *Contents) PrintSummary(head int) {
	glog.Infof("Loaded from %s\n", c.filePath)
	glog.Infof("  %s\n", &c.Fingerprint)

	glog.Infof("  index meta %d count %d\n", c.IndexMeta, len(c.Index))
	for i := 0; i < head && i < len(c.Index); i++ {
		glog.Infof("    index[%d] %+v", i, c.Index[i])
	}
	glog.Infof("  nodes meta %d count %d\n", c.NodesMeta, len(c.Nodes))
	for i := 0; i < head && i < len(c.Nodes); i++ {
		glog.Infof("    nodes[%d] %+v", i, c.Nodes[i])
	}
	glog.Infof("  forward_data_sources meta %d count %d\n", c.ForwardDataSourcesMeta, len(c.ForwardDataSources))
	for i := 0; i < head && i < len(c.ForwardDataSources); i++ {
		glog.Infof("    forward_data_sources[%d] %+v", i, c.ForwardDataSources[i])
	}
	glog.Infof("  reverse_data_sources meta %d count %d\n", c.ReverseDataSourcesMeta, len(c.ReverseDataSources))
	for i := 0; i < head && i < len(c.ReverseDataSources); i++ {
		glog.Infof("    reverse_data_sources[%d] %+v", i, c.ReverseDataSources[i])
	}

	glog.Infof("  forward_weights number_of_elements meta %d count %d\n", c.ForwardWeights.NumOfElements, len(c.ForwardWeights.Values))
	glog.Infof("  forward_weights packed meta %d\n", c.ForwardWeights.PackedMeta)
	for i := 0; i < head && i < len(c.ForwardWeights.Values); i++ {
		glog.Infof("    forward_weights[%d] %d", i, c.ForwardWeights.Values[i])
	}
	glog.Infof("  reverse_weights number_of_elements meta %d count %d\n", c.ReverseWeights.NumOfElements, len(c.ReverseWeights.Values))
	glog.Infof("  reverse_weights packed meta %d\n", c.ReverseWeights.PackedMeta)
	for i := 0; i < head && i < len(c.ReverseWeights.Values); i++ {
		glog.Infof("    reverse_weights[%d] %d", i, c.ReverseWeights.Values[i])
	}
	glog.Infof("  forward_durations number_of_elements meta %d count %d\n", c.ForwardDurations.NumOfElements, len(c.ForwardDurations.Values))
	glog.Infof("  forward_durations packed meta %d\n", c.ForwardDurations.PackedMeta)
	for i := 0; i < head && i < len(c.ForwardDurations.Values); i++ {
		glog.Infof("    forward_durations[%d] %d", i, c.ForwardDurations.Values[i])
	}
	glog.Infof("  reverse_durations number_of_elements meta %d count %d\n", c.ReverseDurations.NumOfElements, len(c.ReverseDurations.Values))
	glog.Infof("  reverse_durations packed meta %d\n", c.ReverseDurations.PackedMeta)
	for i := 0; i < head && i < len(c.ReverseDurations.Values); i++ {
		glog.Infof("    reverse_durations[%d] %d", i, c.ReverseDurations.Values[i])
	}
}

// Validate checks whether the contents valid or not.
func (c *Contents) Validate() error {
	if !c.Fingerprint.IsValid() {
		return fmt.Errorf("invalid fingerprint %v", c.Fingerprint)
	}

	if uint64(c.IndexMeta) != uint64(len(c.Index)) {
		return fmt.Errorf("Index meta not match, count in meta %d, but actual index count %d", c.IndexMeta, len(c.Index))
	}
	if uint64(c.NodesMeta) != uint64(len(c.Nodes)) {
		return fmt.Errorf("Nodes meta not match, count in meta %d, but actual nodes count %d", c.NodesMeta, len(c.Nodes))
	}
	if uint64(c.ForwardDataSourcesMeta) != uint64(len(c.ForwardDataSources)) {
		return fmt.Errorf("forward_data_sources meta not match, count in meta %d, but actual forward_data_sources count %d", c.ForwardDataSourcesMeta, len(c.ForwardDataSources))
	}
	if uint64(c.ReverseDataSourcesMeta) != uint64(len(c.ReverseDataSources)) {
		return fmt.Errorf("reverse_data_sources meta not match, count in meta %d, but actual reverse_data_sources count %d", c.ReverseDataSourcesMeta, len(c.ReverseDataSources))
	}

	if err := c.ForwardWeights.Validate(); err != nil {
		return err
	}
	if err := c.ReverseWeights.Validate(); err != nil {
		return err
	}
	if err := c.ForwardDurations.Validate(); err != nil {
		return err
	}
	if err := c.ReverseDurations.Validate(); err != nil {
		return err
	}

	// sentinel check: https://github.com/Telenav/osrm-backend/blob/6900e30070a4ed3f1ca59004d57010a344cc7c9b/src/extractor/compressed_edge_container.cpp#L395
	indexSentinel := meta.Num(c.Index[len(c.Index)-1])
	if indexSentinel != c.NodesMeta || indexSentinel != c.ForwardDataSourcesMeta || indexSentinel != c.ReverseDataSourcesMeta || indexSentinel != c.ForwardWeights.NumOfElements || indexSentinel != c.ReverseWeights.NumOfElements || indexSentinel != c.ForwardDurations.NumOfElements || indexSentinel != c.ReverseDurations.NumOfElements {
		return fmt.Errorf("index sentinel and meta not match: %d %d %d %d %d %d %d %d",
			indexSentinel, c.NodesMeta, c.ForwardDataSourcesMeta, c.ReverseDataSourcesMeta, c.ForwardWeights.NumOfElements, c.ReverseWeights.NumOfElements, c.ForwardDurations.NumOfElements, c.ReverseDurations.NumOfElements)
	}

	return nil
}

// PostProcess post process the conents once contents loaded if necessary.
func (c *Contents) PostProcess() error {

	if err := c.ForwardWeights.Prune(); err != nil {
		return err
	}
	if err := c.ReverseWeights.Prune(); err != nil {
		return err
	}
	if err := c.ForwardDurations.Prune(); err != nil {
		return err
	}
	if err := c.ReverseDurations.Prune(); err != nil {
		return err
	}

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
