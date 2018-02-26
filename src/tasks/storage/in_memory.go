package storage

import (
	"fmt"
	domain "github.com/Tanibox/tania-server/src/tasks/domain"
	deadlock "github.com/sasha-s/go-deadlock"
	uuid "github.com/satori/go.uuid"
	"time"
)

type TaskEventStorage struct {
	Lock       *deadlock.RWMutex
	TaskEvents []TaskEvent
}

type TaskEvent struct {
	TaskUID uuid.UUID
	Version int
	Event   interface{}
}

func CreateTaskEventStorage() *TaskEventStorage {
	rwMutex := deadlock.RWMutex{}
	deadlock.Opts.DeadlockTimeout = time.Second * 10
	deadlock.Opts.OnPotentialDeadlock = func() {
		fmt.Println("TASK EVENT STORAGE DEADLOCK!")
	}

	return &TaskEventStorage{Lock: &rwMutex}
}

type TaskReadStorage struct {
	Lock        *deadlock.RWMutex
	TaskReadMap map[uuid.UUID]TaskRead
}

func CreateTaskReadStorage() *TaskReadStorage {
	rwMutex := deadlock.RWMutex{}
	deadlock.Opts.DeadlockTimeout = time.Second * 10
	deadlock.Opts.OnPotentialDeadlock = func() {
		fmt.Println("TASK READ STORAGE DEADLOCK!")
	}

	return &TaskReadStorage{TaskReadMap: make(map[uuid.UUID]TaskRead), Lock: &rwMutex}
}

type TaskRead struct {
	Title         string            `json:"title"`
	UID           uuid.UUID         `json:"uid"`
	Description   string            `json:"description"`
	CreatedDate   time.Time         `json:"created_date"`
	DueDate       *time.Time        `json:"due_date, omitempty"`
	CompletedDate *time.Time        `json:"completed_date"`
	CancelledDate *time.Time        `json:"cancelled_date"`
	Priority      string            `json:"priority"`
	Status        string            `json:"status"`
	Domain        string            `json:"domain"`
	DomainDetails domain.TaskDomain `json:"domain_details"`
	Category      string            `json:"category"`
	IsDue         bool              `json:"is_due"`
	AssetID       *uuid.UUID        `json:"asset_id"`
}

func (taskRead *TaskRead) BuildTaskFromTaskRead(t domain.Task) (*domain.Task, error) {

	t.Title = taskRead.Title
	t.UID = taskRead.UID
	t.Description = taskRead.Description
	t.CreatedDate = taskRead.CreatedDate
	t.DueDate = taskRead.DueDate
	t.CompletedDate = taskRead.CompletedDate
	t.Priority = taskRead.Priority
	t.Status = taskRead.Status
	t.Domain = taskRead.Domain
	t.DomainDetails = taskRead.DomainDetails
	t.Category = taskRead.Category
	t.IsDue = taskRead.IsDue
	t.AssetID = taskRead.AssetID

	return &t, nil
}

// CreateTask
func CreateTaskModifiedEvent(uid uuid.UUID, title string, description string, duedate *time.Time, priority string, taskdomain domain.TaskDomain, taskcategory string, assetid *uuid.UUID) (*domain.TaskModified, error) {

	event := domain.TaskModified{
		UID:           uid,
		Title:         title,
		Description:   description,
		Priority:      priority,
		DueDate:       duedate,
		Domain:        taskdomain.Code(),
		DomainDetails: taskdomain,
		Category:      taskcategory,
		AssetID:       assetid,
	}

	return &event, nil
}

func CreateTaskCancelledEvent(uid uuid.UUID, title string, description string, duedate *time.Time, priority string, taskdomain domain.TaskDomain, taskcategory string, assetid *uuid.UUID) *domain.TaskCancelled {

	cancelTime := time.Now()

	event := domain.TaskCancelled{
		UID:           uid,
		Title:         title,
		Description:   description,
		Priority:      priority,
		DueDate:       duedate,
		Domain:        taskdomain.Code(),
		DomainDetails: taskdomain,
		Category:      taskcategory,
		AssetID:       assetid,
		CancelledDate: &cancelTime,
	}

	return &event
}

func CreateTaskCompletedEvent(uid uuid.UUID, title string, description string, duedate *time.Time, priority string, taskdomain domain.TaskDomain, taskcategory string, assetid *uuid.UUID) *domain.TaskCompleted {

	completedTime := time.Now()

	event := domain.TaskCompleted{
		UID:           uid,
		Title:         title,
		Description:   description,
		Priority:      priority,
		DueDate:       duedate,
		Domain:        taskdomain.Code(),
		DomainDetails: taskdomain,
		Category:      taskcategory,
		AssetID:       assetid,
		CompletedDate: &completedTime,
	}

	return &event
}

func CreateTaskDueEvent(uid uuid.UUID, title string, description string, duedate *time.Time, priority string, taskdomain domain.TaskDomain, taskcategory string, assetid *uuid.UUID) *domain.TaskDue {

	event := domain.TaskDue{
		UID:           uid,
		Title:         title,
		Description:   description,
		Priority:      priority,
		DueDate:       duedate,
		Domain:        taskdomain.Code(),
		DomainDetails: taskdomain,
		Category:      taskcategory,
		AssetID:       assetid,
	}

	return &event
}
