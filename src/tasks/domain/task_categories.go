package domain

const (
	TaskCategoryArea        = "area"
	TaskCategoryCrop        = "crop"
	TaskCategoryFinance     = "finance"
	TaskCategoryGeneral     = "general"
	TaskCategoryInventory   = "inventory"
	TaskCategoryNutrient    = "nutrient"
	TaskCategoryPestControl = "pestcontrol"
	TaskCategoryReservoir   = "reservoir"
	TaskCategorySafety      = "safety"
	TaskCategorySanitation  = "sanitation"
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
