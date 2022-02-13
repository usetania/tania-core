package sqlite

import (
	"database/sql"
	"time"

	"github.com/usetania/tania-core/src/user/repository"
	"github.com/usetania/tania-core/src/user/storage"
)

type UserReadRepositorySqlite struct {
	DB *sql.DB
}

func NewUserReadRepositorySqlite(db *sql.DB) repository.UserRead {
	return &UserReadRepositorySqlite{DB: db}
}

func (f *UserReadRepositorySqlite) Save(userRead *storage.UserRead) <-chan error {
	result := make(chan error)

	go func() {
		count := 0

		err := f.DB.QueryRow(`SELECT COUNT(*) FROM USER_READ WHERE UID = ?`, userRead.UID).Scan(&count)
		if err != nil {
			result <- err
		}

		if count > 0 {
			_, err := f.DB.Exec(`UPDATE USER_READ SET
				USERNAME = ?, PASSWORD = ?,
				CREATED_DATE = ?, LAST_UPDATED = ?
				WHERE UID = ?`,
				userRead.Username, userRead.Password,
				userRead.CreatedDate.Format(time.RFC3339), userRead.LastUpdated.Format(time.RFC3339),
				userRead.UID)
			if err != nil {
				result <- err
			}
		} else {
			_, err := f.DB.Exec(`INSERT INTO USER_READ
				(UID, USERNAME, PASSWORD, CREATED_DATE, LAST_UPDATED)
				VALUES (?, ?, ?, ?, ?)`,
				userRead.UID, userRead.Username, userRead.Password,
				userRead.CreatedDate.Format(time.RFC3339), userRead.LastUpdated.Format(time.RFC3339))
			if err != nil {
				result <- err
			}
		}

		result <- nil
		close(result)
	}()

	return result
}
