package poiloader

import "encoding/json"

// Element represent information loaded for point record
type Element struct {
	ID                int64       `json:"id"`
	VendorCode        VendorCode  `json:"vendor_code"`
	VendorPoiID       string      `json:"vendor_poi_id"`
	Lat               float64     `json:"lat"`
	Lon               float64     `json:"lon"`
	NavLat            float64     `json:"nav_lat"`
	NavLon            float64     `json:"nav_lon"`
	DisLat            float64     `json:"dis_lat"`
	DisLon            float64     `json:"dis_lon"`
	MapLinkID         int64       `json:"map_link_id"`
	SideOfStreet      string      `json:"side_of_street"`
	Country           Country     `json:"country"`
	SpaceID           SpaceID     `json:"space_id"`
	AirportCode       string      `json:"airport_code"`
	IsNational        bool        `json:"is_national"`
	IsStateImportance bool        `json:"is_state_importance"`
	IsCityImportance  bool        `json:"is_city_importance"`
	Fax               string      `json:"fax"`
	Phone             json.Number `json:"phone"`
	EncodedPhone      string      `json:"encoded_phone"`
	Email             string      `json:"email"`
	WebURL            string      `json:"web_url"`
	CategoryIDGather  json.Number `json:"category_id_gather"`
	ChainGather       string      `json:"chain_gather"`
	RawCategoryGather string      `json:"raw_category_gather"`
	ChildGather       string      `json:"child_gather"`
	ParentGather      string      `json:"parent_gather"`
	Hilbert           float64     `json:"hilbert"`
	Amenity           Amenity     `json:"amenity"`
}

// Amenity is defined by OSM format: https://wiki.openstreetmap.org/wiki/Tag:amenity%3Dcharging_station
type Amenity string

const (
	ChargingStation Amenity = "charging_station"
)

// Country records IOS code
type Country string

const (
	Usa Country = "USA"
)

// SpaceID is used to define a code for political defined administrative area
type SpaceID string

const (
	UsaCA SpaceID = "USA_CA"
)

// VendorCode defines data source provider
type VendorCode string

const (
	Noel VendorCode = "NOEL"
	Noft VendorCode = "NOFT"
	Nolp VendorCode = "NOLP"
)

// CategoryIDGather records aggregate data for certain type of point
type CategoryIDGather struct {
	Integer *int64
	String  *string
}
