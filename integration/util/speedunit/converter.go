// Package speedunit provides speed unit converters.
package speedunit

const (
	mps2KPH = 3.6
)

// ConvertMPS2KPH converts speed from meter per second to kilometers per hour.
func ConvertMPS2KPH(speed float64) float64 {
	return speed * mps2KPH
}

// ConvertKPH2MPS converts speed from kilometers per hour to meter per second.
func ConvertKPH2MPS(speed float64) float64 {
	return speed / mps2KPH
}
