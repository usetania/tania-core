package domain

import (
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/growth/query"
	"github.com/usetania/tania-core/src/helper/stringhelper"
)

type Crop struct {
	UID          uuid.UUID
	BatchID      string
	Status       CropStatus
	Type         CropType
	Container    CropContainer
	InventoryUID uuid.UUID
	FarmUID      uuid.UUID
	Photos       []CropPhoto

	// Fields to track crop's movement
	InitialArea      InitialArea
	MovedArea        []MovedArea
	HarvestedStorage []HarvestedStorage
	Trash            []Trash

	// Fields to track care crop
	LastFertilized time.Time
	LastPruned     time.Time
	LastPesticided time.Time

	// Notes
	Notes map[uuid.UUID]CropNote

	// Events
	Version            int
	UncommittedChanges []interface{}
}

// CropService handles crop behaviours that needs external interaction to be worked.
type CropService interface {
	FindMaterialByID(uid uuid.UUID) ServiceResult
	FindByBatchID(batchID string) ServiceResult
	FindAreaByID(uid uuid.UUID) ServiceResult
}

// ServiceResult is the container for service result.
type ServiceResult struct {
	Result interface{}
	Error  error
}

// We just recorded two crop status
// because other activities are recurring and parallel.
// So we can't have crop batch with a status, for example `HARVESTED`,
// because not all the crop has harvested.
const (
	CropActive   = "ACTIVE"
	CropArchived = "ARCHIVED"
)

type CropStatus struct {
	Code  string `json:"code"`
	Label string `json:"-"`
}

func CropStatuses() []CropStatus {
	return []CropStatus{
		{Code: CropActive, Label: "Active"},
		{Code: CropArchived, Label: "Archived"},
	}
}

func GetCropStatus(code string) CropStatus {
	for _, v := range CropStatuses() {
		if code == v.Code {
			return v
		}
	}

	return CropStatus{}
}

const (
	CropTypeSeeding = "SEEDING"
	CropTypeGrowing = "GROWING"
)

type CropType struct {
	Code  string
	Label string
}

func CropTypes() []CropType {
	return []CropType{
		{Code: CropTypeSeeding, Label: "Seeding"},
		{Code: CropTypeGrowing, Label: "Growing"},
	}
}

func GetCropType(code string) CropType {
	for _, v := range CropTypes() {
		if code == v.Code {
			return v
		}
	}

	return CropType{}
}

// CropContainerType defines the type of a container.
type CropContainerType interface {
	Code() string
}

// Tray implements CropContainerType
type Tray struct {
	Cell int
}

func (Tray) Code() string { return "TRAY" }

// Pot implements CropContainerType
type Pot struct{}

func (Pot) Code() string { return "POT" }

type CropContainer struct {
	Quantity int
	Type     CropContainerType
}

type InitialArea struct {
	AreaUID         uuid.UUID `json:"area_id"`
	InitialQuantity int       `json:"initial_quantity"`
	CurrentQuantity int       `json:"current_quantity"`
	CreatedDate     time.Time `json:"created_date"`
	LastUpdated     time.Time `json:"last_updated"`

	LastWatered    time.Time `json:"last_watered"`
	LastFertilized time.Time `json:"last_fertilized"`
	LastPruned     time.Time `json:"last_pruned"`
	LastPesticided time.Time `json:"last_pesticided"`
}

type MovedArea struct {
	AreaUID         uuid.UUID `json:"area_id"`
	SourceAreaUID   uuid.UUID `json:"source_area_id"`
	InitialQuantity int       `json:"initial_quantity"`
	CurrentQuantity int       `json:"current_quantity"`
	CreatedDate     time.Time `json:"created_date"`
	LastUpdated     time.Time `json:"last_updated"`

	LastWatered    time.Time `json:"last_watered"`
	LastFertilized time.Time `json:"last_fertilized"`
	LastPruned     time.Time `json:"last_pruned"`
	LastPesticided time.Time `json:"last_pesticided"`
}

type HarvestedStorage struct {
	Quantity             int       `json:"quantity"`
	ProducedGramQuantity float32   `json:"produced_gram_quantity"`
	SourceAreaUID        uuid.UUID `json:"source_area_id"`
	CreatedDate          time.Time `json:"created_date"`
	LastUpdated          time.Time `json:"last_updated"`
}

