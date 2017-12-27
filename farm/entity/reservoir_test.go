package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateReservoir(t *testing.T) {
	// Given
	name := "My reservoir"
	farm, _ := CreateFarm("Farm 1", "This is our farm", "10.00", "11.00", FarmTypeOrganic, "ID", "ID")

	// When
	reservoir, err := CreateReservoir(farm, name)

	// Then
	assert.Nil(t, err)
	assert.NotEqual(t, reservoir, Reservoir{})
}

func TestInvalidCreateReservoir(t *testing.T) {
	// Given
	name := ""
	farm, _ := CreateFarm("Farm 1", "This is our farm", "10.00", "11.00", FarmTypeOrganic, "ID", "ID")

	// When
	_, err := CreateReservoir(farm, name)

	// Then
	assert.Equal(t, err, ReservoirError{ReservoirErrorNameEmptyCode})

	// Given
	name = "asd"

	// When
	_, err = CreateReservoir(farm, name)

	// Then
	assert.Equal(t, err, ReservoirError{ReservoirErrorNameNotEnoughCharacterCode})
}

func TestAttachWaterSource(t *testing.T) {
	// Given
	farm, _ := CreateFarm("Farm 1", "This is our farm", "10.00", "11.00", FarmTypeOrganic, "ID", "ID")

	reservoir1, _ := CreateReservoir(farm, "My Reservoir 1")
	bucket, _ := CreateBucket(100, 50)

	reservoir2, _ := CreateReservoir(farm, "My Reservoir 2")
	tap, _ := CreateTap()

	// When
	err1 := reservoir1.AttachBucket(&bucket)
	err2 := reservoir2.AttachTap(&tap)

	// Then
	val1 := reservoir1.WaterSource
	val2 := reservoir2.WaterSource

	assert.Equal(t, val1, &bucket)
	assert.Nil(t, err1)

	assert.Equal(t, val2, &tap)
	assert.Nil(t, err2)
}

func TestInvalidAttachWaterSource(t *testing.T) {
	// Given
	farm, _ := CreateFarm("Farm 1", "This is our farm", "10.00", "11.00", FarmTypeOrganic, "ID", "ID")
	reservoir, _ := CreateReservoir(farm, "My Reservoir")
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
	farm, _ := CreateFarm("Farm 1", "This is our farm", "10.00", "11.00", FarmTypeOrganic, "ID", "ID")

	reservoir1, _ := CreateReservoir(farm, "My Reservoir 1")
	bucket, _ := CreateBucket(100, 50)
	reservoir1.AttachBucket(&bucket)

	reservoir2, _ := CreateReservoir(farm, "My Reservoir 2")
	tap, _ := CreateTap()
	reservoir2.AttachTap(&tap)

	// When
	val1 := reservoir1.MeasureCondition()
	val2 := reservoir2.MeasureCondition()

	// Then
	assert.Equal(t, val1, float32(1))
	assert.Equal(t, val2, float32(0))
}

func TestChangeTemperature(t *testing.T) {
	// Given
	farm, _ := CreateFarm("Farm 1", "This is our farm", "10.00", "11.00", FarmTypeOrganic, "ID", "ID")
	reservoir, _ := CreateReservoir(farm, "My Reservoir")
	temperature := float32(32)
	ph := float32(4.3)
	ec := float32(23.5)

	// When
	reservoir.ChangeTemperature(temperature, ph, ec)

	// Then
	assert.Equal(t, reservoir.Temperature, temperature)
	assert.Equal(t, reservoir.PH, ph)
	assert.Equal(t, reservoir.EC, ec)
}

func TestInvalidChangeTemperature(t *testing.T) {
	// Given
	farm, _ := CreateFarm("Farm 1", "This is our farm", "10.00", "11.00", FarmTypeOrganic, "ID", "ID")
	reservoir, _ := CreateReservoir(farm, "My Reservoir")
	temperature := float32(32)
	ph1 := float32(-10)
	ec1 := float32(23.5)
	ph2 := float32(4)
	ec2 := float32(-1)

	// When
	err1 := reservoir.ChangeTemperature(temperature, ph1, ec1)
	err2 := reservoir.ChangeTemperature(temperature, ph2, ec2)

	// Then
	assert.Equal(t, err1, ReservoirError{ReservoirErrorPHInvalidCode})
	assert.Equal(t, err2, ReservoirError{ReservoirErrorECInvalidCode})
}
