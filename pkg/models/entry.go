package models

type Entry struct {
	ID     			int    `db:"entry_id"`
	User
	Slug   			string `db:"slug"`
	DisplayName   	string `db:"display_name"`
}
