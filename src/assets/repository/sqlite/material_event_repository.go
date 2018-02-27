package sqlite

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/Tanibox/tania-server/src/assets/repository"
	"github.com/Tanibox/tania-server/src/helper/structhelper"
	uuid "github.com/satori/go.uuid"
)

type MaterialEventRepositorySqlite struct {
	DB *sql.DB
}

func NewMaterialEventRepositorySqlite(db *sql.DB) repository.MaterialEventRepository {
	return &MaterialEventRepositorySqlite{DB: db}
}

func (f *MaterialEventRepositorySqlite) Save(uid uuid.UUID, latestVersion int, events []interface{}) <-chan error {
	result := make(chan error)

	go func() {
		for _, v := range events {
			stmt, err := f.DB.Prepare(`INSERT INTO MATERIAL_EVENT (MATERIAL_UID, VERSION, CREATED_DATE, EVENT) VALUES (?, ?, ?, ?)`)
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
