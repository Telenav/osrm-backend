package builtinio

import (
	"encoding/binary"
	"io"

	"github.com/golang/glog"
)

type writerBinderInt32Slice struct {
	v *[]int32
}

// BindWriterOnInt32Slice binds io.Writer on []int32.
func BindWriterOnInt32Slice(v *[]int32) io.Writer {
	if v == nil {
		return nil
	}
	return &writerBinderInt32Slice{v: v}
}

func (w *writerBinderInt32Slice) Write(p []byte) (int, error) {
	if w.v == nil {
		glog.Fatal("bonded pointer can not be nil")
	}
	if *w.v == nil {
		*w.v = make([]int32, 0)
	}

	var writeLen int
	writeP := p
	for {
		if len(writeP) < int32Bytes {
			break
		}

		var id int32
		id = int32(binary.LittleEndian.Uint32(writeP))

		*w.v = append(*w.v, id)

		writeP = writeP[int32Bytes:]
		writeLen += int32Bytes
	}

	return writeLen, nil
}
