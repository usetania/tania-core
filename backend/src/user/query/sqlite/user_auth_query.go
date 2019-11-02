package sqlite

import (
	"database/sql"
	"time"

	"github.com/Tanibox/tania-core/src/user/query"
	"github.com/Tanibox/tania-core/src/user/storage"
	uuid "github.com/satori/go.uuid"
)

type UserAuthQuerySqlite struct {
	DB *sql.DB
}

func NewUserAuthQuerySqlite(db *sql.DB) query.UserAuthQuery {
	return UserAuthQuerySqlite{DB: db}
}

type userAuthResult struct {
	UserUID      string
	AccessToken  string
	TokenExpires int
	CreatedDate  string
	LastUpdated  string
}

func (s UserAuthQuerySqlite) FindByUserID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		userAuth := storage.UserAuth{}
		rowsData := userAuthResult{}

		err := s.DB.QueryRow(`SELECT USER_UID, ACCESS_TOKEN, TOKEN_EXPIRES, CREATED_DATE, LAST_UPDATED
			FROM USER_AUTH WHERE USER_UID = ?`, uid).Scan(
			&rowsData.UserUID,
			&rowsData.AccessToken,
			&rowsData.TokenExpires,
			&rowsData.CreatedDate,
			&rowsData.LastUpdated,
		)

		if err != nil && err != sql.ErrNoRows {
			result <- query.QueryResult{Error: err}
		}

		if err == sql.ErrNoRows {
			result <- query.QueryResult{Result: userAuth}
		}

		userUID, err := uuid.FromString(rowsData.UserUID)
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

		userAuth = storage.UserAuth{
			UserUID:      userUID,
			AccessToken:  rowsData.AccessToken,
			TokenExpires: rowsData.TokenExpires,
			CreatedDate:  createdDate,
			LastUpdated:  lastUpdated,
		}

		result <- query.QueryResult{Result: userAuth}
		close(result)
	}()

	return result
}
