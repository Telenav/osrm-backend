package osrmtype

import (
	"encoding/binary"
	"fmt"
)

// GeometryID represents OSRM defined Geometry ID.
// C++ Implementation: https://github.com/Telenav/osrm-backend/blob/6283c6074066f98e6d4a9f774f21ea45407c0d52/include/util/typedefs.hpp#L155
type GeometryID struct {
	ID      NodeID // only uses 31 bits in C++ implementation
	Forward bool   // only uses 1 bit in C++ implementation
}

func (g *GeometryID) tryParse(p []byte) error {

	if len(p) < 4 {
		return fmt.Errorf("at least 4 bytes for GeometryID but only got %d bytes", len(p))
	}

	idBytes := []byte{p[0], p[1], p[2], p[3] & 0x7F}
	g.ID = NodeID(binary.LittleEndian.Uint32(idBytes))
	if p[3]&0x80 > 0 {
		g.Forward = true
	} else {
		g.Forward = false
	}
	return nil
}
