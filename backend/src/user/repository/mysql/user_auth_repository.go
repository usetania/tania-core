package mysql

import (
	"database/sql"

	"github.com/usetania/tania-core/src/user/repository"
	"github.com/usetania/tania-core/src/user/storage"
)

type UserAuthRepositoryMysql struct {
	DB *sql.DB
}

func NewUserAuthRepositoryMysql(db *sql.DB) repository.UserAuth {
	return &UserAuthRepositoryMysql{DB: db}
}

func (s *UserAuthRepositoryMysql) Save(userAuth *storage.UserAuth) <-chan error {
	result := make(chan error)

	go func() {
		total := 0

		err := s.DB.QueryRow(`SELECT COUNT(USER_UID)
			FROM USER_AUTH WHERE USER_UID = ?`, userAuth.UserUID.Bytes()).Scan(&total)
		if err != nil {
			result <- err
		}

		if total > 0 {
			_, err := s.DB.Exec(`UPDATE USER_AUTH
				SET ACCESS_TOKEN = ?, TOKEN_EXPIRES = ?, CREATED_DATE = ?, LAST_UPDATED = ?
				WHERE USER_UID = ?`,
				userAuth.AccessToken, userAuth.TokenExpires,
				userAuth.CreatedDate, userAuth.LastUpdated,
				userAuth.UserUID.Bytes())
			if err != nil {
				result <- err
			}
		} else {
			_, err := s.DB.Exec(`INSERT INTO USER_AUTH
				(USER_UID, ACCESS_TOKEN, TOKEN_EXPIRES, CREATED_DATE, LAST_UPDATED)
				VALUES (?,?,?,?,?)`,
				userAuth.UserUID.Bytes(), userAuth.AccessToken, userAuth.TokenExpires,
				userAuth.CreatedDate, userAuth.LastUpdated)
			if err != nil {
				result <- err
			}
		}

		result <- nil
		close(result)
	}()

	return result
}
