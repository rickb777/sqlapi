package require

import (
	. "github.com/onsi/gomega"
	"testing"
)

func TestWrongSizeZero(t *testing.T) {
	g := NewGomegaWithT(t)

	err := ErrWrongSize(0, "wanted %d", 5)
	g.Expect(err.Error()).To(Equal("wanted 5"))
	g.Expect(err.(Sizer).Size()).To(BeEquivalentTo(0))
	g.Expect(IsNotFound(err)).To(BeTrue())
	g.Expect(IsNotUnique(err)).To(BeFalse())
}

func TestWrongSizeMany(t *testing.T) {
	g := NewGomegaWithT(t)

	err := ErrWrongSize(3, "wanted %d", 5)
	g.Expect(err.Error()).To(Equal("wanted 5"))
	g.Expect(err.(Sizer).Size()).To(BeEquivalentTo(3))
	g.Expect(IsNotFound(err)).To(BeFalse())
	g.Expect(IsNotUnique(err)).To(BeTrue())
}
