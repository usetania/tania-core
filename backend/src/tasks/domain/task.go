package domain

import (
	"time"

	"github.com/gofrs/uuid"
)

type TaskService interface {
	FindAreaByID(uid uuid.UUID) ServiceResult
	FindCropByID(uid uuid.UUID) ServiceResult
	FindMaterialByID(uid uuid.UUID) ServiceResult
	FindReservoirByID(uid uuid.UUID) ServiceResult
}

// ServiceResult is the container for service result.
type ServiceResult struct {
	Result interface{}
	Error  error
}

type Task struct {
	UID           uuid.UUID  `json:"uid"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	CreatedDate   time.Time  `json:"created_date"`
	DueDate       *time.Time `json:"due_date,omitempty"`
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

// CreateTask.
func CreateTask(
	ts TaskService,
	title, description, priority, category string,
	duedate *time.Time,
	taskdomain TaskDomain,
	assetid *uuid.UUID,
) (*Task, error) {
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

	err = validateTaskCategory(category)
	if err != nil {
		return &Task{}, err
	}

	err = validateAssetID(ts, assetid, taskdomain.Code())
	if err != nil {
		return &Task{}, err
	}

	uid, err := uuid.NewV4()
	if err != nil {
		return &Task{}, err
	}

	initial := &Task{}

	initial.TrackChange(TaskCreated{
		Title:         title,
		UID:           uid,
		Description:   description,
		CreatedDate:   time.Now(),
		DueDate:       duedate,
		Priority:      priority,
		Status:        TaskStatusCreated,
		Domain:        taskdomain.Code(),
		DomainDetails: taskdomain,
		Category:      category,
		IsDue:         false,
		AssetID:       assetid,
	})

	return initial, nil
}

func (t *Task) ChangeTaskTitle(title string) error {
	if err := validateTaskTitle(title); err != nil {
		return err
	}

	event := TaskTitleChanged{
		UID:   t.UID,
		Title: title,
	}

	t.TrackChange(event)

	return nil
}

func (t *Task) ChangeTaskDescription(description string) (*Task, error) {
	err := validateTaskDescription(description)
	if err != nil {
		return &Task{}, err
	}

	event := TaskDescriptionChanged{
		UID:         t.UID,
		Description: description,
	}

	t.TrackChange(event)

	return t, nil
}

func (t *Task) ChangeTaskDueDate(duedate *time.Time) (*Task, error) {
	err := validateTaskDueDate(duedate)
	if err != nil {
		return &Task{}, err
	}

	event := TaskDueDateChanged{
		UID:     t.UID,
		DueDate: duedate,
	}

	t.TrackChange(event)

	return t, nil
}

func (t *Task) ChangeTaskPriority(priority string) (*Task, error) {
	err := validateTaskPriority(priority)
	if err != nil {
		return &Task{}, err
	}

	event := TaskPriorityChanged{
		UID:      t.UID,
		Priority: priority,
	}

	t.TrackChange(event)

	return t, nil
}

func (t *Task) ChangeTaskCategory(category string) (*Task, error) {
	err := validateTaskCategory(category)
	if err != nil {
		return &Task{}, err
	}

	event := TaskCategoryChanged{
		UID:      t.UID,
		Category: category,
	}

	t.TrackChange(event)

	return t, nil
}

func (t *Task) ChangeTaskDetails(details TaskDomain) (*Task, error) {
	event := TaskDetailsChanged{
		UID:           t.UID,
		DomainDetails: details,
	}

	t.TrackChange(event)

	return t, nil
}

// SetTaskAsDue.
func (t *Task) SetTaskAsDue() {
	t.TrackChange(TaskDue{
		UID: t.UID,
	})
}

// CompleteTask.
func (t *Task) CompleteTask() {
	completedTime := time.Now()

	t.TrackChange(TaskCompleted{
		UID:           t.UID,
		Status:        TaskCompletedCode,
		CompletedDate: &completedTime,
	})
}

// CompleteTask.
func (t *Task) CancelTask() {
	cancelledTime := time.Now()

	t.TrackChange(TaskCancelled{
		UID:           t.UID,
		Status:        TaskCancelledCode,
		CancelledDate: &cancelledTime,
	})
}

// Event Tracking.
func (t *Task) TrackChange(event interface{}) {
	t.UncommittedChanges = append(t.UncommittedChanges, event)
	t.Transition(event)
}

func (t *Task) Transition(event interface{}) {
	switch e := event.(type) {
	case TaskCreated:
		t.Title = e.Title
		t.UID = e.UID
		t.Description = e.Description
		t.CreatedDate = e.CreatedDate
		t.DueDate = e.DueDate
		t.Priority = e.Priority
		t.Status = e.Status
		t.Domain = e.Domain
		t.DomainDetails = e.DomainDetails
		t.Category = e.Category
		t.IsDue = e.IsDue
		t.AssetID = e.AssetID
	case TaskTitleChanged:
		t.Title = e.Title
	case TaskDescriptionChanged:
		t.Description = e.Description
	case TaskDueDateChanged:
		t.DueDate = e.DueDate
	case TaskPriorityChanged:
		t.Priority = e.Priority
	case TaskCategoryChanged:
		t.Category = e.Category
	case TaskDetailsChanged:
		t.DomainDetails = e.DomainDetails
	case TaskCancelled:
		t.CancelledDate = e.CancelledDate
		t.Status = TaskStatusCancelled
	case TaskCompleted:
		t.CompletedDate = e.CompletedDate
		t.Status = TaskStatusCompleted
	case TaskDue:
		t.IsDue = true
	}
}

// Validation

// validateTaskTitle.
func validateTaskTitle(title string) error {
	if title == "" {
		return TaskError{TaskErrorTitleEmptyCode}
	}

	return nil
}

// validateTaskDescription.
func validateTaskDescription(description string) error {
	if description == "" {
		return TaskError{TaskErrorDescriptionEmptyCode}
	}

	return nil
}

// validateTaskDueDate.
func validateTaskDueDate(newdate *time.Time) error {
	if newdate != nil {
		if newdate.Before(time.Now()) {
			return TaskError{TaskErrorDueDateInvalidCode}
		}
	}

	return nil
}

// validateTaskPriority.
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

// validateTaskCategory.
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

// validateAssetID.
func validateAssetID(taskService TaskService, assetid *uuid.UUID, taskdomain string) error {
	if assetid != nil {
		if taskdomain == "" {
			return TaskError{TaskErrorDomainEmptyCode}
		}
		// Find asset in repository
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
