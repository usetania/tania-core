package reservoir

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateWaterSource(t *testing.T) {
	// Given

	// When
	bucket, err1 := CreateBucket(100, 50)
	tap, err2 := CreateTap()

	// Then
	assert.Nil(t, err1)
	assert.NotEqual(t, bucket, Bucket{})
	assert.Equal(t, bucket.Capacity, float32(100))
	assert.Equal(t, bucket.CurrentCapacity, float32(50))

	assert.Nil(t, err2)
	assert.Equal(t, tap, Tap{})
}

func TestCurrentCapacity(t *testing.T) {
	// Given
	bucket, _ := CreateBucket(100, 50)

	// When
	bucket.ChangeCurrentCapacity(85)

	// Then
	assert.Equal(t, bucket.CurrentCapacity, float32(85))
}

func TestInvalidCurrentCapacity(t *testing.T) {
	// Given
	bucket1, _ := CreateBucket(100, 50)
	bucket2, _ := CreateBucket(100, 50)

	// When
	err1 := bucket1.ChangeCurrentCapacity(110)
	err2 := bucket2.ChangeCurrentCapacity(-5)

	// Then
	assert.Equal(t, err1, ReservoirError{ReservoirErrorInvalidCurrentBucketCapacityCode})
	assert.Equal(t, err2, ReservoirError{ReservoirErrorInvalidCurrentBucketCapacityCode})
}

func TestEmptyBucket(t *testing.T) {
	// Given
	bucket, _ := CreateBucket(100, 50)

	// When
	bucket.ChangeCurrentCapacity(19)
	val1 := bucket.IsBucketEmpty()

	bucket.ChangeCurrentCapacity(20)
	val2 := bucket.IsBucketEmpty()

	// Then
	assert.Equal(t, val1, true)
	assert.Equal(t, val2, false)
}
