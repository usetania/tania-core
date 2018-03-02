package sqlite

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/Tanibox/tania-server/src/tasks/domain"
	"github.com/Tanibox/tania-server/src/tasks/query"
	"github.com/Tanibox/tania-server/src/tasks/storage"
	"github.com/Tanibox/tania-server/src/tasks/util/decoder"
	"github.com/mitchellh/mapstructure"
	uuid "github.com/satori/go.uuid"
)

type TaskEventQuerySqlite struct {
	DB *sql.DB
}

func NewTaskEventQuerySqlite(db *sql.DB) query.TaskEventQuery {
	return &TaskEventQuerySqlite{DB: db}
}

func (f *TaskEventQuerySqlite) FindAllByTaskID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		events := []storage.TaskEvent{}

		rows, err := f.DB.Query("SELECT * FROM TASK_EVENT WHERE TASK_UID = ? ORDER BY VERSION ASC", uid)
		if err != nil {
			result <- query.QueryResult{Error: err}
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

			wrapper := query.EventWrapper{}
			json.Unmarshal(rowsData.Event, &wrapper)

			event, err := assertTaskEvent(wrapper)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			taskUID, err := uuid.FromString(rowsData.TaskUID)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			createdDate, err := time.Parse(time.RFC3339, rowsData.CreatedDate)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			events = append(events, storage.TaskEvent{
				TaskUID:     taskUID,
				Version:     rowsData.Version,
				CreatedDate: createdDate,
				Event:       event,
			})
		}

		result <- query.QueryResult{Result: events}
		close(result)
	}()

	return result
}

func assertTaskEvent(wrapper query.EventWrapper) (interface{}, error) {
	mapped := wrapper.EventData.(map[string]interface{})

	f := mapstructure.ComposeDecodeHookFunc(
		decoder.UIDHook(),
		decoder.TimeHook(time.RFC3339),
		decoder.TaskDomainDetailHook(),
	)

	switch wrapper.EventName {
	case "TaskCreated":
		e := domain.TaskCreated{}

		decoder.Decode(f, &mapped, &e)

		return e, nil

	case "TaskModified":
		e := domain.TaskModified{}

		decoder.Decode(f, &mapped, &e)

		return e, nil

	case "TaskCompleted":
		e := domain.TaskCompleted{}

		decoder.Decode(f, &mapped, &e)

		return e, nil

	case "TaskCancelled":
		e := domain.TaskCancelled{}

		decoder.Decode(f, &mapped, &e)

		return e, nil

	case "TaskDue":
		e := domain.TaskDue{}

		decoder.Decode(f, &mapped, &e)

		return e, nil
	}

	return nil, errors.New("Event not decoded succesfully")
}
