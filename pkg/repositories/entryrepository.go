package repositories

import (
	"github.com/atdean/onomatopoedia/pkg/models"
	"github.com/jmoiron/sqlx"
	"log"
)

type EntryRepository struct {
	SqlPool *sqlx.DB
}

func NewEntryRepository(sqlPool *sqlx.DB) *EntryRepository {
	return &EntryRepository {
		SqlPool: sqlPool,
	}
}

func (repo *EntryRepository) GetByID(entryID int) (*models.Entry , error) {
	entry := &models.Entry{}

	queryString := `
		SELECT
			entries.id AS entry_id, entries.user_id, entries.slug, entries.display_name,
			users.username
		FROM entries
		INNER JOIN users ON entries.user_id = users.id
		WHERE entries.id = ?
	`
	row := repo.SqlPool.QueryRowx(queryString, entryID)

	if err := row.StructScan(entry); err != nil {
		log.Println(err)
		return nil, err
	}

	return entry, nil
}

func (repo *EntryRepository) GetBySlug(entrySlug string) (*models.Entry, error) {
	queryString := `
		SELECT
			entries.id AS entry_id, entries.user_id, entries.slug, entries.display_name,
			users.username
		FROM entries
		INNER JOIN users ON entries.user_id = users.id
		WHERE entries.slug = ?
	`
	row := repo.SqlPool.QueryRowx(queryString, entrySlug)

	entry := &models.Entry{}
	if err := row.StructScan(entry); err != nil {
		log.Println(err)
		return nil, err
	}

	return entry, nil
}

func (repo *EntryRepository) GetMostRecent(resultsPerPage int, page int) ([]*models.Entry, error) {
	queryString := `
		SELECT
			entries.id AS entry_id, entries.user_id, entries.slug, entries.display_name,
			users.username
		FROM entries
		INNER JOIN users ON entries.user_id = users.id
		ORDER BY entries.id DESC
		LIMIT ?,?
	`
	rows, err := repo.SqlPool.Queryx(queryString, (page - 1) * resultsPerPage, resultsPerPage)
	defer closeConnection(rows)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	entries := make([]*models.Entry, 0, resultsPerPage)
	for rows.Next() {
		entry := &models.Entry{}
		if err = rows.StructScan(entry); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}