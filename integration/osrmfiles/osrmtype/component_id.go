package osrmtype

import (
	"encoding/binary"
	"fmt"
)

// ComponentID represents Strongly Connected Component ID of an edge-based node.
// C++ Implementation: https://github.com/Telenav/osrm-backend/blob/6283c6074066f98e6d4a9f774f21ea45407c0d52/include/util/typedefs.hpp#L168
type ComponentID struct {
	ID     uint32 // only uses 31 bits in C++ implementation
	IsTiny bool   // only uses 1 bit in C++ implementation
}

const componentIDBytes = 4

func (c *ComponentID) tryParse(p []byte) error {

	if len(p) < componentIDBytes {
		return fmt.Errorf("at least %d bytes for ComponentID but only got %d bytes", componentIDBytes, len(p))
	}

	idBytes := []byte{p[0], p[1], p[2], p[3] & 0x7F}
	c.ID = binary.LittleEndian.Uint32(idBytes)
	if p[3]&0x80 > 0 {
		c.IsTiny = true
	} else {
		c.IsTiny = false
	}
	return nil
}
