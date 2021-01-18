package builtinio

import (
	"encoding/binary"
	"io"

	"github.com/golang/glog"
)

type writerBinderUint32Slice struct {
	v *[]uint32
}

// BindWriterOnUint32Slice binds io.Writer on []uint32.
func BindWriterOnUint32Slice(v *[]uint32) io.Writer {
	if v == nil {
		return nil
	}
	return &writerBinderUint32Slice{v: v}
}

func (w *writerBinderUint32Slice) Write(p []byte) (int, error) {
	if w.v == nil {
		glog.Fatal("bonded pointer can not be nil")
	}
	if *w.v == nil {
		*w.v = make([]uint32, 0)
	}

	var writeLen int
	writeP := p
	for {
		if len(writeP) < uint32Bytes {
			break
		}

		var id uint32
		id = binary.LittleEndian.Uint32(writeP)

		*w.v = append(*w.v, id)

		writeP = writeP[uint32Bytes:]
		writeLen += uint32Bytes
	}

	return writeLen, nil
}
