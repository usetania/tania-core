package mysql

import (
	"database/sql"
	"errors"
	"time"

	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/user/query"
	"github.com/usetania/tania-core/src/user/storage"
	"golang.org/x/crypto/bcrypt"
)

type UserReadQueryMysql struct {
	DB *sql.DB
}

func NewUserReadQueryMysql(db *sql.DB) query.UserRead {
	return UserReadQueryMysql{DB: db}
}

type userReadResult struct {
	UID         []byte
	Username    string
	Password    string
	CreatedDate time.Time
	LastUpdated time.Time
}

func (s UserReadQueryMysql) FindByID(uid uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		userRead := storage.UserRead{}
		rowsData := userReadResult{}

		err := s.DB.QueryRow("SELECT * FROM USER_READ WHERE UID = ?", uid.Bytes()).Scan(
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

		userUID, err := uuid.FromBytes(rowsData.UID)
		if err != nil {
			result <- query.Result{Error: err}
		}

		userRead = storage.UserRead{
			UID:         userUID,
			Username:    rowsData.Username,
			Password:    []byte(rowsData.Password),
			CreatedDate: rowsData.CreatedDate,
			LastUpdated: rowsData.LastUpdated,
		}

		result <- query.Result{Result: userRead}
		close(result)
	}()

	return result
}

func (s UserReadQueryMysql) FindByUsername(username string) <-chan query.Result {
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

		userUID, err := uuid.FromBytes(rowsData.UID)
		if err != nil {
			result <- query.Result{Error: err}
		}

		userRead = storage.UserRead{
			UID:         userUID,
			Username:    rowsData.Username,
			Password:    []byte(rowsData.Password),
			CreatedDate: rowsData.CreatedDate,
			LastUpdated: rowsData.LastUpdated,
		}

		result <- query.Result{Result: userRead}
		close(result)
	}()

	return result
}

func (s UserReadQueryMysql) FindByUsernameAndPassword(username, password string) <-chan query.Result {
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

		userUID, err := uuid.FromBytes(rowsData.UID)
		if err != nil {
			result <- query.Result{Error: err}
		}

		userRead = storage.UserRead{
			UID:         userUID,
			Username:    rowsData.Username,
			Password:    []byte(rowsData.Password),
			CreatedDate: rowsData.CreatedDate,
			LastUpdated: rowsData.LastUpdated,
		}

		result <- query.Result{Result: userRead}
		close(result)
	}()

	return result
}
