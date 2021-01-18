package osrmtype

import (
	"encoding/binary"
)

// TurnIndexBlock represents turn index that can be identify by from,via,to NodeID.
// C++ implementation: https://github.com/Telenav/osrm-backend/blob/6900e30070a4ed3f1ca59004d57010a344cc7c9b/include/extractor/edge_based_graph_factory.hpp#L48-L53
type TurnIndexBlock struct {
	From NodeID
	Via  NodeID
	To   NodeID
}

// TurnIndexBlocks represents slice of TurnIndexBlock.
type TurnIndexBlocks struct {
	Values []TurnIndexBlock

	// The `Write` will be called many times in `io.Copy` due to fixed buffer size.
	// Defaultly the copy buffer size is 32*1024 bytes, see line 391 `copyBuffer()` in https://golang.org/src/io/io.go.
	// Unfortunetly, we expect to parse the packed data which uses turnIndexBlockBytes(4*3=12bytes) for each TurnIndexBlock.
	// We have to cache remain bytes since 32*1024 / 12 != 0.
	// Maximum will store `12-1` bytes.
	unwrittenBytes []byte
}

// constant for TurnIndexBlock
const (
	turnIndexBlockBytes = nodeIDBytes * 3 // 3 * uint32
)

func (t *TurnIndexBlocks) Write(p []byte) (int, error) {

	t.unwrittenBytes = append(t.unwrittenBytes, p...)
	writeLen := len(p)

	writeP := t.unwrittenBytes
	for {

		if len(writeP) < turnIndexBlockBytes {
			break
		}

		var turnIndexBlock TurnIndexBlock
		turnIndexBlock.From = NodeID(binary.LittleEndian.Uint32(writeP))
		turnIndexBlock.Via = NodeID(binary.LittleEndian.Uint32(writeP[nodeIDBytes:]))
		turnIndexBlock.To = NodeID(binary.LittleEndian.Uint32(writeP[nodeIDBytes+nodeIDBytes:]))

		t.Values = append(t.Values, turnIndexBlock)

		writeP = writeP[turnIndexBlockBytes:]
	}

	t.unwrittenBytes = writeP
	return writeLen, nil
}
