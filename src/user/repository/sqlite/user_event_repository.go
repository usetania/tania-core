package sqlite

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/Tanibox/tania-core/src/helper/structhelper"
	"github.com/Tanibox/tania-core/src/user/decoder"
	"github.com/Tanibox/tania-core/src/user/repository"
	uuid "github.com/satori/go.uuid"
)

type UserEventRepositorySqlite struct {
	DB *sql.DB
}

func NewUserEventRepositorySqlite(db *sql.DB) repository.UserEventRepository {
	return &UserEventRepositorySqlite{DB: db}
}

func (f *UserEventRepositorySqlite) Save(uid uuid.UUID, latestVersion int, events []interface{}) <-chan error {
	result := make(chan error)

	go func() {
		for _, v := range events {
			stmt, err := f.DB.Prepare(`INSERT INTO USER_EVENT
				(USER_UID, VERSION, CREATED_DATE, EVENT) VALUES (?, ?, ?, ?)`)
			if err != nil {
				result <- err
			}

			latestVersion++

			e, err := json.Marshal(decoder.EventWrapper{
				EventName: structhelper.GetName(v),
				EventData: v,
			})

			_, err = stmt.Exec(uid, latestVersion, time.Now().Format(time.RFC3339), e)
			if err != nil {
				result <- err
			}
		}

		result <- nil
		close(result)
	}()

	return result
}
