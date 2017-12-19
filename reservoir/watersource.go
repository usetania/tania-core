package reservoir

const (
	ReservoirEmptyBucketPercentage = 0.2
)

// Bucket is value object attached to the Reservoir.waterSource.
type Bucket struct {
	Capacity        float32
	CurrentCapacity float32
}

// Tap is value object attached to the Reservoir.waterSource entity.
type Tap struct {
}

// CreateBucket registers a new Bucket.
func CreateBucket(capacity, currentCapacity float32) (Bucket, error) {
	if capacity <= 0 {
		return Bucket{}, ReservoirError{ReservoirErrorInvalidCapacityCode}
	}

	err := validateCurrentCapacity(capacity, currentCapacity)
	if err != nil {
		return Bucket{}, err
	}

	return Bucket{Capacity: capacity, CurrentCapacity: currentCapacity}, nil
}

// CreateTap registers a new Tab.
func CreateTap() (Tap, error) {
	return Tap{}, nil
}

// ChangeCurrentCapacity changes the amount of water in the Bucket.
func (b *Bucket) ChangeCurrentCapacity(currentCapacity float32) error {
	err := validateCurrentCapacity(b.Capacity, currentCapacity)

	if err != nil {
		return err
	}

	b.CurrentCapacity = currentCapacity

	return nil
}

// IsBucketEmpty is used to check if bucket is empty.
func (b Bucket) IsBucketEmpty() bool {
	if b.CurrentCapacity < b.Capacity*ReservoirEmptyBucketPercentage {
		return true
	}

	return false
}

func validateCurrentCapacity(capacity, currentCapacity float32) error {
	if currentCapacity > capacity || currentCapacity < 0 {
		return ReservoirError{ReservoirErrorInvalidCurrentBucketCapacityCode}
	}

	return nil
}
