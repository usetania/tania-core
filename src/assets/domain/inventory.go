package domain

// type InventoryMaterial struct {
// 	UID       uuid.UUID `json:"uid"`
// 	PlantType PlantType `json:"plant_type"`
// 	Variety   string    `json:"variety"`
// }

// type PlantType interface {
// 	Code() string
// }

// type Vegetable struct{}

// func (v Vegetable) Code() string { return "vegetable" }

// type Fruit struct{}

// func (v Fruit) Code() string { return "fruit" }

// type Herb struct{}

// func (v Herb) Code() string { return "herb" }

// type Flower struct{}

// func (v Flower) Code() string { return "flower" }

// type Tree struct{}

// func (v Tree) Code() string { return "tree" }

// func CreateInventoryMaterial(plantType PlantType, variety string) (InventoryMaterial, error) {
// 	// err := validatePlantType(plantType)
// 	// if err != nil {
// 	// 	return InventoryMaterial{}, err
// 	// }

// 	if variety == "" {
// 		return InventoryMaterial{}, InventoryMaterialError{Code: InventoryMaterialInvalidVariety}
// 	}

// 	uid, err := uuid.NewV4()
// 	if err != nil {
// 		return InventoryMaterial{}, err
// 	}

// 	return InventoryMaterial{
// 		UID:       uid,
// 		PlantType: plantType,
// 		Variety:   variety,
// 	}, nil
// }

func CreateInventoryTools() {}

// func GetPlantTypes() []PlantType {
// 	return []PlantType{
// 		Vegetable{},
// 		Fruit{},
// 		Herb{},
// 		Flower{},
// 		Tree{},
// 	}
// }

// func validatePlantType(plantType PlantType) error {
// 	switch plantType.(type) {
// 	case Vegetable:
// 	case Fruit:
// 	case Herb:
// 	case Flower:
// 	case Tree:
// 	default:
// 		return InventoryMaterialError{Code: InventoryMaterialInvalidPlantType}
// 	}

// 	return nil
// }
