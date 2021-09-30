// Code generated by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/minskylab/ent-hasura/example/basic/ent/note"
)

// Note is the model entity for the Note schema.
type Note struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Title holds the value of the "title" field.
	Title string `json:"title,omitempty"`
	// Content holds the value of the "content" field.
	Content string `json:"content,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the NoteQuery when eager-loading is set.
	Edges NoteEdges `json:"edges"`
}

// NoteEdges holds the relations/edges for other nodes in the graph.
type NoteEdges struct {
	// Authors holds the value of the authors edge.
	Authors []*User `json:"authors,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// AuthorsOrErr returns the Authors value or an error if the edge
// was not loaded in eager-loading.
func (e NoteEdges) AuthorsOrErr() ([]*User, error) {
	if e.loadedTypes[0] {
		return e.Authors, nil
	}
	return nil, &NotLoadedError{edge: "authors"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Note) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case note.FieldID:
			values[i] = new(sql.NullInt64)
		case note.FieldTitle, note.FieldContent:
			values[i] = new(sql.NullString)
		default:
			return nil, fmt.Errorf("unexpected column %q for type Note", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Note fields.
func (n *Note) assignValues(columns []string, values []interface{}) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case note.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			n.ID = int(value.Int64)
		case note.FieldTitle:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field title", values[i])
			} else if value.Valid {
				n.Title = value.String
			}
		case note.FieldContent:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field content", values[i])
			} else if value.Valid {
				n.Content = value.String
			}
		}
	}
	return nil
}

// QueryAuthors queries the "authors" edge of the Note entity.
func (n *Note) QueryAuthors() *UserQuery {
	return (&NoteClient{config: n.config}).QueryAuthors(n)
}

// Update returns a builder for updating this Note.
// Note that you need to call Note.Unwrap() before calling this method if this Note
// was returned from a transaction, and the transaction was committed or rolled back.
func (n *Note) Update() *NoteUpdateOne {
	return (&NoteClient{config: n.config}).UpdateOne(n)
}

// Unwrap unwraps the Note entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (n *Note) Unwrap() *Note {
	tx, ok := n.config.driver.(*txDriver)
	if !ok {
		panic("ent: Note is not a transactional entity")
	}
	n.config.driver = tx.drv
	return n
}

// String implements the fmt.Stringer.
func (n *Note) String() string {
	var builder strings.Builder
	builder.WriteString("Note(")
	builder.WriteString(fmt.Sprintf("id=%v", n.ID))
	builder.WriteString(", title=")
	builder.WriteString(n.Title)
	builder.WriteString(", content=")
	builder.WriteString(n.Content)
	builder.WriteByte(')')
	return builder.String()
}

// Notes is a parsable slice of Note.
type Notes []*Note

func (n Notes) config(cfg config) {
	for _i := range n {
		n[_i].config = cfg
	}
}