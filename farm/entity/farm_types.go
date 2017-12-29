package entity

const (
	FarmTypeOrganic      = "organic"
	FarmTypeHydroponic   = "hydroponic"
	FarmTypeAquaponic    = "aquaponic"
	FarmTypeMushroom     = "mushroom"
	FarmTypeLiveStock    = "livestock"
	FarmTypeFisheries    = "fisheries"
	FarmTypePermaculture = "permaculture"
)

type FarmType struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

func FindAllFarmTypes() []FarmType {
	return []FarmType{
		FarmType{Code: FarmTypeOrganic, Name: "Organic / Soil-Based"},
		FarmType{Code: FarmTypeHydroponic, Name: "Hydroponic"},
		FarmType{Code: FarmTypeAquaponic, Name: "Aquaponic"},
		FarmType{Code: FarmTypeMushroom, Name: "Mushroom"},
		FarmType{Code: FarmTypeLiveStock, Name: "Livestock"},
		FarmType{Code: FarmTypeFisheries, Name: "Fisheries"},
		FarmType{Code: FarmTypePermaculture, Name: "Permaculture"},
	}
}

func FindFarmTypeByCode(code string) (FarmType, error) {
	items := FindAllFarmTypes()

	for _, item := range items {
		if item.Code == code {
			return item, nil
		}
	}

	return FarmType{}, FarmError{FarmErrorInvalidFarmTypeCode}
}
