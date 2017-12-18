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

func TestAttachWaterSource(t *testing.T) {
	// Given
	reservoir1, _ := CreateReservoir("My Reservoir 1", 8, 24.5, 31.8)
	bucket, _ := CreateBucket(100, 50)

	reservoir2, _ := CreateReservoir("My Reservoir 2", 8, 24.5, 31.8)
	tap, _ := CreateTap()

	// When
	err1 := reservoir1.AttachBucket(&bucket)
	err2 := reservoir2.AttachTap(&tap)

	// Then
	val1 := reservoir1.waterSource
	val2 := reservoir2.waterSource

	assert.Equal(t, val1, bucket)
	assert.Nil(t, err1)

	assert.Equal(t, val2, tap)
	assert.Nil(t, err2)
}

func TestInvalidAttachWaterSource(t *testing.T) {
	// Given
	reservoir, _ := CreateReservoir("My Reservoir", 8, 24.5, 31.8)
	bucket1, _ := CreateBucket(100, 50)
	bucket2, _ := CreateBucket(200, 150)
	tap, _ := CreateTap()

	// When
	reservoir.AttachBucket(&bucket1)
	err1 := reservoir.AttachBucket(&bucket2)
	err2 := reservoir.AttachTap(&tap)

	assert.Equal(t, err1, ReservoirError{ReservoirErrorWaterSourceAlreadyAttachedCode})
	assert.Equal(t, err2, ReservoirError{ReservoirErrorWaterSourceAlreadyAttachedCode})
}

func TestMeasureCondition(t *testing.T) {
	// Given
	reservoir1, _ := CreateReservoir("My Reservoir 1", 8, 24.5, 31.8)
	bucket, _ := CreateBucket(100, 50)
	reservoir1.AttachBucket(&bucket)

	reservoir2, _ := CreateReservoir("My Reservoir 2", 10, 21.2, 34.2)
	tap, _ := CreateTap()
	reservoir2.AttachTap(&tap)

	// When
	val1 := reservoir1.MeasureCondition()
	val2 := reservoir2.MeasureCondition()

	// Then
	assert.Equal(t, val1, float32(1))
	assert.Equal(t, val2, float32(0))
}
