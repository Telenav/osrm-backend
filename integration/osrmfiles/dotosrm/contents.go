package dotosrm

import (
	"fmt"
	"io"

	"github.com/Telenav/osrm-backend/integration/osrmfiles/osrmtype"

	"github.com/Telenav/osrm-backend/integration/osrmfiles/meta"
	"github.com/Telenav/osrm-backend/integration/osrmfiles/osrmtype/querynode"

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
	EdgesMeta         meta.Num
	Edges             osrmtype.NodeBasedEdges
	AnnotationsMeta   meta.Num
	Annotations       osrmtype.NodeBasedEdgeAnnotations

	// for internal implementation
	writers  map[string]io.Writer
	filePath string
}

// New creates an empty Contents for `.osrm`.
func New(file string) *Contents {
	c := Contents{}

	c.filePath = file

	// init writers
	c.writers = map[string]io.Writer{}
	c.writers["osrm_fingerprint.meta"] = &c.Fingerprint
	c.writers["/extractor/nodes.meta"] = &c.NodesMeta
	c.writers["/extractor/nodes"] = &c.Nodes
	c.writers["/extractor/barriers.meta"] = &c.BarriersMeta
	c.writers["/extractor/barriers"] = &c.Barriers
	c.writers["/extractor/traffic_lights.meta"] = &c.TrafficLightsMeta
	c.writers["/extractor/traffic_lights"] = &c.TrafficLights
	c.writers["/extractor/edges.meta"] = &c.EdgesMeta
	c.writers["/extractor/edges"] = &c.Edges
	c.writers["/extractor/annotations.meta"] = &c.AnnotationsMeta
	c.writers["/extractor/annotations"] = &c.Annotations

	return &c
}

// PrintSummary prints summary and head lines of contents.
func (c *Contents) PrintSummary(head int) {
	glog.Infof("Loaded from %s\n", c.filePath)
	glog.Infof("  %s\n", &c.Fingerprint)

	glog.Infof("  nodes meta %d count %d\n", c.NodesMeta, len(c.Nodes))
	for i := 0; i < head && i < len(c.Nodes); i++ {
		glog.Infof("    node[%d] %#v", i, c.Nodes[i])
	}

	glog.Infof("  barriers meta %d count %d\n", c.BarriersMeta, len(c.Barriers))
	for i := 0; i < head && i < len(c.Barriers); i++ {
		glog.Infof("    barrier[%d] %v", i, c.Barriers[i])
	}

	glog.Infof("  traffic_lights meta %d count %d\n", c.TrafficLightsMeta, len(c.TrafficLights))
	for i := 0; i < head && i < len(c.TrafficLights); i++ {
		glog.Infof("    traffic_lights[%d] %v", i, c.TrafficLights[i])
	}

	glog.Infof("  edges meta %d count %d\n", c.EdgesMeta, len(c.Edges))
	for i := 0; i < head && i < len(c.Edges); i++ {
		glog.Infof("    edges[%d] %#v", i, c.Edges[i])
	}

	glog.Infof("  annotations meta %d count %d\n", c.AnnotationsMeta, len(c.Annotations))
	for i := 0; i < head && i < len(c.Annotations); i++ {
		glog.Infof("    annotations[%d] %#v", i, c.Annotations[i])
	}

}

// Validate checks whether the contents valid or not.
func (c *Contents) Validate() error {
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
	if uint64(c.EdgesMeta) != uint64(len(c.Edges)) {
		return fmt.Errorf("edges meta not match, count in meta %d, but actual edges count %d", c.EdgesMeta, len(c.Edges))
	}
	if uint64(c.AnnotationsMeta) != uint64(len(c.Annotations)) {
		return fmt.Errorf("annotations meta not match, count in meta %d, but actual annotations count %d", c.AnnotationsMeta, len(c.Annotations))
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

// PostProcess post process the conents once contents loaded if necessary.
func (c *Contents) PostProcess() error {
	return nil // nothing need to do
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
