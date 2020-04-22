// Package unidbpatch implements some utility to optimize handling for UniDB PBF.
package unidbpatch

const (
	unidbWayIDSuffix        = 100
	unidbWayIDSuffixDivisor = int64(1e3)
)

// IsValidWay returns whether the wayID from UniDB is valid or not.
// This validation is for UniDB PBF only.
// UniDB PBF append "100" for all valid ways.
func IsValidWay(wayID int64) bool {
	if wayID < 0 {
		wayID = -wayID
	}

	if (wayID-unidbWayIDSuffix)%unidbWayIDSuffixDivisor == 0 {
		return true
	}
	return false
}

// TrimValidWayIDSuffix trims "100" suffix for UniDB wayID.
func TrimValidWayIDSuffix(wayID int64) int64 {
	if IsValidWay(wayID) {
		return wayID / unidbWayIDSuffixDivisor
	}
	return wayID
}
