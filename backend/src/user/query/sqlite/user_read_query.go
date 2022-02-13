package sqlite

import (
	"database/sql"
	"errors"
	"time"

	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/user/query"
	"github.com/usetania/tania-core/src/user/storage"
	"golang.org/x/crypto/bcrypt"
)

type UserReadQuerySqlite struct {
	DB *sql.DB
}

func NewUserReadQuerySqlite(db *sql.DB) query.UserRead {
	return UserReadQuerySqlite{DB: db}
}

type userReadResult struct {
	UID         string
	Username    string
	Password    string
	CreatedDate string
	LastUpdated string
}

func (s UserReadQuerySqlite) FindByID(uid uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

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

		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			result <- query.Result{Error: err}
		}

		if errors.Is(err, sql.ErrNoRows) {
			result <- query.Result{Result: userRead}
		}

		userUID, err := uuid.FromString(rowsData.UID)
		if err != nil {
			result <- query.Result{Error: err}
		}

		createdDate, err := time.Parse(time.RFC3339, rowsData.CreatedDate)
		if err != nil {
			result <- query.Result{Error: err}
		}

		lastUpdated, err := time.Parse(time.RFC3339, rowsData.LastUpdated)
		if err != nil {
			result <- query.Result{Error: err}
		}

		userRead = storage.UserRead{
			UID:         userUID,
			Username:    rowsData.Username,
			Password:    []byte(rowsData.Password),
			CreatedDate: createdDate,
			LastUpdated: lastUpdated,
		}

		result <- query.Result{Result: userRead}
		close(result)
	}()

	return result
}

func (s UserReadQuerySqlite) FindByUsername(username string) <-chan query.Result {
	result := make(chan query.Result)

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

		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			result <- query.Result{Error: err}
		}

		if errors.Is(err, sql.ErrNoRows) {
			result <- query.Result{Result: userRead}
		}

		userUID, err := uuid.FromString(rowsData.UID)
		if err != nil {
			result <- query.Result{Error: err}
		}

		createdDate, err := time.Parse(time.RFC3339, rowsData.CreatedDate)
		if err != nil {
			result <- query.Result{Error: err}
		}

		lastUpdated, err := time.Parse(time.RFC3339, rowsData.LastUpdated)
		if err != nil {
			result <- query.Result{Error: err}
		}

		userRead = storage.UserRead{
			UID:         userUID,
			Username:    rowsData.Username,
			Password:    []byte(rowsData.Password),
			CreatedDate: createdDate,
			LastUpdated: lastUpdated,
		}

		result <- query.Result{Result: userRead}
		close(result)
	}()

	return result
}

func (s UserReadQuerySqlite) FindByUsernameAndPassword(username, password string) <-chan query.Result {
	result := make(chan query.Result)

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

		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			result <- query.Result{Error: err}
		}

		if errors.Is(err, sql.ErrNoRows) {
			result <- query.Result{Result: userRead}
		}

		err = bcrypt.CompareHashAndPassword([]byte(rowsData.Password), []byte(password))
		if err != nil {
			result <- query.Result{Result: userRead}
		}

		userUID, err := uuid.FromString(rowsData.UID)
		if err != nil {
			result <- query.Result{Error: err}
		}

		createdDate, err := time.Parse(time.RFC3339, rowsData.CreatedDate)
		if err != nil {
			result <- query.Result{Error: err}
		}

		lastUpdated, err := time.Parse(time.RFC3339, rowsData.LastUpdated)
		if err != nil {
			result <- query.Result{Error: err}
		}

		userRead = storage.UserRead{
			UID:         userUID,
			Username:    rowsData.Username,
			Password:    []byte(rowsData.Password),
			CreatedDate: createdDate,
			LastUpdated: lastUpdated,
		}

		result <- query.Result{Result: userRead}
		close(result)
	}()

	return result
}
