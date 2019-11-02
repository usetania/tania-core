package domain

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type TaskService interface {
	FindAreaByID(uid uuid.UUID) ServiceResult
	FindCropByID(uid uuid.UUID) ServiceResult
	FindMaterialByID(uid uuid.UUID) ServiceResult
	FindReservoirByID(uid uuid.UUID) ServiceResult
}

// ServiceResult is the container for service result
type ServiceResult struct {
	Result interface{}
	Error  error
}

type Task struct {
	UID           uuid.UUID  `json:"uid"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	CreatedDate   time.Time  `json:"created_date"`
	DueDate       *time.Time `json:"due_date, omitempty"`
	CompletedDate *time.Time `json:"completed_date"`
	CancelledDate *time.Time `json:"cancelled_date"`
	Priority      string     `json:"priority"`
	Status        string     `json:"status"`
	Domain        string     `json:"domain"`
	DomainDetails TaskDomain `json:"domain_details"`
	Category      string     `json:"category"`
	IsDue         bool       `json:"is_due"`
	AssetID       *uuid.UUID `json:"asset_id"`

	// Events
	Version            int
	UncommittedChanges []interface{}
}

// CreateTask
func CreateTask(taskService TaskService, title string, description string, duedate *time.Time, priority string, taskdomain TaskDomain, taskcategory string, assetid *uuid.UUID) (*Task, error) {
	// add validation

	err := validateTaskTitle(title)
	if err != nil {
		return &Task{}, err
	}

	err = validateTaskDescription(description)
	if err != nil {
		return &Task{}, err
	}

	err = validateTaskDueDate(duedate)
	if err != nil {
		return &Task{}, err
	}

	err = validateTaskPriority(priority)
	if err != nil {
		return &Task{}, err
	}

	err = validateTaskCategory(taskcategory)
	if err != nil {
		return &Task{}, err
	}

	err = validateAssetID(taskService, assetid, taskdomain.Code())
	if err != nil {
		return &Task{}, err
	}

	uid, err := uuid.NewV4()
	if err != nil {
		return &Task{}, err
	}

	initial := &Task{}

	initial.TrackChange(taskService, TaskCreated{
		Title:         title,
		UID:           uid,
		Description:   description,
		CreatedDate:   time.Now(),
		DueDate:       duedate,
		Priority:      priority,
		Status:        TaskStatusCreated,
		Domain:        taskdomain.Code(),
		DomainDetails: taskdomain,
		Category:      taskcategory,
		IsDue:         false,
		AssetID:       assetid,
	})

	return initial, nil
}

func (t *Task) ChangeTaskTitle(taskService TaskService, title string) (*Task, error) {
	err := validateTaskTitle(title)
	if err != nil {
		return &Task{}, err
	}

	event := TaskTitleChanged{
		UID:   t.UID,
		Title: title,
	}

	t.TrackChange(taskService, event)

	return t, nil
}

func (t *Task) ChangeTaskDescription(taskService TaskService, description string) (*Task, error) {
	err := validateTaskDescription(description)
	if err != nil {
		return &Task{}, err
	}

	event := TaskDescriptionChanged{
		UID:         t.UID,
		Description: description,
	}

	t.TrackChange(taskService, event)

	return t, nil
}

func (t *Task) ChangeTaskDueDate(taskService TaskService, duedate *time.Time) (*Task, error) {
	err := validateTaskDueDate(duedate)
	if err != nil {
		return &Task{}, err
	}

	event := TaskDueDateChanged{
		UID:     t.UID,
		DueDate: duedate,
	}

	t.TrackChange(taskService, event)

	return t, nil
}

func (t *Task) ChangeTaskPriority(taskService TaskService, priority string) (*Task, error) {
	err := validateTaskPriority(priority)
	if err != nil {
		return &Task{}, err
	}

	event := TaskPriorityChanged{
		UID:      t.UID,
		Priority: priority,
	}

	t.TrackChange(taskService, event)

	return t, nil
}

func (t *Task) ChangeTaskCategory(taskService TaskService, category string) (*Task, error) {
	err := validateTaskCategory(category)
	if err != nil {
		return &Task{}, err
	}

	event := TaskCategoryChanged{
		UID:      t.UID,
		Category: category,
	}

	t.TrackChange(taskService, event)

	return t, nil
}

func (t *Task) ChangeTaskDetails(taskService TaskService, details TaskDomain) (*Task, error) {

	event := TaskDetailsChanged{
		UID:           t.UID,
		DomainDetails: details,
	}

	t.TrackChange(taskService, event)

	return t, nil
}

// SetTaskAsDue
func (t *Task) SetTaskAsDue(taskService TaskService) {
	t.TrackChange(taskService, TaskDue{
		UID: t.UID,
	})
}

// CompleteTask
func (t *Task) CompleteTask(taskService TaskService) {
	completedTime := time.Now()

	t.TrackChange(taskService, TaskCompleted{
		UID:           t.UID,
		Status:        TaskCompletedCode,
		CompletedDate: &completedTime,
	})
}

// CompleteTask
func (t *Task) CancelTask(taskService TaskService) {
	cancelledTime := time.Now()

	t.TrackChange(taskService, TaskCancelled{
		UID:           t.UID,
		Status:        TaskCancelledCode,
		CancelledDate: &cancelledTime,
	})
}

// Event Tracking

func (state *Task) TrackChange(taskService TaskService, event interface{}) error {
	state.UncommittedChanges = append(state.UncommittedChanges, event)
	err := state.Transition(taskService, event)
	if err != nil {
		return err
	}

	return nil
}

func (state *Task) Transition(taskService TaskService, event interface{}) error {
	switch e := event.(type) {
	case TaskCreated:
		state.Title = e.Title
		state.UID = e.UID
		state.Description = e.Description
		state.CreatedDate = e.CreatedDate
		state.DueDate = e.DueDate
		state.Priority = e.Priority
		state.Status = e.Status
		state.Domain = e.Domain
		state.DomainDetails = e.DomainDetails
		state.Category = e.Category
		state.IsDue = e.IsDue
		state.AssetID = e.AssetID
	case TaskTitleChanged:
		state.Title = e.Title
	case TaskDescriptionChanged:
		state.Description = e.Description
	case TaskDueDateChanged:
		state.DueDate = e.DueDate
	case TaskPriorityChanged:
		state.Priority = e.Priority
	case TaskCategoryChanged:
		state.Category = e.Category
	case TaskDetailsChanged:
		state.DomainDetails = e.DomainDetails
	case TaskCancelled:
		state.CancelledDate = e.CancelledDate
		state.Status = TaskStatusCancelled
	case TaskCompleted:
		state.CompletedDate = e.CompletedDate
		state.Status = TaskStatusCompleted
	case TaskDue:
		state.IsDue = true
	}

	return nil
}

// Validation

// validateTaskTitle
func validateTaskTitle(title string) error {
	if title == "" {
		return TaskError{TaskErrorTitleEmptyCode}
	}
	return nil
}

// validateTaskDescription
func validateTaskDescription(description string) error {
	if description == "" {
		return TaskError{TaskErrorDescriptionEmptyCode}
	}
	return nil
}

// validateTaskDueDate
func validateTaskDueDate(newdate *time.Time) error {
	if newdate != nil {
		if (*newdate).Before(time.Now()) {
			return TaskError{TaskErrorDueDateInvalidCode}
		}
	}
	return nil
}

//validateTaskPriority
func validateTaskPriority(priority string) error {

	if priority == "" {
		return TaskError{TaskErrorPriorityEmptyCode}
	}

	_, err := FindTaskPriorityByCode(priority)
	if err != nil {
		return err
	}

	return nil
}

// validateTaskStatus
func validateTaskStatus(status string) error {

	if status == "" {
		return TaskError{TaskErrorStatusEmptyCode}
	}

	_, err := FindTaskStatusByCode(status)
	if err != nil {
		return err
	}

	return nil
}

// validateTaskCategory
func validateTaskCategory(taskcategory string) error {

	if taskcategory == "" {
		return TaskError{TaskErrorCategoryEmptyCode}
	}

	_, err := FindTaskCategoryByCode(taskcategory)
	if err != nil {
		return err
	}

	return nil
}

// validateAssetID
func validateAssetID(taskService TaskService, assetid *uuid.UUID, taskdomain string) error {

	if assetid != nil {
		if taskdomain == "" {
			return TaskError{TaskErrorDomainEmptyCode}
		}
		//Find asset in repository
		// if not found return error

		switch taskdomain {
		case TaskDomainAreaCode:
			serviceResult := taskService.FindAreaByID(*assetid)

			if serviceResult.Error != nil {
				return serviceResult.Error
			}
		case TaskDomainCropCode:

			serviceResult := taskService.FindCropByID(*assetid)

			if serviceResult.Error != nil {
				return serviceResult.Error
			}
		case TaskDomainInventoryCode:

			serviceResult := taskService.FindMaterialByID(*assetid)

			if serviceResult.Error != nil {
				return serviceResult.Error
			}
		case TaskDomainReservoirCode:
			serviceResult := taskService.FindReservoirByID(*assetid)

			if serviceResult.Error != nil {
				return serviceResult.Error
			}
		default:
			return TaskError{TaskErrorInvalidDomainCode}
		}
	}
	return nil
}
