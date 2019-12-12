package dotosrm

import (
	"archive/tar"
	"fmt"
	"io"
	"os"

	"github.com/Telenav/osrm-backend/integration/osrmfiles/meta"
	"github.com/Telenav/osrm-backend/integration/osrmfiles/querynode"

	"github.com/Telenav/osrm-backend/integration/osrmfiles/fingerprint"
	"github.com/golang/glog"
)

// Contents represents `.osrm` file structure.
type Contents struct {
	Fingerprint fingerprint.Fingerprint
	NodesMeta   meta.Num
	Nodes       querynode.Nodes

	//TODO:

	// for internal implementation
	writers  map[string]io.Writer
	filePath string
}

// Load `.osrm` file to generate a new contents structure.
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

// PrintSummary prints summary of current contents.
func (c *Contents) PrintSummary() {
	glog.Infof("Loaded from %s\n", c.filePath)
	glog.Infof("  %s\n", &c.Fingerprint)
	glog.Infof("  nodes meta %d count %d\n", c.NodesMeta, len(c.Nodes))
}

func new() *Contents {
	c := Contents{}

	// init writers
	c.writers = map[string]io.Writer{}
	c.writers["osrm_fingerprint.meta"] = &c.Fingerprint
	c.writers["/extractor/nodes.meta"] = &c.NodesMeta
	c.writers["/extractor/nodes"] = &c.Nodes

	return &c
}

func (c *Contents) validate() error {
	if !c.Fingerprint.IsValid() {
		return fmt.Errorf("invalid fingerprint %v", c.Fingerprint)
	}
	if uint64(c.NodesMeta) != uint64(len(c.Nodes)) {
		return fmt.Errorf("nodes meta not match, count in meta %d, but actual nodes count %d", c.NodesMeta, len(c.Nodes))
	}

	return nil
}
