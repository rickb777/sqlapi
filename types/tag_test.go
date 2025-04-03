package types

import (
	"os"
	"testing"

	"github.com/rickb777/expect"
)

func TestParseTag(t *testing.T) {
	tagTests := []struct {
		raw string
		tag *Tag
	}{
		{
			TagKey + `:""`,
			nil,
		},
		{
			TagKey + `:"-"`,
			&Tag{Skip: true},
		},
		{
			TagKey + `:"prefixed: true"`,
			&Tag{Prefixed: true},
		},
		{
			TagKey + `:"pk: true"`,
			&Tag{Primary: true, Auto: false},
		},
		{
			TagKey + `:"pk: true, auto: true"`,
			&Tag{Primary: true, Auto: true},
		},
		{
			TagKey + `:"name: foo"`,
			&Tag{Name: "foo"},
		},
		{
			TagKey + `:"type: varchar"`,
			&Tag{Type: "varchar"},
		},
		{
			TagKey + `:"size: 2048"`,
			&Tag{Size: 2048},
		},
		{
			TagKey + `:"default: abc123"`,
			&Tag{Default: "abc123"},
		},
		{
			TagKey + `:"index: fake_index"`,
			&Tag{Index: "fake_index"},
		},
		{
			TagKey + `:"unique: fake_unique_index"`,
			&Tag{Unique: "fake_unique_index"},
		},
		{
			TagKey + `:"fk: alpha.ID, onupdate: setnull, ondelete: setdefault"`,
			&Tag{ForeignKey: "alpha.ID", OnUpdate: "set null", OnDelete: "set default"},
		},
		{
			TagKey + `:"fk: alpha, onupdate: 'set null', ondelete: 'set default'"`,
			&Tag{ForeignKey: "alpha", OnUpdate: "set null", OnDelete: "set default"},
		},
	}

	for _, test := range tagTests {
		expect.Any(ParseTag(test.raw)).Info(test.raw).ToBe(t, test.tag)
	}
}

func TestParseValidation(t *testing.T) {
	tagTests := []struct {
		raw string
		err string
	}{
		{
			TagKey + `:"encode: x"`,
			`unrecognised encode value "x"`,
		},
		{
			TagKey + `:"auto: true"`,
			`auto-increment can only be used on primary keys`,
		},
		{
			TagKey + `:"nk: true, auto: true"`,
			`natural key cannot use auto-increment; auto-increment can only be used on primary keys`,
		},
		{
			TagKey + `:"pk: true, nk: true"`,
			`primary key cannot also be a natural key`,
		},
		{
			TagKey + `:"fk: x.x.x"`,
			`fk value ("x.x.x") must be in 'tablename' or 'tablename.column' form`,
		},
		{
			TagKey + `:"pk: true, fk: a.b"`,
			`foreign key cannot also be a primary key nor a natural key`,
		},
		{
			TagKey + `:"nk: true, fk: a.b"`,
			`foreign key cannot also be a primary key nor a natural key`,
		},
		{
			TagKey + `:"nk: true, index: foo"`,
			`natural key cannot be used with index`,
		},
		{
			TagKey + `:"nk: true, unique: foo"`,
			`natural key should not be used with unique`,
		},
		{
			TagKey + `:"onupdate: x"`,
			`unrecognised onupdate value "x"`,
		},
		{
			TagKey + `:"ondelete: x"`,
			`unrecognised ondelete value "x"`,
		},
		{
			TagKey + `:"onupdate: x, ondelete: y"`,
			`unrecognised onupdate value "x"; unrecognised ondelete value "y"`,
		},
		{
			TagKey + `:"size: -1"`,
			`size cannot be negative (-1)`,
		},
		{
			TagKey + `:"encode: foo"`,
			`unrecognised encode value "foo"`,
		},
	}

	for _, test := range tagTests {
		expect.Error(ParseTag(test.raw)).I(test.raw).ToContain(t, test.err)
	}
}

func TestReadTagsFile(t *testing.T) {
	file := os.TempDir() + "/sqlgen2-test.yaml"
	defer os.Remove(file)

	yml := `
Id:
  pk: true
  auto: true

Foo:
  name: fooish
  type: blob
`

	err := os.WriteFile(file, []byte(yml), 0644)
	expect.Error(err).ToBeNil(t)

	tags, err := ReadTagsFile(file)
	expect.Error(err).ToBeNil(t)
	expect.Map(tags).ToHaveLength(t, 2)

	id := tags["Id"]
	expect.Any(id).ToBe(t, &Tag{Primary: true, Auto: true})

	foo := tags["Foo"]
	expect.Any(foo).ToBe(t, &Tag{Name: "fooish", Type: "blob"})
}
