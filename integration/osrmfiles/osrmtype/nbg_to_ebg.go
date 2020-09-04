package osrmtype

import "encoding/binary"

// NBGToEBG represent mapping between the node based graph u,v nodes and the edge based graph head,tail edge ids.
// Required in the osrm-partition tool to translate from a nbg partition to a ebg partition.
// C++ implementation: https://github.com/Telenav/osrm-backend/blob/0b461183b97de493983ba44749c772719849fd3e/include/extractor/nbg_to_ebg.hpp#L13
type NBGToEBG struct {
	U, V                            NodeID
	ForwardEBGNode, BackwardEBGNode NodeID
}

// NBGToEBGs represents vector of NBGToEBG.
type NBGToEBGs []NBGToEBG

const (
	nBGToEBGBytes = 16
)

func (n *NBGToEBGs) Write(p []byte) (int, error) {

	var writeLen int
	writeP := p
	for {
		if len(writeP) < nBGToEBGBytes {
			break
		}

		var nbgToEBG NBGToEBG
		nbgToEBG.U = NodeID(binary.LittleEndian.Uint32(writeP))
		nbgToEBG.V = NodeID(binary.LittleEndian.Uint32(writeP[4:]))
		nbgToEBG.ForwardEBGNode = NodeID(binary.LittleEndian.Uint32(writeP[8:]))
		nbgToEBG.BackwardEBGNode = NodeID(binary.LittleEndian.Uint32(writeP[12:]))
		*n = append(*n, nbgToEBG)

		writeP = writeP[nBGToEBGBytes:]
		writeLen += nBGToEBGBytes
	}

	return writeLen, nil
}
