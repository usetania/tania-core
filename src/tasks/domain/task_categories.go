package domain

const (
	TaskCategoryArea      = "area"
	TaskCategoryReservoir = "reservoir"
	TaskCategoryCrop      = "crop"
	TaskCategoryInventory = "inventory"
	TaskCategoryDevice    = "device"
	TaskCategoryFinance   = "finance"
)

type TaskCategory struct {
	Code string
	Name string
}

func FindAllTaskCategories() []TaskCategory {
	return []TaskCategory{
		TaskCategory{Code: TaskCategoryArea, Name: "Area"},
		TaskCategory{Code: TaskCategoryReservoir, Name: "Reservoir"},
		TaskCategory{Code: TaskCategoryCrop, Name: "Crop"},
		TaskCategory{Code: TaskCategoryInventory, Name: "Inventory"},
		TaskCategory{Code: TaskCategoryDevice, Name: "Device"},
		TaskCategory{Code: TaskCategoryFinance, Name: "Finance"},
	}
}

func FindTaskCategoryByCode(code string) (TaskCategory, error) {
	items := FindAllTaskCategories()

	for _, item := range items {
		if item.Code == code {
			return item, nil
		}
	}

	return TaskCategory{}, TaskError{TaskErrorInvalidCategoryCode}
}
