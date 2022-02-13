package sqlite

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/tasks/decoder"
	"github.com/usetania/tania-core/src/tasks/query"
	"github.com/usetania/tania-core/src/tasks/storage"
)

type TaskEventQuerySqlite struct {
	DB *sql.DB
}

func NewTaskEventQuerySqlite(db *sql.DB) query.TaskEvent {
	return &TaskEventQuerySqlite{DB: db}
}

func (f *TaskEventQuerySqlite) FindAllByTaskID(uid uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		events := []storage.TaskEvent{}

		rows, err := f.DB.Query("SELECT * FROM TASK_EVENT WHERE TASK_UID = ? ORDER BY VERSION ASC", uid)
		if err != nil {
			result <- query.Result{Error: err}
		}

		rowsData := struct {
			ID          int
			TaskUID     string
			Version     int
			CreatedDate string
			Event       []byte
		}{}

		for rows.Next() {
			rows.Scan(&rowsData.ID, &rowsData.TaskUID, &rowsData.Version, &rowsData.CreatedDate, &rowsData.Event)

			wrapper := decoder.TaskEventWrapper{}
			json.Unmarshal(rowsData.Event, &wrapper)

			taskUID, err := uuid.FromString(rowsData.TaskUID)
			if err != nil {
				result <- query.Result{Error: err}
			}

			createdDate, err := time.Parse(time.RFC3339, rowsData.CreatedDate)
			if err != nil {
				result <- query.Result{Error: err}
			}

			events = append(events, storage.TaskEvent{
				TaskUID:     taskUID,
				Version:     rowsData.Version,
				CreatedDate: createdDate,
				Event:       wrapper.Data,
			})
		}

		result <- query.Result{Result: events}
		close(result)
	}()

	return result
}
