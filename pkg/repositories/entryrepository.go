package repositories

import (
	"database/sql"
	"github.com/atdean/onomatopoedia/pkg/models"
	"log"
)

type EntryRepository struct {
	SqlPool *sql.DB
}

func (repo *EntryRepository) GetByID(entryID int) (*models.Entry, error) {
	entry := &models.Entry{}

	row := repo.SqlPool.QueryRow(`
		SELECT
			entries.id, entries.user_id, entries.slug, entries.display_name
		FROM entries
		WHERE entries.id = ?
	`, entryID)

	if err := row.Scan(&entry.ID, &entry.UserID, &entry.Slug, &entry.DisplayName); err != nil {
		return nil, err
	}

	return entry, nil
}

func (repo *EntryRepository) GetMostRecent(resultsPerPage int, page int) ([]*models.Entry, error) {
	entries := make([]*models.Entry, 0, resultsPerPage)

	rows, err := repo.SqlPool.Query(`
		SELECT
			entries.id, entries.user_id, entries.slug, entries.display_name
		FROM entries
		ORDER BY id DESC
		LIMIT ?,?
	`, (page - 1) * resultsPerPage, resultsPerPage)
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return entries, err
	}

	for rows.Next() {
		entry := &models.Entry{}

		if err := rows.Scan(
			&entry.ID,
			&entry.UserID,
			&entry.Slug,
			&entry.DisplayName,
		); err != nil {
			return entries, err
		}
		entries = append(entries, entry)
	}
	if err = rows.Err(); err != nil {
		return entries, err
	}

	return entries, nil
}