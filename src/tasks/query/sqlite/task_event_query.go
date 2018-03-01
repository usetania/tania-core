package sqlite

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/Tanibox/tania-server/src/tasks/domain"
	"github.com/Tanibox/tania-server/src/tasks/query"
	"github.com/Tanibox/tania-server/src/tasks/storage"
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

	switch wrapper.EventName {
	case "TaskCreated":
		e := domain.TaskCreated{}

		if v, ok := mapped["uid"]; ok {
			uid, err := makeUUID(v)
			if err != nil {
				return nil, err
			}

			e.UID = uid
		}
		if v, ok := mapped["title"]; ok {
			val := v.(string)
			e.Title = val
		}
		if v, ok := mapped["description"]; ok {
			val := v.(string)
			e.Description = val
		}
		if v, ok := mapped["created_date"]; ok {
			val, err := makeTime(v)
			if err != nil {
				return nil, err
			}

			e.CreatedDate = val
		}
		if v, ok := mapped["due_date"]; ok {
			val, err := makeTimePointer(v)
			if err != nil {
				return nil, err
			}

			e.DueDate = val
		}
		if v, ok := mapped["priority"]; ok {
			val := v.(string)
			e.Priority = val
		}
		if v, ok := mapped["status"]; ok {
			val := v.(string)
			e.Status = val
		}
		if v, ok := mapped["domain"]; ok {
			val := v.(string)
			e.Domain = val
		}
		if v, ok := mapped["domain_details"]; ok {
			domainDetails, err := makeDomainDetails(v, e.Domain)
			if err != nil {
				return nil, err
			}

			e.DomainDetails = domainDetails
		}
		if v, ok := mapped["category"]; ok {
			val := v.(string)
			e.Category = val
		}
		if v, ok := mapped["is_due"]; ok {
			val := v.(bool)
			e.IsDue = val
		}
		if v, ok := mapped["asset_id"]; ok {
			uid, err := makeUUID(v)
			if err != nil {
				return nil, err
			}

			e.AssetID = &uid
		}

		return e, nil

	case "TaskModified":
		e := domain.TaskModified{}

		if v, ok := mapped["uid"]; ok {
			uid, err := makeUUID(v)
			if err != nil {
				return nil, err
			}

			e.UID = uid
		}
		if v, ok := mapped["title"]; ok {
			val := v.(string)
			e.Title = val
		}
		if v, ok := mapped["description"]; ok {
			val := v.(string)
			e.Description = val
		}
		if v, ok := mapped["due_date"]; ok {
			val, err := makeTimePointer(v)
			if err != nil {
				return nil, err
			}

			e.DueDate = val
		}
		if v, ok := mapped["priority"]; ok {
			val := v.(string)
			e.Priority = val
		}
		if v, ok := mapped["domain"]; ok {
			val := v.(string)
			e.Domain = val
		}
		if v, ok := mapped["domain_details"]; ok {
			domainDetails, err := makeDomainDetails(v, e.Domain)
			if err != nil {
				return nil, err
			}

			e.DomainDetails = domainDetails
		}
		if v, ok := mapped["category"]; ok {
			val := v.(string)
			e.Category = val
		}
		if v, ok := mapped["asset_id"]; ok {
			uid, err := makeUUID(v)
			if err != nil {
				return nil, err
			}

			e.AssetID = &uid
		}

		return e, nil

	case "TaskCompleted":
		e := domain.TaskCompleted{}

		if v, ok := mapped["uid"]; ok {
			uid, err := makeUUID(v)
			if err != nil {
				return nil, err
			}

			e.UID = uid
		}
		if v, ok := mapped["status"]; ok {
			val := v.(string)
			e.Status = val
		}
		if v, ok := mapped["completed_date"]; ok {
			val, err := makeTimePointer(v)
			if err != nil {
				return nil, err
			}

			e.CompletedDate = val
		}

		return e, nil

	case "TaskCancelled":
		e := domain.TaskCancelled{}

		if v, ok := mapped["uid"]; ok {
			uid, err := makeUUID(v)
			if err != nil {
				return nil, err
			}

			e.UID = uid
		}
		if v, ok := mapped["status"]; ok {
			val := v.(string)
			e.Status = val
		}
		if v, ok := mapped["cancelled_date"]; ok {
			val, err := makeTimePointer(v)
			if err != nil {
				return nil, err
			}

			e.CancelledDate = val
		}

		return e, nil

	case "TaskDue":
		e := domain.TaskDue{}

		if v, ok := mapped["uid"]; ok {
			uid, err := makeUUID(v)
			if err != nil {
				return nil, err
			}

			e.UID = uid
		}

		return e, nil
	}

	return nil, nil
}
