// Package constraint provides types and methods to support foreign-key relationshipd between database
// tables.
//
// Only simple keys are supported, which can be integers, strings or any other suitable type.
// Compound keys are not supported.
package constraint

import (
	"fmt"
	"github.com/rickb777/sqlapi"
)

type Dialect interface {
	Quote(column string) string
}

// Constraint represents data that augments the data-definition SQL statements such as CREATE TABLE.
type Constraint interface {
	// ConstraintSql constructs the CONSTRAINT clause to be included in the CREATE TABLE.
	ConstraintSql(dialect Dialect, name sqlapi.TableName, index int) string

	// Expresses the constraint as a constructor + literals for the API type.
	GoString() string
}

//-------------------------------------------------------------------------------------------------

// Constraints holds constraints.
type Constraints []Constraint

// FkConstraints returns only the foreign key constraints in the Constraints slice.
func (cc Constraints) FkConstraints() FkConstraints {
	list := make(FkConstraints, 0, len(cc))
	for _, c := range cc {
		if fkc, ok := c.(FkConstraint); ok {
			list = append(list, fkc)
		}
	}
	return list
}

//-------------------------------------------------------------------------------------------------

// CheckConstraint holds an expression that refers to table columns and is applied as a precondition
// whenever a table insert, update or delete is attempted. The CheckConstraint expression is in SQL.
type CheckConstraint struct {
	Expression string
}

// ConstraintSql constructs the CONSTRAINT clause to be included in the CREATE TABLE.
func (c CheckConstraint) ConstraintSql(dialect Dialect, name sqlapi.TableName, index int) string {
	return fmt.Sprintf("CONSTRAINT %s_c%d CHECK (%s)", name, index, c.Expression)
}

func (c CheckConstraint) GoString() string {
	panic("not implemented")
}
