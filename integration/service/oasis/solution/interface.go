package solution

import "github.com/Telenav/osrm-backend/integration/api/oasis"

// Generator is the interface that wraps the Generate method
type Generator interface {

	// Generate calculates oasis solutions based on oasis request
	Generate(oasisReq *oasis.Request) (int, []*oasis.Solution, error)
}
