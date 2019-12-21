package repositories

import (
	"errors"
	"fmt"
	"github.com/atdean/onomatopoedia/pkg/models"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
)

type UserRepository struct {
	SqlPool *sqlx.DB
}

func NewUserRepository(sqlPool *sqlx.DB) *UserRepository {
	return &UserRepository{
		SqlPool: sqlPool,
	}
}

func (repo *UserRepository) CreateNewUser(username, email, password string) (*models.User, error) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 12)

	user := &models.User{
		Username: username,
		Email: email,
		Password: string(hashedPassword),
	}

	queryString := `
		INSERT INTO users (username, email, password)
		VALUES(?, ?, ?)
	`
	result, err := repo.SqlPool.Exec(queryString, user.Username, user.Email, user.Password)
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			log.Printf("MySQL Error: %s\n", mySQLError.Message)
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_email_uindex") {
				return nil, fmt.Errorf("Could not insert user: %s\n", mySQLError.Message)
			}
		}
		return nil, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	user.ID = int(userID)

	return user, nil
}

func (repo *UserRepository) GetByID(userID int) (*models.User, error) {
	queryString := `
		SELECT
			id AS user_id, username, email, password
		FROM users
		WHERE id = ?
	`
	row := repo.SqlPool.QueryRowx(queryString, userID)

	user := &models.User{}
	if err := row.StructScan(user); err != nil {
		log.Printf("UserRepository.GetByID: %s\n", err)
		return nil, err
	}

	return user, nil
}

func (repo *UserRepository) GetByUsername(username string) (*models.User, error) {
	queryString := `
		SELECT
			id AS user_id, username, email, password
		FROM users
		WHERE username = ?
	`
	row := repo.SqlPool.QueryRowx(queryString, username)

	user := &models.User{}
	if err := row.StructScan(user); err != nil {
		log.Printf("UserRepository.GetByUsername: %s\n", err)
		return nil, err
	}

	return user, nil
}