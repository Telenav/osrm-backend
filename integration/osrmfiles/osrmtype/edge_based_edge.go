package osrmtype

import (
	"encoding/binary"
	"fmt"
	"math"
)

// EdgeData represents data of edge-expanded graph edge.
// C++ Implementation: https://github.com/Telenav/osrm-backend/blob/3906e34f2079617702a966b317c6431412852dfc/include/extractor/edge_based_edge.hpp#L16
type EdgeData struct {
	TurnID   NodeID
	Weight   EdgeWeight
	Distance EdgeDistance
	Duration EdgeDuration // 30 bits in files to save memory, see https://github.com/Telenav/osrm-backend/blob/3906e34f2079617702a966b317c6431412852dfc/include/extractor/edge_based_edge.hpp#L37
	Forward  bool         // 1 bit
	Backward bool         // 1 bit
}

// EdgeBasedEdge represents Edge-expanded graph edge.
// C++ Implementation: https://github.com/Telenav/osrm-backend/blob/3906e34f2079617702a966b317c6431412852dfc/include/extractor/edge_based_edge.hpp#L13
type EdgeBasedEdge struct {
	Source NodeID
	Target NodeID
	EdgeData
}

// EdgeBasedEdges represents vector of EdgeBasedEdge.
type EdgeBasedEdges []EdgeBasedEdge

const (
	edgeDataBytes      = nodeIDBytes + 3*4                         // 16 bytes
	edgeBasedEdgeBytes = nodeIDBytes + nodeIDBytes + edgeDataBytes // 24 bytes
)

func (e *EdgeBasedEdges) Write(p []byte) (int, error) {
	var writeLen int
	writeP := p
	for {
		if len(writeP) < edgeBasedEdgeBytes {
			break
		}

		var edge EdgeBasedEdge
		if err := edge.tryParse(writeP); err != nil {
			return writeLen, err
		}

		*e = append(*e, edge)

		writeP = writeP[edgeBasedEdgeBytes:]
		writeLen += edgeBasedEdgeBytes
	}

	return writeLen, nil

}

func (e *EdgeBasedEdge) tryParse(p []byte) error {
	if len(p) < edgeBasedEdgeBytes {
		return fmt.Errorf("at least %d bytes for EdgeBasedEdge but only got %d bytes", edgeBasedEdgeBytes, len(p))
	}

	writeP := p

	e.Source = NodeID(binary.LittleEndian.Uint32(writeP))
	writeP = writeP[nodeIDBytes:]

	e.Target = NodeID(binary.LittleEndian.Uint32(writeP))
	writeP = writeP[nodeIDBytes:]

	e.TurnID = NodeID(binary.LittleEndian.Uint32(writeP))
	writeP = writeP[nodeIDBytes:]

	e.Weight = EdgeWeight(binary.LittleEndian.Uint32(writeP))
	writeP = writeP[4:]

	e.Distance = EdgeDistance(math.Float32frombits(binary.LittleEndian.Uint32(writeP)))
	writeP = writeP[4:]

	durationBytes := []byte{writeP[0], writeP[1], writeP[2], writeP[3] & 0x3F}
	e.Duration = EdgeDuration(binary.LittleEndian.Uint32(durationBytes))

	if writeP[3]&0x40 > 0 {
		e.Forward = true
	} else {
		e.Forward = false
	}

	if writeP[3]&0x80 > 0 {
		e.Backward = true
	} else {
		e.Backward = false
	}

	return nil
}
