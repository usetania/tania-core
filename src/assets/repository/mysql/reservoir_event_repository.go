package mysql

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/Tanibox/tania-core/src/assets/decoder"
	"github.com/Tanibox/tania-core/src/assets/repository"
	"github.com/Tanibox/tania-core/src/helper/structhelper"
	"github.com/gofrs/uuid"
)

type ReservoirEventRepositoryMysql struct {
	DB *sql.DB
}

func NewReservoirEventRepositoryMysql(db *sql.DB) repository.ReservoirEvent {
	return &ReservoirEventRepositoryMysql{DB: db}
}

func (f *ReservoirEventRepositoryMysql) Save(uid uuid.UUID, latestVersion int, events []interface{}) <-chan error {
	result := make(chan error)

	go func() {
		for _, v := range events {
			latestVersion++

			stmt, err := f.DB.Prepare(`INSERT INTO RESERVOIR_EVENT
				(RESERVOIR_UID, VERSION, CREATED_DATE, EVENT)
				VALUES (?, ?, ?, ?)`)
			if err != nil {
				result <- err
			}

			e, err := json.Marshal(decoder.EventWrapper{
				EventName: structhelper.GetName(v),
				EventData: v,
			})
			if err != nil {
				panic(err)
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
