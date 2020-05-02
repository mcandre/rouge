package rouge

import (
	"encoding/binary"
	"errors"
)

var BO = binary.BigEndian

var InvalidBufferSize = errors.New("Invalid buffer size (expected a multiple of 3)")

func Uint32sToBytes(xs []uint32) ([]byte, error) {
	bs := []byte{}

	for _, x := range xs {
		le := make([]byte, 4)
		BO.PutUint32(le, x)
		bs = append(bs, le...)
	}

	return bs, nil
}

func BytesToUint32s(bs []byte) ([]uint32, error) {
	if len(bs) % 4 != 0 {
		return nil, InvalidBufferSize
	}

	xs := []uint32{}

	for i := 0; i < len(bs)/4; i++ {
		le := bs[i*4:i*4+4]
		x := BO.Uint32(le)
		xs = append(xs, x)
	}

	return xs, nil
}
