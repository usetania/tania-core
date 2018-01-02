package entity

const (
	AreaLocationOutdoor = "outdoor"
	AreaLocationIndoor  = "indoor"
)

type AreaLocation struct {
	Code string
	Name string
}

func FindAllAreaLocation() []AreaLocation {
	return []AreaLocation{
		AreaLocation{Code: AreaLocationOutdoor, Name: "Field (Outdoor)"},
		AreaLocation{Code: AreaLocationIndoor, Name: "Greenhouse (Indoor)"},
	}
}

func FindAreaLocationByCode(code string) (AreaLocation, error) {
	items := FindAllAreaLocation()

	for _, item := range items {
		if item.Code == code {
			return item, nil
		}
	}

	return AreaLocation{}, AreaError{AreaErrorInvalidAreaLocationCode}
}
