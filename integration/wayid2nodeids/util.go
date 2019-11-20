package wayid2nodeids

func absInt64(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}
