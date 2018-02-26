package sqlite

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/Tanibox/tania-server/src/assets/repository"
	uuid "github.com/satori/go.uuid"
)

type FarmEventRepositorySqlite struct {
	DB *sql.DB
}

func NewFarmEventRepositorySqlite(db *sql.DB) repository.FarmEventRepository {
	return &FarmEventRepositorySqlite{DB: db}
}

func (f *FarmEventRepositorySqlite) Save(uid uuid.UUID, latestVersion int, events []interface{}) <-chan error {
	result := make(chan error)

	go func() {
		stmt, err := f.DB.Prepare(`INSERT INTO FARM_EVENT (FARM_UID, VERSION, CREATED_DATE, EVENTS) VALUES (?, ?, ?, ?)`)
		if err != nil {
			result <- err
			close(result)
		}

		latestVersion++
		em, err := json.Marshal(events)
		if err != nil {
			result <- err
			close(result)
		}

		_, err = stmt.Exec(uid, latestVersion, time.Now().Format(time.RFC3339), em)
		if err != nil {
			result <- err
			close(result)
		}

		result <- nil
		close(result)
	}()

	return result
}
