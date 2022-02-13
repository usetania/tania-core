package server

import (
	"errors"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/usetania/tania-core/src/tasks/domain"
	"github.com/usetania/tania-core/src/tasks/storage"
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
	case domain.TaskTitleChanged:
		// Get TaskRead By UID
		taskReadFromRepo, err := s.getTaskReadFromID(e.UID)
		if err != nil {
			return err
		}

		taskReadFromRepo.Title = e.Title
		taskRead = taskReadFromRepo
	case domain.TaskDescriptionChanged:
		// Get TaskRead By UID
		taskReadFromRepo, err := s.getTaskReadFromID(e.UID)
		if err != nil {
			return err
		}

		taskReadFromRepo.Description = e.Description
		taskRead = taskReadFromRepo
	case domain.TaskPriorityChanged:
		// Get TaskRead By UID
		taskReadFromRepo, err := s.getTaskReadFromID(e.UID)
		if err != nil {
			return err
		}

		taskReadFromRepo.Priority = e.Priority
		taskRead = taskReadFromRepo
	case domain.TaskDueDateChanged:
		// Get TaskRead By UID
		taskReadFromRepo, err := s.getTaskReadFromID(e.UID)
		if err != nil {
			return err
		}

		taskReadFromRepo.DueDate = e.DueDate
		taskRead = taskReadFromRepo
	case domain.TaskCategoryChanged:
		// Get TaskRead By UID
		taskReadFromRepo, err := s.getTaskReadFromID(e.UID)
		if err != nil {
			return err
		}

		taskReadFromRepo.Category = e.Category
		taskRead = taskReadFromRepo
	case domain.TaskDetailsChanged:
		// Get TaskRead By UID
		taskReadFromRepo, err := s.getTaskReadFromID(e.UID)
		if err != nil {
			return err
		}

		taskReadFromRepo.DomainDetails = e.DomainDetails
		taskRead = taskReadFromRepo

	case domain.TaskCompleted:
		// Get TaskRead By UID
		taskReadFromRepo, err := s.getTaskReadFromID(e.UID)
		if err != nil {
			return err
		}

		taskReadFromRepo.CompletedDate = e.CompletedDate
		taskReadFromRepo.Status = domain.TaskStatusCompleted
		taskRead = taskReadFromRepo

	case domain.TaskCancelled:
		// Get TaskRead By UID
		taskReadFromRepo, err := s.getTaskReadFromID(e.UID)
		if err != nil {
			return err
		}

		taskReadFromRepo.CancelledDate = e.CancelledDate
		taskReadFromRepo.Status = domain.TaskStatusCancelled
		taskRead = taskReadFromRepo

	case domain.TaskDue:
		// Get TaskRead By UID
		taskReadFromRepo, err := s.getTaskReadFromID(e.UID)
		if err != nil {
			return err
		}

		taskReadFromRepo.IsDue = true
		taskRead = taskReadFromRepo

	default:
		return errors.New("unknown task event")
	}

	err := <-s.TaskReadRepo.Save(taskRead)
	if err != nil {
		return err
	}

	return nil
}

func (s *TaskServer) getTaskReadFromID(uid uuid.UUID) (*storage.TaskRead, error) {
	readResult := <-s.TaskReadQuery.FindByID(uid)

	taskReadFromRepo, ok := readResult.Result.(storage.TaskRead)

	if taskReadFromRepo.UID != uid {
		return &storage.TaskRead{}, domain.TaskError{Code: domain.TaskErrorTaskNotFoundCode}
	}

	if !ok {
		return &storage.TaskRead{}, echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	return &taskReadFromRepo, nil
}
