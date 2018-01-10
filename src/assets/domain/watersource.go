package domain

const (
	ReservoirEmptyBucketPercentage = 0.2
)

type WaterSource interface {
	Type() string
}

// Bucket is value object attached to the Reservoir.waterSource.
type Bucket struct {
	Capacity float32 `json:"capacity"`
	Volume   float32 `json:"volume"`
}

// Tap is value object attached to the Reservoir.waterSource domain.
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

func (b Bucket) Type() string {
	return "bucket"
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

func (t Tap) Type() string {
	return "tap"
}

func validateVolume(capacity, volume float32) error {
	if volume > capacity || volume < 0 {
		return ReservoirError{ReservoirErrorBucketVolumeInvalidCode}
	}

	return nil
}
