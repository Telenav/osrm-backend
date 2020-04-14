package csvreader

// CompressionType represents compression types that the package is supported.
type CompressionType string

// Supported compression algorithms, either none or snappy at the moment.
const (
	CompressionTypeNone   CompressionType = "none"
	CompressionTypeSnappy                 = "snappy"
)

// Default options
const (
	DefaultMaxReadCount  = 500
	DefaultMaxCacheCount = 100
)

// Options represents the options that can be set when reading csv file.
type Options struct {

	// Compression is the compression type to decode the csv file when reading.
	Compression CompressionType

	// MaxReadCount is the max count when you read from ReadLines for ReadRecords.
	// If <= 0, the DefaultMaxReadCount will be used.
	MaxReadCount int

	// MaxCacheCount is the max cache of internal channel.
	// If <= 0, block channel will be used.
	// If your post actions are heavy, you may want to increase this value for better performance.
	MaxCacheCount int
}

var defaultOptions = Options{
	Compression:   CompressionTypeNone,
	MaxReadCount:  DefaultMaxReadCount,
	MaxCacheCount: DefaultMaxCacheCount,
}

// DefaultOptions returns the default options.
func DefaultOptions() *Options {
	return &defaultOptions
}
