package csvreader

// CompressionType represents compression types that the package is supported.
type CompressionType string

// Supported compression algorithms, either none or snappy at the moment.
const (
	CompressionTypeNone   CompressionType = "none"
	CompressionTypeSnappy                 = "snappy"
)

// Default options by experiment
const (
	DefaultMinReadCount   = 500
	DefaultMaxCacheCount  = 100
	DefaultParallelParser = 3
)

// Options represents the options that can be set when reading csv file.
type Options struct {

	// Compression is the compression type to decode the csv file when reading.
	Compression CompressionType

	// MinReadCount is the minimum count when you read from ReadLines/ReadRecords.
	// If <= 0, the DefaultMinReadCount will be used.
	MinReadCount int

	// MaxCacheCount is the max cache of internal channel.
	// If <= 0, block channel will be used.
	// If your post actions are heavy, you may want to increase this value for better performance.
	MaxCacheCount int

	// ParallelParser is how many parallel parsers to parse from line to records.
	// It ONLY affects for RecordsAsyncReader.
	// It's goal is to improve performance.
	// If <= 0, the DefaultParallelParser will be used.
	ParallelParser int
}

var defaultOptions = Options{
	Compression:    CompressionTypeNone,
	MinReadCount:   DefaultMinReadCount,
	MaxCacheCount:  DefaultMaxCacheCount,
	ParallelParser: DefaultParallelParser,
}

// DefaultOptions returns the default options.
func DefaultOptions() *Options {
	options := defaultOptions // copy to avoid modify the defaultOptions directly
	return &options
}

func validateOptions(options *Options) *Options {
	if options == nil {
		return &defaultOptions
	}
	if options.MinReadCount <= 0 {
		options.MinReadCount = DefaultMinReadCount
	}
	if options.MaxCacheCount <= 0 {
		options.MaxCacheCount = 0 // block channel
	}
	if options.ParallelParser <= 0 {
		options.ParallelParser = DefaultParallelParser
	}
	return options
}
