package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateReservoir(t *testing.T) {
	// Given
	name := "MyReservoir"
	farm, farmErr := CreateFarm("Farm1", FarmTypeOrganic)

	// When
	reservoir, err := CreateReservoir(farm, name)

	// Then
	assert.Nil(t, farmErr)
	assert.Nil(t, err)
	assert.NotEqual(t, reservoir, Reservoir{})
}

func TestInvalidCreateReservoir(t *testing.T) {
	// Given
	farm, farmErr := CreateFarm("Farm1", FarmTypeOrganic)

	reservoirData := []struct {
		farm     Farm
		name     string
		expected ReservoirError
	}{
		{farm, "My<>Reserv", ReservoirError{ReservoirErrorNameAlphanumericOnlyCode}},
		{farm, "MyR", ReservoirError{ReservoirErrorNameNotEnoughCharacterCode}},
		{farm, "MyReservoirNameShouldNotBeMoreThanAHundredLongCharactersMyReservoirNameShouldNotBeMoreThanAHundredLongCharacters", ReservoirError{ReservoirErrorNameExceedMaximunCharacterCode}},
	}

	for _, data := range reservoirData {
		// When
		_, err := CreateReservoir(data.farm, data.name)

		// Then
		assert.Equal(t, data.expected, err)
	}

	// Then
	assert.Nil(t, farmErr)
}

func TestAttachWaterSource(t *testing.T) {
	// Given
	farm, farmErr := CreateFarm("Farm1", FarmTypeOrganic)

	reservoir1, resErr1 := CreateReservoir(farm, "MyReservoir1")
	bucket, bucketErr := CreateBucket(100, 50)

	reservoir2, resErr2 := CreateReservoir(farm, "MyReservoir2")
	tap, tapErr := CreateTap()

	// When
	err1 := reservoir1.AttachBucket(bucket)
	err2 := reservoir2.AttachTap(tap)

	// Then
	val1 := reservoir1.WaterSource
	val2 := reservoir2.WaterSource

	assert.Nil(t, farmErr)
	assert.Nil(t, resErr1)
	assert.Nil(t, resErr2)
	assert.Nil(t, bucketErr)
	assert.Nil(t, tapErr)

	assert.Equal(t, val1, bucket)
	assert.Nil(t, err1)

	assert.Equal(t, val2, tap)
	assert.Nil(t, err2)
}

func TestInvalidAttachWaterSource(t *testing.T) {
	// Given
	farm, farmErr := CreateFarm("Farm1", FarmTypeOrganic)
	reservoir, reservoirErr := CreateReservoir(farm, "MyReservoir")
	bucket1, bucket1Err := CreateBucket(100, 50)
	bucket2, bucket2Err := CreateBucket(200, 150)
	tap, _ := CreateTap()

	// When
	reservoir.AttachBucket(bucket1)
	err1 := reservoir.AttachBucket(bucket2)
	err2 := reservoir.AttachTap(tap)

	// Then
	assert.Nil(t, farmErr)
	assert.Nil(t, reservoirErr)
	assert.Nil(t, bucket1Err)
	assert.Nil(t, bucket2Err)
	assert.Equal(t, err1, ReservoirError{ReservoirErrorWaterSourceAlreadyAttachedCode})
	assert.Equal(t, err2, ReservoirError{ReservoirErrorWaterSourceAlreadyAttachedCode})
}

func TestMeasureCondition(t *testing.T) {
	// Given
	farm, farmErr := CreateFarm("Farm1", FarmTypeOrganic)

	reservoir1, reservoir1Err := CreateReservoir(farm, "MyReservoir1")
	bucket, bucketErr := CreateBucket(100, 50)
	reservoir1.AttachBucket(bucket)

	reservoir2, reservoir2Err := CreateReservoir(farm, "MyReservoir2")
	tap, tapErr := CreateTap()
	reservoir2.AttachTap(tap)

	// When
	val1 := reservoir1.MeasureCondition()
	val2 := reservoir2.MeasureCondition()

	// Then
	assert.Nil(t, farmErr)
	assert.Nil(t, reservoir1Err)
	assert.Nil(t, reservoir2Err)
	assert.Nil(t, bucketErr)
	assert.Nil(t, tapErr)
	assert.Equal(t, val1, float32(1))
	assert.Equal(t, val2, float32(0))
}

func TestChangeTemperature(t *testing.T) {
	// Given
	farm, farmErr := CreateFarm("Farm1", FarmTypeOrganic)
	reservoir, reservoirErr := CreateReservoir(farm, "MyReservoir")
	temperature := float32(32)
	ph := float32(4.3)
	ec := float32(23.5)

	// When
	reservoir.ChangeTemperature(temperature, ph, ec)

	// Then
	assert.Nil(t, farmErr)
	assert.Nil(t, reservoirErr)
	assert.Equal(t, reservoir.Temperature, temperature)
	assert.Equal(t, reservoir.PH, ph)
	assert.Equal(t, reservoir.EC, ec)
}

func TestInvalidChangeTemperature(t *testing.T) {
	// Given
	farm, farmErr := CreateFarm("Farm1", FarmTypeOrganic)
	reservoir, reservoirErr := CreateReservoir(farm, "MyReservoir")
	temperature := float32(32)
	ph1 := float32(-10)
	ec1 := float32(23.5)
	ph2 := float32(4)
	ec2 := float32(-1)

	// When
	err1 := reservoir.ChangeTemperature(temperature, ph1, ec1)
	err2 := reservoir.ChangeTemperature(temperature, ph2, ec2)

	// Then
	assert.Nil(t, farmErr)
	assert.Nil(t, reservoirErr)
	assert.Equal(t, err1, ReservoirError{ReservoirErrorPHInvalidCode})
	assert.Equal(t, err2, ReservoirError{ReservoirErrorECInvalidCode})
}

func TestReservoirCreateRemoveNote(t *testing.T) {
	// Given
	farm, farmErr := CreateFarm("Farm1", FarmTypeOrganic)
	reservoir, reservoirErr := CreateReservoir(farm, "MyReservoir")

	// When
	reservoir.AddNewNote("This is my new note")

	// Then
	assert.Nil(t, farmErr)
	assert.Nil(t, reservoirErr)

	assert.Equal(t, 1, len(reservoir.Notes))
	assert.Equal(t, "This is my new note", reservoir.Notes[0].Content)
	assert.NotNil(t, reservoir.Notes[0].CreatedDate)

	// When
	reservoir.RemoveNote("This is my new note")

	assert.Equal(t, 0, len(reservoir.Notes))
}
