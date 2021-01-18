package dotebg

import (
	"fmt"
	"io"

	"github.com/Telenav/osrm-backend/integration/osrmfiles/osrmtype"

	"github.com/Telenav/osrm-backend/integration/osrmfiles/fingerprint"
	"github.com/Telenav/osrm-backend/integration/osrmfiles/meta"
	"github.com/Telenav/osrm-backend/integration/util/builtinio"
	"github.com/golang/glog"
)

// Contents represents `.osrm.ebg` file structure.
type Contents struct {
	Fingerprint fingerprint.Fingerprint

	NumberOfEdgeBasedNodesMeta meta.Num
	NumberOfEdgeBasedNodes     uint32
	EdgeBasedEdgeListMeta      meta.Num
	EdgeBasedEdgeList          osrmtype.EdgeBasedEdges
	ConnectivityChecksumMeta   meta.Num
	ConnectivityChecksum       uint32

	// for internal implementation
	writers  map[string]io.Writer
	filePath string
}

// New creates an empty Contents for `.osrm.ebg`.
func New(file string) *Contents {
	c := Contents{}

	c.filePath = file

	// init writers
	c.writers = map[string]io.Writer{}
	c.writers["osrm_fingerprint.meta"] = &c.Fingerprint
	c.writers["/common/number_of_edge_based_nodes.meta"] = &c.NumberOfEdgeBasedNodesMeta
	c.writers["/common/number_of_edge_based_nodes"] = builtinio.BindWriterOnUint32(&c.NumberOfEdgeBasedNodes)
	c.writers["/common/edge_based_edge_list.meta"] = &c.EdgeBasedEdgeListMeta
	c.writers["/common/edge_based_edge_list"] = &c.EdgeBasedEdgeList
	c.writers["/common/connectivity_checksum.meta"] = &c.ConnectivityChecksumMeta
	c.writers["/common/connectivity_checksum"] = builtinio.BindWriterOnUint32(&c.ConnectivityChecksum)

	return &c
}

// PrintSummary prints summary and head lines of contents.
func (c *Contents) PrintSummary(head int) {
	glog.Infof("Loaded from %s\n", c.filePath)
	glog.Infof("  %s\n", &c.Fingerprint)

	glog.Infof("  number_of_edge_based_nodes meta %d count 1\n", c.NumberOfEdgeBasedNodesMeta)
	glog.Infof("  number_of_edge_based_nodes: %d", c.NumberOfEdgeBasedNodes)

	glog.Infof("  edge_based_edge_list meta %d count %d\n", c.EdgeBasedEdgeListMeta, len(c.EdgeBasedEdgeList))
	for i := 0; i < head && i < len(c.EdgeBasedEdgeList); i++ {
		glog.Infof("    edge_based_edge_list[%d] %+v", i, c.EdgeBasedEdgeList[i])
	}

	glog.Infof("  connectivity_checksum meta %d count 1\n", c.ConnectivityChecksumMeta)
	glog.Infof("  connectivity_checksum: %d", c.ConnectivityChecksum)

}

// Validate checks whether the contents valid or not.
func (c *Contents) Validate() error {
	if !c.Fingerprint.IsValid() {
		return fmt.Errorf("invalid fingerprint %v", c.Fingerprint)
	}

	if uint64(c.EdgeBasedEdgeListMeta) != uint64(len(c.EdgeBasedEdgeList)) {
		return fmt.Errorf("edge_based_edge_list meta not match, count in meta %d, but actual edge_based_edge_list count %d", c.EdgeBasedEdgeListMeta, len(c.EdgeBasedEdgeList))
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
