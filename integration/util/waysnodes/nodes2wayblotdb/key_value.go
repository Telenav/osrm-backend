package nodes2wayblotdb

import (
	"encoding/binary"

	"github.com/golang/glog"
)

const (
	int64Bytes = 8

	maxKeyBytes   = 2 * int64Bytes // 2 int64
	maxValueBytes = int64Bytes     // 1 int64
)

func key(fromNodeID, toNodeID int64) []byte {

	buf := make([]byte, maxKeyBytes, maxKeyBytes)
	binary.LittleEndian.PutUint64(buf, uint64(fromNodeID))
	binary.LittleEndian.PutUint64(buf[int64Bytes:], uint64(toNodeID))
	return buf
}

func parseKey(buf []byte) (fromNodeID, toNodeID int64) {
	if len(buf) < maxKeyBytes {
		glog.Fatalf("invalid key buf: %v\n", buf)
		return
	}

	fromNodeID = int64(binary.LittleEndian.Uint64(buf))
	toNodeID = int64(binary.LittleEndian.Uint64(buf[int64Bytes:]))
	return
}

func value(wayID int64) []byte {
	buf := make([]byte, maxValueBytes, maxValueBytes)
	binary.LittleEndian.PutUint64(buf, uint64(wayID))
	return buf
}

func parseValue(buf []byte) (wayID int64) {
	wayID = int64(binary.LittleEndian.Uint64(buf))
	return
}
