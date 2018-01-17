package server

import (
	"encoding/json"
	"time"

	domain "github.com/Tanibox/tania-server/src/tasks/domain"
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
		UID			string 	`json:"uid"`
		Description	string		`json:"description"`
		CreatedDate	time.Time	`json:"createddate"`
		DueDate		time.Time	`json:"duedate"`
		Priority	string		`json:"priority"`
		Status		string		`json:"status"`
		TaskType	string		`json:"type"`
		AssetID		string		`json:"assetid"`
	}{
		UID:  st.UID.String(),
		Description: st.Description,
		CreatedDate: st.CreatedDate,
		DueDate: st.DueDate,
		Priority: st.Priority,
		Status: st.Status,
		TaskType: st.TaskType,
		AssetID: st.AssetID.String(),
	})
}
