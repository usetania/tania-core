package domain

const (
	AreaTypeNursery = "nursery"
	AreaTypeGrowing = "growing"
)

type AreaType struct {
	Code string
	Name string
}

func FindAllAreaTypes() []AreaType {
	return []AreaType{
		AreaType{Code: AreaTypeNursery, Name: "Nursery"},
		AreaType{Code: AreaTypeGrowing, Name: "Growing"},
	}
}

func FindAreaTypeByCode(code string) (AreaType, error) {
	items := FindAllAreaTypes()

	for _, item := range items {
		if item.Code == code {
			return item, nil
		}
	}

	return AreaType{}, AreaError{AreaErrorInvalidAreaTypeCode}
}
