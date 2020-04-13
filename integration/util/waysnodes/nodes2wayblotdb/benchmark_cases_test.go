package nodes2wayblotdb

import (
	"math/rand"

	"github.com/Telenav/osrm-backend/integration/util/waysnodes"
)

func randNodesCount() int {
	minNodes := 2
	maxNodes := 10

	n := rand.Intn(maxNodes)
	if n < minNodes {
		n = minNodes
	}
	return n
}

func generateBenchmarkCases(count int) []waysnodes.WayNodes {

	cases := []waysnodes.WayNodes{}

	for i := 0; i < count; i++ {
		c := waysnodes.WayNodes{}
		c.WayID = rand.Int63()
		nodesCount := randNodesCount()
		for j := 0; j < nodesCount; j++ {
			c.NodeIDs = append(c.NodeIDs, rand.Int63())
		}

		cases = append(cases, c)
	}

	return cases
}

var (
	cases100   = generateBenchmarkCases(100)
	cases1000  = generateBenchmarkCases(1000)
	cases10000 = generateBenchmarkCases(10000)
)
