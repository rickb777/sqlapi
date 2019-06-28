package constraint

import "fmt"

// Consequence is the action to be performed after updating or deleting a record constrained by foreign key.
type Consequence string

const (
	// unspecified option is available but its semantics vary by DB vendor, so it's not included here.
	NoAction   Consequence = "no action"
	Restrict   Consequence = "restrict"
	Cascade    Consequence = "cascade"
	SetNull    Consequence = "set null"
	SetDefault Consequence = "set default"
	Delete     Consequence = "delete" // not MySQL
)

// Apply constructs the SQL sub-clause for a consequence of a specified action.
// The prefix is typically arbitrary whitespace.
func (c Consequence) Apply(pfx, action string) string {
	if c == "" {
		return "" // implicitly equivalent to NoAction
	}
	return fmt.Sprintf("%son %s %s", pfx, action, c)
}
