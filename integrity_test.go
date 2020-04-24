package rouge_test

import (
	"github.com/mcandre/rouge"

	"reflect"
	"testing"
	"testing/quick"
)

func TestUint32Marshaling(t *testing.T) {
	symmetricProperty := func(xs []uint32) bool {
		bs, err := rouge.Uint32sToBytes(xs)

		if err != nil {
			t.Error(err)
			return false
		}

		ys, err := rouge.BytesToUint32s(bs)

		if err != nil {
			t.Error(err)
			return false
		}

		return reflect.DeepEqual(ys, xs)
	}

	if err := quick.Check(symmetricProperty, nil); err != nil {
		t.Error(err)
	}
}
