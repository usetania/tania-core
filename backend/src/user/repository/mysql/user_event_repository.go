package mysql

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/helper/structhelper"
	"github.com/usetania/tania-core/src/user/decoder"
	"github.com/usetania/tania-core/src/user/repository"
)

type UserEventRepositoryMysql struct {
	DB *sql.DB
}

func NewUserEventRepositoryMysql(db *sql.DB) repository.UserEvent {
	return &UserEventRepositoryMysql{DB: db}
}

func (f *UserEventRepositoryMysql) Save(uid uuid.UUID, latestVersion int, events []interface{}) <-chan error {
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
			if err != nil {
				result <- err
			}

			_, err = stmt.Exec(uid.Bytes(), latestVersion, time.Now(), e)
			if err != nil {
				result <- err
			}
		}

		result <- nil
		close(result)
	}()

	return result
}
