package sqlite

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/Tanibox/tania-core/src/user/decoder"
	"github.com/Tanibox/tania-core/src/user/query"
	"github.com/Tanibox/tania-core/src/user/storage"
	uuid "github.com/satori/go.uuid"
)

type UserEventQuerySqlite struct {
	DB *sql.DB
}

func NewUserEventQuerySqlite(db *sql.DB) query.UserEventQuery {
	return &UserEventQuerySqlite{DB: db}
}

func (f *UserEventQuerySqlite) FindAllByID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		events := []storage.UserEvent{}

		rows, err := f.DB.Query("SELECT * FROM USER_EVENT WHERE USER_UID = ? ORDER BY VERSION ASC", uid)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		rowsData := struct {
			ID          int
			UserUID     string
			Version     int
			CreatedDate string
			Event       []byte
		}{}

		for rows.Next() {
			rows.Scan(&rowsData.ID, &rowsData.UserUID, &rowsData.Version, &rowsData.CreatedDate, &rowsData.Event)

			wrapper := decoder.UserEventWrapper{}
			err := json.Unmarshal(rowsData.Event, &wrapper)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			userUID, err := uuid.FromString(rowsData.UserUID)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			createdDate, err := time.Parse(time.RFC3339, rowsData.CreatedDate)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			events = append(events, storage.UserEvent{
				UserUID:     userUID,
				Version:     rowsData.Version,
				CreatedDate: createdDate,
				Event:       wrapper.EventData,
			})
		}

		result <- query.QueryResult{Result: events}
		close(result)
	}()

	return result
}
