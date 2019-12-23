package osrm

// RouteResponse represent OSRM api v1 route response.
type RouteResponse struct {
	Code        string      `json:"code"`
	Message     string      `json:"message,omitempty"`
	DataVersion string      `json:"data_version,omitempty"`
	Routes      []*Route    `json:"routes,omitempty"`
	Waypoints   []*Waypoint `json:"waypoints,omitempty"`
}

// Route represents a route through (potentially multiple) waypoints.
type Route struct {
	Distance   float64     `json:"distance"`
	Duration   float64     `json:"duration"`
	Geometry   string      `json:"geometry"`
	Weight     float64     `json:"weight"`
	WeightName string      `json:"weight_name"`
	Legs       []*RouteLeg `json:"legs,omitempty"`
}

// RouteLeg represents a route between two waypoints.
type RouteLeg struct {
	Distance   float64      `json:"distance"`
	Duration   float64      `json:"duration"`
	Weight     float64      `json:"weight"`
	Summary    string       `json:"summary"`
	Steps      []*RouteStep `json:"steps,omitempty"`
	Annotation *Annotation  `json:"annotation,omitempty"`
}

// RouteStep A step consists of a maneuver such as a turn or merge, followed by a distance of travel along a single way to the subsequent step.
type RouteStep struct {
	Distance float64 `json:"distance"`
	Duration float64 `json:"duration"`
	Geometry string  `json:"geometry"`
	Weight   float64 `json:"weight"`
	Name     string  `json:"name"`
	//TODO: others
}

// Waypoint object used to describe waypoint on a route.
type Waypoint struct {
	Name     string     `json:"name"`
	Location [2]float64 `json:"location,omitempty"` // [longitude, latitude]
	Distance float64    `json:"distance"`
	Hint     string     `json:"hint"`
}

// Annotation of the whole route leg with fine-grained information about each segment or node id.
type Annotation struct {
	Distance    []float64 `json:"distance,omitempty"`
	Duration    []float64 `json:"duration,omitempty"`
	DataSources []int     `json:"datasources,omitempty"`
	Nodes       []int64   `json:"nodes,omitempty"`
	Weight      []float64 `json:"weight,omitempty"`
	Speed       []float64 `json:"speed,omitempty"`
	Metadata    *Metadata `json:"metadata,omitempty"`
}

// Metadata related to other annotations
type Metadata struct {
	DataSourceNames []string `json:"datasource_names,omitempty"`
}
