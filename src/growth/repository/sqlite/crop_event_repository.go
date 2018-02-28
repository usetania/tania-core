package sqlite

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/Tanibox/tania-server/src/growth/repository"
	"github.com/Tanibox/tania-server/src/helper/structhelper"
	uuid "github.com/satori/go.uuid"
)

type CropEventRepositorySqlite struct {
	DB *sql.DB
}

func NewCropEventRepositorySqlite(db *sql.DB) repository.CropEventRepository {
	return &CropEventRepositorySqlite{DB: db}
}

func (f *CropEventRepositorySqlite) Save(uid uuid.UUID, latestVersion int, events []interface{}) <-chan error {
	result := make(chan error)

	go func() {
		for _, v := range events {
			stmt, err := f.DB.Prepare(`INSERT INTO CROP_EVENT (CROP_UID, VERSION, CREATED_DATE, EVENT) VALUES (?, ?, ?, ?)`)
			if err != nil {
				result <- err
			}

			latestVersion++

			e, err := json.Marshal(repository.EventWrapper{
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
