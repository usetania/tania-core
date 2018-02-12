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
	DueDate       *time.Time `json:"due_date,omitempty"`
	CompletedDate *time.Time `json:"completed_date"`
	Priority      string     `json:"priority"`
	Status        string     `json:"status"`
	Domain        string     `json:"domain"`
	DomainDetails TaskDomain `json:"domain_details"`
	Category      string     `json:"category"`
	IsDue         bool       `json:"is_due"`
	AssetID       *uuid.UUID `json:"asset_id"`
}

// CreateTask
func CreateTask(taskservice TaskService, title string, description string, duedate *time.Time, priority string, taskdomain TaskDomain, taskcategory string, assetid *uuid.UUID) (Task, error) {
	// add validation

	err := validateTaskTitle(title)
	if err != nil {
		return Task{}, err
	}

	err = validateTaskDueDate(duedate)
	if err != nil {
		return Task{}, err
	}

	err = validateTaskPriority(priority)
	if err != nil {
		return Task{}, err
	}

	err = validateTaskCategory(taskcategory)
	if err != nil {
		return Task{}, err
	}

	err = validateAssetID(taskservice, assetid, taskdomain.Code())
	if err != nil {
		return Task{}, err
	}

	uid, err := uuid.NewV4()
	if err != nil {
		return Task{}, err
	}

	return Task{
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
	}, nil
}

// ChangeDescription
func (t *Task) ChangeTaskDescription(newdescription string) error {

	t.Description = newdescription

	return nil
}

// ChangeDueDate
func (t *Task) ChangeTaskDueDate(newdate *time.Time) error {

	err := validateTaskDueDate(newdate)
	if err != nil {
		return err
	}
	t.DueDate = newdate

	return nil
}

// ChangePriority
func (t *Task) ChangeTaskPriority(newpriority string) error {

	err := validateTaskPriority(newpriority)
	if err != nil {
		return err
	}
	t.Priority = newpriority

	return nil
}

// ChangeStatus
func (t *Task) ChangeTaskStatus(newstatus string) error {

	err := validateTaskStatus(newstatus)
	if err != nil {
		return err
	}
	t.Status = newstatus

	return nil

}

// ChangeCategory
func (t *Task) ChangeTaskCategory(newtaskcategory string) error {

	err := validateTaskCategory(newtaskcategory)
	if err != nil {
		return err
	}
	t.Category = newtaskcategory

	return nil
}

// SetTaskAsDue
func (t *Task) SetTaskAsDue() {
	t.IsDue = true
}

// SetTaskCompletedDate()
func (t *Task) SetTaskCompletedDate() {
	currenttime := time.Now()
	t.CompletedDate = &currenttime
}

// Validation

// validateTaskTitle
func validateTaskTitle(title string) error {
	if title == "" {
		return TaskError{TaskErrorTitleEmptyCode}
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
