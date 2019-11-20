//Package graph defines a Node Based Graph.
//more details refer to https://github.com/Telenav/open-source-spec/blob/master/osrm/doc/understanding_osrm_graph_representation.md#terminology
package graph

//Edge represents NodeBasedEdge structure. It's an directed edge between two nodes.
//https://github.com/Telenav/open-source-spec/blob/master/osrm/doc/understanding_osrm_graph_representation.md#terminology
type Edge struct {
	From NodeID
	To   NodeID
}

// Reverse returns reverse direction edge from original one.
func (e Edge) Reverse() Edge {
	return Edge{From: e.To, To: e.From}
}

// ReverseEdges reverses the edges.
func ReverseEdges(s []Edge) []Edge {
	if len(s) == 0 {
		return s
	}

	for i, j := 0, len(s)-1; i <= j; i, j = i+1, j-1 {
		s[i], s[j] = s[j].Reverse(), s[i].Reverse()
	}
	return s
}
