package shorter

import (
	"math/rand"
	"testing"
)

func TestConvert(t *testing.T) {
	a := rand.Uint64()
	if a > 62*62*62*62 {
		a = 62 * 62 * 62 * 62
	}
	for i := uint64(0); i < a; i++ {
		str := convert10To62(i)
		get, err := convert62To10(str)
		if nil != err || get != i {
			t.Error("convert has error:", err, get, i, str)
			return
		}

	}
}

func BenchmarkConvert(b *testing.B) {
	for i := uint64(0); i < uint64(b.N); i++ {
		str := convert10To62(i)
		get, err := convert62To10(str)
		if nil != err || get != i {
			b.Error("convert has error:", err, get, i, str)
			return
		}
	}
}
