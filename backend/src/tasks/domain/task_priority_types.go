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
		{Code: TaskPriorityUrgent, Name: "Urgent"},
		{Code: TaskPriorityNormal, Name: "Normal"},
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
