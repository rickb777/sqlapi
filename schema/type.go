package schema

import (
	"fmt"

	. "github.com/rickb777/sqlapi/types"
)

type Type struct {
	Name      string // name of source code type.
	PkgPath   string `json:",omitempty" yaml:",omitempty"` // package name (full path)
	PkgName   string `json:",omitempty" yaml:",omitempty"` // package name (short name)
	IsPtr     bool   `json:",omitempty" yaml:",omitempty"`
	IsScanner bool   `json:",omitempty" yaml:",omitempty"`
	IsValuer  bool   `json:",omitempty" yaml:",omitempty"`
	Base      Kind   // underlying source code kind.
}

func (t Type) Tag() string {
	if t.IsPtr {
		return t.Name + "Ptr"
	}
	return t.Name
}

func (t Type) Star() string {
	if t.IsPtr {
		return "*"
	}
	return ""
}

func (t Type) Type() string {
	if len(t.PkgName) > 0 {
		return fmt.Sprintf("%s.%s", t.PkgName, t.Name)
	} else {
		return t.Name
	}
}

func (t Type) NullableValue() string {
	if t.IsPtr {
		switch t.Base {
		case String:
			return "String"
		case Int, Int8, Int16, Int32, Int64,
			Uint, Uint8, Uint16, Uint32, Uint64:
			return "Int64"
		case Float32, Float64:
			return "Float64"
		case Bool:
			return "Bool"
		}
	}
	return ""
}

func (t Type) IsBasicType() bool {
	return t.PkgPath == "" && t.Base.String() == t.Name
}

func (t Type) String() string {
	return fmt.Sprintf("%s%s (%v)", t.Star(), t.Type(), t.Base)
}
