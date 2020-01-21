package osrmfiles

import "io"

// ContentsOperator is the interfaces that wraps Contents' methods.
type ContentsOperator interface {

	// PrintSummary prints summary and head lines of contents.
	PrintSummary(head int)

	// Validate checks whether the contents valid or not.
	Validate() error

	// FindWriter find io.Writer for the specified name, contents can be filled in by the found io.Writer.
	FindWriter(name string) (io.Writer, bool)

	// FilePath returns the file path that stores the contents.
	FilePath() string
}
