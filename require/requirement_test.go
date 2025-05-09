package require

import (
	"fmt"
	"testing"

	"github.com/rickb777/expect"
)

var e0 = fmt.Errorf("foo")

func TestChainErrorIfQueryNotSatisfiedBy_and_ChainErrorIfExecNotSatisfiedBy(t *testing.T) {
	cases := []struct {
		input    error
		req      Requirement
		actual   int64
		expected bool
		sMessage string
		qMessage string
		eMessage string
	}{
		{nil, nil, 0, false, "", "", ""},
		{e0, nil, 0, true, "", "foo", "foo"},

		{nil, None, 0, false, "none", "", ""},
		{nil, None, 5, true, "none", "expected to fetch none but got 5", "expected to change none but changed 5"},
		{nil, One, 5, true, "one", "expected to fetch one but got 5", "expected to change one but changed 5"},
		{nil, Many, 5, false, "many", "", ""},
		{nil, Many, 1, true, "many", "expected to fetch many but got 1", "expected to change many but changed 1"},

		{nil, Exactly(3), 3, false, "exactly 3", "", ""},
		{nil, Exactly(3), 2, true, "exactly 3", "expected to fetch 3 but got 2", "expected to change 3 but changed 2"},

		{nil, AtLeast(3), 5, false, "at least 3", "", ""},
		{nil, AtLeast(3), 2, true, "at least 3", "expected to fetch at least 3 but got 2", "expected to change at least 3 but changed 2"},

		{nil, NoMoreThan(3), 2, false, "no more than 3", "", ""},
		{nil, NoMoreThan(3), 5, true, "no more than 3", "expected to fetch no more than 3 but got 5", "expected to change no more than 3 but changed 5"},
	}

	for _, c := range cases {
		e1 := ChainErrorIfQueryNotSatisfiedBy(c.input, c.req, c.actual)
		if c.req != nil {
			expect.String(c.req.String()).ToBe(t, c.sMessage)
		}

		if c.expected {
			expect.Error(e1).ToContain(t, c.qMessage)
		} else {
			expect.Error(e1).ToBeNil(t)
		}

		e2 := ChainErrorIfExecNotSatisfiedBy(c.input, c.req, c.actual)
		if c.expected {
			expect.Error(e2).ToContain(t, c.eMessage)
		} else {
			expect.Error(e2).ToBeNil(t)
		}
	}
}
