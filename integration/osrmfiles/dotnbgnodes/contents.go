package dotnbgnodes

import (
	"fmt"
	"io"

	"github.com/Telenav/osrm-backend/integration/osrmfiles/fingerprint"
	"github.com/Telenav/osrm-backend/integration/osrmfiles/meta"
	"github.com/Telenav/osrm-backend/integration/osrmfiles/osrmtype"
	"github.com/Telenav/osrm-backend/integration/osrmfiles/osrmtype/packed"
	"github.com/golang/glog"
)

// Contents represents `.osrm.nbg_nodes` file structure.
type Contents struct {
	Fingerprint     fingerprint.Fingerprint
	CoordinatesMeta meta.Num
	Coordinates     osrmtype.Coordinates
	OSMNodeIDs      packed.Uint64Vector

	// for internal implementation
	writers  map[string]io.Writer
	filePath string
}

// New creates an empty Contents for `.osrm.nbg_nodes`.
func New(file string, packedBits uint) *Contents {
	c := Contents{
		OSMNodeIDs: packed.NewUint64Vector(packedBits),
	}

	c.filePath = file

	// init writers
	c.writers = map[string]io.Writer{}
	c.writers["osrm_fingerprint.meta"] = &c.Fingerprint
	c.writers["/common/nbn_data/coordinates.meta"] = &c.CoordinatesMeta
	c.writers["/common/nbn_data/coordinates"] = &c.Coordinates
	c.writers["/common/nbn_data/osm_node_ids/number_of_elements.meta"] = &c.OSMNodeIDs.NumOfElements
	c.writers["/common/nbn_data/osm_node_ids/packed.meta"] = &c.OSMNodeIDs.PackedMeta
	c.writers["/common/nbn_data/osm_node_ids/packed"] = &c.OSMNodeIDs

	return &c
}

// PrintSummary prints summary and head lines of contents.
func (c *Contents) PrintSummary(head int) {
	glog.Infof("Loaded from %s\n", c.filePath)
	glog.Infof("  %s\n", &c.Fingerprint)

	glog.Infof("  coordinates meta %d count\n", c.CoordinatesMeta)
	for i := 0; i < head && i < len(c.Coordinates); i++ {
		glog.Infof("    coordinate[%d] %v", i, c.Coordinates[i])
	}

	glog.Infof("  osm_node_ids number_of_elements meta %d count\n", c.OSMNodeIDs.NumOfElements)
	glog.Infof("  osm_node_ids packed meta %d count\n", c.OSMNodeIDs.PackedMeta)
	for i := 0; i < head && i < len(c.OSMNodeIDs.Values); i++ {
		glog.Infof("    osm_node_ids[%d] %v", i, c.OSMNodeIDs.Values[i])
	}
}

// Validate checks whether the contents valid or not.
func (c *Contents) Validate() error {
	if !c.Fingerprint.IsValid() {
		return fmt.Errorf("invalid fingerprint %v", c.Fingerprint)
	}
	if uint64(c.CoordinatesMeta) != uint64(len(c.Coordinates)) {
		return fmt.Errorf("coordinates meta not match, count in meta %d, but actual coordinates count %d", c.CoordinatesMeta, len(c.Coordinates))
	}
	if err := c.OSMNodeIDs.Validate(); err != nil {
		return err
	}

	return nil
}

// PostProcess post process the conents once contents loaded if necessary.
func (c *Contents) PostProcess() error {
	return c.OSMNodeIDs.Prune()
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
