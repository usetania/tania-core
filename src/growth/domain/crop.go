package domain

import (
	"strings"
	"time"

	"github.com/Tanibox/tania-server/src/growth/query"
	"github.com/Tanibox/tania-server/src/helper/stringhelper"
	uuid "github.com/satori/go.uuid"
)

type Crop struct {
	UID          uuid.UUID
	BatchID      string
	Status       CropStatus
	Type         CropType
	Container    CropContainer
	InventoryUID uuid.UUID
	FarmUID      uuid.UUID
	Photo        CropPhoto

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

// CropService handles crop behaviours that needs external interaction to be worked
type CropService interface {
	FindMaterialByID(uid uuid.UUID) ServiceResult
	FindByBatchID(batchID string) ServiceResult
	FindAreaByID(uid uuid.UUID) ServiceResult
}

// ServiceResult is the container for service result
type ServiceResult struct {
	Result interface{}
	Error  error
}

// We just recorded two crop status
// because other activities are recurring and parallel.
// So we can't have crop batch with a status, for example `HARVESTED`,
// because not all the crop has harvested
const (
	CropActive = "ACTIVE"
	CropEnd    = "END"
)

type CropStatus struct {
	Code  string `json:"code"`
	Label string `json:"-"`
}

func CropStatuses() []CropStatus {
	return []CropStatus{
		{Code: CropActive, Label: "Active"},
		{Code: CropEnd, Label: "End"},
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

// CropContainerType defines the type of a container
type CropContainerType interface {
	Code() string
}

// Tray implements CropContainerType
type Tray struct {
	Cell int
}

func (t Tray) Code() string { return "TRAY" }

// Pot implements CropContainerType
type Pot struct{}

func (p Pot) Code() string { return "POT" }

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
	AreaUID         uuid.UUID
	SourceAreaUID   uuid.UUID
	InitialQuantity int
	CurrentQuantity int
	CreatedDate     time.Time
	LastUpdated     time.Time

	LastWatered    time.Time
	LastFertilized time.Time
	LastPruned     time.Time
	LastPesticided time.Time
}

type HarvestedStorage struct {
	Quantity             int
	ProducedGramQuantity float32
	SourceAreaUID        uuid.UUID
	CreatedDate          time.Time
	LastUpdated          time.Time
}

type Trash struct {
	Quantity      int
	SourceAreaUID uuid.UUID
	CreatedDate   time.Time
	LastUpdated   time.Time
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
	Filename string `json:"filename"`
	MimeType string `json:"mime_type"`
	Size     int    `json:"size"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
}

func (state *Crop) TrackChange(event interface{}) {
	state.UncommittedChanges = append(state.UncommittedChanges, event)
	state.Transition(event)
}

func (state *Crop) Transition(event interface{}) {
	switch e := event.(type) {
	case CropBatchCreated:
		state.UID = e.UID
		state.BatchID = e.BatchID
		state.Status = e.Status
		state.Type = e.Type
		state.Container = e.Container
		state.InventoryUID = e.InventoryUID
		state.InitialArea = InitialArea{
			AreaUID:         e.InitialAreaUID,
			InitialQuantity: e.Quantity,
			CurrentQuantity: e.Quantity,
			CreatedDate:     e.CreatedDate,
			LastUpdated:     e.CreatedDate,
		}
		state.FarmUID = e.FarmUID
	case CropBatchMoved:
		if state.InitialArea.AreaUID == e.SrcAreaUID {
			state.InitialArea.CurrentQuantity -= e.Quantity
		}

		for i, v := range state.MovedArea {
			if v.AreaUID == e.SrcAreaUID {
				state.MovedArea[i].CurrentQuantity -= e.Quantity
			}
		}

		isDstExist := false
		for _, v := range state.MovedArea {
			if v.AreaUID == e.DstAreaUID {
				isDstExist = true
			}
		}

		if isDstExist {
			for i, v := range state.MovedArea {
				if v.AreaUID == e.DstAreaUID {
					state.MovedArea[i].CurrentQuantity += e.Quantity
					state.MovedArea[i].LastUpdated = e.MovedDate
				}
			}
		} else {
			state.MovedArea = append(state.MovedArea, MovedArea{
				AreaUID:         e.DstAreaUID,
				SourceAreaUID:   e.SrcAreaUID,
				InitialQuantity: e.Quantity,
				CurrentQuantity: e.Quantity,
				CreatedDate:     e.MovedDate,
				LastUpdated:     e.MovedDate,
			})
		}
	case CropBatchWatered:
		if state.InitialArea.AreaUID == e.AreaUID {
			state.InitialArea.LastWatered = e.WateringDate
		}

		for i, v := range state.MovedArea {
			if v.AreaUID == e.AreaUID {
				state.MovedArea[i].LastWatered = e.WateringDate
			}
		}
	}
}

func CreateCropBatch(
	cropService CropService,
	areaUID uuid.UUID,
	cropType string,
	inventoryUID uuid.UUID,
	quantity int, containerType CropContainerType) (*Crop, error) {

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

	containerCell := 0
	switch v := containerType.(type) {
	case Tray:
		containerCell = v.Cell
	case Pot:

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
		UID:             uid,
		BatchID:         batchID,
		Status:          GetCropStatus(CropActive),
		Type:            ct,
		Container:       cropContainer,
		ContainerType:   containerType.Code(),
		ContainerCell:   containerCell,
		InventoryUID:    inv.UID,
		VarietyName:     inv.Name,
		PlantType:       inv.MaterialSeedPlantTypeCode,
		CreatedDate:     createdDate,
		InitialAreaUID:  area.UID,
		InitialAreaName: area.Name,
		Quantity:        quantity,
		FarmUID:         area.FarmUID,
	})

	return initial, nil
}

func (c *Crop) MoveToArea(cropService CropService, sourceAreaUID uuid.UUID, destinationAreaUID uuid.UUID, quantity int) error {
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
	if srcArea.Type == "SEEDING" && dstArea.Type == "GROWING" {
		isValidMoveRules = true
	} else if srcArea.Type == "SEEDING" && dstArea.Type == "SEEDING" {
		isValidMoveRules = true
	} else if srcArea.Type == "GROWING" && dstArea.Type == "GROWING" {
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
	c.TrackChange(CropBatchMoved{
		UID:           c.UID,
		BatchID:       c.BatchID,
		ContainerType: c.Container.Type.Code(),
		Quantity:      quantity,
		SrcAreaUID:    srcArea.UID,
		SrcAreaName:   srcArea.Name,
		SrcAreaType:   srcArea.Type,
		DstAreaUID:    dstArea.UID,
		DstAreaName:   dstArea.Name,
		DstAreaType:   dstArea.Type,
		MovedDate:     time.Now(),
	})

	return nil
}

func (c *Crop) Harvest(
	cropService CropService,
	sourceAreaUID uuid.UUID,
	harvestType string,
	producedQuantity float32,
	producedUnit ProducedUnit) error {

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

	// Produced Quantity always converted to gram
	totalProduced := producedQuantity
	if producedUnit.Code == Kg {
		totalProduced = producedQuantity * 1000
	}

	// Process //
	// If harvestType All, then empty the quantity in the area because it has been all harvested
	// Else if harvestType Partial, then we assume that the quantity of moved plant is 0
	harvestedQuantity := 0
	if ht.Code == HarvestTypeAll {
		if c.InitialArea.AreaUID == sourceAreaUID {
			harvestedQuantity = c.InitialArea.CurrentQuantity
			c.InitialArea.CurrentQuantity = 0
		}
		for i, v := range c.MovedArea {
			if v.AreaUID == sourceAreaUID {
				harvestedQuantity = c.MovedArea[i].CurrentQuantity
				c.MovedArea[i].CurrentQuantity = 0
			}
		}
	}

	// Check source area existance. If already exist, then just update it
	isExist := false
	for i, v := range c.HarvestedStorage {
		if v.SourceAreaUID == sourceAreaUID {
			c.HarvestedStorage[i].Quantity += harvestedQuantity
			c.HarvestedStorage[i].LastUpdated = time.Now()
			isExist = true
		}
	}

	if !isExist {
		hs := HarvestedStorage{
			Quantity:      harvestedQuantity,
			SourceAreaUID: sourceAreaUID,
			CreatedDate:   time.Now(),
			LastUpdated:   time.Now(),
		}
		c.HarvestedStorage = append(c.HarvestedStorage, hs)
	}

	// Calculate the produced harvest
	for i, v := range c.HarvestedStorage {
		if v.SourceAreaUID == srcArea.UID {
			c.HarvestedStorage[i].ProducedGramQuantity += totalProduced
		}
	}

	return nil
}

func (c *Crop) Dump(cropService CropService, sourceAreaUID uuid.UUID, quantity int) error {
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

	if quantity <= 0 {
		return CropError{Code: CropDumpErrorInvalidQuantity}
	}

	// Process //
	// Check source area existance. If already exist, then just update it
	isExist := false
	for i, v := range c.Trash {
		if v.SourceAreaUID == srcArea.UID {
			c.Trash[i].Quantity += quantity
			c.Trash[i].LastUpdated = time.Now()
			isExist = true
		}
	}

	if !isExist {
		t := Trash{
			Quantity:      quantity,
			SourceAreaUID: srcArea.UID,
			CreatedDate:   time.Now(),
			LastUpdated:   time.Now(),
		}
		c.Trash = append(c.Trash, t)
	}

	// Reduce the quantity in the area because it has been dumped
	if c.InitialArea.AreaUID == sourceAreaUID {
		c.InitialArea.CurrentQuantity -= quantity
	}
	for i, v := range c.MovedArea {
		if v.AreaUID == sourceAreaUID {
			c.MovedArea[i].CurrentQuantity -= quantity
		}
	}

	return nil
}

func (c *Crop) Fertilize() error {
	// c.LastFertilized = time.Now()

	return nil
}

func (c *Crop) Prune() error {
	// c.LastPruned = time.Now()

	return nil
}

func (c *Crop) Pesticide() error {
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

	c.Type = ct

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

	c.Container = CropContainer{
		Quantity: quantity,
		Type:     containerType,
	}

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

	c.InventoryUID = inventory.UID
	c.BatchID = batchID

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

	cropNote := CropNote{
		UID:         uid,
		Content:     content,
		CreatedDate: time.Now(),
	}

	if len(c.Notes) == 0 {
		c.Notes = make(map[uuid.UUID]CropNote)
	}

	c.Notes[uid] = cropNote

	return nil
}

func (c *Crop) RemoveNote(uid uuid.UUID) error {
	if uid == (uuid.UUID{}) {
		return CropError{Code: CropNoteErrorNotFound}
	}

	found := false
	for _, v := range c.Notes {
		if v.UID == uid {
			delete(c.Notes, uid)
			found = true
		}
	}

	if !found {
		return CropError{Code: CropNoteErrorNotFound}
	}

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

func generateBatchID(cropService CropService, inventory query.CropMaterialQueryResult, createdDate time.Time) (string, error) {
	// Generate Batch ID
	// Format the date to become daymonth format like 25jan
	dateFormat := strings.ToLower(createdDate.Format("2Jan"))

	// Get variety name and split it to slice
	varietySlice := strings.Fields(inventory.Name)
	varietyFormat := ""
	for _, v := range varietySlice {
		// 	// For every value, get only the first three characters
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
	serviceResult := cropService.FindByBatchID(batchID)
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
