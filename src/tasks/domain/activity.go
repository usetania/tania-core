package domain

import (
	uuid "github.com/satori/go.uuid"
	"math"
	"strconv"
)

type Activity interface {
	Type() string
}

// SeedActivity
type SeedActivity struct {
}

func (sa SeedActivity) Type() string {
	return ActivityTypeSeed
}

func CreateSeedActivity() (SeedActivity, error) {
	return SeedActivity{}, nil
}

// FertilizeActivity
type FertilizeActivity struct {
}

func (sa FertilizeActivity) Type() string {
	return ActivityTypeFertilize
}

func CreateFertilizeActivity() (FertilizeActivity, error) {
	return FertilizeActivity{}, nil
}

// PruneActivity
type PruneActivity struct {
}

func (sa PruneActivity) Type() string {
	return ActivityTypePrune
}

func CreatePruneActivity() (PruneActivity, error) {
	return PruneActivity{}, nil
}

// PesticideActivity
type PesticideActivity struct {
}

func (sa PesticideActivity) Type() string {
	return ActivityTypePesticide
}

func CreatePesticideActivity() (PesticideActivity, error) {
	return PesticideActivity{}, nil
}

// MoveToAreaActivity
type MoveToAreaActivity struct {
	SourceAreaID      uuid.UUID `json:"source_area_id"`
	DestinationAreaID uuid.UUID `json:"destination_area_id"`
	Quantity          float32   `json:"quantity"`
}

func (sa MoveToAreaActivity) Type() string {
	return ActivityTypeMoveToArea
}

func CreateMoveToAreaActivity(taskService TaskService, source string, dest string, qnt string) (MoveToAreaActivity, error) {

	err := validateAreaID(source)
	if err != nil {
		return MoveToAreaActivity{}, TaskError{TaskErrorActivitySourceInvalidCode}
	}
	src_id, _ := uuid.FromString(source)
	err = validateAssetID(taskService, src_id, TaskTypeArea)
	if err != nil {
		return MoveToAreaActivity{}, TaskError{TaskErrorActivitySourceInvalidCode}
	}

	err = validateAreaID(dest)
	if err != nil {
		return MoveToAreaActivity{}, TaskError{TaskErrorActivityDestinationInvalidCode}
	}
	dest_id, _ := uuid.FromString(dest)
	err = validateAssetID(taskService, dest_id, TaskTypeArea)
	if err != nil {
		return MoveToAreaActivity{}, TaskError{TaskErrorActivityDestinationInvalidCode}
	}

	quantity64, err := strconv.ParseFloat(qnt, 32)
	if err != nil || math.IsNaN(quantity64) {
		return MoveToAreaActivity{}, TaskError{TaskErrorActivityQuantityInvalidCode}
	}
	quantity32 := float32(quantity64)
	err = validateQuantity(quantity32)
	if err != nil {
		return MoveToAreaActivity{}, err
	}
	return MoveToAreaActivity{SourceAreaID: src_id, DestinationAreaID: dest_id, Quantity: quantity32}, nil
}

// DumpActivity
type DumpActivity struct {
	SourceAreaID uuid.UUID `json:"source_area_id"`
	Quantity     float32   `json:"quantity"`
}

func (sa DumpActivity) Type() string {
	return ActivityTypeDump
}

func CreateDumpActivity(taskService TaskService, source string, qnt string) (DumpActivity, error) {

	err := validateAreaID(source)
	if err != nil {
		return DumpActivity{}, TaskError{TaskErrorActivitySourceInvalidCode}
	}
	src_id, _ := uuid.FromString(source)
	err = validateAssetID(taskService, src_id, TaskTypeArea)
	if err != nil {
		return DumpActivity{}, TaskError{TaskErrorActivitySourceInvalidCode}
	}

	quantity64, err := strconv.ParseFloat(qnt, 32)
	if err != nil || math.IsNaN(quantity64) {
		return DumpActivity{}, TaskError{TaskErrorActivityQuantityInvalidCode}
	}
	quantity32 := float32(quantity64)
	err = validateQuantity(quantity32)
	if err != nil {
		return DumpActivity{}, err
	}
	return DumpActivity{SourceAreaID: src_id, Quantity: quantity32}, nil
}

// HarvestActivity
type HarvestActivity struct {
	SourceAreaID uuid.UUID `json:"source_area_id"`
	Quantity     float32   `json:"quantity"`
}

func (sa HarvestActivity) Type() string {
	return ActivityTypeHarvest
}

func CreateHarvestActivity(taskService TaskService, source string, qnt string) (HarvestActivity, error) {

	err := validateAreaID(source)
	if err != nil {
		return HarvestActivity{}, TaskError{TaskErrorActivitySourceInvalidCode}
	}
	src_id, _ := uuid.FromString(source)
	err = validateAssetID(taskService, src_id, TaskTypeArea)
	if err != nil {
		return HarvestActivity{}, TaskError{TaskErrorActivitySourceInvalidCode}
	}

	quantity64, err := strconv.ParseFloat(qnt, 32)
	if err != nil || math.IsNaN(quantity64) {
		return HarvestActivity{}, TaskError{TaskErrorActivityQuantityInvalidCode}
	}
	quantity32 := float32(quantity64)
	err = validateQuantity(quantity32)
	if err != nil {
		return HarvestActivity{}, err
	}
	return HarvestActivity{SourceAreaID: src_id, Quantity: quantity32}, nil
}

//Validation
func validateAreaID(id string) error {
	_, err := uuid.FromString(id)
	if err != nil {
		return err
	}
	return nil
}
func validateQuantity(qnt float32) error {
	if qnt < 0 {
		return TaskError{TaskErrorActivityQuantityInvalidCode}
	}

	return nil
}
