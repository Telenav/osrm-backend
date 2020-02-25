package nametable

import (
	"encoding/binary"
	"fmt"
)

// blockReference is designed for reference block in IndexedData.
// C++ implementation: https://github.com/Telenav/osrm-backend/blob/b24b8a085dc10bea279ffb352049330beae23791/include/util/indexed_data.hpp#L51
type blockReference struct {
	offset     uint32
	descriptor uint32
}

const (
	blockSize = 16

	blockReferenceOffsetBytes     = 4
	blockReferenceDescriptorBytes = 4
	blockReferenceBytes           = blockReferenceOffsetBytes + blockReferenceDescriptorBytes
)

func (i *IndexedData) assembleBlockReferences() error {

	for i.BlocksBuffer.Len() > 0 {

		nextBytes := i.BlocksBuffer.Next(blockReferenceBytes)
		if len(nextBytes) < blockReferenceBytes {
			return fmt.Errorf("remain bytes %d insufficient for a blockReference, requires %d bytes", nextBytes, blockReferenceBytes)
		}

		var b blockReference
		b.offset = binary.LittleEndian.Uint32(nextBytes)
		b.descriptor = binary.LittleEndian.Uint32(nextBytes[blockReferenceOffsetBytes:])
		i.blocks = append(i.blocks, b)
	}

	return nil
}
