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
	assert.NotEqual(t, bucket, Bucket{})
	assert.Nil(t, err1)

	assert.Equal(t, tap, Tap{})
	assert.Nil(t, err2)
}
