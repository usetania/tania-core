package mathhelper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEqualFloat32(t *testing.T) {
	// Given

	// When
	val1 := IsEqual(float32(10.0001), float32(10.00009))
	val2 := IsEqual(float32(10.0005), float32(10.00001))
	val3 := IsEqual(float32(10.00009), float32(10.0001))
	val4 := IsEqual(float32(10.00001), float32(10.0005))

	// Then
	assert.Equal(t, val1, true)
	assert.Equal(t, val2, false)
	assert.Equal(t, val3, true)
	assert.Equal(t, val4, false)
}
