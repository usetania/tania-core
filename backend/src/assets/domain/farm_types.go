package domain

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
		{Code: FarmTypeOrganic, Name: "Organic / Soil-Based"},
		{Code: FarmTypeHydroponic, Name: "Hydroponic"},
		{Code: FarmTypeAquaponic, Name: "Aquaponic"},
		{Code: FarmTypeMushroom, Name: "Mushroom"},
		{Code: FarmTypeLiveStock, Name: "Livestock"},
		{Code: FarmTypeFisheries, Name: "Fisheries"},
		{Code: FarmTypePermaculture, Name: "Permaculture"},
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
