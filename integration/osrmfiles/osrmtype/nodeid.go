// Package osrmtype declares OSRM defined types, e.g. NodeID, EdgeID, etc.
// C++ Implementation: https://github.com/Telenav/osrm-backend/blob/6283c6074066f98e6d4a9f774f21ea45407c0d52/include/util/typedefs.hpp#L72
package osrmtype

import "encoding/binary"

// NodeID represents OSRM defined Node ID.
// C++ Implementation: https://github.com/Telenav/osrm-backend/blob/6283c6074066f98e6d4a9f774f21ea45407c0d52/include/util/typedefs.hpp#L72
type NodeID uint32

// NodeIDs represents vector of NodeID.
type NodeIDs []NodeID

// SpecialNodeID represents invalid NodeID.
// C++ Implementation: https://github.com/Telenav/osrm-backend/blob/6283c6074066f98e6d4a9f774f21ea45407c0d52/include/util/typedefs.hpp#L102
const (
	SpecialNodeID  = ^NodeID(0)
	InvalidNodeID  = SpecialNodeID
	MaxValidNodeID = InvalidNodeID - 1 // 2^32-1
)

const nodeIDBytes = 4 // uint32

func (n *NodeIDs) Write(p []byte) (int, error) {

	var writeLen int
	writeP := p
	for {
		if len(writeP) < nodeIDBytes {
			break
		}

		var id NodeID
		id = NodeID(binary.LittleEndian.Uint32(writeP))

		*n = append(*n, id)

		writeP = writeP[nodeIDBytes:]
		writeLen += nodeIDBytes
	}

	return writeLen, nil
}
