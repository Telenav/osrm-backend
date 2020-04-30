package osrmtype

// Waypoint object used to describe waypoint used in route or table.
type Waypoint struct {
	Name     string     `json:"name"`
	Location [2]float64 `json:"location,omitempty"` // [longitude, latitude]
	Distance float64    `json:"distance"`
	Hint     string     `json:"hint"`
}
