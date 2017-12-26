package reservoir

import (
	"testing"

	"github.com/Tanibox/tania-server/helper/mathhelper"
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
	assert.InDelta(t, bucket.Capacity, float32(100), mathhelper.EPSILON)
	assert.InDelta(t, bucket.Volume, float32(50), mathhelper.EPSILON)

	assert.Nil(t, err2)
	assert.Equal(t, tap, Tap{})
}

func TestVolume(t *testing.T) {
	// Given
	bucket, _ := CreateBucket(100, 50)

	// When
	bucket.ChangeVolume(85)

	// Then
	assert.Equal(t, bucket.Volume, float32(85))
}

func TestInvalidVolume(t *testing.T) {
	// Given
	bucket1, _ := CreateBucket(100, 50)
	bucket2, _ := CreateBucket(100, 50)

	// When
	err1 := bucket1.ChangeVolume(110)
	err2 := bucket2.ChangeVolume(-5)

	// Then
	assert.Equal(t, err1, ReservoirError{ReservoirErrorBucketVolumeInvalidCode})
	assert.Equal(t, err2, ReservoirError{ReservoirErrorBucketVolumeInvalidCode})
}

func TestEmptyBucket(t *testing.T) {
	// Given
	bucket, _ := CreateBucket(100, 50)

	// When
	bucket.ChangeVolume(19)
	val1 := bucket.IsBucketEmpty()

	bucket.ChangeVolume(20)
	val2 := bucket.IsBucketEmpty()

	// Then
	assert.Equal(t, val1, true)
	assert.Equal(t, val2, false)
}
