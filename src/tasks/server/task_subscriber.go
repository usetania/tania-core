package server

import (
	"errors"
	"net/http"

	"github.com/Tanibox/tania-server/src/tasks/domain"
	"github.com/Tanibox/tania-server/src/tasks/storage"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
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
			log.Error(domain.TaskError{domain.TaskErrorTaskNotFoundCode})
		}
		if !ok {
			log.Error(echo.NewHTTPError(http.StatusBadRequest, "Internal server error"))
		}

		taskReadFromRepo.Title = e.Title
		taskReadFromRepo.Description = e.Description
		taskReadFromRepo.Priority = e.Priority
		taskReadFromRepo.DueDate = e.DueDate
		taskReadFromRepo.DomainDetails = e.DomainDetails
		taskReadFromRepo.Category = e.Category
		taskReadFromRepo.AssetID = e.AssetID
		taskRead = &taskReadFromRepo

	case domain.TaskCompleted:

		// Get TaskRead By UID
		readResult := <-s.TaskReadQuery.FindByID(e.UID)

		taskReadFromRepo, ok := readResult.Result.(storage.TaskRead)

		if taskReadFromRepo.UID != e.UID {
			log.Error(domain.TaskError{domain.TaskErrorTaskNotFoundCode})
		}
		if !ok {
			log.Error(echo.NewHTTPError(http.StatusBadRequest, "Internal server error"))
		}

		taskReadFromRepo.CompletedDate = e.CompletedDate
		taskReadFromRepo.Status = domain.TaskStatusCompleted
		taskRead = &taskReadFromRepo

	case domain.TaskCancelled:

		// Get TaskRead By UID
		readResult := <-s.TaskReadQuery.FindByID(e.UID)

		taskReadFromRepo, ok := readResult.Result.(storage.TaskRead)

		if taskReadFromRepo.UID != e.UID {
			log.Error(domain.TaskError{domain.TaskErrorTaskNotFoundCode})
		}
		if !ok {
			log.Error(echo.NewHTTPError(http.StatusBadRequest, "Internal server error"))
		}

		taskReadFromRepo.CancelledDate = e.CancelledDate
		taskReadFromRepo.Status = domain.TaskStatusCancelled
		taskRead = &taskReadFromRepo

	case domain.TaskDue:

		// Get TaskRead By UID
		readResult := <-s.TaskReadQuery.FindByID(e.UID)

		taskReadFromRepo, ok := readResult.Result.(storage.TaskRead)

		if taskReadFromRepo.UID != e.UID {
			log.Error(domain.TaskError{domain.TaskErrorTaskNotFoundCode})
		}
		if !ok {
			log.Error(echo.NewHTTPError(http.StatusBadRequest, "Internal server error"))
		}

		taskReadFromRepo.IsDue = true
		taskRead = &taskReadFromRepo

	default:
		log.Error(errors.New("Unknown task event"))
	}

	err := <-s.TaskReadRepo.Save(taskRead)
	if err != nil {
		log.Error(err)
	}

	return nil
}
