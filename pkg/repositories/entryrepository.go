package repositories

import (
	"database/sql"
	"github.com/atdean/onomatopoedia/pkg/models"
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