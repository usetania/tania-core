package domain

type InventoryMaterial struct {
	PlantType PlantType
	Variety   string
}

type PlantType interface {
	Code() string
}

type Vegetable struct{}

func (v Vegetable) Code() string { return "vegetable" }

type Fruit struct{}

func (v Fruit) Code() string { return "fruit" }

type Herb struct{}

func (v Herb) Code() string { return "herb" }

type Flower struct{}

func (v Flower) Code() string { return "flower" }

type Tree struct{}

func (v Tree) Code() string { return "tree" }

func CreateInventoryMaterial(plantType PlantType) (InventoryMaterial, error) {
	err := validatePlantType(plantType)
	if err != nil {
		return InventoryMaterial{}, err
	}

	return InventoryMaterial{PlantType: plantType}, nil
}

func CreateInventoryTools() {}

func validatePlantType(plantType PlantType) error {
	switch plantType.(type) {
	case Vegetable:
	case Fruit:
	case Herb:
	case Flower:
	case Tree:
	default:
		return InventoryMaterialError{Code: InventoryMaterialInvalidPlantType}
	}

	return nil
}
