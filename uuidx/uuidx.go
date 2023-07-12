package uuidx

import (
	"encoding/hex"
	"unsafe"

	"github.com/google/uuid"
)

// Using global variable for gomonkey patch.
// Gomonkey function patch feature is not allowed on Mac.
var New = func() string {
	id := uuid.New()
	dst := make([]byte, 32)
	hex.Encode(dst, id[:])
	return bytes2string(dst)
}

func bytes2string(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
