package constraint

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/rickb777/sqlapi"
	"github.com/rickb777/sqlapi/dialect"
	"github.com/rickb777/sqlapi/schema"
	"github.com/rickb777/sqlapi/util"
)

// FkConstraints holds foreign key constraints.
type FkConstraints []FkConstraint

// Find finds the FkConstraint on a specified column.
func (cc FkConstraints) Find(fkColumn string) (FkConstraint, bool) {
	for _, fkc := range cc {
		if fkc.ForeignKeyColumn == fkColumn {
			return fkc, true
		}
	}
	return FkConstraint{}, false
}

//-------------------------------------------------------------------------------------------------

// Reference holds a table + column reference used by constraints.
// The table name should not include any schema or other prefix.
// The column may be blank, but the functionality is then reduced
// (there will be insufficient metadata to use Relationship methods).
type Reference struct {
	TableName string
	Column    string // only one column is supported
}

//-------------------------------------------------------------------------------------------------

// FkConstraint holds a pair of references and their update/delete consequences.
// ForeignKeyColumn is the 'owner' of the constraint.
type FkConstraint struct {
	ForeignKeyColumn string // only one column is supported
	Parent           Reference
	Update, Delete   Consequence
}

// FkConstraintOn constructs a foreign key constraint in a fluent style.
func FkConstraintOn(column string) FkConstraint {
	return FkConstraint{ForeignKeyColumn: column}
}

// FkConstraintOfField constructs a foreign key constraint from a struct field.
func FkConstraintOfField(field *schema.Field) FkConstraint {
	tags := field.GetTags()
	tbl, col := tags.ParentReference()
	return FkConstraint{
		ForeignKeyColumn: field.SqlName,
		Parent: Reference{
			TableName: tbl,
			Column:    col,
		},
		Update: Consequence(tags.OnUpdate),
		Delete: Consequence(tags.OnDelete),
	}
}

// RefersTo sets the parent reference. The column may be blank.
func (c FkConstraint) RefersTo(tableName string, column string) FkConstraint {
	c.Parent = Reference{tableName, column}
	return c
}

// OnUpdate sets the update consequence.
func (c FkConstraint) OnUpdate(consequence Consequence) FkConstraint {
	c.Update = consequence
	return c
}

// OnDelete sets the delete consequence.
func (c FkConstraint) OnDelete(consequence Consequence) FkConstraint {
	c.Delete = consequence
	return c
}

// ConstraintSql constructs the CONSTRAINT clause to be included in the CREATE TABLE.
func (c FkConstraint) ConstraintSql(q dialect.Quoter, name sqlapi.TableName, index int) string {
	return baseConstraintSql(q, name, index, c.sql(q, name.Prefix), "", "")
}

// Column constructs the foreign key clause needed to configure the database.
func (c FkConstraint) sql(q dialect.Quoter, prefix string) string {
	column := ""
	if c.Parent.Column != "" {
		column = " (" + q.Quote(c.Parent.Column) + ")"
	}
	return fmt.Sprintf("foreign key (%s) references %s%s%s%s",
		q.Quote(c.ForeignKeyColumn), q.Quote(prefix+c.Parent.TableName), column,
		c.Update.Apply(" ", "update"),
		c.Delete.Apply(" ", "delete"))
}

func (c FkConstraint) GoString() string {
	return fmt.Sprintf(`constraint.FkConstraint{"%s", constraint.Reference{"%s", "%s"}, "%s", "%s"}`,
		c.ForeignKeyColumn, c.Parent.TableName, c.Parent.Column, c.Update, c.Delete)
}

//func (c FkConstraint) AlterTable() AlterTable {
//	return AlterTable{c.Child.TableName, c.ConstraintSql(0)}
//}

// NoCascade changes both the Update and Delete consequences to NoAction.
func (c FkConstraint) NoCascade() FkConstraint {
	c.Update = NoAction
	c.Delete = NoAction
	return c
}

// RelationshipWith constructs the Relationship that is expressed by the parent reference in
// the FkConstraint and the child's foreign key.
//
// The table names do not include any prefix.
func (c FkConstraint) RelationshipWith(child sqlapi.TableName) Relationship {
	return Relationship{
		Parent: c.Parent,
		Child:  Reference{child.Name, c.ForeignKeyColumn},
	}
}

//-------------------------------------------------------------------------------------------------

// Relationship represents a parent-child relationship.
// Only simple keys are supported (compound keys are not supported).
type Relationship struct {
	Parent, Child Reference
}

// IdsUnusedAsForeignKeys finds all the primary keys in the parent table that have no foreign key
// in the dependent (child) table. The table tbl provides the database or transaction handle; either
// the parent or the child table can be used for thi purpose.
func (rel Relationship) IdsUnusedAsForeignKeys(tbl sqlapi.Table) (util.Int64Set, error) {
	if rel.Parent.Column == "" || rel.Child.Column == "" {
		return nil, errors.Errorf("IdsUnusedAsForeignKeys requires the column names to be specified")
	}

	// TODO benchmark two candidates and choose the better
	// http://stackoverflow.com/questions/3427353/sql-statement-question-how-to-retrieve-records-of-a-table-where-the-primary-ke?rq=1
	//	s := fmt.Sprintf(
	//		`SELECT a.%s
	//			FROM %s a
	//			WHERE NOT EXISTS (
	//   				SELECT 1 FROM %s b
	//   				WHERE %s.%s = %s.%s
	//			)`,
	//		primary.ForeignKeyColumn, primary.TableName, foreign.TableName, primary.TableName, primary.ForeignKeyColumn, foreign.TableName, foreign.ForeignKeyColumn)

	// http://stackoverflow.com/questions/13108587/selecting-primary-keys-that-does-not-has-foreign-keys-in-another-table
	pfx := tbl.Name().Prefix
	s := fmt.Sprintf(
		`SELECT a.%s
			FROM %s%s a
			LEFT OUTER JOIN %s%s b ON a.%s = b.%s
			WHERE b.%s IS null`,
		rel.Parent.Column,
		pfx, rel.Parent.TableName,
		pfx, rel.Child.TableName,
		rel.Parent.Column, rel.Child.Column,
		rel.Child.Column)
	return fetchIds(tbl, s)
}

// IdsUsedAsForeignKeys finds all the primary keys in the parent table that have at least one foreign key
// in the dependent (child) table.
func (rel Relationship) IdsUsedAsForeignKeys(tbl sqlapi.Table) (util.Int64Set, error) {
	if rel.Parent.Column == "" || rel.Child.Column == "" {
		return nil, errors.Errorf("IdsUsedAsForeignKeys requires the column names to be specified")
	}

	pfx := tbl.Name().Prefix
	s := fmt.Sprintf(
		`SELECT DISTINCT a.%s AS Id
			FROM %s%s a
			INNER JOIN %s%s b ON a.%s = b.%s`,
		rel.Parent.Column,
		pfx, rel.Parent.TableName,
		pfx, rel.Child.TableName,
		rel.Parent.Column, rel.Child.Column)
	return fetchIds(tbl, s)
}

func fetchIds(tbl sqlapi.Table, query string) (util.Int64Set, error) {
	rows, err := tbl.Query(query)
	if err != nil {
		return nil, tbl.Database().LogIfError(errors.Wrap(err, query))
	}
	defer rows.Close()

	set := util.NewInt64Set()
	for rows.Next() {
		var id int64
		rows.Scan(&id)
		set.Add(id)
	}
	return set, tbl.Database().LogIfError(rows.Err())
}
