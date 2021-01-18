package dotebgnodes

import (
	"fmt"
	"io"

	"github.com/Telenav/osrm-backend/integration/osrmfiles/fingerprint"
	"github.com/Telenav/osrm-backend/integration/osrmfiles/meta"
	"github.com/Telenav/osrm-backend/integration/osrmfiles/osrmtype"
	"github.com/golang/glog"
)

// Contents represents `.osrm.ebg_nodes` file structure.
type Contents struct {
	Fingerprint fingerprint.Fingerprint

	AnnotationsMeta meta.Num
	Annotations     osrmtype.NodeBasedEdgeAnnotations
	NodesMeta       meta.Num
	Nodes           osrmtype.EdgeBasedNodes

	// for internal implementation
	writers  map[string]io.Writer
	filePath string
}

// New creates an empty Contents for `.osrm.ebg_nodes`.
func New(file string) *Contents {
	c := Contents{}

	c.filePath = file

	// init writers
	c.writers = map[string]io.Writer{}
	c.writers["osrm_fingerprint.meta"] = &c.Fingerprint
	c.writers["/common/ebg_node_data/annotations.meta"] = &c.AnnotationsMeta
	c.writers["/common/ebg_node_data/annotations"] = &c.Annotations
	c.writers["/common/ebg_node_data/nodes.meta"] = &c.NodesMeta
	c.writers["/common/ebg_node_data/nodes"] = &c.Nodes

	return &c
}

// PrintSummary prints summary and head lines of contents.
func (c *Contents) PrintSummary(head int) {
	glog.Infof("Loaded from %s\n", c.filePath)
	glog.Infof("  %s\n", &c.Fingerprint)

	glog.Infof("  edge_based_nodes meta %d count %d\n", c.NodesMeta, len(c.Nodes))
	for i := 0; i < head && i < len(c.Nodes); i++ {
		glog.Infof("    edge_based_nodes[%d] %+v", i, c.Nodes[i])
	}
	glog.Infof("  annotations meta %d count %d\n", c.AnnotationsMeta, len(c.Annotations))
	for i := 0; i < head && i < len(c.Annotations); i++ {
		glog.Infof("    annotations[%d] %+v", i, c.Annotations[i])
	}
}

// Validate checks whether the contents valid or not.
func (c *Contents) Validate() error {
	if !c.Fingerprint.IsValid() {
		return fmt.Errorf("invalid fingerprint %v", c.Fingerprint)
	}

	if uint64(c.NodesMeta) != uint64(len(c.Nodes)) {
		return fmt.Errorf("edge_based_nodes meta not match, count in meta %d, but actual edge_based_nodes count %d", c.NodesMeta, len(c.Nodes))
	}
	if uint64(c.AnnotationsMeta) != uint64(len(c.Annotations)) {
		return fmt.Errorf("annotations meta not match, count in meta %d, but actual annotations count %d", c.AnnotationsMeta, len(c.Annotations))
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
