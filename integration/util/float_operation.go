package util

const (
	epsilon4Float64 float64 = 0.00000001
	epsilon4Float32 float64 = 0.00001
)

// Float64Equal compares two float64 numbers and returns true if they equal to each other
// with precision limit of 0.00000001
// More information about float numbers:
// Floating Point Numbers https://users.cs.fiu.edu/~downeyt/cop2400/float.htm
// What Every Computer Scientist Should Know About Floating-Point Arithmetic https://docs.oracle.com/cd/E19957-01/806-3568/ncg_goldberg.html
func Float64Equal(a, b float64) bool {
	return floatEqual(a, b, epsilon4Float64)
}

// Float32Equal compares two float32 numbers and returns true if they equal to each other
// with precision limit of 0.00001
func Float32Equal(a, b float32) bool {
	return floatEqual(float64(a), float64(b), epsilon4Float32)
}

func floatEqual(a, b, epsilon float64) bool {
	if (a-b) < epsilon4Float64 && (b-a) < epsilon4Float64 {
		return true
	}
	return false
}
