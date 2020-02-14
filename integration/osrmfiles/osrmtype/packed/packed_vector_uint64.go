// Package packed implements the C++ `class PackedVector`.
// Instead of pack data in memory and decode value during retrieve,
//   this implementation is much more simpler: only parse the data when load from file.
// In another word, it parse the packed data from file but stores in memory as normal slice.
// Maybe we can implements a real `PackedVector` later which will be more complex but saves memory.
// C++ implementation: https://github.com/Telenav/osrm-backend/blob/6283c6074066f98e6d4a9f774f21ea45407c0d52/include/util/packed_vector.hpp#L92
package packed

import (
	"encoding/binary"
	"fmt"

	"github.com/golang/glog"

	"github.com/Telenav/osrm-backend/integration/osrmfiles/meta"
)

// Uint64Vector represents the `PackedOSMIDs`. It stores slice of OSMID(uint64), but raw data comes from packed data.
// C++ implementation: https://github.com/Telenav/osrm-backend/blob/6283c6074066f98e6d4a9f774f21ea45407c0d52/include/extractor/packed_osm_ids.hpp#L14
type Uint64Vector struct {
	NumOfElements meta.Num
	PackedMeta    meta.Num
	Values        []uint64

	bits       uint // how many bits per Value in packed data, [1,64]
	blockWords uint // number of words per block
	blockBytes uint // number of bytes per block

	// The `Write` will be called many times in `io.Copy` due to fixed buffer size.
	// Defaultly the copy buffer size is 32*1024 bytes, see line 391 `copyBuffer()` in https://golang.org/src/io/io.go.
	// Unfortunetly, we expect to parse the packed data which uses Bits*8 as a block.
	// We have to cache remain bytes since most of time 32*1024 / (Bits*8) != 0.
	// Maximum will store `Bits*8-1` bytes.
	unwrittenBytes []byte

	totalPackedBytes uint64
}

const (
	// https://github.com/Project-OSRM/osrm-backend/blob/9234b2ae76bdbbb91cbb51142bfc0ee1252c4abd/include/util/packed_vector.hpp#L99
	maxBits = 64

	// wordBytes/wordBits represents underlying storage bytes/bits. packed_vector always use uint64 as the storage.
	// - https://github.com/Telenav/osrm-backend/blob/6283c6074066f98e6d4a9f774f21ea45407c0d52/include/util/packed_vector.hpp#L94
	// - https://github.com/Telenav/osrm-backend/blob/6283c6074066f98e6d4a9f774f21ea45407c0d52/include/util/packed_vector.hpp#L102
	wordBytes = 8  // 8 bytes per uint64
	wordBits  = 64 // 64 bits per uint64
)

// NewUint64Vector creates a new packd Uint64Vector with fixed bits setting.
func NewUint64Vector(bits uint) Uint64Vector {
	if bits > maxBits || bits == 0 {
		glog.Fatalf("invalid bits %d, only allows [1,%d]", bits, maxBits)
	}

	u := Uint64Vector{
		bits: bits,
	}

	// number of words per block
	// https://github.com/Telenav/osrm-backend/blob/6283c6074066f98e6d4a9f774f21ea45407c0d52/include/util/packed_vector.hpp#L110
	u.blockWords = u.bits
	u.blockBytes = u.bits * wordBytes

	u.unwrittenBytes = make([]byte, 0, u.blockBytes)

	return u
}

func (u *Uint64Vector) Write(p []byte) (int, error) {

	var writeLen int
	writeP := p
	for {
		if len(u.unwrittenBytes) > 0 { // consume last unwritten bytes first
			retrieveBytes := int(u.blockBytes) - len(u.unwrittenBytes)
			if len(writeP) < retrieveBytes {
				u.unwrittenBytes = append(u.unwrittenBytes, writeP...)
				writeLen += len(writeP)
				break
			}

			u.unwrittenBytes = append(u.unwrittenBytes, writeP[:retrieveBytes]...)

			n, err := u.parseBlock(u.unwrittenBytes)
			if err != nil {
				return n, err
			}
			if n != int(u.blockBytes) {
				glog.Fatalf("expect parse block bytes %d but got %d", u.blockBytes, n)
			}
			u.unwrittenBytes = u.unwrittenBytes[:0]

			writeP = writeP[retrieveBytes:]
			writeLen += retrieveBytes
		}

		if len(writeP) < int(u.blockBytes) { // insufficient for a block
			u.unwrittenBytes = append(u.unwrittenBytes, writeP...)
			writeLen += len(writeP)
			break
		}

		n, err := u.parseBlock(writeP[:u.blockBytes])
		if err != nil {
			return writeLen, err
		}
		if n != int(u.blockBytes) {
			glog.Fatalf("expect parse block bytes %d but got %d", u.blockBytes, n)
		}
		writeP = writeP[u.blockBytes:]
		writeLen += n
	}

	u.totalPackedBytes += uint64(writeLen)
	return writeLen, nil
}

