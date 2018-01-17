package domain

const (
	AreaTypeSeeding = "seeding"
	AreaTypeGrowing = "growing"
)

type AreaType struct {
	Code string
	Name string
}

func FindAllAreaTypes() []AreaType {
	return []AreaType{
		AreaType{Code: AreaTypeSeeding, Name: "Seeding"},
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
