package osrmv1

// RouteResponse represent OSRM api v1 route response.
type RouteResponse struct {
	Code        string  `json:"code"`
	Message     string  `json:"message"`
	DataVersion string  `json:"data_version"`
	Routes      []Route `json:"routes"`
}

// Route represents a route through (potentially multiple) waypoints.
type Route struct {
	Distance   float32    `json:"distance"`
	Duration   float32    `json:"duration"`
	Geometry   string     `json:"geometry"`
	Weight     float32    `json:"weight"`
	WeightName string     `json:"weight_name"`
	Legs       []RouteLeg `json:"legs"`
}

// RouteLeg represents a route between two waypoints.
type RouteLeg struct {
	Distance float32 `json:"distance"`
	Duration float32 `json:"duration"`
	Weight   float32 `json:"weight"`
	Summary  string  `json:"summary"`
}
