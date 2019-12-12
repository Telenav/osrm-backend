package meta

import (
	"encoding/binary"
	"fmt"
)

// Num represents number of elements in `.osrm[.xxx]` files.
type Num uint64

const numBytes = 8 // only need 8 bytes for uint64

func (n *Num) Write(p []byte) (int, error) {
	if len(p) < numBytes {
		return 0, fmt.Errorf("byte array len %d insufficient", len(p))
	}
	*n = Num(binary.LittleEndian.Uint64(p))
	return numBytes, nil
}
