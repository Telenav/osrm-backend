package solution

const (
	StatusOrigAndDestIsNotReachable    = 1 // StatusOrigAndDestIsNotReachable means orig could not reach destination with current energy capacity
	StatusNoNeedCharge                 = 2 // StatusNoNeedCharge means orig could reach destination with current energy capacity
	StatusChargeForSingleTime          = 3 // StatusChargeForSingleTime means orig could reach destination with single charge
	StatusChargeForMultipleTime        = 4 // StatusChargeForMultipleTime means orig could reach destination with multiple charge
	StatusFailedToCalculateRoute       = 5 // StatusFailedToCalculateRoute means failed to calculate route between Orig and destination
	StatusFailedToGenerateChargeResult = 6 // StatusFailedToGenerateChargeResult means failed to generate charge result
	StatusIncorrectRequest             = 7 // StatusIncorrectRequest means incorrect request parameters
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
