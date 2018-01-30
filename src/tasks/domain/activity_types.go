package domain

const (
	ActivityTypeSeed       = "seed"
	ActivityTypeFertilize  = "fertilize"
	ActivityTypePrune      = "prune"
	ActivityTypePesticide  = "pesticide"
	ActivityTypeMoveToArea = "movetoarea"
	ActivityTypeDump       = "dump"
	ActivityTypeHarvest    = "harvest"
)

type ActivityType struct {
	Code string
	Name string
}

func FindAllActivityTypes() []ActivityType {
	return []ActivityType{
		ActivityType{Code: ActivityTypeSeed, Name: "Seed"},
		ActivityType{Code: ActivityTypeFertilize, Name: "Fertilize"},
		ActivityType{Code: ActivityTypePrune, Name: "Prune"},
		ActivityType{Code: ActivityTypePesticide, Name: "Pesticide"},
		ActivityType{Code: ActivityTypeMoveToArea, Name: "MoveToArea"},
		ActivityType{Code: ActivityTypeDump, Name: "Dump"},
		ActivityType{Code: ActivityTypeHarvest, Name: "Harvest"},
	}
}

func FindActivityTypeByCode(code string) (ActivityType, error) {
	items := FindAllActivityTypes()

	for _, item := range items {
		if item.Code == code {
			return item, nil
		}
	}

	return ActivityType{}, TaskError{TaskErrorActivityTypeInvalidCode}
}
