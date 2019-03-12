package where

import (
	"bytes"
	"fmt"
	"github.com/rickb777/sqlapi/dialect"
	"reflect"
	"strings"
)

const (
	WhereAdverb = "WHERE "
	HavingVerb  = "HAVING "
)

// Expression is an element in a WHERE clause. Expressions may be nested in various ways.
type Expression interface {
	fmt.Stringer
	And(Expression) Clause
	Or(Expression) Clause
	Build() (string, []interface{})
}

// Where constructs the sql clause beginning "WHERE ...". It will contain '?' style placeholders;
// these need to be passed through the relevant dialect ReplacePlaceholders processing.
func Where(wh Expression) (string, []interface{}) {
	return Build(WhereAdverb, wh)
}

// Having constructs the sql clause beginning "HAVING ...". It will contain '?' style placeholders;
// these need to be passed through the relevant dialect ReplacePlaceholders processing.
func Having(wh Expression) (string, []interface{}) {
	return Build(HavingVerb, wh)
}

// Build constructs the sql clause beginning with some verb/adverb. It will contain '?' style placeholders;
// these need to be passed through the relevant dialect ReplacePlaceholders processing.
func Build(adverb string, wh Expression) (string, []interface{}) {
	if wh == nil {
		return "", nil
	}
	sql, args := wh.Build()
	if sql == "" {
		return "", nil
	}
	return adverb + sql, args
}

//-------------------------------------------------------------------------------------------------

type not struct {
	expression Expression
}

func (not not) Build() (string, []interface{}) {
	sql, args := not.expression.Build()
	if sql == "" {
		return "", nil
	}
	return "NOT (" + sql + ")", args
}

func (not not) String() string {
	sql, args := not.Build()
	return insertLiteralValues(sql, args)
}

//-------------------------------------------------------------------------------------------------

// Condition is a simple condition such as an equality test.
type Condition struct {
	Column, Predicate string
	Args              []interface{}
}

func (cl Condition) Build() (string, []interface{}) {
	sql := dialect.Quote(cl.Column) + cl.Predicate

	var args []interface{}
	for _, arg := range cl.Args {
		value := reflect.ValueOf(arg)
		switch value.Kind() {
		case reflect.Array, reflect.Slice:
			for j := 0; j < value.Len(); j++ {
				args = append(args, value.Index(j).Interface())
			}

		default:
			args = append(args, arg)
		}
	}
	return sql, args
}

func (cl Condition) String() string {
	sql, args := cl.Build()
	return insertLiteralValues(sql, args)
}

//-------------------------------------------------------------------------------------------------

// Clause is a compound expression.
type Clause struct {
	wheres      []Expression
	conjunction string
}

func (wh Clause) Build() (string, []interface{}) {
	if len(wh.wheres) == 0 {
		return "", nil
	}

	var sqls []string
	var args []interface{}

	for _, where := range wh.wheres {
		sql, a2 := where.Build()
		if len(sql) > 0 {
			sqls = append(sqls, "("+sql+")")
			args = append(args, a2...)
		}
	}

	sql := strings.Join(sqls, wh.conjunction)
	return sql, args
}

func (wh Clause) String() string {
	sql, args := wh.Build()
	return insertLiteralValues(sql, args)
}

//-------------------------------------------------------------------------------------------------

func insertLiteralValues(sql string, args []interface{}) string {
	buf := &bytes.Buffer{}
	idx := 0
	for _, r := range sql {
		if r == '?' && idx < len(args) {
			v := args[idx]
			t := reflect.TypeOf(v)
			switch t.Kind() {
			case reflect.Bool,
				reflect.Int,
				reflect.Int8,
				reflect.Int16,
				reflect.Int32,
				reflect.Int64,
				reflect.Uint,
				reflect.Uint8,
				reflect.Uint16,
				reflect.Uint32,
				reflect.Uint64,
				reflect.Uintptr,
				reflect.Float32,
				reflect.Float64:
				buf.WriteString(fmt.Sprintf(`%v`, v))
			default:
				buf.WriteString(fmt.Sprintf(`'%v'`, v))
			}
			idx++
		} else {
			buf.WriteRune(r)
		}
	}
	return buf.String()
}
