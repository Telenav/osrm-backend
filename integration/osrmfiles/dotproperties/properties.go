// Package dotproperties implements parser for `.osrm.properties` file.
// The `struct ProfileProperties` has been written into `.osrm.properties` file by take its address and `sizeof()` its size directly,
// which strongly related on C++ memory alignment on platform.
// In this Golang implementation, we assume that on 64 bits machine, which has minimum alignment 4 bytes, and struct should be n*8 bytes.
// It may not work if out of this memory alignment convention. E.g. it doesn't support 32 bits machine.
// Be aware of any change of `struct ProfileProperties` and memory alignment of your compiler/platform if anything goes wrong.
package dotproperties

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
)

// Properties represents the LUA profile properties that affect run time queries.
// C++ implementation: https://github.com/Telenav/osrm-backend/blob/8e094d1b90a31c090cdafa9854985343ea1e15c3/include/extractor/profile_properties.hpp#L23
// Be aware that the parsing should be strictly follow the memory layout of the `struct ProfileProperties` in C++ implementation,
// which means the parsing code should be changed also if anything changed in the C++ `struct ProfileProperties`,
// otherwise error will occrred.
// Also, the parsing code only considers 64 bits machine, DOESN'T support 32 bits machine.
type Properties struct {
	TrafficSignalPenalty   int32 // NOTE that it's different from `traffic_light_penality`, actually it's been DEPRECATED already.
	UTurnPenalty           int32
	MaxSpeedForMapMatching float64

	// 4 bool together takes 4 bytes
	ContinueStraightAtWaypoint bool
	UseTurnRestrictions        bool
	LeftHandDriving            bool // DEPRECATED, see https://github.com/Telenav/osrm-backend/blob/8e094d1b90a31c090cdafa9854985343ea1e15c3/include/extractor/profile_properties.hpp#L129
	FallbackToDuration         bool

	WeightName        string
	ClassNames        []string
	ExcludableClasses []uint8

	// uint + 2 bool together takes 8 bytes on 64 bits machine
	WeightPrecision         uint
	ForceSplitEdges         bool
	CallTaglessNodeFunction bool
}

const (
	propertiesBytes = 2344 // the C++ implementation layout actual size

	// https://github.com/Telenav/osrm-backend/blob/8e094d1b90a31c090cdafa9854985343ea1e15c3/include/extractor/profile_properties.hpp#L132
	weightNameBytes       = 256
	classNameBytes        = 256
	classNamesCount       = 8
	excluableClassesCount = 8

	lastPaddingBytes = 4 // last padding bytes due to C++ memory alignment
)

func (pr *Properties) Write(p []byte) (int, error) {
	if len(p) != propertiesBytes {
		return 0, fmt.Errorf("lua properties should always has %d bytes on 64-bit machine. Check whether anything changed on `struct ProfileProperties`, or memory alignment on machine", propertiesBytes)
	}

	var writeLen int
	writeP := p

	// first 16 bytes
	pr.TrafficSignalPenalty = int32(binary.LittleEndian.Uint32(writeP))
	pr.UTurnPenalty = int32(binary.LittleEndian.Uint32(writeP[4:]))
	pr.MaxSpeedForMapMatching = math.Float64frombits(binary.LittleEndian.Uint64(writeP[8:]))
	writeP = writeP[16:]
	writeLen += 16

	// 4 bool/4 bytes
	if writeP[0] > 0 {
		pr.ContinueStraightAtWaypoint = true
	}
	if writeP[1] > 0 {
		pr.UseTurnRestrictions = true
	}
	if writeP[2] > 0 {
		pr.LeftHandDriving = true
	}
	if writeP[3] > 0 {
		pr.FallbackToDuration = true
	}
	writeP = writeP[4:]
	writeLen += 4

	// names
	weightNameValidBytes := bytes.IndexByte(writeP[:weightNameBytes], 0)
	if weightNameValidBytes > 0 {
		pr.WeightName = string(writeP[:weightNameValidBytes])
	} else {
		pr.WeightName = string(writeP[:weightNameBytes])
	}
	writeP = writeP[weightNameBytes:]
	writeLen += weightNameBytes
	for i := 0; i < classNamesCount; i++ {
		classNameValidBytes := bytes.IndexByte(writeP[:classNameBytes], 0)
		if classNameValidBytes > 0 {
			pr.ClassNames = append(pr.ClassNames, string(writeP[:classNameValidBytes]))
		} else if classNameValidBytes == 0 {
			pr.ClassNames = append(pr.ClassNames, "")
		} else {
			pr.ClassNames = append(pr.ClassNames, string(writeP[:classNameBytes]))
		}
		writeP = writeP[classNameBytes:]
		writeLen += classNameBytes
	}
	for i := 0; i < excluableClassesCount; i++ {
		pr.ExcludableClasses = append(pr.ExcludableClasses, uint8(writeP[i]))
	}
	writeP = writeP[excluableClassesCount:]
	writeLen += excluableClassesCount

	// last 8 bytes
	pr.WeightPrecision = uint(binary.LittleEndian.Uint32(writeP))
	if writeP[4] > 0 {
		pr.ForceSplitEdges = true
	}
	if writeP[5] > 0 {
		pr.CallTaglessNodeFunction = true
	}
	writeP = writeP[8:]
	writeLen += 8

	// extra 4 bytes for memory align on 64 bits machine
	writeP = writeP[lastPaddingBytes:]
	writeLen += lastPaddingBytes

	if writeLen != propertiesBytes {
		return writeLen, fmt.Errorf("expect write %d bytes but actually %d bytes. Check whether anything changed on `struct ProfileProperties`, or memory alignment on machine", propertiesBytes, writeLen)
	}
	return writeLen, nil
}
