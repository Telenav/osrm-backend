package dotturnweightpenalties

import (
	"fmt"
	"io"

	"github.com/Telenav/osrm-backend/integration/util/builtinio"

	"github.com/Telenav/osrm-backend/integration/osrmfiles/fingerprint"
	"github.com/Telenav/osrm-backend/integration/osrmfiles/meta"
	"github.com/golang/glog"
)

// Contents represents `.osrm.turn_weight_penalties` file structure.
type Contents struct {
	Fingerprint fingerprint.Fingerprint

	TurnWeights     []int16
	TurnWeightsMeta meta.Num

	// for internal implementation
	writers  map[string]io.Writer
	filePath string
}

// New creates an empty Contents for `.osrm.turn_weight_penalties`.
func New(file string) *Contents {
	c := Contents{}

	c.filePath = file

	// init writers
	c.writers = map[string]io.Writer{}
	c.writers["osrm_fingerprint.meta"] = &c.Fingerprint
	c.writers["/common/turn_penalty/weight"] = builtinio.BindWriterOnInt16Slice(&c.TurnWeights)
	c.writers["/common/turn_penalty/weight.meta"] = &c.TurnWeightsMeta
	return &c
}

// PrintSummary prints summary and head lines of contents.
func (c *Contents) PrintSummary(head int) {
	glog.Infof("Loaded from %s\n", c.filePath)
	glog.Infof("  %s\n", &c.Fingerprint)

	glog.Infof("  turn_weights meta %d count %d\n", c.TurnWeightsMeta, len(c.TurnWeights))
	for i := 0; i < head && i < len(c.TurnWeights); i++ {
		glog.Infof("    turn_weights[%d] %+v", i, c.TurnWeights[i])
	}
}

// Validate checks whether the contents valid or not.
func (c *Contents) Validate() error {
	if !c.Fingerprint.IsValid() {
		return fmt.Errorf("invalid fingerprint %v", c.Fingerprint)
	}

	if uint64(c.TurnWeightsMeta) != uint64(len(c.TurnWeights)) {
		return fmt.Errorf("turn_weights meta not match, count in meta %d, but actual count %d", c.TurnWeightsMeta, len(c.TurnWeights))
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
