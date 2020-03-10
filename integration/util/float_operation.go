package util

const epsilon float64 = 0.00000001

// FloatEquals compares two float64 numbers and returns true if they equal to each other
// with precision limit of 0.00000001
// More information about float numbers:
// Floating Point Numbers https://users.cs.fiu.edu/~downeyt/cop2400/float.htm
// What Every Computer Scientist Should Know About Floating-Point Arithmetic https://docs.oracle.com/cd/E19957-01/806-3568/ncg_goldberg.html
func FloatEquals(a, b float64) bool {
	if (a-b) < epsilon && (b-a) < epsilon {
		return true
	}
	return false
}
