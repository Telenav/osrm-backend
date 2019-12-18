package osrmtype

// NameID represents OSRM defined Name ID.
// C++ Implementation: https://github.com/Telenav/osrm-backend/blob/6283c6074066f98e6d4a9f774f21ea45407c0d52/include/util/typedefs.hpp#L74
type NameID uint32

// AnnotationID represents OSRM defined Annotation ID.
// C++ Implementation: https://github.com/Telenav/osrm-backend/blob/6283c6074066f98e6d4a9f774f21ea45407c0d52/include/util/typedefs.hpp#L75
type AnnotationID uint32

// EdgeWeight represents weight of this Edge.
// C++ Implementation: https://github.com/Telenav/osrm-backend/blob/6283c6074066f98e6d4a9f774f21ea45407c0d52/include/util/typedefs.hpp#L76
type EdgeWeight int32

// EdgeDuration represents duration of this Edge.
// C++ Implementation: https://github.com/Telenav/osrm-backend/blob/6283c6074066f98e6d4a9f774f21ea45407c0d52/include/util/typedefs.hpp#L77
type EdgeDuration int32

// EdgeDistance represents distance of this Edge.
// C++ Implementation: https://github.com/Telenav/osrm-backend/blob/6283c6074066f98e6d4a9f774f21ea45407c0d52/include/util/typedefs.hpp#L78
type EdgeDistance float32

// LaneDescriptionID represents OSRM defined Lane Description ID.
// C++ Implementation: https://github.com/Telenav/osrm-backend/blob/6283c6074066f98e6d4a9f774f21ea45407c0d52/include/util/typedefs.hpp#L90
type LaneDescriptionID uint16

// ClassData represents OSRM defined ClassData.
// C++ Implementation: https://github.com/Telenav/osrm-backend/blob/6283c6074066f98e6d4a9f774f21ea45407c0d52/include/extractor/class_data.hpp#L15
type ClassData uint8

// invalid definitions
const (
	InvalidLaneDescriptionID = ^LaneDescriptionID(0)
)
