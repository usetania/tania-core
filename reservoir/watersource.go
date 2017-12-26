package reservoir

const (
	ReservoirEmptyBucketPercentage = 0.2
)

// Bucket is value object attached to the Reservoir.waterSource.
type Bucket struct {
	Capacity float32
	Volume   float32
}

// Tap is value object attached to the Reservoir.waterSource entity.
type Tap struct {
}

// CreateBucket registers a new Bucket.
func CreateBucket(capacity, volume float32) (Bucket, error) {
	if capacity <= 0 {
		return Bucket{}, ReservoirError{ReservoirErrorBucketCapacityInvalidCode}
	}

	err := validateVolume(capacity, volume)
	if err != nil {
		return Bucket{}, err
	}

	return Bucket{Capacity: capacity, Volume: volume}, nil
}

// CreateTap registers a new Tab.
func CreateTap() (Tap, error) {
	return Tap{}, nil
}

// ChangeVolume changes the amount of water in the Bucket.
func (b *Bucket) ChangeVolume(volume float32) error {
	err := validateVolume(b.Capacity, volume)

	if err != nil {
		return err
	}

	b.Volume = volume

	return nil
}

// IsBucketEmpty is used to check if bucket is empty.
func (b Bucket) IsBucketEmpty() bool {
	if b.Volume < b.Capacity*ReservoirEmptyBucketPercentage {
		return true
	}

	return false
}

func validateVolume(capacity, volume float32) error {
	if volume > capacity || volume < 0 {
		return ReservoirError{ReservoirErrorBucketCapacityInvalidCode}
	}

	return nil
}
