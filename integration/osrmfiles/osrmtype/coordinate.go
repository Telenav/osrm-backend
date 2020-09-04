package osrmtype

import "encoding/binary"

// FixedLat defines FixedLatitude with COORDINATE_PRECISION = 1e6
// C++ Implementation: https://github.com/Telenav/osrm-backend/blob/6283c6074066f98e6d4a9f774f21ea45407c0d52/include/util/coordinate.hpp#L66
// C++ implementation toFixed()/toFloat() see https://github.com/Telenav/osrm-backend/blob/6283c6074066f98e6d4a9f774f21ea45407c0d52/include/util/coordinate.hpp#L88
type FixedLat int32

// FixedLon defines FixedLatitude with COORDINATE_PRECISION = 1e6
// C++ Implementation: https://github.com/Telenav/osrm-backend/blob/6283c6074066f98e6d4a9f774f21ea45407c0d52/include/util/coordinate.hpp#L66
// C++ implementation toFixed()/toFloat() see https://github.com/Telenav/osrm-backend/blob/6283c6074066f98e6d4a9f774f21ea45407c0d52/include/util/coordinate.hpp#L88
type FixedLon int32

// Coordinate represents Coordinate structure.
// C++ implementation: https://github.com/Telenav/osrm-backend/blob/6283c6074066f98e6d4a9f774f21ea45407c0d52/include/util/coordinate.hpp#L185
type Coordinate struct {
	FixedLon
	FixedLat
}

// Coordinates represents Coordinate vector structure.
type Coordinates []Coordinate

const (
	FixedLonBytes   = 4 // int32
	FixedLatBytes   = 4 // int32
	CoordinateBytes = FixedLonBytes + FixedLatBytes
)

func (c *Coordinates) Write(p []byte) (int, error) {

	var writeLen int
	writeP := p
	for {
		if len(writeP) < CoordinateBytes {
			break
		}

		var coordinate Coordinate
		coordinate.FixedLon = FixedLon(binary.LittleEndian.Uint32(writeP))
		coordinate.FixedLat = FixedLat(binary.LittleEndian.Uint32(writeP[FixedLonBytes:]))

		*c = append(*c, coordinate)

		writeP = writeP[CoordinateBytes:]
		writeLen += CoordinateBytes
	}

	return writeLen, nil
}
