package osrmtype

// TravelMode defines travel mode.
// C++ Implementation: https://github.com/Telenav/osrm-backend/blob/6283c6074066f98e6d4a9f774f21ea45407c0d52/include/extractor/travel_mode.hpp#L41
type TravelMode uint8

// C++ Implementation: https://github.com/Telenav/osrm-backend/blob/6283c6074066f98e6d4a9f774f21ea45407c0d52/include/extractor/travel_mode.hpp#L43
const (
	TravelModeInaccessible TravelMode = 0
	TravelModeDriving      TravelMode = 1
	TravelModeCycling      TravelMode = 2
	TravelModeWalking      TravelMode = 3
	TravelModeFerry        TravelMode = 4
	TravelModeTrain        TravelMode = 5
	TravelModePushingBike  TravelMode = 6
	// FIXME only for testbot.lua
	TravelModeStepsUp   = 8
	TravelModeStepsDown = 9
	TravelModeRiverUp   = 10
	TravelModeRiverDown = 11
	TravelModeRoute     = 12
)
