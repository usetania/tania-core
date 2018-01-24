package domain

const (
	TaskStatusInProgress = "inprogress"
	TaskStatusCancelled  = "cancelled"
	TaskStatusComplete   = "completed"
)

type TaskStatus struct {
	Code string
	Name string
}

func FindAllTaskStatus() []TaskStatus {
	return []TaskStatus{
		TaskStatus{Code: TaskStatusInProgress, Name: "In Progress"},
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
