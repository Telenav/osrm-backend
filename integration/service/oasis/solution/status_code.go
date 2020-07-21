package solution

const (
	StatusOrigAndDestIsNotReachable = iota
	StatusNoNeedCharge
	StatusChargeForSingleTime
	StatusChargeForMultipleTime
	StatusFailedToCalculateRoute
	StatusFailedToGenerateChargeResult
	StatusIncorrectRequest
)

var statusText = map[int]string{
	StatusOrigAndDestIsNotReachable:    "Orig could not reach destination with current energy capacity",
	StatusNoNeedCharge:                 "Orig could reach destination with current energy capacity",
	StatusChargeForSingleTime:          "Orig could reach destination with single charge",
	StatusChargeForMultipleTime:        "Orig could reach destination with multiple charge",
	StatusFailedToCalculateRoute:       "Failed to calculate route between Orig and destination",
	StatusFailedToGenerateChargeResult: "Failed to generate charge result",
	StatusIncorrectRequest:             "Incorrect request parameters",
}

// StatusText returns a text for the solution status code. It returns the empty
// string if the code is unknown.
func StatusText(code int) string {
	return statusText[code]
}
