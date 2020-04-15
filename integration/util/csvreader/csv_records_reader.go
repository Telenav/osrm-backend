package csvreader

import (
	"strings"
	"sync"
	"time"

	"github.com/golang/glog"
)

// RecordsAsyncReader represent async reader to read csv lines.
type RecordsAsyncReader struct {
	l *LinesAsyncReader

	recordsChan chan [][]string
}

// NewRecordsAsyncReader creates a new async reader to read csv records.
func NewRecordsAsyncReader(file string, options *Options) *RecordsAsyncReader {

	r := RecordsAsyncReader{}
	r.l = NewLinesAsyncReader(file, options)
	r.recordsChan = make(chan [][]string, options.MaxCacheCount)

	return &r
}

// Err returns the first non-EOF error that was encountered by the async reader.
func (r *RecordsAsyncReader) Err() error {
	return r.l.Err()
}

// Start starts the async reader to read from file and returns immediatelly.
func (r *RecordsAsyncReader) Start() {

	go func() {
		defer close(r.recordsChan) // close records channel until all parser exited

		wg := sync.WaitGroup{}
		wg.Add(r.l.options.ParallelParser)
		for i := 0; i < r.l.options.ParallelParser; i++ {
			go func() {
				r.parser()
				wg.Done()
			}()
		}
		wg.Wait()
	}()
	r.l.Start()
}

// ReadRecords will block if no data to return.
// It returns minimum MinReadCount csv records(except last read) and true per call if data available.
// If the second value returns false, the async reading process has been closed, and you'd better to check whether any Err() occurs.
// It's safe to call it by multiple goroutines at the same time.
func (r *RecordsAsyncReader) ReadRecords() ([][]string, bool) {
	records, ok := <-r.recordsChan
	return records, ok
}

func (r *RecordsAsyncReader) parser() {

	startTime := time.Now()

	var readLines, parsedRecords int

	recordsPerTrans := r.l.options.MinReadCount
	records := make([][]string, 0, recordsPerTrans)

	for {
		lines, ok := r.l.ReadLines()
		if !ok {
			break
		}
		readLines++

		for _, line := range lines {
			record := strings.Split(line, ",")
			if len(record) == 0 {
				continue // silently ignore empty record
			}
			records = append(records, record)
			parsedRecords++
		}

		if len(records) >= recordsPerTrans {
			r.recordsChan <- records
			records = make([][]string, 0, recordsPerTrans)
		}
	}
	if len(records) > 0 {
		r.recordsChan <- records
		records = nil
	}

	glog.V(2).Infof("Parsed %d records from %d lines, takes %f seconds.\n", parsedRecords, readLines, time.Now().Sub(startTime).Seconds())
}
