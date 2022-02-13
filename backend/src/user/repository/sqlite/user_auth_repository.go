package sqlite

import (
	"database/sql"
	"time"

	"github.com/usetania/tania-core/src/user/repository"
	"github.com/usetania/tania-core/src/user/storage"
)

type UserAuthRepositorySqlite struct {
	DB *sql.DB
}

func NewUserAuthRepositorySqlite(db *sql.DB) repository.UserAuth {
	return &UserAuthRepositorySqlite{DB: db}
}

func (s *UserAuthRepositorySqlite) Save(userAuth *storage.UserAuth) <-chan error {
	result := make(chan error)

	go func() {
		total := 0

		err := s.DB.QueryRow(`SELECT COUNT(USER_UID)
			FROM USER_AUTH WHERE USER_UID = ?`, userAuth.UserUID).Scan(&total)
		if err != nil {
			result <- err
		}

		if total > 0 {
			_, err := s.DB.Exec(`UPDATE USER_AUTH
				SET ACCESS_TOKEN = ?, TOKEN_EXPIRES = ?, CREATED_DATE = ?, LAST_UPDATED = ?
				WHERE USER_UID = ?`,
				userAuth.AccessToken, userAuth.TokenExpires,
				userAuth.CreatedDate.Format(time.RFC3339), userAuth.LastUpdated.Format(time.RFC3339),
				userAuth.UserUID)
			if err != nil {
				result <- err
			}
		} else {
			_, err := s.DB.Exec(`INSERT INTO USER_AUTH
				(USER_UID, ACCESS_TOKEN, TOKEN_EXPIRES, CREATED_DATE, LAST_UPDATED)
				VALUES (?,?,?,?,?)`,
				userAuth.UserUID, userAuth.AccessToken, userAuth.TokenExpires,
				userAuth.CreatedDate.Format(time.RFC3339), userAuth.LastUpdated.Format(time.RFC3339))
			if err != nil {
				result <- err
			}
		}

		result <- nil
		close(result)
	}()

	return result
}
