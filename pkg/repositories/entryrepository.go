package repositories

import (
	"github.com/atdean/onomatopoedia/pkg/models"
	"github.com/jmoiron/sqlx"
	"log"
)

type EntryRepository struct {
	SqlPool *sqlx.DB
}

func (repo *EntryRepository) GetByID(entryID int) (*models.Entry, error) {
	entry := &models.Entry{}

	queryString := `
		SELECT
			entries.id, entries.user_id, entries.slug, entries.display_name
		FROM entries
		WHERE entries.id = ?
	`
	row := repo.SqlPool.QueryRowx(queryString, entryID)

	if err := row.StructScan(entry); err != nil {
		return nil, err
	}

	return entry, nil
}

func (repo *EntryRepository) GetMostRecent(resultsPerPage int, page int) ([]*models.Entry, error) {
	entries := make([]*models.Entry, 0, resultsPerPage)

	queryString := `
		SELECT
			entries.id, entries.user_id, entries.slug, entries.display_name
		FROM entries
		ORDER BY id DESC
		LIMIT ?,?
	`
	rows, err := repo.SqlPool.Queryx(queryString, (page - 1) * resultsPerPage, resultsPerPage)
	defer func() {
		if err := rows.Close(); err != nil {
			log.Fatalln("Could not close SQL connection.")
		}
	}()
	if err != nil {
		log.Println(err)
		return entries, err
	}

	for rows.Next() {
		entry := &models.Entry{}
		if err = rows.StructScan(entry); err != nil {
			return entries, err
		}
		entries = append(entries, entry)
	}
	if err = rows.Err(); err != nil {
		return entries, err
	}

	return entries, nil
}