package server

import (
	"fmt"
	"github.com/Tanibox/tania-server/src/tasks/domain"
	"github.com/Tanibox/tania-server/src/tasks/storage"
	"github.com/labstack/echo"
	"net/http"
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
	case domain.TaskModified:

		// Get TaskRead By UID
		readResult := <-s.TaskReadQuery.FindByID(e.UID)

		taskReadFromRepo, ok := readResult.Result.(storage.TaskRead)

		if taskReadFromRepo.UID != e.UID {
			return domain.TaskError{domain.TaskErrorTaskNotFoundCode}
		}
		if !ok {
			return echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
		}
		fmt.Println(taskReadFromRepo)

		taskReadFromRepo.Title = e.Title
		taskReadFromRepo.Description = e.Description
		taskReadFromRepo.Priority = e.Priority
		taskReadFromRepo.DueDate = e.DueDate
		taskReadFromRepo.DomainDetails = e.DomainDetails
		taskReadFromRepo.Category = e.Category
		taskReadFromRepo.AssetID = e.AssetID
		taskRead = &taskReadFromRepo

	case domain.TaskCompleted:
	case domain.TaskCancelled:
	case domain.TaskDue:
	default:
		fmt.Println("Unknown Task Event")
	}

	fmt.Println("Saving from task event")
	fmt.Println(taskRead)
	err := <-s.TaskReadRepo.Save(taskRead)
	if err != nil {
		return err
	}

	return nil
}
