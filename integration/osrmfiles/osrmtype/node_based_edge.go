package osrmtype

import (
	"encoding/binary"
	"fmt"
	"math"
)

// NodeBasedEdge represents a segment(connect 2 OSM nodes) with direction.
// Terminology: https://github.com/Telenav/open-source-spec/blob/master/osrm/doc/understanding_osrm_graph_representation.md#terminology
// C++ Implementation: https://github.com/Telenav/osrm-backend/blob/6283c6074066f98e6d4a9f774f21ea45407c0d52/include/extractor/node_based_edge.hpp#L92
type NodeBasedEdge struct {
	Source       NodeID                      // 4 bytes in .osrm file
	Target       NodeID                      // 4 bytes in .osrm file
	Weight       EdgeWeight                  // 4 bytes in .osrm file
	Duration     EdgeDuration                // 4 bytes in .osrm file
	Distance     EdgeDistance                // 4 bytes in .osrm file
	GeometryID                               // 4 bytes in .osrm file
	AnnotationID                             // 4 bytes in .osrm file
	Flags        NodeBasedEdgeClassification // 4 bytes in .osrm file
}

// NodeBasedEdges represents vector of NodeBasedEdge.
type NodeBasedEdges []NodeBasedEdge

// NodeBasedEdgeClassification describing the class of the road.
// C++ Implementation: https://github.com/Telenav/osrm-backend/blob/6283c6074066f98e6d4a9f774f21ea45407c0d52/include/extractor/node_based_edge.hpp#L20
type NodeBasedEdgeClassification struct {
	Forward    bool // 1 bit in .osrm file
	Backward   bool // 1 bit in .osrm file
	IsSplit    bool // 1 bit in .osrm file
	Roundabout bool // 1 bit in .osrm file
	Circular   bool // 1 bit in .osrm file
	Startpoint bool // 1 bit in .osrm file
	Restricted bool // 1 bit in .osrm file
	// still 1 bit reserved
	RoadClassification              // 2 bytes in .osrm file
	HighwayTurnClassification uint8 // 4 bits in .osrm file
	AccessTurnClassification  uint8 // 4 bits in .osrm file
}

// NodeBasedEdgeAnnotation is annotative data, used in parts in guidance generation, in parts during navigation (classes) but
// mostly for annotation of edges. The entry can be shared between multiple edges and usually
// describes features present on OSM ways. This is the place to put specific data that you want to
// see as part of the API output but that does not influence navigation.
// C++ Implemnetation: https://github.com/Telenav/osrm-backend/blob/6283c6074066f98e6d4a9f774f21ea45407c0d52/include/extractor/node_based_edge.hpp#L66
type NodeBasedEdgeAnnotation struct {
	NameID                 // 4 bytes in .osrm file
	LaneDescriptionID      // 2 bytes in .osrm file
	ClassData              // 1 byte in .osrm file
	TravelMode             // 4 bits in .osrm file
	IsLeftHandDriving bool // 1 bit in .osrm file
}

// NodeBasedEdgeAnnotations represents vector of NodeBasedEdgeAnnotation.
type NodeBasedEdgeAnnotations []NodeBasedEdgeAnnotation

const (
	nodeBasedEdgeBytes           = 32 // 4 * 8
	nodeBasedEdgeAnnotationBytes = 8
)

func (n *NodeBasedEdges) Write(p []byte) (int, error) {

	var writeLen int
	writeP := p
	for {
		if len(writeP) < nodeBasedEdgeBytes {
			break
		}

		var edge NodeBasedEdge
		edge.Source = NodeID(binary.LittleEndian.Uint32(writeP))
		edge.Target = NodeID(binary.LittleEndian.Uint32(writeP[4:]))
		edge.Weight = EdgeWeight(binary.LittleEndian.Uint32(writeP[8:]))
		edge.Duration = EdgeDuration(binary.LittleEndian.Uint32(writeP[12:]))
		edge.Distance = EdgeDistance(
			math.Float32frombits(binary.LittleEndian.Uint32(writeP[16:])))
		if err := edge.GeometryID.tryParse(writeP[20:]); err != nil {
			return writeLen, err
		}
		edge.AnnotationID = AnnotationID(binary.LittleEndian.Uint32(writeP[24:]))
		if err := edge.Flags.tryParse(writeP[28:]); err != nil {
			return writeLen, err
		}

		*n = append(*n, edge)

		writeP = writeP[nodeBasedEdgeBytes:]
		writeLen += nodeBasedEdgeBytes
	}

	return writeLen, nil
}

func (n *NodeBasedEdgeClassification) tryParse(p []byte) error {

	if len(p) < 4 {
		return fmt.Errorf("at least 4 bytes for NodeBasedEdgeClassification but only got %d bytes", len(p))
	}

	if p[0]&0x01 > 0 {
		n.Forward = true
	}
	if p[0]&0x02 > 0 {
		n.Backward = true
	}
	if p[0]&0x04 > 0 {
		n.IsSplit = true
	}
	if p[0]&0x08 > 0 {
		n.Roundabout = true
	}
	if p[0]&0x10 > 0 {
		n.Circular = true
	}
	if p[0]&0x20 > 0 {
		n.Startpoint = true
	}
	if p[0]&0x40 > 0 {
		n.Restricted = true
	}
	if err := n.RoadClassification.tryParse(p[1:]); err != nil {
		return err
	}
	n.HighwayTurnClassification = p[3] & 0x0F
	n.AccessTurnClassification = (p[3] & 0xF0) >> 4

	return nil
}

func (n *NodeBasedEdgeAnnotations) Write(p []byte) (int, error) {

	var writeLen int
	writeP := p
	for {
		if len(writeP) < nodeBasedEdgeAnnotationBytes {
			break
		}

		var annotation NodeBasedEdgeAnnotation
		annotation.NameID = NameID(binary.LittleEndian.Uint32(writeP))
		annotation.LaneDescriptionID = LaneDescriptionID(binary.LittleEndian.Uint16(writeP[4:]))
		annotation.ClassData = ClassData(writeP[6])
		annotation.TravelMode = TravelMode(writeP[7] & 0x0F)
		if writeP[7]&0x10 > 0 {
			annotation.IsLeftHandDriving = true
		}

		*n = append(*n, annotation)

		writeP = writeP[nodeBasedEdgeAnnotationBytes:]
		writeLen += nodeBasedEdgeAnnotationBytes
	}

	return writeLen, nil
}
