package nametable

// sum2bits calculates summation of 16 2-bit values using SWAR
// https://github.com/Telenav/open-source-spec/blob/master/osrm/doc/sum2bits-swar.md
// C++ implementation: https://github.com/Telenav/osrm-backend/blob/b24b8a085dc10bea279ffb352049330beae23791/include/util/indexed_data.hpp#L94
func sum2Bits(value uint32) uint32 {
	value = (value >> 2 & 0x33333333) + (value & 0x33333333)
	value = (value >> 4 & 0x0f0f0f0f) + (value & 0x0f0f0f0f)
	value = (value >> 8 & 0x00ff00ff) + (value & 0x00ff00ff)
	return (value >> 16 & 0x0000ffff) + (value & 0x0000ffff)
}
