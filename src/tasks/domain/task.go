package domain

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type TaskService interface {
	FindAreaByID(uid uuid.UUID) ServiceResult
	FindCropByID(uid uuid.UUID) ServiceResult
}

// ServiceResult is the container for service result
type ServiceResult struct {
	Result interface{}
	Error  error
}

type Task struct {
	UID          uuid.UUID  `json:"uid"`
	Description  string     `json:"description"`
	CreatedDate  time.Time  `json:"created_date"`
	DueDate      *time.Time `json:"due_date,omitempty"`
	Priority     string     `json:"priority"`
	Status       string     `json:"status"`
	TaskType     string     `json:"type"`
	IsDue        bool       `json:"is_due"`
	AssetID      uuid.UUID  `json:"asset_id"`
	TaskActivity Activity   `json:"Activity"`
}

// CreateTask
func CreateTask(task_service TaskService, description string, due_date *time.Time, priority string, tasktype string, asset_id string, activity Activity) (Task, error) {
	// add validation

	err := validateTaskDueDate(due_date)
	if err != nil {
		return Task{}, err
	}

	err = validateTaskPriority(priority)
	if err != nil {
		return Task{}, err
	}

	err = validateTaskType(tasktype)
	if err != nil {
		return Task{}, err
	}
	if asset_id == "" {
		return Task{}, TaskError{TaskErrorAssetIDEmptyCode}
	}
	asset, err := uuid.FromString(asset_id)
	if err != nil {
		return Task{}, TaskError{TaskErrorInvalidAssetIDCode}
	}

	err = validateAssetID(task_service, asset, tasktype)
	if err != nil {
		return Task{}, err
	}

	uid, err := uuid.NewV4()
	if err != nil {
		return Task{}, err
	}

	return Task{
		UID:          uid,
		Description:  description,
		CreatedDate:  time.Now(),
		DueDate:      due_date,
		Priority:     priority,
		Status:       TaskStatusCreated,
		TaskType:     tasktype,
		IsDue:        false,
		AssetID:      asset,
		TaskActivity: activity,
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
func (t *Task) ChangeTaskType(newtasktype string) error {

	err := validateTaskType(newtasktype)
	if err != nil {
		return err
	}
	t.TaskType = newtasktype

	return nil
}

//SetTaskAsDue
func (t *Task) SetTaskAsDue() {
	t.IsDue = true
}

// Validation

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

// validateTaskType
func validateTaskType(tasktype string) error {

	if tasktype == "" {
		return TaskError{TaskErrorTypeEmptyCode}
	}

	_, err := FindTaskTypeByCode(tasktype)
	if err != nil {
		return err
	}

	return nil
}

// validateAssetID
func validateAssetID(taskService TaskService, asset_id uuid.UUID, tasktype string) error {

	if tasktype == "" {
		return TaskError{TaskErrorTypeEmptyCode}
	}
	//Find asset in repository
	// if not found return error

	switch tasktype {
	case TaskTypeArea:
		serviceResult := taskService.FindAreaByID(asset_id)

		if serviceResult.Error != nil {
			return serviceResult.Error
		}
	case TaskTypeCrop:

		serviceResult := taskService.FindCropByID(asset_id)

		if serviceResult.Error != nil {
			return serviceResult.Error
		}
	default:
	}

	return nil
}
