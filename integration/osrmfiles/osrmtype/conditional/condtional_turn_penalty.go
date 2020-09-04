// Package conditional implements the C++ `struct ConditionalTurnPenalty` that used to store conditonal restrictions which has been compiled by `osrm-extract`.
// C++ implementation: https://github.com/Telenav/osrm-backend/blob/f45ab75cf9eb57cb9c857ea564beb95be0523968/include/extractor/conditional_turn_penalty.hpp#L17
package conditional

import (
	"encoding/binary"
	"fmt"

	"github.com/Telenav/osrm-backend/integration/osrmfiles/meta"
	"github.com/Telenav/osrm-backend/integration/osrmfiles/osrmtype"
	"github.com/golang/glog"
)

// TurnPenalties represents vector of TurnPenalty.
type TurnPenalties struct {
	TurnPenalties []TurnPenalty

	turnPenaltiesCount meta.Num

	// The `Write` will be called many times in `io.Copy` due to fixed buffer size.
	// Defaultly the copy buffer size is 32*1024 bytes, see line 391 `copyBuffer()` in https://golang.org/src/io/io.go.
	// So we don't know whether we have sufficient data in parsing.
	// We have to cache remain bytes if any fail occurs, and try next time to see whether data is sufficient.
	unwritten []byte

	totalParsedBytes uint64
}

// New creates conditional turn penalties.
func New() TurnPenalties {
	return TurnPenalties{
		TurnPenalties: []TurnPenalty{},
		unwritten:     []byte{},
	}
}

// TurnPenalty represents the C++ `struct ConditionalTurnPenalty`.
// C++ implementation: https://github.com/Telenav/osrm-backend/blob/f45ab75cf9eb57cb9c857ea564beb95be0523968/include/extractor/conditional_turn_penalty.hpp#L17
type TurnPenalty struct {
	TurnOffset uint64
	osrmtype.Coordinate

	conditionsCount meta.Num

	Conditions []OpeningHours
}

func newTurnPenalty() TurnPenalty {
	return TurnPenalty{
		Conditions: []OpeningHours{},
	}
}

// Validate validates whether expected bytes were processed by parsing successfully.
func (t TurnPenalties) Validate(expectedBytes uint64) error {
	if expectedBytes != t.totalParsedBytes {
		return fmt.Errorf("bytes not match, expect %d but actually %d", expectedBytes, t.totalParsedBytes)
	}
	if len(t.unwritten) > 0 {
		return fmt.Errorf("%d bytes have not been parsed", len(t.unwritten))
	}

	return nil
}

func (t *TurnPenalties) Write(p []byte) (int, error) {

	writeLen := len(p) // always return len(p) since all data will be cached even they may have not been parsed

	t.unwritten = append(t.unwritten, p...)
	writeP := t.unwritten

	if t.totalParsedBytes == 0 { // parse count of conditional turn penalties at the beginning
		n, err := t.turnPenaltiesCount.Write(writeP)
		if err != nil {
			glog.Warning(err)
			return writeLen, nil // cached all unwritten data
		}
		writeP = writeP[n:]
		t.totalParsedBytes += uint64(n)
	}

	for {
		if len(writeP) == 0 {
			break
		}

		if len(t.TurnPenalties) == int(t.turnPenaltiesCount) { // all data has been parsed
			break
		}

		// write each conditional turn penalty
		// the data will be cached if it fails, because it mostly caused by insufficient data and will be succeed in next time.
		tp := newTurnPenalty()
		n, err := tp.Write(writeP)
		if err != nil {
			glog.V(2).Infof("New turn penalty error: \"%v\" due to incomplete go-lang IO buffer, will continue Write in next round.", err)
			break
		}
		t.TurnPenalties = append(t.TurnPenalties, tp)

		writeP = writeP[n:]
		t.totalParsedBytes += uint64(n)
	}

	t.unwritten = writeP
	return writeLen, nil
}

const (
	turnOffsetBytes = 8 // uint64
)

func (t *TurnPenalty) Write(p []byte) (int, error) {
	if len(p) < turnOffsetBytes+osrmtype.CoordinateBytes {
		return 0, fmt.Errorf("%T byte array len %d insufficient, requires at least %d bytes", t, len(p), turnOffsetBytes+osrmtype.CoordinateBytes)
	}

	var writeLen int
	writeP := p

	t.TurnOffset = binary.LittleEndian.Uint64(writeP)
	writeP = writeP[turnOffsetBytes:]
	writeLen += turnOffsetBytes

	// the order is different with normal coordinate serialization, i.e., "Lat Lon" vs. "Lon Lat"
	t.Coordinate.FixedLat = osrmtype.FixedLat(binary.LittleEndian.Uint32(writeP))
	t.Coordinate.FixedLon = osrmtype.FixedLon(binary.LittleEndian.Uint32(writeP[osrmtype.FixedLatBytes:]))
	writeP = writeP[osrmtype.CoordinateBytes:]
	writeLen += osrmtype.CoordinateBytes

	n, err := t.conditionsCount.Write(writeP)
	if err != nil {
		return 0, err
	}
	writeP = writeP[n:]
	writeLen += n

	// write opening hours
	for i := 0; i < int(t.conditionsCount); i++ {
		oh := newOpenningHours()
		n, err := oh.Write(writeP)
		if err != nil {
			return 0, err
		}
		t.Conditions = append(t.Conditions, oh)

		writeP = writeP[n:]
		writeLen += n
	}

	return writeLen, nil
}
