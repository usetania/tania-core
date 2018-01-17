package domain

const (
	TaskTypeArea = "area"
	TaskTypeReservoir = "reservoir"
	TaskTypeCrop = "crop"
	TaskTypeInventory = "inventory"
	TaskTypeDevice = "device"
	TaskTypeFinance = "finance"
)

type TaskType struct {
	Code string
	Name string
}

func FindAllTaskTypes() []TaskType {
	return []TaskType{
		TaskType{Code: TaskTypeArea, Name: "Area"},
		TaskType{Code: TaskTypeReservoir, Name: "Reservoir"},
		TaskType{Code: TaskTypeCrop, Name: "Crop"},
		TaskType{Code: TaskTypeInventory, Name: "Inventory"},
		TaskType{Code: TaskTypeDevice, Name: "Device"},
		TaskType{Code: TaskTypeFinance, Name: "Finance"},
	}
}

func FindTaskTypeByCode(code string) (TaskType, error) {
	items := FindAllTaskTypes()

	for _, item := range items {
		if item.Code == code {
			return item, nil
		}
	}

	return TaskType{}, TaskError{TaskErrorInvalidTypeCode}
}
