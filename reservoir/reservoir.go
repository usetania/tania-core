// Package reservoir provides the operation that farm owner or his/her staff
// can do with the reservoir in a farm
package reservoir

// ReservoirError is a custom error from Go built-in error
type ReservoirError struct {
	code int
}

func (e ReservoirError) Error() string {
	switch e.code {
	case ReservoirErrorEmptyNameCode:
		return "Reservoir name is required."
	case ReservoirErrorInvalidPHCode:
		return "Reservoir pH value is invalid."
	case ReservoirErrorInvalidECCode:
		return "Reservoir EC value is invalid."
	case ReservoirErrorInvalidCapacityCode:
		return "Reservoir bucket is invalid."
	case ReservoirErrorBucketAlreadyAttachedCode:
		return "Reservoir water source is already attached."
	default:
		return "Unrecognized Reservoir Error Code"
	}
}

const (
	ReservoirErrorEmptyNameCode = iota
	ReservoirErrorInvalidPHCode
	ReservoirErrorInvalidECCode
	ReservoirErrorInvalidCapacityCode
	ReservoirErrorBucketAlreadyAttachedCode
)

// Reservoir is entity that provides the operation that farm owner or his/her staff
// can do with the reservoir in a farm
type Reservoir struct {
	UID         string
	Name        string
	PH          float32
	EC          float32
	Temperature float32

	waterSource interface{}
}

// Bucket is value object attached to the Reservoir.waterSource
type Bucket struct {
	Capacity float32
}

// Tap is value object attached to the Reservoir.waterSource entity
type Tap struct {
}

// CreateReservoir registers a new Reservoir
func CreateReservoir(name string, ph, ec, temperature float32) (Reservoir, error) {
	err := validateName(name)
	if err != nil {
		return Reservoir{}, err
	}

	err = validatePH(ph)
	if err != nil {
		return Reservoir{}, err
	}

	err = validateEC(ph)
	if err != nil {
		return Reservoir{}, err
	}

	return Reservoir{
		Name:        name,
		PH:          ph,
		EC:          ec,
		Temperature: temperature,
	}, nil
}

// CreateBucket registers a new Bucket
func CreateBucket(capacity float32) (Bucket, error) {
	if capacity <= 0 {
		return Bucket{}, ReservoirError{ReservoirErrorInvalidCapacityCode}
	}

	return Bucket{Capacity: capacity}, nil
}

// CreateTap registers a new Tab
func CreateTap() (Tap, error) {
	return Tap{}, nil
}

// AttachBucket attach Bucket value object to Reservoir.waterSource
func (r *Reservoir) AttachBucket(bucket *Bucket) error {
	if r.IsAttachedToWaterSource() {
		return ReservoirError{ReservoirErrorBucketAlreadyAttachedCode}
	}

	r.waterSource = *bucket
	return nil
}

// AttachTap attach Tap value object to Reservoir.waterSource
func (r *Reservoir) AttachTap(tap *Tap) error {
	if r.IsAttachedToWaterSource() {
		return ReservoirError{ReservoirErrorBucketAlreadyAttachedCode}
	}

	r.waterSource = *tap
	return nil
}

// IsAttachedToTap is used to check if Reservoir is attached to Tap WaterSource
func (r Reservoir) IsAttachedToTap() bool {
	_, ok := r.waterSource.(Tap)
	return ok
}

// IsAttachedToBucket is used to check if Reservoir is attached to Bucket WaterSource
func (r Reservoir) IsAttachedToBucket() bool {
	_, ok := r.waterSource.(Bucket)
	return ok
}

// IsAttachedToWaterSource is used to check if Reservoir is attached to WaterSource
func (r Reservoir) IsAttachedToWaterSource() bool {
	return r.waterSource != nil
}

// MeasureCondition will measure the Reservoir condition based on its properties
func (r Reservoir) MeasureCondition() float32 {
	if !r.IsAttachedToBucket() {
		// We can't measure non bucket reservoir
		return 0
	}

	// Do measure here
	return 1
}

// ChangeInformation is used to change Reservoir information. All arguments are optional
func (r *Reservoir) ChangeInformation(name string, ph, ec, temperature float32) error {
	err := validatePH(ph)
	if err != nil {
		return err
	}

	err = validateEC(ec)
	if err != nil {
		return err
	}

	if name != "" {
		r.Name = name
	}
	if ph != 0 {
		r.PH = ph
	}
	if ec != 0 {
		r.EC = ec
	}
	if temperature != 0 {
		r.Temperature = temperature
	}

	return nil
}

func validateName(name string) error {
	if name == "" {
		return ReservoirError{ReservoirErrorEmptyNameCode}
	}

	return nil
}

func validatePH(ph float32) error {
	if ph < 0 {
		return ReservoirError{ReservoirErrorInvalidPHCode}
	}

	return nil
}

func validateEC(ec float32) error {
	if ec <= 0 {
		return ReservoirError{ReservoirErrorInvalidECCode}
	}

	return nil
}