type Trash struct {
	Quantity      int       `json:"quantity"`
	SourceAreaUID uuid.UUID `json:"source_area_id"`
	CreatedDate   time.Time `json:"created_date"`
	LastUpdated   time.Time `json:"last_updated"`
}

const (
	HarvestTypeAll     = "ALL"
	HarvestTypePartial = "PARTIAL"
)

type HarvestType struct {
	Code  string
	Label string
}

func HarvestTypes() []HarvestType {
	return []HarvestType{
		{Code: HarvestTypeAll, Label: "All"},
		{Code: HarvestTypePartial, Label: "Partial"},
	}
}

func GetHarvestType(code string) HarvestType {
	for _, v := range HarvestTypes() {
		if v.Code == code {
			return v
		}
	}

	return HarvestType{}
}

const (
	Kg = "Kg"
	Gr = "Gr"
)

type ProducedUnit struct {
	Code  string
	Label string
}

func ProducedUnits() []ProducedUnit {
	return []ProducedUnit{
		{Code: Kg, Label: "kg"},
		{Code: Gr, Label: "gr"},
	}
}

func GetProducedUnit(code string) ProducedUnit {
	for _, v := range ProducedUnits() {
		if code == v.Code {
			return v
		}
	}

	return ProducedUnit{}
}

type CropNote struct {
	UID         uuid.UUID `json:"uid"`
	Content     string    `json:"content"`
	CreatedDate time.Time `json:"created_date"`
}

type CropPhoto struct {
	UID         uuid.UUID `json:"uid"`
	Filename    string    `json:"filename"`
	MimeType    string    `json:"mime_type"`
	Size        int       `json:"size"`
	Width       int       `json:"width"`
	Height      int       `json:"height"`
	Description string    `json:"description"`
}

func (c *Crop) TrackChange(event interface{}) {
	c.UncommittedChanges = append(c.UncommittedChanges, event)
	c.Transition(event)
}

