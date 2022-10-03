package mathhelper_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/usetania/tania-core/src/helper/mathhelper"
)

func TestIsEqualFloat32(t *testing.T) {
	t.Parallel()
	// Given
	// When
	val1 := mathhelper.IsEqual(float32(10.0001), float32(10.00009))
	val2 := mathhelper.IsEqual(float32(10.0005), float32(10.00001))
	val3 := mathhelper.IsEqual(float32(10.00009), float32(10.0001))
	val4 := mathhelper.IsEqual(float32(10.00001), float32(10.0005))

	// Then
	assert.Equal(t, val1, true)
	assert.Equal(t, val2, false)
	assert.Equal(t, val3, true)
	assert.Equal(t, val4, false)
}
