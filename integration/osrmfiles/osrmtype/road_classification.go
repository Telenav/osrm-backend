package osrmtype

import (
	"fmt"
)

// RoadClassification describing the class of the road.
// C++ implementation: https://github.com/Telenav/osrm-backend/blob/6283c6074066f98e6d4a9f774f21ea45407c0d52/include/extractor/road_classification.hpp#L56
type RoadClassification struct {
	MotowayClass      bool  // 1 bit in .osrm file
	LinkClass         bool  // 1 bit in .osrm file
	MaybeIgnored      bool  // 1 bit in .osrm file
	RoadPriorityClass uint8 // 5 bits in .osrm file
	NumberOfLanes     uint8 // 8 bits in .osrm file
}

// Priorities are used to distinguish between how likely a turn is in comparison to a different
// road.
// C++ Implementation: https://github.com/Telenav/osrm-backend/blob/6283c6074066f98e6d4a9f774f21ea45407c0d52/include/extractor/road_classification.hpp#L21
const (
	// Top priority Road
	RoadPriorityClassMotoway     uint8 = 0
	RoadPriorityClassMotowayLink uint8 = 1
	// Second highest priority
	RoadPriorityClassTrunk     uint8 = 2
	RoadPriorityClassTrunkLink uint8 = 3
	// Main roads
	RoadPriorityClassPrimary       uint8 = 4
	RoadPriorityClassPrimaryLink   uint8 = 5
	RoadPriorityClassSecondary     uint8 = 6
	RoadPriorityClassSecondaryLink uint8 = 7
	RoadPriorityClassTertiary      uint8 = 8
	RoadPriorityClassTertiaryLink  uint8 = 9
	// Residential Categories
	RoadPriorityClassMainResidential uint8 = 10
	RoadPriorityClassSideResidential uint8 = 11
	RoadPriorityClassAlley           uint8 = 12
	RoadPriorityClassParking         uint8 = 13
	// Link Category
	RoadPriorityClassLinkRoad     uint8 = 14
	RoadPriorityClassUnclassified uint8 = 15
	// Bike Accessible
	RoadPriorityClassBikePath uint8 = 16
	// Walk Accessible
	RoadPriorityClassFootPath uint8 = 18
	// Link types are usually not considered in forks, unless amongst each other.
	// a road simply offered for connectivity. Will be ignored in forks/other decisions. Always
	// considered non-obvious to continue on
	RoadPriorityClassConnectivity uint8 = 31
)

func (r *RoadClassification) tryParse(p []byte) error {

	if len(p) < 2 {
		return fmt.Errorf("at least 2 bytes for RoadClassification but only got %d bytes", len(p))
	}

	if p[0]&0x01 > 0 {
		r.MotowayClass = true
	}
	if p[0]&0x02 > 0 {
		r.MaybeIgnored = true
	}
	if p[0]&0x04 > 0 {
		r.LinkClass = true
	}

	r.RoadPriorityClass = (p[0] & 0xF8) >> 3
	r.NumberOfLanes = p[1]
	return nil
}
