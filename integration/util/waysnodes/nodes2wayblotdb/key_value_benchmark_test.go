package nodes2wayblotdb

import "testing"

func BenchmarkGenerateKey(b *testing.B) {

	cases := cases1000
	for i := 0; i < b.N; i++ {
		for _, c := range cases {
			for i := 0; i < len(c.NodeIDs)-1; i++ {
				key(c.NodeIDs[i], c.NodeIDs[i+1]) // only generate, don't need to use
			}
		}
	}
}

func BenchmarkGenerateValue(b *testing.B) {

	cases := cases1000
	for i := 0; i < b.N; i++ {
		for _, c := range cases {
			value(c.WayID) // only generate, don't need to use
		}
	}
}
