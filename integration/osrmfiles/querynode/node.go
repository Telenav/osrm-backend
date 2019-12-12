// Package querynode implements QueryNode structure in .osrm file.
// C++ implementation https://github.com/Telenav/osrm-backend/blob/master-telenav/include/extractor/query_node.hpp
package querynode

import (
	"encoding/binary"
)

// Node represents QueryNode structure.
// C++ implementation https://github.com/Telenav/osrm-backend/blob/6283c6074066f98e6d4a9f774f21ea45407c0d52/include/extractor/query_node.hpp#L16
type Node struct {
	// COORDINATE_PRECISION = 1e6
	// C++ implementation toFixed()/toFloat() see https://github.com/Telenav/osrm-backend/blob/6283c6074066f98e6d4a9f774f21ea45407c0d52/include/util/coordinate.hpp#L88
	FixedLon int32
	FixedLat int32

	ID uint64
}

// Nodes represents QueryNode vector structure.
type Nodes []Node

const (
	fixedLonBytes = 4
	fixedLatBytes = 4
	idBytes       = 8
	nodeBytes     = fixedLonBytes + fixedLatBytes + idBytes
)

func (n *Nodes) Write(p []byte) (int, error) {

	var writeLen int
	writeP := p
	for {
		if len(writeP) < nodeBytes {
			break
		}

		var node Node
		node.FixedLon = int32(binary.LittleEndian.Uint32(writeP))
		node.FixedLat = int32(binary.LittleEndian.Uint32(writeP[fixedLonBytes:]))
		node.ID = binary.LittleEndian.Uint64(writeP[fixedLonBytes+fixedLatBytes:])

		*n = append(*n, node)

		writeP = writeP[nodeBytes:]
		writeLen += nodeBytes
	}

	return writeLen, nil
}
