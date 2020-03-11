package osrmtype

import (
	"encoding/binary"
)

// CompressedNodeBasedGraphEdge represent edge in compressed node based graph.
// C++ implementation: https://github.com/Telenav/osrm-backend/blob/9e063361797c1905021040d4ed3c452116fedfc6/include/extractor/compressed_node_based_graph_edge.hpp#L12
type CompressedNodeBasedGraphEdge struct {
	Source NodeID // 4 bytes in .osrm.cnbg file
	Target NodeID // 4 bytes in .osrm.cnbg file
}

// CompressedNodeBasedGraphEdges represents vector of CompressedNodeBasedGraphEdge.
type CompressedNodeBasedGraphEdges []CompressedNodeBasedGraphEdge

const (
	compressedNodeBasedGraphEdgeBytes = 8
)

func (c *CompressedNodeBasedGraphEdges) Write(p []byte) (int, error) {

	var writeLen int
	writeP := p
	for {
		if len(writeP) < compressedNodeBasedGraphEdgeBytes {
			break
		}

		var edge CompressedNodeBasedGraphEdge
		edge.Source = NodeID(binary.LittleEndian.Uint32(writeP))
		edge.Target = NodeID(binary.LittleEndian.Uint32(writeP[4:]))
		*c = append(*c, edge)

		writeP = writeP[compressedNodeBasedGraphEdgeBytes:]
		writeLen += compressedNodeBasedGraphEdgeBytes
	}

	return writeLen, nil
}
