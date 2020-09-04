package conditional

import (
	"encoding/binary"
	"fmt"

	"github.com/Telenav/osrm-backend/integration/osrmfiles/meta"
)

// OpeningHours represent "opening hours" format http://wiki.openstreetmap.org/wiki/Key:opening_hours
// C++ implementation: https://github.com/Telenav/osrm-backend/blob/f45ab75cf9eb57cb9c857ea564beb95be0523968/include/util/opening_hours.hpp#L24
type OpeningHours struct {
	Modifier uint32 // enum, see https://github.com/Telenav/osrm-backend/blob/f45ab75cf9eb57cb9c857ea564beb95be0523968/include/util/opening_hours.hpp#L26

	timesCount meta.Num
	Times      []TimeSpan

	weekdaysCount meta.Num
	Weekdays      []WeekdayRange

	monthdaysCount meta.Num
	Monthdays      []MonthdayRange
}

func newOpenningHours() OpeningHours {
	return OpeningHours{
		Times:     []TimeSpan{},
		Weekdays:  []WeekdayRange{},
		Monthdays: []MonthdayRange{},
	}
}

// TimeEvent represents `enum Event: unsigned char` type in C++ definition.
// C++ implementation: https://github.com/Telenav/osrm-backend/blob/f45ab75cf9eb57cb9c857ea564beb95be0523968/include/util/opening_hours.hpp#L37
type TimeEvent uint8

// TimeSpan represents `struct TimeSpan` in C++ definition.
// C++ implementation: https://github.com/Telenav/osrm-backend/blob/f45ab75cf9eb57cb9c857ea564beb95be0523968/include/util/opening_hours.hpp#L59
type TimeSpan struct {
	FromEvent   TimeEvent
	FromMinutes int32 // https://github.com/Telenav/osrm-backend/blob/f45ab75cf9eb57cb9c857ea564beb95be0523968/include/util/opening_hours.hpp#L48

	ToEvent   TimeEvent
	ToMinutes int32 // https://github.com/Telenav/osrm-backend/blob/f45ab75cf9eb57cb9c857ea564beb95be0523968/include/util/opening_hours.hpp#L48
}

// WeekdayRange represents `struct WeekdayRange` in C++ definition.
// C++ implementation: https://github.com/Telenav/osrm-backend/blob/f45ab75cf9eb57cb9c857ea564beb95be0523968/include/util/opening_hours.hpp#L92
type WeekdayRange struct {
	Weekydays         int
	OvernightWeekdays int
}

// Monthday represents `struct Monthday` in C++ definition.
// C++ implementation: https://github.com/Telenav/osrm-backend/blob/f45ab75cf9eb57cb9c857ea564beb95be0523968/include/util/opening_hours.hpp#L113
type Monthday struct {
	Year  int
	Month int8
	Day   int8
}

// MonthdayRange represents `struct MonthdayRange` in C++ definition.
// C++ implementation: https://github.com/Telenav/osrm-backend/blob/f45ab75cf9eb57cb9c857ea564beb95be0523968/include/util/opening_hours.hpp#L129
type MonthdayRange struct {
	From Monthday
	To   Monthday
}

const (
	modifierBytes = 4 // 4 bytes enum

	timeEventBytes    = 1
	timeReservedBytes = 3 // 3 bytes reserved between event ~ minutes
	timeMinutesBytes  = 4
	timeSpanBytes     = (timeEventBytes + timeEventBytes + timeMinutesBytes) * 2 // *2: from, to

	weekdayBytes          = 4
	overnightWeekdayBytes = 4
	weekdaysRangeBytes    = weekdayBytes + overnightWeekdayBytes

	monthdayYearBytes     = 4
	monthdayMonthBytes    = 1
	monthdayDayBytes      = 1
	monthdayReservedBytes = 2 // 2 bytes reserved at the end of Monthday
	monthdayBytes         = monthdayYearBytes + monthdayMonthBytes + monthdayDayBytes + monthdayReservedBytes
	monthdaysRangeBytes   = monthdayBytes * 2 // *2: from, to
)

