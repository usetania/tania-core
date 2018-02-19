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
		Title         string            `json:"title"`
		UID           uuid.UUID         `json:"uid"`
		Description   string            `json:"description"`
		CreatedDate   time.Time         `json:"created_date"`
		DueDate       *time.Time        `json:"due_date, omitempty"`
		CompletedDate *time.Time        `json:"completed_date"`
		Priority      string            `json:"priority"`
		Status        string            `json:"status"`
		Domain        string            `json:"domain"`
		DomainDetails domain.TaskDomain `json:"domain_details"`
		Category      string            `json:"category"`
		IsDue         bool              `json:"is_due"`
		AssetID       *uuid.UUID        `json:"asset_id"`

		// Events
		Version            int           `json:"int"`
		UncommittedChanges []interface{} `json:"events"`
	}{
		Title:              st.Title,
		UID:                st.UID,
		Description:        st.Description,
		CreatedDate:        st.CreatedDate,
		DueDate:            st.DueDate,
		CompletedDate:      st.CompletedDate,
		Priority:           st.Priority,
		Status:             st.Status,
		Domain:             st.DomainDetails.Code(),
		DomainDetails:      st.DomainDetails,
		Category:           st.Category,
		IsDue:              st.IsDue,
		AssetID:            st.AssetID,
		Version:            st.Version,
		UncommittedChanges: st.UncommittedChanges,
	})
}
