package domain

const (
	TaskStatusCreated = "created"
	TaskStatusInProgress = "inprogress"
	TaskStatusCancelled = "cancelled"
	TaskStatusComplete = "complete"
)

type TaskStatus struct {
	Code string
	Name string
}

func FindAllTaskStatus() []TaskStatus {
	return []TaskStatus{
		TaskStatus{Code: TaskStatusCreated, Name: "Created"},
		TaskStatus{Code: TaskStatusInProgress, Name: "In Progress"},
		TaskStatus{Code: TaskStatusCancelled, Name: "cancelled"},
		TaskStatus{Code: TaskStatusComplete, Name: "Complete"},
	}
}

func FindTaskStatusByCode(code string) (TaskStatus, error) {
	items := FindAllTaskStatus()

	for _, item := range items {
		if item.Code == code {
			return item, nil
		}
	}

	return TaskStatus{}, TaskError{TaskErrorInvalidStatusCode}
}