func (o *OpeningHours) Write(p []byte) (int, error) {
	if len(p) < modifierBytes {
		return 0, fmt.Errorf("%T byte array len %d insufficient, requires at least %d bytes", o, len(p), modifierBytes)
	}

	var writeLen int
	writeP := p

	// write modifier
	o.Modifier = binary.LittleEndian.Uint32(writeP)
	writeP = writeP[modifierBytes:]
	writeLen += modifierBytes

	n, err := o.timesCount.Write(writeP)
	if err != nil {
		return 0, err
	}
	writeP = writeP[n:]
	writeLen += n

	// write times
	for i := 0; i < int(o.timesCount); i++ {
		ts := TimeSpan{}
		n, err := ts.Write(writeP)
		if err != nil {
			return 0, err
		}
		o.Times = append(o.Times, ts)

		writeP = writeP[n:]
		writeLen += n
	}

	n, err = o.weekdaysCount.Write(writeP)
	if err != nil {
		return 0, err
	}
	writeP = writeP[n:]
	writeLen += n

	// write weekdays
	for i := 0; i < int(o.weekdaysCount); i++ {
		wr := WeekdayRange{}
		n, err := wr.Write(writeP)
		if err != nil {
			return 0, err
		}
		o.Weekdays = append(o.Weekdays, wr)

		writeP = writeP[n:]
		writeLen += n
	}

	n, err = o.monthdaysCount.Write(writeP)
	if err != nil {
		return 0, err
	}
	writeP = writeP[n:]
	writeLen += n

	// write monthdays
	for i := 0; i < int(o.monthdaysCount); i++ {
		mr := MonthdayRange{}
		n, err := mr.Write(writeP)
		if err != nil {
			return 0, err
		}
		o.Monthdays = append(o.Monthdays, mr)

		writeP = writeP[n:]
		writeLen += n
	}

	return writeLen, nil
}

func (t *TimeSpan) Write(p []byte) (int, error) {
	if len(p) < timeSpanBytes {
		return 0, fmt.Errorf("%T byte array len %d insufficient, requires at least %d bytes", t, len(p), timeSpanBytes)
	}

	var writeLen int
	writeP := p

	// from
	t.FromEvent = TimeEvent(writeP[0])
	writeP = writeP[timeEventBytes:]
	writeLen += timeEventBytes

	writeP = writeP[timeReservedBytes:]
	writeLen += timeReservedBytes

	t.FromMinutes = int32(binary.LittleEndian.Uint32(writeP))
	writeP = writeP[timeMinutesBytes:]
	writeLen += timeMinutesBytes

	// to
	t.ToEvent = TimeEvent(writeP[0])
	writeP = writeP[timeEventBytes:]
	writeLen += timeEventBytes

	writeP = writeP[timeReservedBytes:]
	writeLen += timeReservedBytes

	t.ToMinutes = int32(binary.LittleEndian.Uint32(writeP))
	writeP = writeP[timeMinutesBytes:]
	writeLen += timeMinutesBytes

	return writeLen, nil
}

func (w *WeekdayRange) Write(p []byte) (int, error) {
	if len(p) < weekdaysRangeBytes {
		return 0, fmt.Errorf("%T byte array len %d insufficient, requires at least %d bytes", w, len(p), weekdaysRangeBytes)
	}

	var writeLen int
	writeP := p

	w.Weekydays = int(binary.LittleEndian.Uint32(writeP))
	writeP = writeP[weekdayBytes:]
	writeLen += weekdayBytes

	w.OvernightWeekdays = int(binary.LittleEndian.Uint32(writeP))
	writeP = writeP[overnightWeekdayBytes:]
	writeLen += overnightWeekdayBytes

	return writeLen, nil
}

func (m *Monthday) Write(p []byte) (int, error) {
	if len(p) < monthdayBytes {
		return 0, fmt.Errorf("%T byte array len %d insufficient, requires at least %d bytes", m, len(p), monthdayBytes)
	}

	var writeLen int
	writeP := p

	m.Year = int(binary.LittleEndian.Uint32(writeP))
	writeP = writeP[monthdayYearBytes:]
	writeLen += monthdayYearBytes

	m.Month = int8(writeP[0])
	writeP = writeP[monthdayMonthBytes:]
	writeLen += monthdayMonthBytes

	m.Day = int8(writeP[0])
	writeP = writeP[monthdayDayBytes:]
	writeLen += monthdayDayBytes

	writeP = writeP[monthdayReservedBytes:]
	writeLen += monthdayReservedBytes

	return writeLen, nil

}

func (m *MonthdayRange) Write(p []byte) (int, error) {
	if len(p) < monthdaysRangeBytes {
		return 0, fmt.Errorf("%T byte array len %d insufficient, requires at least %d bytes", m, len(p), monthdaysRangeBytes)
	}

	var writeLen int
	writeP := p

	n, err := m.From.Write(writeP)
	if err != nil {
		return 0, err
	}
	writeP = writeP[n:]
	writeLen += n

	n, err = m.To.Write(writeP)
	if err != nil {
		return 0, err
	}
	writeP = writeP[n:]
	writeLen += n

	return writeLen, nil
}
