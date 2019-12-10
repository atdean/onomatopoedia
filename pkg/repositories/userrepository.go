package repositories

import (
	"github.com/atdean/onomatopoedia/pkg/models"
	"github.com/jmoiron/sqlx"
	"log"
)

type UserRepository struct {
	SqlPool *sqlx.DB
}

func NewUserRepository(sqlPool *sqlx.DB) *UserRepository {
	return &UserRepository{
		SqlPool: sqlPool,
	}
}

func (repo *UserRepository) GetByID(userID int) (*models.User, error) {
	user := &models.User{}

	queryString := `
		SELECT
			id AS user_id, username, email, password
		FROM users
		WHERE id = ?
	`
	row := repo.SqlPool.QueryRowx(queryString, userID)

	if err := row.StructScan(user); err != nil {
		log.Printf("UserRepository.GetByID: %s\n", err)
		return nil, err
	}

	return user, nil
}

func (repo *UserRepository) GetByUsername(username string) (*models.User, error) {
	user := &models.User{}

	queryString := `
		SELECT
			id AS user_id, username, email, password
		FROM users
		WHERE username = ?
	`
	row := repo.SqlPool.QueryRowx(queryString, username)

	if err := row.StructScan(user); err != nil {
		log.Printf("UserRepository.GetByUsername: %s\n", err)
		return nil, err
	}

	return user, nil
}