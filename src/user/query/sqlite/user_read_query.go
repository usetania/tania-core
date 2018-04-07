package sqlite

import (
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/Tanibox/tania-core/src/user/query"
	"github.com/Tanibox/tania-core/src/user/storage"
	uuid "github.com/satori/go.uuid"
)

type UserReadQuerySqlite struct {
	DB *sql.DB
}

func NewUserReadQuerySqlite(db *sql.DB) query.UserReadQuery {
	return UserReadQuerySqlite{DB: db}
}

type userReadResult struct {
	UID         string
	Username    string
	Password    string
	CreatedDate string
	LastUpdated string
}

func (s UserReadQuerySqlite) FindByID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		userRead := storage.UserRead{}
		rowsData := userReadResult{}

		err := s.DB.QueryRow("SELECT * FROM USER_READ WHERE UID = ?", uid).Scan(
			&rowsData.UID,
			&rowsData.Username,
			&rowsData.Password,
			&rowsData.CreatedDate,
			&rowsData.LastUpdated,
		)

		if err != nil && err != sql.ErrNoRows {
			result <- query.QueryResult{Error: err}
		}

		if err == sql.ErrNoRows {
			result <- query.QueryResult{Result: userRead}
		}

		userUID, err := uuid.FromString(rowsData.UID)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		createdDate, err := time.Parse(time.RFC3339, rowsData.CreatedDate)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		lastUpdated, err := time.Parse(time.RFC3339, rowsData.LastUpdated)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		userRead = storage.UserRead{
			UID:         userUID,
			Username:    rowsData.Username,
			Password:    []byte(rowsData.Password),
			CreatedDate: createdDate,
			LastUpdated: lastUpdated,
		}

		result <- query.QueryResult{Result: userRead}
		close(result)
	}()

	return result
}

func (s UserReadQuerySqlite) FindByUsername(username string) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		userRead := storage.UserRead{}
		rowsData := userReadResult{}

		err := s.DB.QueryRow("SELECT * FROM USER_READ WHERE USERNAME = ?", username).Scan(
			&rowsData.UID,
			&rowsData.Username,
			&rowsData.Password,
			&rowsData.CreatedDate,
			&rowsData.LastUpdated,
		)

		if err != nil && err != sql.ErrNoRows {
			result <- query.QueryResult{Error: err}
		}

		if err == sql.ErrNoRows {
			result <- query.QueryResult{Result: userRead}
		}

		userUID, err := uuid.FromString(rowsData.UID)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		createdDate, err := time.Parse(time.RFC3339, rowsData.CreatedDate)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		lastUpdated, err := time.Parse(time.RFC3339, rowsData.LastUpdated)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		userRead = storage.UserRead{
			UID:         userUID,
			Username:    rowsData.Username,
			Password:    []byte(rowsData.Password),
			CreatedDate: createdDate,
			LastUpdated: lastUpdated,
		}

		result <- query.QueryResult{Result: userRead}
		close(result)
	}()

	return result
}

func (s UserReadQuerySqlite) FindByUsernameAndPassword(username, password string) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		userRead := storage.UserRead{}
		rowsData := userReadResult{}

		err := s.DB.QueryRow(`SELECT * FROM USER_READ
			WHERE USERNAME = ?`, username).Scan(
			&rowsData.UID,
			&rowsData.Username,
			&rowsData.Password,
			&rowsData.CreatedDate,
			&rowsData.LastUpdated,
		)

		if err != nil && err != sql.ErrNoRows {
			result <- query.QueryResult{Error: err}
		}

		if err == sql.ErrNoRows {
			result <- query.QueryResult{Result: userRead}
		}

		err = bcrypt.CompareHashAndPassword([]byte(rowsData.Password), []byte(password))
		if err != nil {
			result <- query.QueryResult{Result: userRead}
		}

		userUID, err := uuid.FromString(rowsData.UID)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		createdDate, err := time.Parse(time.RFC3339, rowsData.CreatedDate)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		lastUpdated, err := time.Parse(time.RFC3339, rowsData.LastUpdated)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		userRead = storage.UserRead{
			UID:         userUID,
			Username:    rowsData.Username,
			Password:    []byte(rowsData.Password),
			CreatedDate: createdDate,
			LastUpdated: lastUpdated,
		}

		result <- query.QueryResult{Result: userRead}
		close(result)
	}()

	return result
}
