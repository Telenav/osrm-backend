package nearbychargestation

// Response for search service
type Response struct {
	Status       Status    `json:"status"`
	ResponseTime int       `json:"response_time"`
	Results      []*Result `json:"results"`
}

// Status for search response
type Status struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Result for search places
type Result struct {
	ID        string `json:"id"`
	Place     Place  `json:"place"`
	Distance  int    `json:"distance"`
	Facets    Facet  `json:"facets"`
	DetailURL string `json:"detail_url"`
}

// Place contains name, phone, address information
type Place struct {
	Address []*Address `json:"address"`
}

// Address contains coordinate information
type Address struct {
	GeoCoordinate  Coordinate    `json:"geo_coordinates"`
	NavCoordinates []*Coordinate `json:"nav_coordinates"`
}

// Coordinate for specific location
type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Facet contains open hours, charge station status, nearby information
type Facet struct {
	EVConnectors EVConnector `json:"ev_connectors"`
}

// EVConnector contains charge station status
type EVConnector struct {
	TotalNumber     int               `json:"total_number"`
	ConnectorCounts []*ConnectorCount `json:"connector_counts"`
}

// ConnectorCount contains charge level and related count
type ConnectorCount struct {
	Level int `json:"level"`
	Total int `total:"total"`
}
