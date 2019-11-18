//Package nodebasededge defines NodeBasedEdge structure.
//more details refer to https://github.com/Telenav/open-source-spec/blob/master/osrm/doc/understanding_osrm_graph_representation.md#terminology
package nodebasededge

//Edge represent NodeBasedEdge structure. It's an directed edge between two nodes.
//https://github.com/Telenav/open-source-spec/blob/master/osrm/doc/understanding_osrm_graph_representation.md#terminology
type Edge struct {
	FromNode int64
	ToNode   int64
}

// Reverse returns reverse direction edge from original one.
func (e Edge) Reverse() Edge {
	return Edge{FromNode: e.ToNode, ToNode: e.FromNode}
}
