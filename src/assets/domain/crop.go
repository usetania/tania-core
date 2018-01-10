package domain

import uuid "github.com/satori/go.uuid"

type Crop struct {
	CurrentArea Area
	Type        CropType
}

type CropType interface {
	Code() string
}
type Seeding struct {
}
type Growing struct {
}

func (s Seeding) Code() string {
	return "seeding"
}

func (g Growing) Code() string {
	return "growing"
}

func CreateCropBatch(area Area, cropType CropType) (Crop, error) {
	if area.UID == (uuid.UUID{}) {
		return Crop{}, CropError{Code: CropErrorInvalidArea}
	}

	return Crop{
		CurrentArea: area,
		Type:        cropType,
	}, nil
}
