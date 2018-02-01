package domain

import "errors"

const (
	MaterialTypeSeedCode                = "SEED"
	MaterialTypeGrowingMediumCode       = "GROWING_MEDIUM"
	MaterialTypeAgrochemicalCode        = "AGROCHEMICAL"
	MaterialTypeLabelAndCropSupportCode = "LABEL_AND_CROP_SUPPORT"
	MaterialTypeSeedingContainerCode    = "SEEDING_CONTAINER"
	MaterialTypePostHarvestSupplyCode   = "POST_HARVEST_SUPPLY"
	MaterialTypeOtherCode               = "OTHER"
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
	Code  string `json:"code"`
	Label string `json:"label"`
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

type MaterialTypeAgrochemical struct {
	ChemicalType ChemicalType
}

func (mt MaterialTypeAgrochemical) Code() string {
	return MaterialTypeAgrochemicalCode
}

const (
	ChemicalTypeDisinfectant = "DISINFECTANT"
	ChemicalTypeFertilizer   = "FERTILIZER"
	ChemicalTypeHormone      = "HORMONE"
	ChemicalTypeManure       = "MANURE"
	ChemicalTypePesticide    = "PESTICIDE"
)

type ChemicalType struct {
	Code  string
	Label string
}

func ChemicalTypes() []ChemicalType {
	return []ChemicalType{
		{Code: ChemicalTypeDisinfectant, Label: "Disinfectant and Sanitizer"},
		{Code: ChemicalTypeFertilizer, Label: "Fertilizer"},
		{Code: ChemicalTypeHormone, Label: "Hormone and Growth Agent"},
		{Code: ChemicalTypeManure, Label: "Manure"},
		{Code: ChemicalTypePesticide, Label: "Pesticide"},
	}
}

func GetChemicalType(code string) ChemicalType {
	for _, v := range ChemicalTypes() {
		if v.Code == code {
			return v
		}
	}

	return ChemicalType{}
}

func CreateMaterialTypeAgrochemical(chemicalType string) (MaterialTypeAgrochemical, error) {
	ct := GetChemicalType(chemicalType)
	if ct == (ChemicalType{}) {
		return MaterialTypeAgrochemical{}, errors.New("options wrong")
	}

	return MaterialTypeAgrochemical{ct}, nil
}

type MaterialTypeGrowingMedium struct {
}

func (mt MaterialTypeGrowingMedium) Code() string {
	return MaterialTypeGrowingMediumCode
}

type MaterialTypeLabelAndCropSupport struct {
}

func (mt MaterialTypeLabelAndCropSupport) Code() string {
	return MaterialTypeLabelAndCropSupportCode
}

type MaterialTypeSeedingContainer struct {
	ContainerType ContainerType
}

func (mt MaterialTypeSeedingContainer) Code() string {
	return MaterialTypeSeedingContainerCode
}

const (
	ContainerTypeTray = "TRAY"
	ContainerTypePot  = "POT"
)

type ContainerType struct {
	Code  string `json:"code"`
	Label string `json:"label"`
}

func ContainerTypes() []ContainerType {
	return []ContainerType{
		{Code: ContainerTypeTray, Label: "Tray"},
		{Code: ContainerTypePot, Label: "Pot"},
	}
}

func GetContainerType(code string) ContainerType {
	for _, v := range ContainerTypes() {
		if v.Code == code {
			return v
		}
	}

	return ContainerType{}
}

func CreateMaterialTypeSeedingContainer(containerType string) (MaterialTypeSeedingContainer, error) {
	ct := GetContainerType(containerType)
	if ct == (ContainerType{}) {
		return MaterialTypeSeedingContainer{}, errors.New("options wrong")
	}

	return MaterialTypeSeedingContainer{ContainerType: ct}, nil
}

type MaterialTypePostHarvestSupply struct {
}

func (mt MaterialTypePostHarvestSupply) Code() string {
	return MaterialTypePostHarvestSupplyCode
}
