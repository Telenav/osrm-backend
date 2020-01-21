package dotosrmdottimestamp

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"os"

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

// Load `.osrm.timestamp` file to generate a new contents structure.
func Load(file string) (*Contents, error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	glog.V(2).Infof("open %s succeed.\n", file)

	contents := new()

	// Open and iterate through the files in the archive.
	tr := tar.NewReader(f)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			glog.Fatal(err)
		}
		glog.V(1).Infof("%s\n", hdr.Name)
		writer, found := contents.writers[hdr.Name]
		if !found {
			glog.Warningf("unrecognized content in tar: %s", hdr.Name)
			continue
		}

		if _, err := io.Copy(writer, tr); err != nil {
			glog.Fatal(err)
		}
	}

	// validate loaded contents
	if err := contents.validate(); err != nil {
		return nil, err
	}

	contents.filePath = file
	return contents, nil
}

func new() *Contents {
	c := Contents{}

	// init writers
	c.writers = map[string]io.Writer{}
	c.writers["osrm_fingerprint.meta"] = &c.Fingerprint
	c.writers["/common/timestamp.meta"] = &c.TimestampMeta
	c.writers["/common/timestamp"] = &c.Timestamp

	return &c
}

func (c *Contents) validate() error {
	if !c.Fingerprint.IsValid() {
		return fmt.Errorf("invalid fingerprint %v", c.Fingerprint)
	}
	if uint64(c.TimestampMeta) != uint64(c.Timestamp.Len()) {
		return fmt.Errorf("timestamp meta not match, count in meta %d, but actual timestamp bytse count %d", c.TimestampMeta, c.Timestamp.Len())
	}
	return nil
}

// PrintSummary prints summary and head lines of current contents.
func (c *Contents) PrintSummary(head int) {
	//TODO:
	glog.Infof("Loaded from %s\n", c.filePath)
	glog.Infof("  %s\n", &c.Fingerprint)

	glog.Infof("  timestamp(e.g. data_version) meta %d count %d\n", c.TimestampMeta, c.Timestamp.Len())
	glog.Infof("  timestamp(e.g. data_version) %v\n", c.Timestamp)
}
