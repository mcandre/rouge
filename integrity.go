package rouge

import (
	"encoding/binary"
	"errors"
)

// BO denotes an internal representation of PCM data.
var BO = binary.BigEndian

// ErrorInvalidBufferSize indicates a configuration error.
var ErrorInvalidBufferSize = errors.New("Invalid buffer size (expected a multiple of 4)")

// Uint32sToBytes applies BO encoding onto unsigned 32-bit integers.
func Uint32sToBytes(xs []uint32) ([]byte, error) {
	bs := []byte{}

	for _, x := range xs {
		le := make([]byte, 4)
		BO.PutUint32(le, x)
		bs = append(bs, le...)
	}

	return bs, nil
}

// BytesToUint32s applies BO decoding from byte arrays.
func BytesToUint32s(bs []byte) ([]uint32, error) {
	if len(bs)%4 != 0 {
		return nil, ErrorInvalidBufferSize
	}

	xs := []uint32{}

	for i := 0; i < len(bs)/4; i++ {
		le := bs[i*4 : i*4+4]
		x := BO.Uint32(le)
		xs = append(xs, x)
	}

	return xs, nil
}
