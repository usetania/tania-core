package repository

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/Tanibox/tania-server/src/assets/domain"
	"github.com/Tanibox/tania-server/src/assets/storage"
	"github.com/Tanibox/tania-server/src/helper/structhelper"
	uuid "github.com/satori/go.uuid"
)

type ReservoirEventRepository interface {
	Save(uid uuid.UUID, latestVersion int, events []interface{}) <-chan error
}

type ReservoirEventRepositoryInMemory struct {
	Storage *storage.ReservoirEventStorage
}

func NewReservoirEventRepositoryInMemory(s *storage.ReservoirEventStorage) ReservoirEventRepository {
	return &ReservoirEventRepositoryInMemory{Storage: s}
}

type ReservoirEventRepositorySqlite struct {
	DB *sql.DB
}

func NewReservoirEventRepositorySqlite(db *sql.DB) ReservoirEventRepository {
	return &ReservoirEventRepositorySqlite{DB: db}
}

func NewReservoirFromHistory(events []storage.ReservoirEvent) *domain.Reservoir {
	state := &domain.Reservoir{}
	for _, v := range events {
		state.Transition(v.Event)
		state.Version++
	}
	return state
}

// Save is to save
func (f *ReservoirEventRepositoryInMemory) Save(uid uuid.UUID, latestVersion int, events []interface{}) <-chan error {
	result := make(chan error)

	go func() {
		f.Storage.Lock.Lock()
		defer f.Storage.Lock.Unlock()

		for _, v := range events {
			latestVersion++
			f.Storage.ReservoirEvents = append(f.Storage.ReservoirEvents, storage.ReservoirEvent{
				ReservoirUID: uid,
				Version:      latestVersion,
				Event:        v,
			})
		}

		result <- nil

		close(result)
	}()

	return result
}

func (f *ReservoirEventRepositorySqlite) Save(uid uuid.UUID, latestVersion int, events []interface{}) <-chan error {
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

			e, err := json.Marshal(EventWrapper{
				EventName: structhelper.GetName(v),
				EventData: v,
			})

			if err != nil {
				panic(err)
			}

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
