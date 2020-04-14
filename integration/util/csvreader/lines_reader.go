// Package csvreader is dedicated to read huge csv file as fast as possible.
package csvreader

import (
	"bufio"
	"io"
	"os"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/golang/snappy"
)

// LinesAsyncReader represent async reader to read csv lines.
type LinesAsyncReader struct {
	file    string
	options Options

	linesChan chan []string

	errMutex sync.RWMutex
	err      error
}

// NewLinesAsyncReader creates a new async reader to read lines.
func NewLinesAsyncReader(file string, options *Options) *LinesAsyncReader {

	l := LinesAsyncReader{}
	l.file = file
	if options == nil {
		options = &defaultOptions
	}
	l.options = *options
	l.linesChan = make(chan []string, options.MaxCacheCount)

	return &l
}

// Err returns the first non-EOF error that was encountered by the async reader.
func (l *LinesAsyncReader) Err() error {
	l.errMutex.RLock()
	defer l.errMutex.RUnlock()
	return l.err
}

// Start starts the async reader to read from file and returns immediatelly.
func (l *LinesAsyncReader) Start() {
	go l.reader()
}

// ReadLines will block if no data to return.
// It returns maximum MaxReadCount lines and true per call if data available.
// If the second value returns false, the async reading process has been closed, and you'd better to check whether any Err() occurs.
// It's safe to call it by multiple goroutines at the same time.
func (l *LinesAsyncReader) ReadLines() ([]string, bool) {
	ss, ok := <-l.linesChan
	return ss, ok
}

func (l *LinesAsyncReader) setError(err error) {
	l.errMutex.Lock()
	defer l.errMutex.Unlock()
	l.err = err
}

func (l *LinesAsyncReader) newCompressedReader(r io.Reader) io.Reader {
	if l.options.Compression == CompressionTypeSnappy {
		return snappy.NewReader(r)
	}
	return r
}

func (l *LinesAsyncReader) reader() {
	defer close(l.linesChan)

	startTime := time.Now()

	f, err := os.Open(l.file)
	if err != nil {
		l.setError(err)
		return
	}
	defer f.Close()

	var total int

	linesPerTrans := l.options.MaxReadCount
	lines := make([]string, 0, linesPerTrans)

	scanner := bufio.NewScanner(l.newCompressedReader(f))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		total++

		if len(lines) >= linesPerTrans {
			l.linesChan <- lines
			lines = make([]string, 0, linesPerTrans)
		}
	}
	if len(lines) > 0 { // remain lines
		l.linesChan <- lines
		lines = nil
	}
	if scanner.Err() != nil {
		l.setError(err)
		return
	}

	glog.V(2).Infof("Read %d lines, takes %f seconds.\n", total, time.Now().Sub(startTime).Seconds())
}