// parse a complete block
func (u *Uint64Vector) parseBlock(p []byte) (int, error) {
	if len(p) != int(u.blockBytes) {
		return 0, fmt.Errorf("invalid block bytes %d, expect %d", len(p), u.blockBytes)
	}

	var cachedLower uint64 // last remain bits, stores on lower
	var cachedBits uint    // should always < u.bits

	var writeLen int
	writeP := p
	for {
		if len(writeP) < wordBytes {
			break
		}

		word := uint64(binary.LittleEndian.Uint64(writeP))

		// construct first value: consume with last cached
		firstConsumeBits := u.bits - cachedBits
		upper := (word & (^(1 << firstConsumeBits))) << cachedBits
		newValue := upper | cachedLower
		u.Values = append(u.Values, newValue)
		word = word >> firstConsumeBits

		consumedBits := firstConsumeBits
		for wordBits-consumedBits >= u.bits { // more values may still available
			newValue = word & (^(1 << u.bits))
			u.Values = append(u.Values, newValue)
			word = word >> u.bits
			consumedBits += u.bits
		}
		cachedLower = word
		cachedBits = wordBits - consumedBits
		if cachedLower>>cachedBits != 0 {
			glog.Fatalf("expect cacheBits %d, but got cache 0x%x", cachedBits, cachedLower)
		}

		writeP = writeP[wordBytes:]
		writeLen += wordBytes

		if writeLen%int(u.blockBytes) == 0 && (cachedLower != 0 || cachedBits != 0) {
			// a full block has been prased, then nothing should be cached anymore
			glog.Fatalf("parsed len %d reached block, expect no cache anymore but got cache 0x%x cachedBits %d, block words %d, block bytes %d",
				writeLen, cachedLower, cachedBits, u.blockWords, u.blockBytes)
		}
	}
	return writeLen, nil
}

// Validate checks whether the Uint64Vector, which parsed from packed data, valid or not.
func (u *Uint64Vector) Validate() error {

	if len(u.unwrittenBytes) != 0 {
		return fmt.Errorf("expect consumed all packed data, but len(unwrittenBytes) %d", len(u.unwrittenBytes))
	}

	if uint64(u.NumOfElements) != uint64(len(u.Values)) {
		return fmt.Errorf("packed uint64 vector meta not match, number_of_elements in meta %d, but actual count %d", u.NumOfElements, len(u.Values))
	}
	if u.totalPackedBytes != uint64(u.PackedMeta*wordBytes) {
		return fmt.Errorf("packed uint64 vector packed meta not match, packed size in meta %d bytes %d, but actual bytes %d", u.PackedMeta, u.PackedMeta*wordBytes, u.totalPackedBytes)
	}

	totalElementsBits := uint64(u.NumOfElements) * uint64(u.bits)
	totalPackedBits := uint64(u.PackedMeta * wordBits)
	if totalElementsBits > totalPackedBits {
		return fmt.Errorf("packed uint64 vector bits not match, expect totalElementsBits %d <= totalPackedBits %d",
			totalElementsBits, totalPackedBits)
	}

	return nil
}

// Prune prunes values due to last possible half empty packed word.
func (u *Uint64Vector) Prune() error {
	if uint64(u.NumOfElements) == uint64(len(u.Values)) {
		return nil //nothing need to do
	}

	for i := u.NumOfElements; uint64(i) < uint64(len(u.Values)); i++ {
		// if u.Values more than u.NumOfElements, all these extra should be 0
		if u.Values[i] != 0 {
			return fmt.Errorf("expect extra Values[%d] == 0, but got %d", i, u.Values[i])
		}
	}
	u.Values = u.Values[:u.NumOfElements]
	return nil
}
