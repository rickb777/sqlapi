package schema

import (
	. "github.com/onsi/gomega"
	"github.com/rickb777/sqlapi/util"
	"testing"
)

func TestDistinctTypes(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		list     FieldList
		expected util.StringSet
	}{
		{FieldList{id}, util.NewStringSet(i64.Tag())},
		{FieldList{id, id, id}, util.NewStringSet(i64.Tag())},
		{FieldList{id, category}, util.NewStringSet(i64.Tag(), cat.Tag())},
		{FieldList{id,
			category,
			name,
			qual,
			diff,
			age,
			bmi,
			active,
			labels,
			fave,
			avatar,
			updated}, util.NewStringSet(i64.Tag(), boo.Tag(), cat.Tag(), str.Tag(), spt.Tag(), ipt.Tag(), upt.Tag(), fpt.Tag(), bgi.Tag(), sli.Tag(), bys.Tag(), tim.Tag())},
	}
	for _, c := range cases {
		s := c.list.DistinctTypes()
		g.Expect(util.NewStringSet(s...)).To(Equal(c.expected))
	}
}
