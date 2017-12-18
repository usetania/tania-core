package reservoir

import (
	"testing"
)

func TestCreateReservoir(t *testing.T) {
	// Given
	name := "My reservoir"
	ph := float32(10.0)
	ec := float32(12.34)
	temperature := float32(27.0)

	// When
	_, err := CreateReservoir(name, ph, ec, temperature)

	// Then
	if err != nil {
		t.Error(err)
	}
}

func TestAbnormalCreateReservoir(t *testing.T) {
	// Given
	name := ""

	// When
	_, err := CreateReservoir(name, 0, 0, 0)

	// Then
	if err != ReservoirErrorEmptyName {
		t.Error("Expected error return: ", ReservoirErrorEmptyName)
	}

	// Given
	name = "My Reservoir"
	ph := float32(-10)

	// When
	_, err = CreateReservoir(name, ph, 0, 0)

	// Then
	if err != ReservoirErrorInvalidPH {
		t.Error("Expected error return: ", ReservoirErrorInvalidPH)
	}

	// Given
	name = "My Reservoir"
	ec := float32(0)
	ec2 := float32(-10)

	// When
	_, err = CreateReservoir(name, 0, ec, 0)
	_, err2 := CreateReservoir(name, 0, ec2, 0)

	// Then
	if err != ReservoirErrorInvalidEC {
		t.Error("Expected error return: ", ReservoirErrorInvalidEC)
	}
	if err2 != ReservoirErrorInvalidEC {
		t.Error("Expected error return: ", ReservoirErrorInvalidEC)
	}
}

func TestCreateWaterSource(t *testing.T) {
	// Given

	// When
	bucket, err := CreateBucket(100)
	tap, err := CreateTap()

	// Then
	if bucket == (Bucket{}) && err != nil {
		t.Error("Expected error return: ", nil, ", found: ", err)
	}
	if tap == (Tap{}) && err != nil {
		t.Error("Expected error return: ", nil, ", found: ", err)
	}
}

func TestAttachBucket(t *testing.T) {
	// Given
	reservoir, err := CreateReservoir("My Reservoir", 8, 24.5, 31.8)
	bucket, err := CreateBucket(100)

	// When
	err = reservoir.AttachBucket(&bucket)

	// Then
	val := reservoir.waterSource

	if val == nil || err != nil {
		t.Error("Expected error: ", nil, ", found: ", err)
	}
}
