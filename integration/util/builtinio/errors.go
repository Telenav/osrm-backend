package builtinio

import (
	"fmt"
)

func newInsufficiantBytesError(expectBytes, gotBytes int, typeName string) error {
	return fmt.Errorf("expect %d bytes for type %s but only got %d bytes", expectBytes, typeName, gotBytes)
}
