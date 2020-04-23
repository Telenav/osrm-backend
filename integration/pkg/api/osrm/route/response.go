package route

import (
	"encoding/json"
)

// Response represent OSRM api v1 route response.
type Response struct {
	Code        string      `json:"code"`
	Message     string      `json:"message,omitempty"`
	DataVersion string      `json:"data_version,omitempty"`
	Routes      []*Route    `json:"routes,omitempty"`
	Waypoints   []*Waypoint `json:"waypoints,omitempty"`
}

// Route represents a route through (potentially multiple) waypoints.
type Route struct {
	Distance   float64 `json:"distance"`
	Duration   float64 `json:"duration"`
	Geometry   string  `json:"geometry"`
	Weight     float64 `json:"weight"`
	WeightName string  `json:"weight_name"`
	Legs       []*Leg  `json:"legs,omitempty"`
}

// Leg represents a route between two waypoints.
type Leg struct {
	Distance   float64     `json:"distance"`
	Duration   float64     `json:"duration"`
	Weight     float64     `json:"weight"`
	Summary    string      `json:"summary"`
	Steps      []*Step     `json:"steps,omitempty"`
	Annotation *Annotation `json:"annotation,omitempty"`
}

// Step A step consists of a maneuver such as a turn or merge, followed by a distance of travel along a single way to the subsequent step.
type Step struct {
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
	Ways        []int64   `json:"ways,omitempty"` // NOT osrm original
	Weight      []float64 `json:"weight,omitempty"`
	Speed       []float64 `json:"speed,omitempty"`
	Metadata    *Metadata `json:"metadata,omitempty"`
}

// Metadata related to other annotations
type Metadata struct {
	DataSourceNames []string `json:"datasource_names,omitempty"`
}

// UnmarshalJSON unmarshal Annotation JSON.
// Note that it's a temp solution for the string type nodes.
func (a *Annotation) UnmarshalJSON(bs []byte) error {

	// A temporary type that uses json.Number instead of int64
	var tmp struct {
		Distance    []float64     `json:"distance,omitempty"`
		Duration    []float64     `json:"duration,omitempty"`
		DataSources []int         `json:"datasources,omitempty"`
		Nodes       []json.Number `json:"nodes,omitempty"`
		Ways        []int64       `json:"ways,omitempty"` // NOT osrm original
		Weight      []float64     `json:"weight,omitempty"`
		Speed       []float64     `json:"speed,omitempty"`
		Metadata    *Metadata     `json:"metadata,omitempty"`
	}

	// Unmarshal into the temporary
	if err := json.Unmarshal(bs, &tmp); err != nil {
		return err
	}

	// Convert each json.Number into an int64
	a.Nodes = make([]int64, len(tmp.Nodes))
	for i := range tmp.Nodes {
		v, err := tmp.Nodes[i].Int64()
		if err != nil {
			return err
		}
		a.Nodes[i] = v
	}

	// Save other fields
	a.Distance = tmp.Distance
	a.Duration = tmp.Duration
	a.DataSources = tmp.DataSources
	a.Ways = tmp.Ways
	a.Weight = tmp.Weight
	a.Speed = tmp.Speed
	a.Metadata = tmp.Metadata

	return nil
}
