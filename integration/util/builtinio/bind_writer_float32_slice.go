package builtinio

import (
	"encoding/binary"
	"io"
	"math"

	"github.com/golang/glog"
)

type writerBinderFloat32Slice struct {
	v *[]float32
}

// BindWriterOnFloat32Slice binds io.Writer on []float32.
func BindWriterOnFloat32Slice(v *[]float32) io.Writer {
	if v == nil {
		return nil
	}
	return &writerBinderFloat32Slice{v: v}
}

func (w *writerBinderFloat32Slice) Write(p []byte) (int, error) {
	if w.v == nil {
		glog.Fatal("bonded pointer can not be nil")
	}
	if *w.v == nil {
		*w.v = make([]float32, 0)
	}

	var writeLen int
	writeP := p
	for {
		if len(writeP) < float32Bytes {
			break
		}

		v := math.Float32frombits(binary.LittleEndian.Uint32(writeP))

		*w.v = append(*w.v, v)

		writeP = writeP[float32Bytes:]
		writeLen += float32Bytes
	}

	return writeLen, nil
}
