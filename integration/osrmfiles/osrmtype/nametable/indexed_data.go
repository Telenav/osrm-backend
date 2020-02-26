package nametable

import (
	"encoding/binary"
	"fmt"

	"github.com/golang/glog"
)

// blockReference is designed for reference block in IndexedData.
// C++ implementation: https://github.com/Telenav/osrm-backend/blob/b24b8a085dc10bea279ffb352049330beae23791/include/util/indexed_data.hpp#L51
type blockReference struct {
	offset     uint32
	descriptor uint32
}

const (
	blockSize           = 16 // https://github.com/Telenav/osrm-backend/blob/b24b8a085dc10bea279ffb352049330beae23791/include/extractor/name_table.hpp#L52
	blockContainingSize = blockSize + 1

	twoBitsMax = 0x3

	blockReferenceOffsetBytes     = 4
	blockReferenceDescriptorBytes = 4
	blockReferenceBytes           = blockReferenceOffsetBytes + blockReferenceDescriptorBytes
)

// C++ implementation: https://github.com/Telenav/osrm-backend/blob/b24b8a085dc10bea279ffb352049330beae23791/include/util/indexed_data.hpp#L332
func (i *IndexedData) at(id uint32) ([]byte, error) {

	blockIndex := id / blockContainingSize
	blockInternalIndex := id % blockContainingSize

	if int(blockIndex) >= len(i.blocks) {
		return nil, fmt.Errorf("too big block index %d from id %d, total blocks count %d", blockIndex, id, len(i.blocks))
	}

	first := i.blocks[blockIndex].offset
	dataEnd := uint32(len(i.values))
	if blockIndex+1 < uint32(len(i.blocks)) {
		dataEnd = i.blocks[blockIndex+1].offset
	}

	descriptor := i.blocks[blockIndex].descriptor
	dataStart := first + sum2Bits(descriptor)

	for j := uint32(0); j < blockInternalIndex; j++ {
		byteLen := descriptor & twoBitsMax
		dataStart, first = i.varAdvance(byteLen, dataStart, first)
		descriptor = descriptor >> 2 // each 2-bits represents a length
	}

	if blockInternalIndex < blockSize {
		dataEnd, _ = i.varAdvance(descriptor&twoBitsMax, dataStart, first)
	}

	if dataStart >= uint32(len(i.values)) || dataEnd >= uint32(len(i.values)) {
		return nil, fmt.Errorf("invalid data start,end(%d,%d) position, len(values)=%d", dataStart, dataEnd, len(i.values))
	}

	return i.values[dataStart:dataEnd], nil
}

// C++ implementation: https://github.com/Telenav/osrm-backend/blob/b24b8a085dc10bea279ffb352049330beae23791/include/util/indexed_data.hpp#L70
func (i *IndexedData) varAdvance(byteLen, dataStart, lengthStart uint32) (uint32, uint32) {
	if byteLen > twoBitsMax { // byteLen should only in 2-bits
		glog.Fatalf("invalid byteLen %d", byteLen)
	}

	if byteLen == 1 {
		dataStart += uint32(i.values[lengthStart])
		lengthStart++
	} else if byteLen == 2 {
		dataStart += uint32(i.values[lengthStart])
		lengthStart++
		dataStart += uint32(i.values[lengthStart]) << 8
		lengthStart++
	} else if byteLen == 3 {
		dataStart += uint32(i.values[lengthStart])
		lengthStart++
		dataStart += uint32(i.values[lengthStart]) << 8
		lengthStart++
		dataStart += uint32(i.values[lengthStart]) << 16
		lengthStart++
	}
	return dataStart, lengthStart
}

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
