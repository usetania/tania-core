package mysql

import (
	"database/sql"
	"time"

	"github.com/Tanibox/tania-core/src/user/query"
	"github.com/Tanibox/tania-core/src/user/storage"
	uuid "github.com/satori/go.uuid"
)

type UserAuthQueryMysql struct {
	DB *sql.DB
}

func NewUserAuthQueryMysql(db *sql.DB) query.UserAuthQuery {
	return UserAuthQueryMysql{DB: db}
}

type userAuthResult struct {
	UserUID      []byte
	AccessToken  string
	TokenExpires int
	CreatedDate  time.Time
	LastUpdated  time.Time
}

func (s UserAuthQueryMysql) FindByUserID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		userAuth := storage.UserAuth{}
		rowsData := userAuthResult{}

		err := s.DB.QueryRow(`SELECT USER_UID, ACCESS_TOKEN, TOKEN_EXPIRES, CREATED_DATE, LAST_UPDATED
			FROM USER_AUTH WHERE USER_UID = ?`, uid.Bytes()).Scan(
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

		userUID, err := uuid.FromBytes(rowsData.UserUID)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		userAuth = storage.UserAuth{
			UserUID:      userUID,
			AccessToken:  rowsData.AccessToken,
			TokenExpires: rowsData.TokenExpires,
			CreatedDate:  rowsData.CreatedDate,
			LastUpdated:  rowsData.LastUpdated,
		}

		result <- query.QueryResult{Result: userAuth}
		close(result)
	}()

	return result
}
