// Package osrmtype declares OSRM defined types, e.g. NodeID, EdgeID, etc.
// C++ Implementation: https://github.com/Telenav/osrm-backend/blob/6283c6074066f98e6d4a9f774f21ea45407c0d52/include/util/typedefs.hpp#L72
package osrmtype

import "github.com/Telenav/osrm-backend/integration/util/builtinio"

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

	nodeIDBytes = 4 // uint32

)

func (n *NodeIDs) Write(p []byte) (int, error) {

	var u []uint32
	nWrite, err := builtinio.BindWriterOnUint32Slice(&u).Write(p)
	if err != nil {
		return nWrite, err
	}

	for _, v := range u {
		*n = append(*n, NodeID(v))
	}

	return nWrite, nil
}
