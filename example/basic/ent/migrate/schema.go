// Code generated by entc, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// NotesColumns holds the columns for the "notes" table.
	NotesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "title", Type: field.TypeString},
		{Name: "content", Type: field.TypeString},
	}
	// NotesTable holds the schema information for the "notes" table.
	NotesTable = &schema.Table{
		Name:       "notes",
		Columns:    NotesColumns,
		PrimaryKey: []*schema.Column{NotesColumns[0]},
	}
	// UsersColumns holds the columns for the "users" table.
	UsersColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "email", Type: field.TypeString, Unique: true},
		{Name: "name", Type: field.TypeString},
	}
	// UsersTable holds the schema information for the "users" table.
	UsersTable = &schema.Table{
		Name:       "users",
		Columns:    UsersColumns,
		PrimaryKey: []*schema.Column{UsersColumns[0]},
	}
	// UserNotesColumns holds the columns for the "user_notes" table.
	UserNotesColumns = []*schema.Column{
		{Name: "user_id", Type: field.TypeInt},
		{Name: "note_id", Type: field.TypeInt},
	}
	// UserNotesTable holds the schema information for the "user_notes" table.
	UserNotesTable = &schema.Table{
		Name:       "user_notes",
		Columns:    UserNotesColumns,
		PrimaryKey: []*schema.Column{UserNotesColumns[0], UserNotesColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "user_notes_user_id",
				Columns:    []*schema.Column{UserNotesColumns[0]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "user_notes_note_id",
				Columns:    []*schema.Column{UserNotesColumns[1]},
				RefColumns: []*schema.Column{NotesColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		NotesTable,
		UsersTable,
		UserNotesTable,
	}
)

func init() {
	UserNotesTable.ForeignKeys[0].RefTable = UsersTable
	UserNotesTable.ForeignKeys[1].RefTable = NotesTable
}
