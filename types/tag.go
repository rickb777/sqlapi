package types

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
)

const TagKey = "sql"

// Tag stores the parsed data from the tag string in
// a struct field. These are all optional.
type Tag struct {
	Name       string `json:",omitempty" yaml:"name"`     // explicit column name
	Type       string `json:",omitempty" yaml:"type"`     // explicit column type (SQL syntax)
	Default    string `json:",omitempty" yaml:"default"`  // default SQL value
	Prefixed   bool   `json:",omitempty" yaml:"prefixed"` // use struct nesting to name the column
	Primary    bool   `json:",omitempty" yaml:"pk"`       // is a primary key
	Natural    bool   `json:",omitempty" yaml:"nk"`       // is a natural key so a unique index will be added automatically
	Auto       bool   `json:",omitempty" yaml:"auto"`     // is auto-incremented
	Index      string `json:",omitempty" yaml:"index"`    // the name of an index
	Unique     string `json:",omitempty" yaml:"unique"`   // the name of a unique index
	ForeignKey string `json:",omitempty" yaml:"fk"`       // relationship to another table
	OnUpdate   string `json:",omitempty" yaml:"onupdate"` // what to do on update (no action, cascade, delete, restrict, set null, set default)
	OnDelete   string `json:",omitempty" yaml:"ondelete"` // what to do on delete
	Size       int    `json:",omitempty" yaml:"size"`     // storage size
	Encode     string `json:",omitempty" yaml:"encode"`   // used for struct types: one of json | text | driver
	Skip       bool   `json:",omitempty" yaml:"skip"`     // ignore the field
	// TODO Check      string `yaml:"check"` // specify SQL constraint checks
}

func (tag *Tag) ParentReference() (string, string) {
	if tag == nil || tag.ForeignKey == "" {
		return "", ""
	}
	ss := strings.Split(tag.ForeignKey, ".")
	if len(ss) < 2 {
		return ss[0], ""
	}
	return ss[0], ss[1]
}

//-------------------------------------------------------------------------------------------------

func inSet(s string, set ...string) bool {
	for _, x := range set {
		if s == x {
			return true
		}
	}
	return false
}

func normalizeConsequence(consequence string) string {
	if consequence == "setnull" {
		return "set null"
	}
	if consequence == "setdefault" {
		return "set default"
	}
	return consequence
}

func normalize(tag *Tag) *Tag {
	tag.OnUpdate = normalizeConsequence(tag.OnUpdate)
	tag.OnDelete = normalizeConsequence(tag.OnDelete)
	return tag
}

func (tag *Tag) validate() error {
	sep := ""
	buf := &bytes.Buffer{}

	if tag.Primary && tag.Natural {
		io.WriteString(buf, sep)
		io.WriteString(buf, "primary key cannot also be a natural key")
		sep = "; "
	}

	if tag.Auto && tag.Natural {
		io.WriteString(buf, sep)
		io.WriteString(buf, "natural key cannot use auto-increment")
		sep = "; "
	}

	if tag.Auto && !tag.Primary {
		io.WriteString(buf, sep)
		io.WriteString(buf, "auto-increment can only be used on primary keys")
		sep = "; "
	}

	if tag.Natural && tag.Index != "" {
		io.WriteString(buf, sep)
		io.WriteString(buf, "natural key cannot be used with index")
		sep = "; "
	}

	if tag.Natural && tag.Unique != "" {
		io.WriteString(buf, sep)
		io.WriteString(buf, "natural key should not be used with unique")
		sep = "; "
	}

	if tag.Size < 0 {
		io.WriteString(buf, sep)
		fmt.Fprintf(buf, "size cannot be negative (%d)", tag.Size)
		sep = "; "
	}

	if !inSet(tag.Encode, "", "json", "text", "driver") {
		io.WriteString(buf, sep)
		fmt.Fprintf(buf, "unrecognised encode value %q", tag.Encode)
		sep = "; "
	}

	if tag.ForeignKey != "" {
		if tag.Primary || tag.Natural {
			io.WriteString(buf, sep)
			io.WriteString(buf, "foreign key cannot also be a primary key nor a natural key")
			sep = "; "
		}

		slice := strings.Split(tag.ForeignKey, ".")
		if len(slice) < 1 || 2 < len(slice) {
			io.WriteString(buf, sep)
			fmt.Fprintf(buf, "fk value (%q) must be in 'tablename' or 'tablename.column' form", tag.ForeignKey)
			sep = "; "
		}
	}

	if !inSet(tag.OnUpdate, "", "cascade", "delete", "restrict", "set null", "set default") {
		io.WriteString(buf, sep)
		fmt.Fprintf(buf, "unrecognised onupdate value %q", tag.OnUpdate)
		sep = "; "
	}

	if !inSet(tag.OnDelete, "", "cascade", "delete", "restrict", "set null", "set default") {
		io.WriteString(buf, sep)
		fmt.Fprintf(buf, "unrecognised ondelete value %q", tag.OnDelete)
		sep = "; "
	}

	if buf.Len() > 0 {
		return errors.Errorf(buf.String())
	}
	return nil
}

var zero = Tag{}

func (t *Tag) checkZero() *Tag {
	if *t == zero {
		return nil
	}
	return t
}

// ParseTag parses a tag string from the struct
// field and unmarshals into a Tag struct.
func ParseTag(raw string) (*Tag, error) {
	var tag = new(Tag)

	raw = strings.TrimSpace(topAndTail(strings.TrimSpace(raw)))
	structTag := reflect.StructTag(raw)
	value := strings.TrimSpace(structTag.Get(TagKey))

	if value == "-" {
		tag.Skip = true
		return tag, nil
	}

	// wrap the string in curly braces so that we can use the Yaml parser.
	yamlValue := fmt.Sprintf("{ %s }", value)

	// unmarshals the Yaml formatted string into the Tag structure.
	var err = yaml.Unmarshal([]byte(yamlValue), tag)
	if err != nil {
		return tag, errors.Wrapf(err, "parse tag YAML %q", raw)
	}

	normalize(tag)
	return tag.checkZero(), tag.validate()
}

//-------------------------------------------------------------------------------------------------

type Tags map[string]*Tag

func (tags Tags) String() string {
	b := &bytes.Buffer{}
	for n, t := range tags {
		fmt.Fprintf(b, "%-10s: %+v\n", n, t)
	}
	return b.String()
}

func topAndTail(s string) string {
	last := len(s) - 1
	if len(s) >= 2 && s[0] == s[last] {
		return s[1:last]
	}
	return s
}

func ReadTagsFile(file string) (Tags, error) {
	if file == "" {
		return nil, nil
	}

	f, err := os.Open(file)
	if err != nil {
		return nil, errors.Wrapf(err, "opening tags file %s", file)
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, errors.Wrapf(err, "reading tags file %s", file)
	}

	tags := make(Tags)

	err = yaml.Unmarshal(b, tags)
	if err != nil {
		return nil, errors.Wrapf(err, "parsing YAML tags file %s", file)
	}

	//DevInfo("tags from %s\n%s\n", file, tags)
	return tags, nil
}
