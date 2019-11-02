package mysql

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/Tanibox/tania-core/src/tasks/decoder"

	"github.com/Tanibox/tania-core/src/tasks/query"
	"github.com/Tanibox/tania-core/src/tasks/storage"
	uuid "github.com/satori/go.uuid"
)

type TaskEventQueryMysql struct {
	DB *sql.DB
}

func NewTaskEventQueryMysql(db *sql.DB) query.TaskEventQuery {
	return &TaskEventQueryMysql{DB: db}
}

func (f *TaskEventQueryMysql) FindAllByTaskID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		events := []storage.TaskEvent{}

		rows, err := f.DB.Query("SELECT * FROM TASK_EVENT WHERE TASK_UID = ? ORDER BY VERSION ASC", uid.Bytes())
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		rowsData := struct {
			ID          int
			TaskUID     []byte
			Version     int
			CreatedDate time.Time
			Event       []byte
		}{}

		for rows.Next() {
			rows.Scan(&rowsData.ID, &rowsData.TaskUID, &rowsData.Version, &rowsData.CreatedDate, &rowsData.Event)

			wrapper := decoder.TaskEventWrapper{}
			json.Unmarshal(rowsData.Event, &wrapper)

			taskUID, err := uuid.FromBytes(rowsData.TaskUID)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			events = append(events, storage.TaskEvent{
				TaskUID:     taskUID,
				Version:     rowsData.Version,
				CreatedDate: rowsData.CreatedDate,
				Event:       wrapper.Data,
			})
		}

		result <- query.QueryResult{Result: events}
		close(result)
	}()

	return result
}
