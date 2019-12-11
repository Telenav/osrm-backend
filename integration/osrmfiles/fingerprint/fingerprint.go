// Package fingerprint implements OSRM fingerprint structure of `.osrm` and `.osrm.xxx` files.
// C++ implementation see https://github.com/Telenav/osrm-backend/blob/master-telenav/include/util/fingerprint.hpp
package fingerprint

import (
	"fmt"
)

// Fingerprint represents OSRM fingerprint in `.osrm` files.
type Fingerprint struct {
	magicNumber  [4]byte
	majorVersion uint8
	minorVersion uint8
	patchVersion uint8
	checksum     uint8 // CRC8 of the previous bytes to ensure the fingerprint is not damaged
}

func (f *Fingerprint) Write(p []byte) (n int, err error) {
	if len(p) < 8 {
		return 0, fmt.Errorf("less than 8 bytes")
	}

	writeLen := 0
	for i := 0; i < len(f.magicNumber); i++ {
		f.magicNumber[i] = p[i]
	}
	writeLen += len(f.magicNumber)

	f.majorVersion = uint8(p[writeLen])
	writeLen++
	f.minorVersion = uint8(p[writeLen])
	writeLen++
	f.patchVersion = uint8(p[writeLen])
	writeLen++
	f.checksum = uint8(p[writeLen])
	writeLen++

	return writeLen, nil
}

// IsValid verifies that the fingerprint has the expected magic number, and the checksum is correct.
func (f *Fingerprint) IsValid() bool {

	if !isValidMagicNumber(f.magicNumber) {
		return false
	}

	// TODO: CRC8 checksum verify

	return true
}

func (f *Fingerprint) String() string {
	return fmt.Sprintf("%c%c%c%c v%d.%d.%d",
		f.magicNumber[0], f.magicNumber[1], f.magicNumber[2], f.magicNumber[3],
		f.majorVersion, f.minorVersion, f.patchVersion)
}

func isValidMagicNumber(m [4]byte) bool {

	predefinedMagicNumber := []byte{'O', 'S', 'R', 'N'}

	if len(m) != len(predefinedMagicNumber) {
		return false
	}

	for i := 0; i < len(predefinedMagicNumber); i++ {
		if m[i] != predefinedMagicNumber[i] {
			return false
		}
	}

	return true
}
