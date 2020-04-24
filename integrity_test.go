package rouge_test

import (
	"github.com/mcandre/rouge"

	"reflect"
	"testing"
	"testing/quick"
)

func TestIntMarshaling(t *testing.T) {
	symmetricProperty := func(xs []int) bool {
		bs, err := rouge.IntsToBytes(xs)

		if err != nil {
			if err == rouge.IntegerElementOutOfBounds {
				return true
			}

			t.Error(err)
		}

		ys, err := rouge.BytesToInts(bs)

		if err != nil {
			t.Error(err)
		}

		return reflect.DeepEqual(ys, xs)
	}

	if err := quick.Check(symmetricProperty, nil); err != nil {
		t.Error(err)
	}
}
