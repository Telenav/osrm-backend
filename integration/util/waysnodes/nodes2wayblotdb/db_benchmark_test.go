package nodes2wayblotdb

import (
	"math/rand"
	"os"
	"testing"

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

func benchmarkWrite(b *testing.B, cases []waysnodes.WayNodes) {

	// Open a new DB.
	db, err := Open("bench-write.db", false)
	if err != nil {
		b.Error(err)
	}
	defer os.Remove(db.db.Path()) // remove temp db file after test.

	for i := 0; i < b.N; i++ {
		for _, c := range cases {
			err = db.Write(c.WayID, c.NodeIDs)
			if err != nil {
				b.Error(err)
			}
		}
	}

	//b.Log(db.Statistics())

	// Close database to release the file lock.
	if err := db.Close(); err != nil {
		b.Error(err)
	}
}

func BenchmarkWriteDBCases100(b *testing.B) {
	benchmarkWrite(b, cases100)
}
func BenchmarkWriteDBCases1000(b *testing.B) {
	benchmarkWrite(b, cases1000)
}
func BenchmarkWriteDBCases10000(b *testing.B) {
	benchmarkWrite(b, cases10000)
}

func benchmarkBatchWriteDB(b *testing.B, cases []waysnodes.WayNodes) {
	// Open a new DB.
	db, err := Open("bench-write.db", false)
	if err != nil {
		b.Error(err)
	}
	defer os.Remove(db.db.Path()) // remove temp db file after test.

	for i := 0; i < b.N; i++ {
		err = db.BatchWrite(cases)
		if err != nil {
			b.Error(err)
		}
	}

	//b.Log(db.Statistics())

	// Close database to release the file lock.
	if err := db.Close(); err != nil {
		b.Error(err)
	}
}

func BenchmarkBatchWriteDBCases100(b *testing.B) {
	benchmarkBatchWriteDB(b, cases100)
}

func BenchmarkBatchWriteDBCases1000(b *testing.B) {
	benchmarkBatchWriteDB(b, cases1000)
}
func BenchmarkBatchWriteDBCases10000(b *testing.B) {
	benchmarkBatchWriteDB(b, cases10000)
}
