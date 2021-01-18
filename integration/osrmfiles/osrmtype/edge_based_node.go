package osrmtype

import (
	"encoding/binary"
	"fmt"
)

// EdgeBasedNode represents a specific direction (forward or backward) of an segment.
// Terminology: https://github.com/Telenav/open-source-spec/blob/master/osrm/doc/understanding_osrm_graph_representation.md#terminology
// C++ Implementation: https://github.com/Telenav/osrm-backend/blob/6283c6074066f98e6d4a9f774f21ea45407c0d52/include/extractor/edge_based_node.hpp#L11
type EdgeBasedNode struct {
	GeometryID        // uses 4 bytes in C++ implementation
	ComponentID       // uses 4 bytes in C++ implementation
	AnnotationID      // only uses 31 bits in C++ implementation
	Segregated   bool // only uses 1 bit in C++ implementation
}

// EdgeBasedNodes represents vector of EdgeBasedNode.
type EdgeBasedNodes []EdgeBasedNode

const (
	annotationIDAndSegregatedBytes = 4
	edgeBasedNodeBytes             = geometryIDBytes + componentIDBytes + annotationIDAndSegregatedBytes // 12 bytes
)

func (e *EdgeBasedNode) tryParse(p []byte) error {
	if len(p) < edgeBasedNodeBytes {
		return fmt.Errorf("at least %d bytes for EdgeBasedNode but only got %d bytes", edgeBasedNodeBytes, len(p))
	}

	writeP := p

	if err := e.GeometryID.tryParse(writeP); err != nil {
		return err
	}
	writeP = writeP[geometryIDBytes:]

	if err := e.ComponentID.tryParse(writeP); err != nil {
		return err
	}
	writeP = writeP[componentIDBytes:]

	annotationIDBytes := []byte{writeP[0], writeP[1], writeP[2], writeP[3] & 0x7F}
	e.AnnotationID = AnnotationID(binary.LittleEndian.Uint32(annotationIDBytes))
	if writeP[3]&0x80 > 0 {
		e.Segregated = true
	} else {
		e.Segregated = false
	}

	return nil
}

func (e *EdgeBasedNodes) Write(p []byte) (int, error) {

	var writeLen int
	writeP := p
	for {
		if len(writeP) < edgeBasedNodeBytes {
			break
		}

		var n EdgeBasedNode
		if err := n.tryParse(writeP); err != nil {
			return writeLen, err
		}

		*e = append(*e, n)

		writeP = writeP[edgeBasedNodeBytes:]
		writeLen += edgeBasedNodeBytes
	}

	return writeLen, nil
}
