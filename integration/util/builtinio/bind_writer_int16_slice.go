package builtinio

import (
	"encoding/binary"
	"io"

	"github.com/golang/glog"
)

type writerBinderInt16Slice struct {
	v *[]int16
}

// BindWriterOnInt16Slice binds io.Writer on []int16.
func BindWriterOnInt16Slice(v *[]int16) io.Writer {
	if v == nil {
		return nil
	}
	return &writerBinderInt16Slice{v: v}
}

func (w *writerBinderInt16Slice) Write(p []byte) (int, error) {
	if w.v == nil {
		glog.Fatal("bonded pointer can not be nil")
	}
	if *w.v == nil {
		*w.v = make([]int16, 0)
	}

	var writeLen int
	writeP := p
	for {
		if len(writeP) < int16Bytes {
			break
		}

		v := int16(binary.LittleEndian.Uint16(writeP))

		*w.v = append(*w.v, v)

		writeP = writeP[int16Bytes:]
		writeLen += int16Bytes
	}

	return writeLen, nil
}
