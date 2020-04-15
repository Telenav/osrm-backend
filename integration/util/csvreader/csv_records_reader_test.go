package csvreader_test

import (
	"fmt"
	"log"

	"github.com/Telenav/osrm-backend/integration/util/csvreader"
)

func ExampleRecordsAsyncReader() {

	r := csvreader.NewRecordsAsyncReader("test.csv", nil)
	r.Start()

	for {
		records, ok := r.ReadRecords()
		if !ok {
			break
		}

		for _, record := range records {
			// TODO: use record
			fmt.Println(record)
		}
	}

	if err := r.Err(); err != nil {
		log.Fatal(err)
	}
}
