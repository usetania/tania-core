package domain

const (
	TaskCategoryArea        = "AREA"
	TaskCategoryCrop        = "CROP"
	TaskCategoryFinance     = "FINANCE"
	TaskCategoryGeneral     = "GENERAL"
	TaskCategoryInventory   = "INVENTORY"
	TaskCategoryNutrient    = "NUTRIENT"
	TaskCategoryPestControl = "PESTCONTROL"
	TaskCategoryReservoir   = "RESERVOIR"
	TaskCategorySafety      = "SAFETY"
	TaskCategorySanitation  = "SANITATION"
)

type TaskCategory struct {
	Code string
	Name string
}

func FindAllTaskCategories() []TaskCategory {
	return []TaskCategory{
		TaskCategory{Code: TaskCategoryArea, Name: "Area"},
		TaskCategory{Code: TaskCategoryCrop, Name: "Crop"},
		TaskCategory{Code: TaskCategoryFinance, Name: "Finance"},
		TaskCategory{Code: TaskCategoryGeneral, Name: "General"},
		TaskCategory{Code: TaskCategoryInventory, Name: "Inventory"},
		TaskCategory{Code: TaskCategoryNutrient, Name: "Nutrient"},
		TaskCategory{Code: TaskCategoryPestControl, Name: "Pest Control"},
		TaskCategory{Code: TaskCategoryReservoir, Name: "Reservoir"},
		TaskCategory{Code: TaskCategorySafety, Name: "Safety"},
		TaskCategory{Code: TaskCategorySanitation, Name: "Sanitation"},
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