func (c *Crop) Transition(event interface{}) {
	switch e := event.(type) {
	case CropBatchCreated:
		c.UID = e.UID
		c.BatchID = e.BatchID
		c.Status = e.Status
		c.Type = e.Type
		c.Container = e.Container
		c.InventoryUID = e.InventoryUID
		c.InitialArea = InitialArea{
			AreaUID:         e.InitialAreaUID,
			InitialQuantity: e.Quantity,
			CurrentQuantity: e.Quantity,
			CreatedDate:     e.CreatedDate,
			LastUpdated:     e.CreatedDate,
		}
		c.FarmUID = e.FarmUID

	case CropBatchInventoryChanged:
		c.InventoryUID = e.InventoryUID
		c.BatchID = e.BatchID

	case CropBatchTypeChanged:
		c.Type = e.Type

	case CropBatchContainerChanged:
		c.Container = e.Container
		c.InitialArea.CurrentQuantity = e.Container.Quantity
		c.InitialArea.InitialQuantity = e.Container.Quantity

	case CropBatchMoved:
		if c.InitialArea.AreaUID == e.SrcAreaUID {
			ia, ok := e.UpdatedSrcArea.(InitialArea)
			if ok {
				c.InitialArea = ia
			}
		}

		for i, v := range c.MovedArea {
			ma, ok := e.UpdatedSrcArea.(MovedArea)
			if ok {
				if v.AreaUID == ma.AreaUID {
					c.MovedArea[i] = ma
				}
			}
		}

		if c.InitialArea.AreaUID == e.DstAreaUID {
			ia, ok := e.UpdatedDstArea.(InitialArea)
			if ok {
				c.InitialArea = ia
			}
		}

		isFound := false

		for i, v := range c.MovedArea {
			ma, ok := e.UpdatedDstArea.(MovedArea)

			if ok {
				if v.AreaUID == ma.AreaUID {
					c.MovedArea[i] = ma
					isFound = true
				}
			}
		}

		if !isFound {
			ma, ok := e.UpdatedDstArea.(MovedArea)
			if ok {
				c.MovedArea = append(c.MovedArea, ma)
			}
		}

	case CropBatchHarvested:
		isFound := false

		for i, v := range c.HarvestedStorage {
			if v.SourceAreaUID == e.UpdatedHarvestedStorage.SourceAreaUID {
				c.HarvestedStorage[i] = e.UpdatedHarvestedStorage
				isFound = true
			}
		}

		if !isFound {
			c.HarvestedStorage = append(c.HarvestedStorage, e.UpdatedHarvestedStorage)
		}

		if e.HarvestedAreaCode == "INITIAL_AREA" {
			ha := e.HarvestedArea.(InitialArea)
			c.InitialArea = ha
		} else if e.HarvestedAreaCode == "MOVED_AREA" {
			ma := e.HarvestedArea.(MovedArea)

			for i, v := range c.MovedArea {
				if v.AreaUID == ma.AreaUID {
					c.MovedArea[i] = ma
				}
			}
		}

		c.Status = GetCropStatus(e.CropStatus)

	case CropBatchDumped:
		isFound := false

		for i, v := range c.Trash {
			if v.SourceAreaUID == e.UpdatedTrash.SourceAreaUID {
				c.Trash[i] = e.UpdatedTrash
				isFound = true
			}
		}

		if !isFound {
			c.Trash = append(c.Trash, e.UpdatedTrash)
		}

		if e.DumpedAreaCode == "INITIAL_AREA" {
			da := e.DumpedArea.(InitialArea)
			c.InitialArea = da
		} else if e.DumpedAreaCode == "MOVED_AREA" {
			da := e.DumpedArea.(MovedArea)

			for i, v := range c.MovedArea {
				if v.AreaUID == da.AreaUID {
					c.MovedArea[i] = da
				}
			}
		}

		c.Status = GetCropStatus(e.CropStatus)

	case CropBatchWatered:
		if c.InitialArea.AreaUID == e.AreaUID {
			c.InitialArea.LastWatered = e.WateringDate
		}

		for i, v := range c.MovedArea {
			if v.AreaUID == e.AreaUID {
				c.MovedArea[i].LastWatered = e.WateringDate
			}
		}

	case CropBatchNoteCreated:
		if len(c.Notes) == 0 {
			c.Notes = make(map[uuid.UUID]CropNote)
		}

		c.Notes[e.UID] = CropNote{
			UID:         e.UID,
			Content:     e.Content,
			CreatedDate: e.CreatedDate,
		}

	case CropBatchNoteRemoved:
		delete(c.Notes, e.UID)

	case CropBatchPhotoCreated:
		c.Photos = append(c.Photos, CropPhoto{
			UID:         e.UID,
			Filename:    e.Filename,
			MimeType:    e.MimeType,
			Size:        e.Size,
			Width:       e.Width,
			Height:      e.Height,
			Description: e.Description,
		})
	}
}

func CreateCropBatch(
	cropService CropService,
	areaUID uuid.UUID,
	cropType string,
	inventoryUID uuid.UUID,
	quantity int, containerType CropContainerType) (*Crop, error,
) {
	serviceResult := cropService.FindAreaByID(areaUID)
	if serviceResult.Error != nil {
		return nil, serviceResult.Error
	}

	area := serviceResult.Result.(query.CropAreaQueryResult)

	ct := GetCropType(cropType)
	if ct == (CropType{}) {
		return nil, CropError{Code: CropErrorInvalidCropType}
	}

	serviceResult = cropService.FindMaterialByID(inventoryUID)
	if serviceResult.Error != nil {
		return nil, serviceResult.Error
	}

	inv := serviceResult.Result.(query.CropMaterialQueryResult)

	err := validateContainer(quantity, containerType)
	if err != nil {
		return nil, err
	}

	cropContainer := CropContainer{
		Quantity: quantity,
		Type:     containerType,
	}

	createdDate := time.Now()

	batchID, err := generateBatchID(cropService, inv, createdDate)
	if err != nil {
		return nil, err
	}

	uid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	initial := &Crop{}

	initial.TrackChange(CropBatchCreated{
		UID:            uid,
		BatchID:        batchID,
		Status:         GetCropStatus(CropActive),
		Type:           ct,
		Container:      cropContainer,
		InventoryUID:   inv.UID,
		CreatedDate:    createdDate,
		InitialAreaUID: area.UID,
		Quantity:       quantity,
		FarmUID:        area.FarmUID,
	})

	return initial, nil
}

