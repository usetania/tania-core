package domain

import (
	uuid "github.com/satori/go.uuid"
)

type Activity interface {
	Type() string
}

// SeedActivity
type SeedActivity struct {
}

func (sa SeedActivity) Type() string {
	return "seed"
}

func CreateSeedActivity() (SeedActivity, error) {
	return SeedActivity{}, nil
}

// FertilizeActivity
type FertilizeActivity struct {
}

func (sa FertilizeActivity) Type() string {
	return "fertilize"
}

func CreateFertilizeActivity() (FertilizeActivity, error) {
	return FertilizeActivity{}, nil
}

// PruneActivity
type PruneActivity struct {
}

func (sa PruneActivity) Type() string {
	return "prune"
}

func CreatePruneActivity() (PruneActivity, error) {
	return PruneActivity{}, nil
}

// PesticideActivity
type PesticideActivity struct {
}

func (sa PesticideActivity) Type() string {
	return "pesticide"
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
	return "movetoarea"
}

func CreateMoveToAreaActivity(source string, dest string, qnt float32) (MoveToAreaActivity, error) {

	src_id, err := uuid.FromString(source)
	if err != nil {
		return MoveToAreaActivity{}, err
	}

	dest_id, err := uuid.FromString(source)
	if err != nil {
		return MoveToAreaActivity{}, err
	}

	err = validateAreaID(src_id)
	if err != nil {
		return MoveToAreaActivity{}, err
	}
	err = validateAreaID(dest_id)
	if err != nil {
		return MoveToAreaActivity{}, err
	}
	err = validateQuantity(qnt)
	if err != nil {
		return MoveToAreaActivity{}, err
	}
	return MoveToAreaActivity{SourceAreaID: src_id, DestinationAreaID: dest_id, Quantity: qnt}, nil
}

// DumpActivity
type DumpActivity struct {
	SourceAreaID uuid.UUID `json:"source_area_id"`
	Quantity     float32   `json:"quantity"`
}

func (sa DumpActivity) Type() string {
	return "dump"
}

func CreateDumpActivity(source string, qnt float32) (DumpActivity, error) {

	src_id, err := uuid.FromString(source)
	if err != nil {
		return DumpActivity{}, err
	}

	err = validateAreaID(src_id)
	if err != nil {
		return DumpActivity{}, err
	}
	err = validateQuantity(qnt)
	if err != nil {
		return DumpActivity{}, err
	}
	return DumpActivity{SourceAreaID: src_id, Quantity: qnt}, nil
}

// HarvestActivity
type HarvestActivity struct {
	SourceAreaID uuid.UUID `json:"source_area_id"`
	Quantity     float32   `json:"quantity"`
}

func (sa HarvestActivity) Type() string {
	return "harvest"
}

func CreateHarvestActivity(source string, qnt float32) (HarvestActivity, error) {

	src_id, err := uuid.FromString(source)
	if err != nil {
		return HarvestActivity{}, err
	}

	err = validateAreaID(src_id)
	if err != nil {
		return HarvestActivity{}, err
	}
	err = validateQuantity(qnt)
	if err != nil {
		return HarvestActivity{}, err
	}
	return HarvestActivity{SourceAreaID: src_id, Quantity: qnt}, nil
}

//Validation
func validateAreaID(id uuid.UUID) error {
	//TODO

	return nil
}
func validateQuantity(qnt float32) error {
	if qnt < 0 {
		return TaskError{TaskErrorActivityQuantityInvalid}
	}

	return nil
}
