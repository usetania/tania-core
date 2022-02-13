package domain

const (
	TaskStatusCreated   = "CREATED"
	TaskStatusCancelled = "CANCELLED"
	TaskStatusCompleted = "COMPLETED"
)

type TaskStatus struct {
	Code string
	Name string
}

func FindAllTaskStatus() []TaskStatus {
	return []TaskStatus{
		{Code: TaskStatusCreated, Name: "Created"},
		{Code: TaskStatusCancelled, Name: "Cancelled"},
		{Code: TaskStatusCompleted, Name: "Completed"},
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
