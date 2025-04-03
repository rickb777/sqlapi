package require

import (
	"testing"

	"github.com/rickb777/expect"
)

func TestWrongSizeZero(t *testing.T) {
	err := ErrWrongSize(0, "wanted %d", 5)
	expect.Error(err).ToContain(t, "wanted 5")
	expect.Number(err.(Sizer).Size()).ToBe(t, 0)
	expect.Bool(IsNotFound(err)).ToBeTrue(t)
	expect.Bool(IsNotUnique(err)).ToBeFalse(t)
}

func TestWrongSizeMany(t *testing.T) {
	err := ErrWrongSize(3, "wanted %d", 5)
	expect.Error(err).ToContain(t, "wanted 5")
	expect.Number(err.(Sizer).Size()).ToBe(t, 3)
	expect.Bool(IsNotFound(err)).ToBeFalse(t)
	expect.Bool(IsNotUnique(err)).ToBeTrue(t)
}
