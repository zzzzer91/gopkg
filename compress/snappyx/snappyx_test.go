package snappyx

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var bs = func() []byte {
	f, err := os.Open("../testdata/Isaac.Newton-Opticks.txt")
	if err != nil {
		panic(err)
	}
	bs, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	return bs
}()

func TestEncode(t *testing.T) {
	rst := Encode(bs)
	t.Log(len(rst))
	s, err := DecodeToString(rst)
	assert.Nil(t, err)
	t.Log(len(s))
	assert.Equal(t, string(bs), s)
}

func BenchmarkEncode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Encode(bs)
	}
}

func BenchmarkDecodeToString(b *testing.B) {
	rst := Encode(bs)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := DecodeToString(rst)
		assert.Nil(b, err)
	}
}
