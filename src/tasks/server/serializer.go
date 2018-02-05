package server

import (
	"encoding/json"
	"time"

	domain "github.com/Tanibox/tania-server/src/tasks/domain"
	uuid "github.com/satori/go.uuid"
)

type SimpleTask domain.Task

func MapToSimpleTask(Tasks []domain.Task) []SimpleTask {
	TaskList := make([]SimpleTask, len(Tasks))

	for i, Task := range Tasks {
		TaskList[i] = SimpleTask(Task)
	}

	return TaskList
}

func (st SimpleTask) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		UID         uuid.UUID  `json:"uid"`
		Description string     `json:"description"`
		CreatedDate time.Time  `json:"created_date"`
		DueDate     *time.Time `json:"due_date, omitempty"`
		Priority    string     `json:"priority"`
		Status      string     `json:"status"`
		TaskType    string     `json:"type"`
		IsDue       bool       `json:"is_due"`
		AssetID     uuid.UUID  `json:"asset_id"`
	}{
		UID:         st.UID,
		Description: st.Description,
		CreatedDate: st.CreatedDate,
		DueDate:     st.DueDate,
		Priority:    st.Priority,
		Status:      st.Status,
		TaskType:    st.TaskType,
		IsDue:       st.IsDue,
		AssetID:     st.AssetID,
	})
}
