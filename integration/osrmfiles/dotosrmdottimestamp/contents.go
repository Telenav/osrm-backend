package dotosrmdottimestamp

import (
	"bytes"
	"fmt"
	"io"

	"github.com/Telenav/osrm-backend/integration/osrmfiles/fingerprint"
	"github.com/Telenav/osrm-backend/integration/osrmfiles/meta"
	"github.com/golang/glog"
)

// Contents represents `.osrm.timestamp` file structure.
type Contents struct {
	Fingerprint   fingerprint.Fingerprint
	TimestampMeta meta.Num
	Timestamp     bytes.Buffer

	// for internal implementation
	writers  map[string]io.Writer
	filePath string
}

// New creates an empty Contents for `.osrm.timestamp`.
func New(file string) *Contents {
	c := Contents{}

	c.filePath = file

	// init writers
	c.writers = map[string]io.Writer{}
	c.writers["osrm_fingerprint.meta"] = &c.Fingerprint
	c.writers["/common/timestamp.meta"] = &c.TimestampMeta
	c.writers["/common/timestamp"] = &c.Timestamp

	return &c
}

// Validate checks whether the contents valid or not.
func (c *Contents) Validate() error {
	if !c.Fingerprint.IsValid() {
		return fmt.Errorf("invalid fingerprint %v", c.Fingerprint)
	}
	if uint64(c.TimestampMeta) != uint64(c.Timestamp.Len()) {
		return fmt.Errorf("timestamp meta not match, count in meta %d, but actual timestamp bytse count %d", c.TimestampMeta, c.Timestamp.Len())
	}
	return nil
}

// PrintSummary prints summary and head lines of contents.
func (c *Contents) PrintSummary(head int) {
	glog.Infof("Loaded from %s\n", c.filePath)
	glog.Infof("  %s\n", &c.Fingerprint)

	glog.Infof("  timestamp(a.k.a. data_version) meta %d count %d\n", c.TimestampMeta, c.Timestamp.Len())
	if head > 0 {
		glog.Infof("  timestamp(a.k.a. data_version) %v\n", c.Timestamp.String())
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
