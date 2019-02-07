package schema

import (
	"fmt"
	. "github.com/rickb777/sqlapi/types"
	"sort"
	"strings"
)

type SqlEncode int

// List of vendor-specific keywords
const (
	ENCNONE SqlEncode = iota
	ENCJSON
	ENCTEXT
	ENCDRIVER // SQL driver uses Scan() & Value() to encode & decode
)

type SqlToken int

// List of vendor-specific keywords
const (
	AUTO_INCREMENT SqlToken = iota
	PRIMARY_KEY
)

type TableDescription struct {
	Type string
	Name string

	Fields  FieldList
	Index   []*Index `json:",omitempty"`
	Primary *Field   `json:",omitempty"` // compound primaries are not supported
}

type Node struct {
	Name   string
	Type   Type
	Parent *Node `json:",omitempty"`
}

type Field struct {
	Node
	SqlName string
	Encode  SqlEncode `json:",omitempty"`
	Tags    *Tag      `json:",omitempty"`
}

type Index struct {
	Name   string
	Unique bool

	Fields FieldList
}

func (t *TableDescription) HasLastInsertId() bool {
	return t.Primary != nil && t.Primary.Type.Base.IsInteger()
}

func (t *TableDescription) HasPrimaryKey() bool {
	return t.Primary != nil
}

func (t *TableDescription) SafePrimary() Field {
	if t.Primary != nil {
		return *t.Primary
	}
	return Field{}
}

func (t *TableDescription) NumColumnNames(withAuto bool) int {
	num := 0
	for _, f := range t.Fields {
		if withAuto || !f.Tags.Auto {
			num++
		}
	}
	return num
}

func (table *TableDescription) ColumnNames(withAuto bool) Identifiers {
	if withAuto {
		return table.Fields.SqlNames()
	}
	return table.Fields.NoAuto().SqlNames()
}

func (t *TableDescription) SimpleFields() FieldList {
	list := make(FieldList, 0, len(t.Fields))
	for _, f := range t.Fields {
		if f.Encode == ENCNONE && f.IsExported() {
			switch f.Type.Base {
			case String, // Bool is not provided
				Int, Int8, Int16, Int32, Int64,
				Uint, Uint8, Uint16, Uint32, Uint64,
				Float32, Float64:
				list = append(list, f)
			}
		}
	}
	return list
}

//-------------------------------------------------------------------------------------------------

func (f *Field) IsExported() bool {
	name0 := f.Name[0]
	return 'A' <= name0 && name0 <= 'Z'
}

//-------------------------------------------------------------------------------------------------

func (i *Index) UniqueStr() string {
	if i.Unique {
		return "UNIQUE "
	}
	return ""
}

func (i *Index) JoinedNames(sep string) string {
	return i.Fields.Names().MkString(sep)
}

func (i *Index) Columns() string {
	return i.Fields.SqlNames().MkString(",")
}

func (i *Index) Single() bool {
	return len(i.Fields) == 1
}

//-------------------------------------------------------------------------------------------------

// Parts gets the node containment chain as a sequence of names of parts.
func (node *Node) Parts() []string {
	d := 0
	for n := node; n != nil; n = n.Parent {
		d++
	}

	p := make([]string, d)
	d--
	for n := node; n != nil; n = n.Parent {
		p[d] = n.Name
		d--
	}
	return p
}

func (node *Node) JoinParts(delta int, sep string) string {
	parts := node.Parts()
	if delta > 0 {
		parts = parts[:len(parts)-delta]
	}
	return strings.Join(parts, sep)
}

//-------------------------------------------------------------------------------------------------

type FieldList []*Field

func (list FieldList) DistinctTypes() []Type {
	types := NewTypeSet()

	for _, field := range list {
		types.Add(field.Type)
	}

	slice := types.ToSlice()
	sort.Slice(slice, func(i, j int) bool { return slice[i].Tag() < slice[j].Tag() })
	return slice
}

func (list FieldList) FormalParams() Identifiers {
	parts := make(Identifiers, len(list))
	for i, field := range list {
		parts[i] = fmt.Sprintf(`%s %s`, strings.ToLower(field.Name), field.Type.Type())
	}
	return parts
}

func (list FieldList) WhereClauses() Identifiers {
	parts := make(Identifiers, len(list))
	for i, field := range list {
		parts[i] = fmt.Sprintf(`where.Eq(%q, %s)`, field.SqlName, strings.ToLower(field.Name))
	}
	return parts
}

func (list FieldList) Names() Identifiers {
	ids := make(Identifiers, len(list))
	for i, field := range list {
		ids[i] = field.Name
	}
	return ids
}

func (list FieldList) SqlNames() Identifiers {
	ids := make(Identifiers, len(list))
	for i, field := range list {
		ids[i] = field.SqlName
	}
	return ids
}

func (list FieldList) Filter(predicate func(*Field) bool) FieldList {
	filtered := make(FieldList, 0, len(list))
	for _, field := range list {
		if predicate(field) {
			filtered = append(filtered, field)
		}
	}
	return filtered
}

// Pointers returns only the fields that have pointer types.
func (list FieldList) Pointers() FieldList {
	return list.Filter(func(field *Field) bool {
		return field.Type.IsPtr
	})
}

// Pointers returns only the fields that have non-pointer types.
func (list FieldList) NoPointers() FieldList {
	return list.Filter(func(field *Field) bool {
		return !field.Type.IsPtr
	})
}

// NoSkips returns only the fields without the skip flag set.
func (list FieldList) NoSkips() FieldList {
	return list.Filter(func(field *Field) bool {
		return field.Tags == nil || !field.Tags.Skip
	})
}

// NoPrimary returns all the fields except any marked primary.
func (list FieldList) NoPrimary() FieldList {
	return list.Filter(func(field *Field) bool {
		return field.Tags == nil || !field.Tags.Primary
	})
}

// NoPrimary returns all the fields except any marked auto-increment.
func (list FieldList) NoAuto() FieldList {
	return list.Filter(func(field *Field) bool {
		return field.Tags == nil || !field.Tags.Auto
	})
}
