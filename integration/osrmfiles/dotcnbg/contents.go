package dotcnbg

import (
	"fmt"
	"io"

	"github.com/Telenav/osrm-backend/integration/osrmfiles/fingerprint"
	"github.com/Telenav/osrm-backend/integration/osrmfiles/meta"
	"github.com/Telenav/osrm-backend/integration/osrmfiles/osrmtype"
	"github.com/golang/glog"
)

// Contents represents `.osrm.cnbg` file structure.
type Contents struct {
	Fingerprint                       fingerprint.Fingerprint
	CompressedNodeBasedGraphEdgesMeta meta.Num
	osrmtype.CompressedNodeBasedGraphEdges

	// for internal implementation
	writers  map[string]io.Writer
	filePath string
}

// New creates an empty Contents for `.osrm.cnbg`.
func New(file string) *Contents {
	c := Contents{}

	c.filePath = file

	// init writers
	c.writers = map[string]io.Writer{}
	c.writers["osrm_fingerprint.meta"] = &c.Fingerprint
	c.writers["/extractor/cnbg.meta"] = &c.CompressedNodeBasedGraphEdgesMeta
	c.writers["/extractor/cnbg"] = &c.CompressedNodeBasedGraphEdges

	return &c
}

// PrintSummary prints summary and head lines of contents.
func (c *Contents) PrintSummary(head int) {
	glog.Infof("Loaded from %s\n", c.filePath)
	glog.Infof("  %s\n", &c.Fingerprint)

	glog.Infof("  compressed node based graph meta %d count %d\n", c.CompressedNodeBasedGraphEdgesMeta, len(c.CompressedNodeBasedGraphEdges))
	for i := 0; i < head && i < len(c.CompressedNodeBasedGraphEdges); i++ {
		glog.Infof("    CompressedNodeBasedGraphEdges[%d] %v", i, c.CompressedNodeBasedGraphEdges[i])
	}

}

// Validate checks whether the contents valid or not.
func (c *Contents) Validate() error {
	if !c.Fingerprint.IsValid() {
		return fmt.Errorf("invalid fingerprint %v", c.Fingerprint)
	}
	if uint64(c.CompressedNodeBasedGraphEdgesMeta) != uint64(len(c.CompressedNodeBasedGraphEdges)) {
		return fmt.Errorf("CompressedNodeBasedGraphEdges meta not match, count in meta %d, but actual count %d", c.CompressedNodeBasedGraphEdgesMeta, len(c.CompressedNodeBasedGraphEdges))
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
