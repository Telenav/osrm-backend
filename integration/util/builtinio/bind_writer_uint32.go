// Package builtinio provides io.Writer binding on built-in types.
package builtinio

import (
	"encoding/binary"
	"io"

	"github.com/golang/glog"
)

type writerBinderUint32 struct {
	v *uint32
}

// BindWriterOnUint32 binds io.Writer on values of built-in type uint32.
func BindWriterOnUint32(v *uint32) io.Writer {
	if v == nil {
		return nil
	}
	return &writerBinderUint32{v: v}
}

func (w *writerBinderUint32) Write(p []byte) (int, error) {
	if w.v == nil {
		glog.Fatal("bonded pointer can not be nil")
	}

	if len(p) < uint32Bytes {
		return 0, newInsufficiantBytesError(uint32Bytes, len(p), "uint32")
	}
	*w.v = binary.LittleEndian.Uint32(p)
	writeLen := uint32Bytes

	return writeLen, nil
}
