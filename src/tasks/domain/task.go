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
func CreateTask(taskservice TaskService, title string, description string, duedate *time.Time, priority string, taskdomain TaskDomain, taskcategory string, assetid *uuid.UUID) (*Task, error) {
	// add validation

	err := validateTaskTitle(title)
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

	err = validateAssetID(taskservice, assetid, taskdomain.Code())
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
		Category:      taskcategory,
		IsDue:         false,
		AssetID:       assetid,
	})

	return initial, nil
}

// ChangeTitle
func (t *Task) ChangeTaskTitle(newtitle string) error {

	err := validateTaskTitle(newtitle)
	if err != nil {
		return err
	}

	t.Title = newtitle

	return nil
}

// ChangeDescription
func (t *Task) ChangeTaskDescription(newdescription string) error {

	err := validateTaskDescription(newdescription)
	if err != nil {
		return err
	}

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

// ChangeTaskAssetID
func (t *Task) ChangeTaskAssetID(taskService TaskService, newasset *uuid.UUID) error {

	err := validateAssetID(taskService, newasset, t.Domain)
	if err != nil {
		return err
	}
	t.AssetID = newasset
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

// Chnage TaskDomainDetails
func (t *Task) ChangeTaskDomainDetails(newdetails TaskDomain) {
	t.DomainDetails = newdetails
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

// Event Tracking

func (state *Task) TrackChange(event interface{}) {
	state.UncommittedChanges = append(state.UncommittedChanges, event)
	state.Transition(event)
}

func (state *Task) Transition(event interface{}) {
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
	case TaskModified:
		state.Title = e.Title
		state.Description = e.Description
		state.DueDate = e.DueDate
		state.Priority = e.Priority
		state.DomainDetails = e.DomainDetails
		state.Category = e.Category
		state.AssetID = e.AssetID
	case TaskCancelled:
		state.CancelledDate = e.CancelledDate
		state.Status = TaskStatusCancelled
	case TaskCompleted:
		state.CompletedDate = e.CompletedDate
		state.Status = TaskStatusCompleted
	}
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
