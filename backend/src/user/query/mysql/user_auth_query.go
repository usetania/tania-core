package mysql

import (
	"database/sql"
	"errors"
	"time"

	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/user/query"
	"github.com/usetania/tania-core/src/user/storage"
)

type UserAuthQueryMysql struct {
	DB *sql.DB
}

func NewUserAuthQueryMysql(db *sql.DB) query.UserAuth {
	return UserAuthQueryMysql{DB: db}
}

type userAuthResult struct {
	UserUID      []byte
	AccessToken  string
	TokenExpires int
	CreatedDate  time.Time
	LastUpdated  time.Time
}

func (s UserAuthQueryMysql) FindByUserID(uid uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

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

		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			result <- query.Result{Error: err}
		}

		if errors.Is(err, sql.ErrNoRows) {
			result <- query.Result{Result: userAuth}
		}

		userUID, err := uuid.FromBytes(rowsData.UserUID)
		if err != nil {
			result <- query.Result{Error: err}
		}

		userAuth = storage.UserAuth{
			UserUID:      userUID,
			AccessToken:  rowsData.AccessToken,
			TokenExpires: rowsData.TokenExpires,
			CreatedDate:  rowsData.CreatedDate,
			LastUpdated:  rowsData.LastUpdated,
		}

		result <- query.Result{Result: userAuth}
		close(result)
	}()

	return result
}