func (c *Crop) MoveToArea(cropService CropService, sourceAreaUID, destinationAreaUID uuid.UUID, quantity int) error {
	// Validate //
	// Check if source area is exist in DB
	serviceResult := cropService.FindAreaByID(sourceAreaUID)
	if serviceResult.Error != nil {
		return serviceResult.Error
	}

	srcArea, ok := serviceResult.Result.(query.CropAreaQueryResult)
	if !ok {
		return CropError{Code: CropMoveToAreaErrorInvalidSourceArea}
	}

	if srcArea.UID == (uuid.UUID{}) {
		return CropError{Code: CropMoveToAreaErrorSourceAreaNotFound}
	}

	// Check if destination area is exist in DB
	serviceResult = cropService.FindAreaByID(destinationAreaUID)
	if serviceResult.Error != nil {
		return serviceResult.Error
	}

	dstArea, ok := serviceResult.Result.(query.CropAreaQueryResult)
	if !ok {
		return CropError{Code: CropMoveToAreaErrorInvalidDestinationArea}
	}

	if dstArea.UID == (uuid.UUID{}) {
		return CropError{Code: CropMoveToAreaErrorDestinationAreaNotFound}
	}

	// Check if movement rules for area type is valid
	isValidMoveRules := false
	if (srcArea.Type == "SEEDING" && dstArea.Type == "GROWING") ||
		(srcArea.Type == "SEEDING" && dstArea.Type == "SEEDING") ||
		(srcArea.Type == "GROWING" && dstArea.Type == "GROWING") {
		isValidMoveRules = true
	}

	if !isValidMoveRules {
		return CropError{Code: CropMoveToAreaErrorInvalidAreaRules}
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
	isValidSrcArea := false
	isValidQuantity := false

	if c.InitialArea.AreaUID == srcArea.UID {
		isValidSrcArea = true
		isValidQuantity = (c.InitialArea.CurrentQuantity - quantity) >= 0
	}

	for _, v := range c.MovedArea {
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

	// Process //
	movedDate := time.Now()

	var updatedSrcArea interface{}

	var updatedSrcAreaCode string

	if c.InitialArea.AreaUID == srcArea.UID {
		ia := c.InitialArea
		ia.CurrentQuantity -= quantity

		updatedSrcAreaCode = "INITIAL_AREA"
		updatedSrcArea = ia
	}

	for _, v := range c.MovedArea {
		if v.AreaUID == srcArea.UID {
			ma := v
			ma.CurrentQuantity -= quantity

			updatedSrcAreaCode = "MOVED_AREA"
			updatedSrcArea = ma
		}
	}

	var updatedDstArea interface{}

	updatedDstAreaCode := ""
	isDstFoundInInitial := false
	isDstFoundInMoved := false

	if c.InitialArea.AreaUID == dstArea.UID {
		ia := c.InitialArea
		ia.CurrentQuantity += quantity
		ia.LastUpdated = movedDate

		updatedDstArea = ia
		updatedDstAreaCode = "INITIAL_AREA"
		isDstFoundInInitial = true
	}

	// If Destination not found in the Initial Area,
	// then we will look at it next in the Moved Area
	if !isDstFoundInInitial {
		for _, v := range c.MovedArea {
			if v.AreaUID == dstArea.UID {
				da := v
				da.CurrentQuantity += quantity
				da.LastUpdated = movedDate

				updatedDstArea = da
				updatedDstAreaCode = "MOVED_AREA"
				isDstFoundInMoved = true
			}
		}
	}

	if !isDstFoundInInitial && !isDstFoundInMoved {
		updatedDstArea = MovedArea{
			AreaUID:         dstArea.UID,
			SourceAreaUID:   srcArea.UID,
			InitialQuantity: quantity,
			CurrentQuantity: quantity,
			CreatedDate:     movedDate,
			LastUpdated:     movedDate,
		}
		updatedDstAreaCode = "MOVED_AREA"
	}

	// Process //
	c.TrackChange(CropBatchMoved{
		UID:                c.UID,
		Quantity:           quantity,
		SrcAreaUID:         srcArea.UID,
		DstAreaUID:         dstArea.UID,
		MovedDate:          movedDate,
		UpdatedSrcArea:     updatedSrcArea,
		UpdatedSrcAreaCode: updatedSrcAreaCode,
		UpdatedDstArea:     updatedDstArea,
		UpdatedDstAreaCode: updatedDstAreaCode,
	})

	return nil
}

func (c *Crop) Harvest(
	cropService CropService,
	sourceAreaUID uuid.UUID,
	harvestType string,
	producedQuantity float32,
	producedUnit ProducedUnit,
	notes string,
) error {
	// Validate //
	// Check if source area is exist in DB
	serviceResult := cropService.FindAreaByID(sourceAreaUID)
	if serviceResult.Error != nil {
		return serviceResult.Error
	}

	srcArea, ok := serviceResult.Result.(query.CropAreaQueryResult)
	if !ok {
		return CropError{Code: CropHarvestErrorInvalidSourceArea}
	}

	if srcArea == (query.CropAreaQueryResult{}) {
		return CropError{Code: CropHarvestErrorSourceAreaNotFound}
	}

	// Check if area is already set in the crop
	isAreaValid := false
	if c.InitialArea.AreaUID == sourceAreaUID {
		isAreaValid = true
	}

	for _, v := range c.MovedArea {
		if v.AreaUID == sourceAreaUID {
			isAreaValid = true
		}
	}

	if !isAreaValid {
		return CropError{Code: CropHarvestErrorSourceAreaNotFound}
	}

	if srcArea.Type != "GROWING" {
		return CropError{Code: CropHarvestErrorInvalidSourceArea}
	}

	ht := GetHarvestType(harvestType)
	if ht == (HarvestType{}) {
		return CropError{Code: CropHarvestErrorInvalidHarvestType}
	}

	// Process //
	harvestDate := time.Now()

	// If harvestType All, then empty the quantity in the area because it has been all harvested
	// Else if harvestType Partial, then we assume that the quantity of moved plant is 0
	harvestedQuantity := 0

	var harvestedArea interface{}

	harvestedAreaCode := ""

	if ht.Code == HarvestTypeAll {
		if c.InitialArea.AreaUID == srcArea.UID {
			ia := c.InitialArea

			if ia.CurrentQuantity <= 0 {
				return CropError{Code: CropHarvestErrorNotEnoughQuantity}
			}

			harvestedQuantity = ia.CurrentQuantity
			ia.CurrentQuantity = 0

			harvestedArea = ia
			harvestedAreaCode = "INITIAL_AREA"
		}

		for _, v := range c.MovedArea {
			if v.AreaUID == srcArea.UID {
				ma := v

				if ma.CurrentQuantity <= 0 {
					return CropError{Code: CropHarvestErrorNotEnoughQuantity}
				}

				harvestedQuantity = ma.CurrentQuantity
				ma.CurrentQuantity = 0

				harvestedArea = ma
				harvestedAreaCode = "MOVED_AREA"
			}
		}
	}

	// Check source area existence. If already exist, then just update it
	harvestedStorage := HarvestedStorage{}
	isExist := false

	for _, v := range c.HarvestedStorage {
		if v.SourceAreaUID == srcArea.UID {
			harvestedStorage = v
			harvestedStorage.Quantity += harvestedQuantity
			harvestedStorage.LastUpdated = harvestDate

			isExist = true
		}
	}

	if !isExist {
		harvestedStorage.Quantity = harvestedQuantity
		harvestedStorage.SourceAreaUID = srcArea.UID
		harvestedStorage.CreatedDate = harvestDate
		harvestedStorage.LastUpdated = harvestDate
	}

	// Calculate the produced harvest
	// Produced Quantity always converted to gram
	totalProduced := producedQuantity
	if producedUnit.Code == Kg {
		totalProduced = producedQuantity * 1000
	}

	harvestedStorage.ProducedGramQuantity += totalProduced

	// Check all the quantity in InitialArea and MovedArea,
	// if its all empty then crop status is marked to archive
	initialAreaEmpty := false
	if harvestedAreaCode == "INITIAL_AREA" {
		initialAreaEmpty = true
	} else if c.InitialArea.CurrentQuantity == 0 {
		initialAreaEmpty = true
	}

	movedAreaEmpty := true

	if harvestedAreaCode == "MOVED_AREA" {
		for _, v := range c.MovedArea {
			ha := harvestedArea.(MovedArea)

			if v.AreaUID == ha.AreaUID {
				ha.CurrentQuantity = 0
			}

			if ha.CurrentQuantity != 0 {
				movedAreaEmpty = false
			}
		}
	} else {
		for _, v := range c.MovedArea {
			if v.CurrentQuantity != 0 {
				movedAreaEmpty = false
			}
		}
	}

	status := CropActive
	if initialAreaEmpty && movedAreaEmpty {
		status = CropArchived
	}

	// Process //
	c.TrackChange(CropBatchHarvested{
		UID:                     c.UID,
		CropStatus:              status,
		HarvestType:             ht.Code,
		HarvestedQuantity:       harvestedQuantity,
		ProducedGramQuantity:    totalProduced,
		UpdatedHarvestedStorage: harvestedStorage,
		HarvestedArea:           harvestedArea,
		HarvestedAreaCode:       harvestedAreaCode,
		HarvestDate:             harvestDate,
		Notes:                   notes,
	})

	return nil
}

func (c *Crop) Dump(cropService CropService, sourceAreaUID uuid.UUID, quantity int, notes string) error {
	// Validate //
	// Check if source area is exist in DB
	serviceResult := cropService.FindAreaByID(sourceAreaUID)
	if serviceResult.Error != nil {
		return serviceResult.Error
	}

	srcArea, ok := serviceResult.Result.(query.CropAreaQueryResult)
	if !ok {
		return CropError{Code: CropDumpErrorInvalidSourceArea}
	}

	if srcArea == (query.CropAreaQueryResult{}) {
		return CropError{Code: CropDumpErrorSourceAreaNotFound}
	}

	// Check if area is already set in the crop
	isAreaValid := false
	isQuantityValid := true

	if c.InitialArea.AreaUID == sourceAreaUID {
		isAreaValid = true

		if (c.InitialArea.CurrentQuantity - quantity) < 0 {
			isQuantityValid = false
		}
	}

	for _, v := range c.MovedArea {
		if v.AreaUID == sourceAreaUID {
			isAreaValid = true

			if (v.CurrentQuantity - quantity) < 0 {
				isQuantityValid = false
			}
		}
	}

	if !isAreaValid {
		return CropError{Code: CropDumpErrorSourceAreaNotFound}
	}

	if !isQuantityValid {
		return CropError{Code: CropDumpErrorNotEnoughQuantity}
	}

	if quantity <= 0 {
		return CropError{Code: CropDumpErrorInvalidQuantity}
	}

	// Check source area existence. If already exist, then just update it
	var dumpedArea interface{}

	updatedTrash := Trash{}
	dumpDate := time.Now()

	isExist := false

	for i, v := range c.Trash {
		if v.SourceAreaUID == srcArea.UID {
			updatedTrash = v
			updatedTrash.Quantity = c.Trash[i].Quantity + quantity
			isExist = true
		}
	}

	if !isExist {
		updatedTrash.Quantity = quantity
		updatedTrash.SourceAreaUID = srcArea.UID
		updatedTrash.CreatedDate = dumpDate
		updatedTrash.LastUpdated = dumpDate
	}

	// Reduce the quantity in the area because it has been dumped
	dumpedAreaCode := ""

	if c.InitialArea.AreaUID == srcArea.UID {
		ia := c.InitialArea
		ia.CurrentQuantity -= quantity
		ia.LastUpdated = dumpDate

		dumpedArea = ia
		dumpedAreaCode = "INITIAL_AREA"
	}

	for _, v := range c.MovedArea {
		if v.AreaUID == srcArea.UID {
			ma := v
			ma.CurrentQuantity -= quantity
			ma.LastUpdated = dumpDate

			dumpedArea = ma
			dumpedAreaCode = "MOVED_AREA"
		}
	}

	// Check all the quantity in InitialArea and MovedArea,
	// if its all empty then crop status is marked to archive
	initialAreaEmpty := false

	if dumpedAreaCode == "INITIAL_AREA" {
		ia := dumpedArea.(InitialArea)

		if ia.CurrentQuantity == 0 {
			initialAreaEmpty = true
		}
	} else if c.InitialArea.CurrentQuantity == 0 {
		initialAreaEmpty = true
	}

	movedAreaEmpty := true

	if dumpedAreaCode == "MOVED_AREA" {
		for _, v := range c.MovedArea {
			ha := dumpedArea.(MovedArea)

			if v.AreaUID == ha.AreaUID {
				ha.CurrentQuantity = 0
			}

			if ha.CurrentQuantity != 0 {
				movedAreaEmpty = false
			}
		}
	} else {
		for _, v := range c.MovedArea {
			if v.CurrentQuantity != 0 {
				movedAreaEmpty = false
			}
		}
	}

	status := CropActive
	if initialAreaEmpty && movedAreaEmpty {
		status = CropArchived
	}

	// Process //
	c.TrackChange(CropBatchDumped{
		UID:            c.UID,
		CropStatus:     status,
		Quantity:       quantity,
		UpdatedTrash:   updatedTrash,
		DumpedArea:     dumpedArea,
		DumpedAreaCode: dumpedAreaCode,
		DumpDate:       time.Now(),
		Notes:          notes,
	})

	return nil
}

func (*Crop) Fertilize() error {
	// c.LastFertilized = time.Now()
	return nil
}

func (*Crop) Prune() error {
	// c.LastPruned = time.Now()
	return nil
}

func (*Crop) Pesticide() error {
	// c.LastPesticided = time.Now()
	return nil
}

func (c *Crop) Water(cropService CropService, sourceAreaUID uuid.UUID, wateringDate time.Time) error {
	serviceResult := cropService.FindAreaByID(sourceAreaUID)
	if serviceResult.Error != nil {
		return serviceResult.Error
	}

	srcArea, ok := serviceResult.Result.(query.CropAreaQueryResult)
	if !ok {
		return CropError{Code: CropDumpErrorInvalidSourceArea}
	}

	if srcArea == (query.CropAreaQueryResult{}) {
		return CropError{Code: CropDumpErrorSourceAreaNotFound}
	}

	if wateringDate.IsZero() {
		return CropError{Code: CropWaterErrorInvalidWateringDate}
	}

	c.TrackChange(CropBatchWatered{
		UID:           c.UID,
		BatchID:       c.BatchID,
		ContainerType: c.Container.Type.Code(),
		AreaUID:       srcArea.UID,
		AreaName:      srcArea.Name,
		WateringDate:  wateringDate,
	})

	return nil
}

func (c *Crop) ChangeCropType(cropType string) error {
	ct := GetCropType(cropType)
	if ct == (CropType{}) {
		return CropError{Code: CropErrorInvalidCropType}
	}

	c.TrackChange(CropBatchTypeChanged{
		UID:  c.UID,
		Type: ct,
	})

	return nil
}

func (c *Crop) ChangeCropStatus(cropStatus string) error {
	cs := GetCropStatus(cropStatus)
	if cs == (CropStatus{}) {
		return CropError{Code: CropErrorInvalidCropStatus}
	}

	c.Status = cs

	return nil
}

func (c *Crop) ChangeContainer(quantity int, containerType CropContainerType) error {
	err := validateContainer(quantity, containerType)
	if err != nil {
		return err
	}

	if c.InitialArea.CurrentQuantity != c.InitialArea.InitialQuantity {
		return CropError{Code: CropContainerErrorInvalidType}
	}

	c.TrackChange(CropBatchContainerChanged{
		UID: c.UID,
		Container: CropContainer{
			Quantity: quantity,
			Type:     containerType,
		},
	})

	return nil
}

func (c *Crop) ChangeInventory(cropService CropService, inventoryUID uuid.UUID) error {
	serviceResult := cropService.FindMaterialByID(inventoryUID)

	if serviceResult.Error != nil {
		return serviceResult.Error
	}

	inventory := serviceResult.Result.(query.CropMaterialQueryResult)

	batchID, err := generateBatchID(cropService, inventory, c.InitialArea.CreatedDate)
	if err != nil {
		return err
	}

	c.TrackChange(CropBatchInventoryChanged{
		UID:          c.UID,
		InventoryUID: inventory.UID,
		BatchID:      batchID,
	})

	return nil
}

func (c *Crop) AddNewNote(content string) error {
	if content == "" {
		return CropError{Code: CropNoteErrorInvalidContent}
	}

	uid, err := uuid.NewV4()
	if err != nil {
		return err
	}

	c.TrackChange(CropBatchNoteCreated{
		UID:         uid,
		CropUID:     c.UID,
		Content:     content,
		CreatedDate: time.Now(),
	})

	return nil
}

func (c *Crop) RemoveNote(uid uuid.UUID) error {
	if uid == (uuid.UUID{}) {
		return CropError{Code: CropNoteErrorNotFound}
	}

	found := CropNote{}

	for _, v := range c.Notes {
		if v.UID == uid {
			found = v
		}
	}

	if found == (CropNote{}) {
		return CropError{Code: CropNoteErrorNotFound}
	}

	c.TrackChange(CropBatchNoteRemoved{
		UID:         found.UID,
		CropUID:     c.UID,
		Content:     found.Content,
		CreatedDate: found.CreatedDate,
	})

	return nil
}

func (c *Crop) AddPhoto(filename, mimeType string, size, width, height int, description string) error {
	if filename == "" {
		return CropError{CropErrorPhotoInvalidFilename}
	}

	if mimeType == "" {
		return CropError{CropErrorPhotoInvalidMimeType}
	}

	if size <= 0 {
		return CropError{CropErrorPhotoInvalidSize}
	}

	if description == "" {
		return CropError{CropErrorPhotoInvalidDescription}
	}

	uid, err := uuid.NewV4()
	if err != nil {
		return err
	}

	c.TrackChange(CropBatchPhotoCreated{
		UID:         uid,
		CropUID:     c.UID,
		Filename:    filename,
		MimeType:    mimeType,
		Size:        size,
		Width:       width,
		Height:      height,
		Description: description,
	})

	return nil
}

// CalculateDaysSinceSeeding will find how long since its been seeded
// It basically tell use the days since this crop is created.
func (c Crop) CalculateDaysSinceSeeding() int {
	now := time.Now()

	diff := now.Sub(c.InitialArea.CreatedDate)

	days := int(diff.Hours()) / 24

	return days
}

func generateBatchID(cs CropService, inventory query.CropMaterialQueryResult, createdDate time.Time) (string, error) {
	// Generate Batch ID
	// Format the date to become daymonth format like 25jan
	dateFormat := strings.ToLower(createdDate.Format("2Jan"))

	// Get variety name and split it to slice
	varietySlice := strings.Fields(inventory.Name)
	varietyFormat := ""

	for _, v := range varietySlice {
		// For every value, get only the first three characters
		var format string

		if len(v) > 3 {
			format = strings.ToLower(v[0:3])
		} else {
			format = strings.ToLower(v)
		}

		varietyFormat = stringhelper.Join(varietyFormat, format, "-")
	}

	// Join that variety and date
	batchID := stringhelper.Join(varietyFormat, dateFormat)

	// Validate Uniqueness of Batch ID.
	serviceResult := cs.FindByBatchID(batchID)
	if serviceResult.Error != nil {
		return "", serviceResult.Error
	}

	return batchID, nil
}

func validateContainer(quantity int, containerType CropContainerType) error {
	if quantity <= 0 {
		return CropError{Code: CropContainerErrorInvalidQuantity}
	}

	switch v := containerType.(type) {
	case Tray:
		if v.Cell <= 0 {
			return CropError{Code: CropContainerErrorInvalidTrayCell}
		}
	case Pot:
	default:
		return CropError{Code: CropContainerErrorInvalidType}
	}

	return nil
}
