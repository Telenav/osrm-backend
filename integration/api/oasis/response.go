package oasis

import "github.com/Telenav/osrm-backend/integration/api/search/nearbychargestation"

// Response for oasis service
type Response struct {
	Code      string      `json:"code"`
	Message   string      `json:"message,omitempty"`
	Solutions []*Solution `json:"solutions,omitempty"`
}

// Solution contains recommended charge stations
type Solution struct {
	Distance       float64          `json:"distance"`
	Duration       float64          `json:"duration"`
	RemainingRage  float64          `json:"estimate_remaining_range"`
	Weight         float64          `json:"weight"`
	WeightName     string           `json:"weight_name"`
	ChargeStations []*ChargeStation `json:"charge_stations"`
}

// ChargeStation contains location, time and energy level, could be used as waypoints for routing request
type ChargeStation struct {
	Address     []*nearbychargestation.Address `json:"address"`
	WaitTime    float64                        `json:"wait_time"`
	ChargeTime  float64                        `json:"charge_time"`
	ChargeRange float64                        `json:"charge_range"`
	DetailURL   string                         `json:"detail_url"`
}
