package csvreader_test

import (
	"fmt"
	"log"
	"sync"

	"github.com/Telenav/osrm-backend/integration/util/csvreader"
)

func ExampleLinesAsyncReader() {

	l := csvreader.NewLinesAsyncReader("test.csv", nil)
	l.Start()

	for {
		lines, ok := l.ReadLines()
		if !ok {
			break
		}

		for _, line := range lines {
			// TODO: use line
			fmt.Println(line)
		}
	}

	if err := l.Err(); err != nil {
		log.Fatal(err)
	}
}

func ExampleLinesAsyncReader_snappyCompressed() {

	options := csvreader.DefaultOptions()
	options.Compression = csvreader.CompressionTypeSnappy

	l := csvreader.NewLinesAsyncReader("test.csv.snappy", options)
	l.Start()

	for {
		lines, ok := l.ReadLines()
		if !ok {
			break
		}

		for _, line := range lines {
			// TODO: use line
			fmt.Println(line)
		}
	}

	if err := l.Err(); err != nil {
		log.Fatal(err)
	}
}

func ExampleLinesAsyncReader_parallelConsuming() {

	l := csvreader.NewLinesAsyncReader("test.csv", nil)
	l.Start()

	parallel := 5
	wg := sync.WaitGroup{}
	wg.Add(parallel)
	for i := 0; i < parallel; i++ {
		go func() {
			for {
				lines, ok := l.ReadLines()
				if !ok {
					break
				}

				for _, line := range lines {
					// TODO: use line
					fmt.Println(line)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()

	if err := l.Err(); err != nil {
		log.Fatal(err)
	}
}
