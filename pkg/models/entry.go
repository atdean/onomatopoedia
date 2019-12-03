package models

// Entry represents a correlation to a single "word" in the "entries" table
type Entry struct {
	ID     int    `db:"id"`
	UserID int    `db:"user_id"`
	Slug   string `db:"slug"`
	Word   string `db:"word"`
}
