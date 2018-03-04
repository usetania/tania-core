package server

import (
	"errors"
	"net/http"

	"github.com/Tanibox/tania-server/src/tasks/domain"
	"github.com/Tanibox/tania-server/src/tasks/storage"
	"github.com/labstack/echo"
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
			log.Error(err)
		}

		taskReadFromRepo.Title = e.Title
		taskRead = taskReadFromRepo
	case domain.TaskDescriptionChanged:

		// Get TaskRead By UID
		taskReadFromRepo, err := s.getTaskReadFromID(e.UID)
		if err != nil {
			log.Error(err)
		}

		taskReadFromRepo.Description = e.Description
		taskRead = taskReadFromRepo
	case domain.TaskPriorityChanged:

		// Get TaskRead By UID
		taskReadFromRepo, err := s.getTaskReadFromID(e.UID)
		if err != nil {
			log.Error(err)
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
			log.Error(err)
		}

		taskReadFromRepo.Category = e.Category
		taskRead = &taskReadFromRepo
	case domain.TaskDetailsChanged:

		// Get TaskRead By UID
		taskReadFromRepo, err := s.getTaskReadFromID(e.UID)
		if err != nil {
			log.Error(err)
		}

		taskReadFromRepo.DomainDetails = e.DomainDetails
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

func (s *TaskServer) getTaskReadFromID(uid uuid.UUID) (*storage.TaskRead, error) {

	readResult := <-s.TaskReadQuery.FindByID(uid)

	taskReadFromRepo, ok := readResult.Result.(storage.TaskRead)

	if taskReadFromRepo.UID != uid {
		return &storage.TaskRead{}, domain.TaskError{domain.TaskErrorTaskNotFoundCode}
	}
	if !ok {
		return &storage.TaskRead{}, echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	} else {
		return &taskReadFromRepo, nil
	}
}
