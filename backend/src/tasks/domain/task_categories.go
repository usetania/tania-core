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
		{Code: TaskCategoryArea, Name: "Area"},
		{Code: TaskCategoryCrop, Name: "Crop"},
		{Code: TaskCategoryFinance, Name: "Finance"},
		{Code: TaskCategoryGeneral, Name: "General"},
		{Code: TaskCategoryInventory, Name: "Inventory"},
		{Code: TaskCategoryNutrient, Name: "Nutrient"},
		{Code: TaskCategoryPestControl, Name: "Pest Control"},
		{Code: TaskCategoryReservoir, Name: "Reservoir"},
		{Code: TaskCategorySafety, Name: "Safety"},
		{Code: TaskCategorySanitation, Name: "Sanitation"},
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
