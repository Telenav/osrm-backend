package nodes2wayblotdb

import (
	"flag"
	"os"
	"testing"

	"github.com/Telenav/osrm-backend/integration/util/waysnodes"
)

var testFlags struct {
	printDBStats bool
}

func init() {
	flag.BoolVar(&testFlags.printDBStats, "print-db-stats", false, "Print DB stats.")
}

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

	if testFlags.printDBStats {
		b.Log(db.Statistics())
	}

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

	if testFlags.printDBStats {
		b.Log(db.Statistics())
	}

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
