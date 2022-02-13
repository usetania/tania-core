package mysql

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/user/decoder"
	"github.com/usetania/tania-core/src/user/query"
	"github.com/usetania/tania-core/src/user/storage"
)

type UserEventQueryMysql struct {
	DB *sql.DB
}

func NewUserEventQueryMysql(db *sql.DB) query.UserEvent {
	return &UserEventQueryMysql{DB: db}
}

func (f *UserEventQueryMysql) FindAllByID(uid uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		events := []storage.UserEvent{}

		rows, err := f.DB.Query("SELECT * FROM USER_EVENT WHERE USER_UID = ? ORDER BY VERSION ASC", uid.Bytes())
		if err != nil {
			result <- query.Result{Error: err}
		}

		rowsData := struct {
			ID          int
			UserUID     []byte
			Version     int
			CreatedDate time.Time
			Event       []byte
		}{}

		for rows.Next() {
			rows.Scan(&rowsData.ID, &rowsData.UserUID, &rowsData.Version, &rowsData.CreatedDate, &rowsData.Event)

			wrapper := decoder.UserEventWrapper{}

			err := json.Unmarshal(rowsData.Event, &wrapper)
			if err != nil {
				result <- query.Result{Error: err}
			}

			userUID, err := uuid.FromBytes(rowsData.UserUID)
			if err != nil {
				result <- query.Result{Error: err}
			}

			events = append(events, storage.UserEvent{
				UserUID:     userUID,
				Version:     rowsData.Version,
				CreatedDate: rowsData.CreatedDate,
				Event:       wrapper.EventData,
			})
		}

		result <- query.Result{Result: events}
		close(result)
	}()

	return result
}
