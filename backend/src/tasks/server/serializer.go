package server

import (
	"github.com/usetania/tania-core/src/tasks/domain"
	"github.com/usetania/tania-core/src/tasks/storage"
)

func MapTaskToTaskRead(task *domain.Task) *storage.TaskRead {
	taskRead := &storage.TaskRead{
		Title:         task.Title,
		UID:           task.UID,
		Description:   task.Description,
		CreatedDate:   task.CreatedDate,
		DueDate:       task.DueDate,
		CompletedDate: task.CompletedDate,
		CancelledDate: task.CancelledDate,
		Priority:      task.Priority,
		Status:        task.Status,
		Domain:        task.Domain,
		DomainDetails: task.DomainDetails,
		Category:      task.Category,
		IsDue:         task.IsDue,
		AssetID:       task.AssetID,
	}

	return taskRead
}
