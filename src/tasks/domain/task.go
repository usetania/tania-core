package entity

import (
	"github.com/Tanibox/tania-server/src/helper/validationhelper"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Task struct {
	UID			uuid.UUID 	`json:"uid"`
	Description	string		`json:"description"`
	CreatedDate	time.Time	`json:"createddate"`
	DueDate		time.Time	`json:"duedate"`
	Priority	string		`json:"priority"`
	Status		string		`json:"status"`
	TaskType	string		`json:"type"`
	ParentUID	uuid.UUID	`json:"parentuid"`
}

// CreateTask
func CreateTask(description string, duedate time.Time, priority string, status string, tasktype string, parentuid string) (Task, error) {
	// add validation

	err := validateTaskDescription(description)
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
	
	err = validateTaskStatus(status)
	if err != nil {
		return Task{}, err
	}
	
	err = validateTaskType(tasktype)
	if err != nil {
		return Task{}, err
	}
	
	err = validateParentUID(parentuid)
	if err != nil {
		return Task{}, err
	}

	uid, err := uuid.NewV4()
	if err != nil {
		return Task{}, err
	}
	parent, err := uuid.FromString(parentuid)
	if err != nil {
		// throw error
	}

	return Task {
		UID:			uid,
		Description:	description,
		CreatedDate:	time.Now(),
		DueDate:		duedate,
		Priority:		priority,
		Status:			status,
		TaskType:		tasktype,
		ParentUID:		parent,
	}, nil
}

// ChangeDescription
func (t *Task) ChangeTaskDescription (newdescription string) error {

	err := validateTaskDescription(newdescription)
	if err != nil {
		return err
	}
	t.Description = newdescription

	return nil
}

// ChangeDueDate
func (t *Task) ChangeTaskDueDate (newdate time.Time) error {

	err := validateTaskDueDate(newdate)
	if err != nil {
		return err
	}
	t.DueDate = newdate

	return nil
}

// ChangePriority
func (t *Task) ChangeTaskPriority (newpriority string) error {

	err := validateTaskPriority(newpriority)
	if err != nil {
		return err
	}
	t.Priority = newpriority

	return nil
}

// ChangeStatus
func (t *Task) ChangeTaskStatus (newstatus string) error {

	err := validateTaskPriority(newstatus)
	if err != nil {
		return err
	}
	t.Status = newstatus

	return nil

}

// ChangeCategory
func (t *Task) ChangeTaskType (newtasktype string) error {

	err := validateTaskPriority(newtasktype)
	if err != nil {
		return err
	}
	t.TaskType = newtasktype

	return nil
}


// Validation 

// validateTaskDescription
func validateTaskDescription (description string) error {
	if description == "" {
		// return error
	}
	if !validationhelper.IsAlphanumeric(description) {
		// return error
	}
	return nil
}

// validateTaskDueDate
func validateTaskDueDate (newdate time.Time) error {
	if newdate.Before(time.Now()) {
		//return error
	}
	return nil
}

//validateTaskPriority
func validateTaskPriority (priority string) error {

	_, err := FindTaskPriorityByCode(priority)
	if err != nil {
		return err
	}

	return nil
}

// validateTaskStatus
func validateTaskStatus (status string) error {

	_, err := FindTaskStatusByCode(status)
	if err != nil {
		return err
	}

	return nil
}

// validateTaskType
func validateTaskType (tasktype string) error {

	_, err := FindTaskTypeByCode(tasktype)
	if err != nil {
		return err
	}

	return nil
}

// validateParentUID
func validateParentUID (parentuid string) error {
	if parentuid == "" {
		// return error
	}

	//Find parent in repository
	// if not found return error

	return nil
}