// Code generated by entc, DO NOT EDIT.

package note

const (
	// Label holds the string label denoting the note type in the database.
	Label = "note"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldTitle holds the string denoting the title field in the database.
	FieldTitle = "title"
	// FieldContent holds the string denoting the content field in the database.
	FieldContent = "content"
	// EdgeAuthors holds the string denoting the authors edge name in mutations.
	EdgeAuthors = "authors"
	// Table holds the table name of the note in the database.
	Table = "notes"
	// AuthorsTable is the table that holds the authors relation/edge. The primary key declared below.
	AuthorsTable = "user_notes"
	// AuthorsInverseTable is the table name for the User entity.
	// It exists in this package in order to avoid circular dependency with the "user" package.
	AuthorsInverseTable = "users"
)

// Columns holds all SQL columns for note fields.
var Columns = []string{
	FieldID,
	FieldTitle,
	FieldContent,
}

var (
	// AuthorsPrimaryKey and AuthorsColumn2 are the table columns denoting the
	// primary key for the authors relation (M2M).
	AuthorsPrimaryKey = []string{"user_id", "note_id"}
)

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}