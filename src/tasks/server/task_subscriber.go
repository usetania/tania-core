package server

import (
	"github.com/Tanibox/tania-server/src/tasks/domain"
	"github.com/Tanibox/tania-server/src/tasks/storage"
)

func (s *TaskServer) SaveToTaskReadModel(event interface{}) error {
	taskRead := &storage.TaskRead{}

	switch e := event.(type) {
	case domain.TaskCreated:

		taskRead.Title = e.Title
		taskRead.UID = e.UID
		taskRead.Description = e.Description
		taskRead.CreatedDate = e.CreatedDate
		taskRead.DueDate = e.DueDate
		taskRead.Priority = e.Priority
		taskRead.Status = e.Status
		taskRead.Domain = e.Domain
		taskRead.DomainDetails = e.DomainDetails
		taskRead.Category = e.Category
		taskRead.IsDue = e.IsDue
		taskRead.AssetID = e.AssetID

	}

	err := <-s.TaskReadRepo.Save(taskRead)
	if err != nil {
		return err
	}

	return nil
}
