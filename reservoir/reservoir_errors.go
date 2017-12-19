package reservoir

const (
	ReservoirErrorEmptyNameCode = iota
	ReservoirErrorNotEnoughCharacterCode
	ReservoirErrorInvalidPHCode
	ReservoirErrorInvalidECCode
	ReservoirErrorInvalidCapacityCode
	ReservoirErrorWaterSourceAlreadyAttachedCode
	ReservoirErrorInvalidCurrentBucketCapacityCode
)

// ReservoirError is a custom error from Go built-in error
type ReservoirError struct {
	code int
}

func (e ReservoirError) Error() string {
	switch e.code {
	case ReservoirErrorEmptyNameCode:
		return "Reservoir name is required."
	case ReservoirErrorNotEnoughCharacterCode:
		return "Not enough character on Reservoir Name"
	case ReservoirErrorInvalidPHCode:
		return "Reservoir pH value is invalid."
	case ReservoirErrorInvalidECCode:
		return "Reservoir EC value is invalid."
	case ReservoirErrorInvalidCapacityCode:
		return "Reservoir bucket is invalid."
	case ReservoirErrorWaterSourceAlreadyAttachedCode:
		return "Reservoir water source is already attached."
	case ReservoirErrorInvalidCurrentBucketCapacityCode:
		return "Current Reservoir bucket capacity is invalid."
	default:
		return "Unrecognized Reservoir Error Code"
	}
}
