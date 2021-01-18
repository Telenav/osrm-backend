package dotturnpenaltiesindex

import (
	"fmt"
	"io"

	"github.com/Telenav/osrm-backend/integration/osrmfiles/fingerprint"
	"github.com/Telenav/osrm-backend/integration/osrmfiles/osrmtype"
	"github.com/golang/glog"
)

// Contents represents `.osrm.turn_penalties_index` file structure.
type Contents struct {
	Fingerprint fingerprint.Fingerprint

	TurnIndexes osrmtype.TurnIndexBlocks

	// for internal implementation
	writers  map[string]io.Writer
	filePath string
}

// New creates an empty Contents for `.osrm.turn_penalties_index`.
func New(file string) *Contents {
	c := Contents{}

	c.filePath = file

	// init writers
	c.writers = map[string]io.Writer{}
	c.writers["osrm_fingerprint.meta"] = &c.Fingerprint
	c.writers["/extractor/turn_index"] = &c.TurnIndexes
	return &c
}

// PrintSummary prints summary and head lines of contents.
func (c *Contents) PrintSummary(head int) {
	glog.Infof("Loaded from %s\n", c.filePath)
	glog.Infof("  %s\n", &c.Fingerprint)

	glog.Infof("  turn_indexes count %d\n", len(c.TurnIndexes.Values))
	for i := 0; i < head && i < len(c.TurnIndexes.Values); i++ {
		glog.Infof("    turn_indexes[%d] %+v", i, c.TurnIndexes.Values[i])
	}
}

// Validate checks whether the contents valid or not.
func (c *Contents) Validate() error {
	if !c.Fingerprint.IsValid() {
		return fmt.Errorf("invalid fingerprint %v", c.Fingerprint)
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
