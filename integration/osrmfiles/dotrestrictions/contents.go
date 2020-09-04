package dotrestrictions

import (
	"fmt"
	"io"

	"github.com/Telenav/osrm-backend/integration/osrmfiles/osrmtype/conditional"

	"github.com/Telenav/osrm-backend/integration/osrmfiles/fingerprint"
	"github.com/Telenav/osrm-backend/integration/osrmfiles/meta"
	"github.com/golang/glog"
)

// Contents represents `.osrm.restrictions` file structure.
type Contents struct {
	Fingerprint                  fingerprint.Fingerprint
	ConditionalTurnPenaltiesMeta meta.Num
	ConditionalTurnPenalties     conditional.TurnPenalties

	// for internal implementation
	writers  map[string]io.Writer
	filePath string
}

// New creates an empty Contents for `.osrm.restrictions`.
func New(file string) *Contents {
	c := Contents{}

	c.filePath = file

	// init writers
	c.writers = map[string]io.Writer{}
	c.writers["osrm_fingerprint.meta"] = &c.Fingerprint
	c.writers["/common/conditional_restrictions.meta"] = &c.ConditionalTurnPenaltiesMeta
	c.writers["/common/conditional_restrictions"] = &c.ConditionalTurnPenalties

	return &c
}

// PrintSummary prints summary and head lines of contents.
func (c *Contents) PrintSummary(head int) {
	glog.Infof("Loaded from %s\n", c.filePath)
	glog.Infof("  %s\n", &c.Fingerprint)

	glog.Infof("  ConditionalTurnPenalties meta %d (bytes)\n", c.ConditionalTurnPenaltiesMeta)
	glog.Infof("  ConditionalTurnPenalties count %d\n", len(c.ConditionalTurnPenalties.TurnPenalties))
	for i := 0; i < head && i < len(c.ConditionalTurnPenalties.TurnPenalties); i++ {
		glog.Infof("    ConditionalTurnPenalties[%d] %+v", i, c.ConditionalTurnPenalties.TurnPenalties[i])
	}
}

// Validate checks whether the contents valid or not.
func (c *Contents) Validate() error {
	if !c.Fingerprint.IsValid() {
		return fmt.Errorf("invalid fingerprint %v", c.Fingerprint)
	}

	return c.ConditionalTurnPenalties.Validate(uint64(c.ConditionalTurnPenaltiesMeta))
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
