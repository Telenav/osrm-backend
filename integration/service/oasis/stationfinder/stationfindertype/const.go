package stationfindertype

import (
	"math"

	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/common"
)

// todo @codebear801 move this file to internal/common

// OrigLocationIDStr defines name for orig
const OrigLocationIDStr string = "orig_location"
const OrigLocationID common.PlaceID = math.MaxInt64 - 1

// DestLocationIDStr defines name for dest
const DestLocationIDStr string = "dest_location"
const DestLocationID common.PlaceID = math.MaxInt64 - 2

// InvalidPlaceID defines name for InvalidPlaceID
const InvalidPlaceIDStr = "invalid_place_id"
const InvalidPlaceID = math.MaxInt64
