package domain

import "errors"

const (
	MaterialTypeSeedCode                = "MATERIAL_SEED"
	MaterialTypeGrowingMediumCode       = "MATERIAL_GROWING_MEDIUM"
	MaterialTypeAgrochemicalCode        = "MATERIAL_AGROCHEMICAL"
	MaterialTypeLabelAndCropSupportCode = "MATERIAL_LABEL_AND_CROP_SUPPORT"
	MaterialTypeSeedingContainerCode    = "MATERIAL_SEEDING_CONTAINER"
	MaterialTypePostHarvestSupplyCode   = "MATERIAL_POST_HARVEST_SUPPLY"
	MaterialTypeOtherCode               = "MATERIAL_OTHER"
)

type MaterialType interface {
	Code() string
}

type MaterialTypeSeed struct {
	PlantType PlantType
}

func (mt MaterialTypeSeed) Code() string {
	return MaterialTypeSeedCode
}

func CreateMaterialTypeSeed(plantType string) (MaterialTypeSeed, error) {
	pt := GetPlantType(plantType)
	if pt == (PlantType{}) {
		return MaterialTypeSeed{}, errors.New("options wrong")
	}

	return MaterialTypeSeed{pt}, nil
}

type PlantType struct {
	Code  string
	Label string
}

const (
	PlantTypeVegetable = "VEGETABLE"
	PlantTypeFruit     = "FRUIT"
	PlantTypeHerb      = "HERB"
	PlantTypeFlower    = "FLOWER"
	PlantTypeTree      = "TREE"
)

func PlantTypes() []PlantType {
	return []PlantType{
		{Code: PlantTypeVegetable, Label: "Vegetable"},
		{Code: PlantTypeFruit, Label: "Fruit"},
		{Code: PlantTypeHerb, Label: "Herb"},
		{Code: PlantTypeFlower, Label: "Flower"},
		{Code: PlantTypeTree, Label: "Tree"},
	}
}

func GetPlantType(code string) PlantType {
	for _, v := range PlantTypes() {
		if v.Code == code {
			return v
		}
	}

	return PlantType{}
}
