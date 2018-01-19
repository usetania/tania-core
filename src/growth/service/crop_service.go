package service

import (
	"strings"

	"github.com/Tanibox/tania-server/src/growth/domain"
	"github.com/Tanibox/tania-server/src/growth/query"
	"github.com/Tanibox/tania-server/src/helper/stringhelper"
	uuid "github.com/satori/go.uuid"
)

type CropService struct {
	CropQuery              query.CropQuery
	InventoryMaterialQuery query.InventoryMaterialQuery
}

func (s CropService) ChangeInventory(crop *domain.Crop, inventoryUID uuid.UUID) error {
	result := s.InventoryMaterialQuery.FindByID(inventoryUID)

	inv, ok := <-result.Result.(service.CropInventory)
	if !ok {
		return CropError{Code: CropInventoryErrorInvalidInventory}
	}

	if inv == (domain.InventoryMaterial{}) {
		return CropError{Code: CropInventoryErrorNotFound}
	}

	// Generate Batch ID
	// Format the date to become daymonth format like 25jan
	dateFormat := strings.ToLower(crop.CreatedDate.Format("2Jan"))

	// Get variety name and split it to slice
	varietySlice := strings.Fields(inventory.Variety)
	varietyFormat := ""
	for _, v := range varietySlice {
		// For every value, get only the first three characters
		format := ""
		if len(v) > 3 {
			format = strings.ToLower(string(v[0:3]))
		} else {
			format = strings.ToLower(string(v))
		}

		varietyFormat = stringhelper.Join(varietyFormat, format, "-")
	}

	// Join that variety and date
	batchID := stringhelper.Join(varietyFormat, dateFormat)

	// Validate Uniqueness of Batch ID.
	resultQuery := <-s.CropQuery.FindByBatchID(batchID)
	if resultQuery.Error != nil {
		return resultQuery.Error
	}

	cropFound, ok := resultQuery.Result.(domain.Crop)
	if !ok {
		return domain.CropError{Code: domain.CropErrorInvalidBatchID}
	}
	if cropFound.UID != (uuid.UUID{}) {
		return domain.CropError{Code: domain.CropErrorBatchIDAlreadyCreated}
	}

	crop.Inventory = inventory
	crop.BatchID = batchID

	return nil
}
