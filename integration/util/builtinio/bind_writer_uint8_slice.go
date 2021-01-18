package builtinio

import (
	"io"

	"github.com/golang/glog"
)

type writerBinderUint8Slice struct {
	v *[]uint8
}

// BindWriterOnUint8Slice binds io.Writer on []uint8.
func BindWriterOnUint8Slice(v *[]uint8) io.Writer {
	if v == nil {
		return nil
	}
	return &writerBinderUint8Slice{v: v}
}

func (w *writerBinderUint8Slice) Write(p []byte) (int, error) {
	if w.v == nil {
		glog.Fatal("bonded pointer can not be nil")
	}
	if *w.v == nil {
		*w.v = make([]uint8, 0)
	}

	*w.v = append(*w.v, p...)
	return len(p), nil
}
