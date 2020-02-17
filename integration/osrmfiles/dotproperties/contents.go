package dotproperties

import (
	"fmt"
	"io"

	"github.com/Telenav/osrm-backend/integration/osrmfiles/fingerprint"
	"github.com/Telenav/osrm-backend/integration/osrmfiles/meta"
	"github.com/Telenav/osrm-backend/integration/osrmfiles/osrmtype/profile"
	"github.com/golang/glog"
)

// Contents represents `.osrm.properties` file structure.
type Contents struct {
	Fingerprint    fingerprint.Fingerprint
	PropertiesMeta meta.Num
	profile.Properties

	// for internal implementation
	writers  map[string]io.Writer
	filePath string
}

// New creates an empty Contents for `.osrm.properties`.
func New(file string) *Contents {
	c := Contents{}

	c.filePath = file

	// init writers
	c.writers = map[string]io.Writer{}
	c.writers["osrm_fingerprint.meta"] = &c.Fingerprint
	c.writers["/common/properties.meta"] = &c.PropertiesMeta
	c.writers["/common/properties"] = &c.Properties

	return &c
}

// Validate checks whether the contents valid or not.
func (c *Contents) Validate() error {
	if !c.Fingerprint.IsValid() {
		return fmt.Errorf("invalid fingerprint %v", c.Fingerprint)
	}
	if c.PropertiesMeta != 1 { // always only have 1 lua properies stored
		return fmt.Errorf("invalid properties.meta %d, expect always 1", c.PropertiesMeta)
	}

	return nil
}

// PostProcess post process the conents once contents loaded if necessary.
func (c *Contents) PostProcess() error {
	return nil // nothing need to do
}

// PrintSummary prints summary and head lines of contents.
func (c *Contents) PrintSummary(head int) {
	glog.Infof("Loaded from %s\n", c.filePath)
	glog.Infof("  %s\n", &c.Fingerprint)

	glog.Infof("  propertiesMeta meta %d\n", c.PropertiesMeta)
	if head > 0 {
		glog.Infof("  properties %#v\n", c.Properties)
	}
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
