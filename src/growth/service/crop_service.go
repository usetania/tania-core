package service

import (
	"strings"

	"github.com/Tanibox/tania-server/src/growth/domain"
	"github.com/Tanibox/tania-server/src/growth/query"
	"github.com/Tanibox/tania-server/src/helper/stringhelper"
	uuid "github.com/satori/go.uuid"
)

type CropService struct {
	AreaQuery              query.AreaQuery
	CropQuery              query.CropQuery
	InventoryMaterialQuery query.InventoryMaterialQuery
}

func (s CropService) ChangeInventory(crop *domain.Crop, inventoryUID uuid.UUID) error {
	result := s.InventoryMaterialQuery.FindByID(inventoryUID)

	inv, ok := <-result.Result.(query.CropInventory)
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

func (s CropService) MoveToArea(crop *domain.Crop, sourceAreaUID uuid.UUID, destinationAreaUID uuid.UUID, quantity int) error {
	// Validate //
	// Check if source area is exist in DB
	result := s.AreaQuery.FindByID(sourceAreaUID)
	srcArea, ok := result.Result.(query.CropArea)
	if !ok {
		return CropError{Code: CropAreaErrorInvalidSourceArea}
	}

	if srcArea == (query.CropArea{}) {
		return CropError{Code: CropAreaErrorSourceAreaNotFound}
	}

	// Check if destination area is exist in DB
	result = s.AreaQuery.FindByID(destinationAreaUID)
	dstArea, ok  result.Result.(query.CropArea)
	if !ok {
		return Croperror{Code: CropAreaErrorInvalidDestinationArea}
	}

	dstArea, ok = result.Result.(query.CropArea)
	if !ok {
		return CropError{Code: CropAreaErrorInvalidDestinationArea}
	}

	if dstArea == (query.CropArea{}) {
		return CropError{Code: CropAreaErrorDestinationAreaNotFound}
	}

	// Check if movement rules for area type is valid
	isValidMoveRules = false
	if srcArea.Type == "seeding" && dstArea.Type == "growing" {
		isValidMoveRules = true
	} else if srcArea.Type == "seeding" && dstArea.Type == "seeding" {
		isValidMoveRules = true
	} else if srcArea.Type == "growing" && dstArea.Type == "growing" {
		isValidMoveRules = true
	}

	if !isValidMoveRules {
		return CropError{Code: CropMoveToAreaErrorInvalidArea}
	}

	// source and destination area cannot be the same
	if srcArea.UID == dstArea.UID {
		return CropError{Code: CropMoveToAreaErrorCannotBeSame}
	}

	// Quantity to be moved cannot be empty
	if quantity <= 0 {
		return CropError{Code: CropMoveToAreaErrorInvalidQuantity}
	}

	// Check validity of the source area input and the quantity to the existing crop source area.
	isValidSrcArea = false
	isValidQuantity = false
	if crop.InitialArea.AreaUID == srcArea.UID {
		isValidSrcArea = true
		isValidQuantity = (crop.InitialArea.CurrentQuantity - quantity) >= 0
	}

	for i, v := range crop.MovedArea {
		if v.AreaUID == srcArea.UID {
			isValidSrcArea = true
			isValidQuantity = (v.CurrentQuantity - quantity) >= 0
		}
	}

	if !isValidSrcArea {
		return CropError{Code: CropMoveToAreaErrorInvalidExistingArea}
	}
	if !isValidQuantity {
		return CropError{Code: CropMoveToAreaErrorInvalidQuantity}
	}

	// Check existance of the destination area input to the existing crop destination area.
	isDstExist = false
	existingDst := MovedArea{}
	for i, v := range crop.MovedArea {
		if v.AreaUID == dstArea.UID {
			isDstExist = true
			existingDst = v
		}
	}

	// Process //
	if crop.InitialArea.AreaUID == srcArea.UID {
		crop.InitialArea.CurrentQuantity -= quantity
	}

	for i, v := range crop.MovedArea {
		if v.AreaUID == srcArea.UID {
			crop.MovedArea[i].CurrentQuantity -= quantity
		}
	}

	if isDstExist {
		for i, v := range crop.MovedArea {
			if v.AreaUID == dstArea.UID {
				crop.MovedArea[i].CurrentQuantity += quantity
			}
		}
	} else {
		crop.MovedArea = append(crop.MovedArea, MovedArea{
			AreaUID: dstArea.UID,
			SourceAreaUID: srcArea.UID,
			InitialQuantity: quantity,
			CurrentQuantity: quantity,
			Date: time.Now(),
		})
	}

	return nil
}

func (c *Crop) Harvest(sourceAreaUID uuid.UUID, quantity int) error {
	return nil
}

func (c *Crop) Dump(sourceAreaUID uuid.UUID, quantity int) error {
	return nil
}
