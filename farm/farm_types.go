package farm

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
	Code string
	Name string
}

func FindAllFarmTypes() []FarmType {
	var types []FarmType

	types = append(types, FarmType{Code: FarmTypeOrganic, Name: "Organic / Soil-Based"})
	types = append(types, FarmType{Code: FarmTypeHydroponic, Name: "Hydroponic"})
	types = append(types, FarmType{Code: FarmTypeAquaponic, Name: "Aquaponic"})
	types = append(types, FarmType{Code: FarmTypeMushroom, Name: "Mushroom"})
	types = append(types, FarmType{Code: FarmTypeLiveStock, Name: "Livestock"})
	types = append(types, FarmType{Code: FarmTypeFisheries, Name: "Fisheries"})
	types = append(types, FarmType{Code: FarmTypePermaculture, Name: "Permaculture"})

	return types
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
