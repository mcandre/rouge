package rouge

import (
	"encoding/binary"
	"errors"
	"math"
)

var BO = binary.LittleEndian
var IntegerElementOutOfBounds = errors.New("Integer element out of bounds")
var InvalidBufferSize = errors.New("Invalid buffer size (expected a multiple of 3)")

func IntsToBytes(xs []int) ([]byte, error) {
	if len(xs) == 0 {
		return []byte{}, nil
	}

	var bs []byte

	for _, x := range xs {
		if x > math.MaxUint32 {
			return nil, IntegerElementOutOfBounds
		}

		le := make([]byte, 4)
		y := uint32(x)
		BO.PutUint32(le, y)
		bs = append(bs, le...)
	}

	return bs, nil
}

func BytesToInts(bs []byte) ([]int, error) {
	if len(bs) == 0 {
		return []int{}, nil
	}

	if len(bs) % 4 != 0 {
		return nil, InvalidBufferSize
	}

	var xs []int

	for i := 0; i < len(bs)/4; i++ {
		le := bs[i*4:i*4+4]
		y := BO.Uint32(le)
		x := int(y)
		xs = append(xs, x)
	}

	return xs, nil
}
