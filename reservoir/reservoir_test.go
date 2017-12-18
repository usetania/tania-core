package reservoir

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateReservoir(t *testing.T) {
	// Given
	name := "My reservoir"
	ph := float32(10.0)
	ec := float32(12.34)
	temperature := float32(27.0)

	// When
	reservoir, err := CreateReservoir(name, ph, ec, temperature)

	// Then
	assert.Nil(t, err)
	assert.NotEqual(t, reservoir, Reservoir{})
}

func TestAbnormalCreateReservoir(t *testing.T) {
	// Given
	name := ""

	// When
	_, err := CreateReservoir(name, 0, 0, 0)

	// Then
	assert.Equal(t, err, ReservoirError{ReservoirErrorEmptyNameCode})

	// Given
	name = "My Reservoir"
	ph := float32(-10)

	// When
	_, err = CreateReservoir(name, ph, 0, 0)

	// Then
	assert.Equal(t, err, ReservoirError{ReservoirErrorInvalidPHCode})

	// Given
	name = "My Reservoir"
	ec := float32(0)
	ec2 := float32(-10)

	// When
	_, err = CreateReservoir(name, 0, ec, 0)
	_, err2 := CreateReservoir(name, 0, ec2, 0)

	// Then
	assert.Equal(t, err, ReservoirError{ReservoirErrorInvalidECCode})
	assert.Equal(t, err2, ReservoirError{ReservoirErrorInvalidECCode})

}

func TestCreateWaterSource(t *testing.T) {
	// Given

	// When
	bucket, err1 := CreateBucket(100)
	tap, err2 := CreateTap()

	// Then
	assert.NotEqual(t, bucket, Bucket{})
	assert.Nil(t, err1)

	assert.Equal(t, tap, Tap{})
	assert.Nil(t, err2)
}

func TestAttachBucket(t *testing.T) {
	// Given
	reservoir, err := CreateReservoir("My Reservoir", 8, 24.5, 31.8)
	bucket, err := CreateBucket(100)

	// When
	err = reservoir.AttachBucket(&bucket)

	// Then
	val := reservoir.waterSource

	assert.Equal(t, val, bucket)
	assert.Nil(t, err)
}
