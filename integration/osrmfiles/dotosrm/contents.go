package dotosrm

import (
	"archive/tar"
	"fmt"
	"io"
	"os"

	"github.com/Telenav/osrm-backend/integration/osrmfiles/osrmtype"

	"github.com/Telenav/osrm-backend/integration/osrmfiles/meta"
	"github.com/Telenav/osrm-backend/integration/osrmfiles/querynode"

	"github.com/Telenav/osrm-backend/integration/osrmfiles/fingerprint"
	"github.com/golang/glog"
)

// Contents represents `.osrm` file structure.
type Contents struct {
	Fingerprint       fingerprint.Fingerprint
	NodesMeta         meta.Num
	Nodes             querynode.Nodes
	BarriersMeta      meta.Num
	Barriers          osrmtype.NodeIDs
	TrafficLightsMeta meta.Num
	TrafficLights     osrmtype.NodeIDs

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
	glog.Infof("  barriers meta %d count %d\n", c.BarriersMeta, len(c.Barriers))
	glog.Infof("  traffic_lights meta %d count %d\n", c.TrafficLightsMeta, len(c.TrafficLights))

}

func new() *Contents {
	c := Contents{}

	// init writers
	c.writers = map[string]io.Writer{}
	c.writers["osrm_fingerprint.meta"] = &c.Fingerprint
	c.writers["/extractor/nodes.meta"] = &c.NodesMeta
	c.writers["/extractor/nodes"] = &c.Nodes
	c.writers["/extractor/barriers.meta"] = &c.BarriersMeta
	c.writers["/extractor/barriers"] = &c.Barriers
	c.writers["/extractor/traffic_lights.meta"] = &c.TrafficLightsMeta
	c.writers["/extractor/traffic_lights"] = &c.TrafficLights

	return &c
}

func (c *Contents) validate() error {
	if !c.Fingerprint.IsValid() {
		return fmt.Errorf("invalid fingerprint %v", c.Fingerprint)
	}
	if uint64(c.NodesMeta) != uint64(len(c.Nodes)) {
		return fmt.Errorf("nodes meta not match, count in meta %d, but actual nodes count %d", c.NodesMeta, len(c.Nodes))
	}
	if uint64(c.BarriersMeta) != uint64(len(c.Barriers)) {
		return fmt.Errorf("barriers meta not match, count in meta %d, but actual barriers count %d", c.BarriersMeta, len(c.Barriers))
	}
	if uint64(c.TrafficLightsMeta) != uint64(len(c.TrafficLights)) {
		return fmt.Errorf("traffic_lights meta not match, count in meta %d, but actual traffic_lights count %d", c.TrafficLightsMeta, len(c.TrafficLights))
	}

	// check relationship between nodes and barriers/traffic_lights
	// https://github.com/Telenav/open-source-spec/blob/master/osrm/doc/osrm-toolchain-files/map.osrm.md
	if uint64(c.NodesMeta) > uint64(osrmtype.MaxValidNodeID) {
		return fmt.Errorf("too big nodes meta %d, osrm NodeID will overflow", c.NodesMeta)
	}
	if len(c.Barriers) > 0 {
		maxBarrierNodeID := c.Barriers[len(c.Barriers)-1]
		if uint64(c.NodesMeta) <= uint64(maxBarrierNodeID) {
			return fmt.Errorf("too big barrier NodeID %d for nodes meta %d", maxBarrierNodeID, c.NodesMeta)
		}
	}
	if len(c.TrafficLights) > 0 {
		maxTrafficLightNodeID := c.TrafficLights[len(c.TrafficLights)-1]
		if uint64(c.NodesMeta) <= uint64(maxTrafficLightNodeID) {
			return fmt.Errorf("too big traffic_light NodeID %d for nodes meta %d", maxTrafficLightNodeID, c.NodesMeta)
		}
	}

	return nil
}
