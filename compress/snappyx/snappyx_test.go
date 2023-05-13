package snappyx

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	f, err := os.Open("./testdata/Isaac.Newton-Opticks.txt")
	assert.Nil(t, err)
	bs, err := io.ReadAll(f)
	assert.Nil(t, err)
	t.Log(len(bs))
	rst := Encode(bs)
	t.Log(len(rst))
	s, err := DecodeToString(rst)
	assert.Nil(t, err)
	t.Log(len(s))
	assert.Equal(t, string(bs), s)
}
