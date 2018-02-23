package repository

import (
	"database/sql"
	"encoding/json"

	"github.com/Tanibox/tania-server/src/assets/domain"
	"github.com/Tanibox/tania-server/src/assets/storage"
	uuid "github.com/satori/go.uuid"
)

type FarmEventRepository interface {
	Save(uid uuid.UUID, latestVersion int, events []interface{}) <-chan error
}

type FarmEventRepositoryInMemory struct {
	Storage *storage.FarmEventStorage
}

func NewFarmEventRepositoryInMemory(s *storage.FarmEventStorage) FarmEventRepository {
	return &FarmEventRepositoryInMemory{Storage: s}
}

type FarmEventRepositorySqlite struct {
	DB *sql.DB
}

func NewFarmEventRepositorySqlite(db *sql.DB) FarmEventRepository {
	return &FarmEventRepositorySqlite{DB: db}
}

func NewFarmFromHistory(events []storage.FarmEvent) *domain.Farm {
	state := &domain.Farm{}
	for _, v := range events {
		state.Transition(v.Event)
		state.Version++
	}
	return state
}

// Save is to save
func (f *FarmEventRepositoryInMemory) Save(uid uuid.UUID, latestVersion int, events []interface{}) <-chan error {
	result := make(chan error)

	go func() {
		f.Storage.Lock.Lock()
		defer f.Storage.Lock.Unlock()

		for _, v := range events {
			latestVersion++
			f.Storage.FarmEvents = append(f.Storage.FarmEvents, storage.FarmEvent{
				FarmUID: uid,
				Version: latestVersion,
				Event:   v,
			})
		}

		result <- nil

		close(result)
	}()

	return result
}

func (f *FarmEventRepositorySqlite) Save(uid uuid.UUID, latestVersion int, events []interface{}) <-chan error {
	result := make(chan error)

	go func() {
		stmt, err := f.DB.Prepare(`INSERT INTO FARM_EVENT (FARM_UID, VERSION, EVENTS) VALUES (?, ?, ?)`)
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

		_, err = stmt.Exec(uid, latestVersion, em)
		if err != nil {
			result <- err
			close(result)
		}

		result <- nil
		close(result)
	}()

	return result
}
