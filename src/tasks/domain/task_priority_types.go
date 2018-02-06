package domain

const (
	TaskPriorityUrgent = "URGENT"
	TaskPriorityNormal = "NORMAL"
)

type TaskPriority struct {
	Code string
	Name string
}

func FindAllTaskPriority() []TaskPriority {
	return []TaskPriority{
		TaskPriority{Code: TaskPriorityUrgent, Name: "Urgent"},
		TaskPriority{Code: TaskPriorityNormal, Name: "Normal"},
	}
}

func FindTaskPriorityByCode(code string) (TaskPriority, error) {
	items := FindAllTaskPriority()

	for _, item := range items {
		if item.Code == code {
			return item, nil
		}
	}

	return TaskPriority{}, TaskError{TaskErrorInvalidPriorityCode}
}
