package domain

const (
	TaskStatusCreated   = "created"
	TaskStatusCancelled = "cancelled"
	TaskStatusComplete  = "completed"
)

type TaskStatus struct {
	Code string
	Name string
}

func FindAllTaskStatus() []TaskStatus {
	return []TaskStatus{
		TaskStatus{Code: TaskStatusCreated, Name: "Created"},
		TaskStatus{Code: TaskStatusCancelled, Name: "Cancelled"},
		TaskStatus{Code: TaskStatusComplete, Name: "Completed"},
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
