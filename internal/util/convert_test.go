package util

import (
	"fmt"
	"testing"
)

var table = []struct {
	in  uint64
	out string
}{{
	in:  0,
	out: "0 B",
}, {
	in:  1,
	out: "1 B",
}, {
	in:  1024,
	out: "1.0 kB",
}, {
	in:  4792,
	out: "4.8 kB",
}, {
	in:  10240,
	out: "10.2 kB",
},
}

func TestConvertBytes(t *testing.T) {
	for _, v := range table {
		if out := ByteCountSI(v.in); out != v.out {
			t.Errorf("ConvertBytes(%d) = %s, want %s", v.in, out, v.out)
		}
	}
}

func BenchmarkConvertBytes(b *testing.B) {
	for _, v := range table {
		b.Run(fmt.Sprintf("input_size_%d", v.in), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ByteCountSI(v.in)
			}
		})
	}
}
