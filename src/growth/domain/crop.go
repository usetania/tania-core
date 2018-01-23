package domain

import (
	"strings"
	"time"

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
	CreatedDate  time.Time
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
}

// CropService handles crop behaviours that needs external interaction to be worked
type CropService interface {
	FindInventoryMaterialByID(uid uuid.UUID) ServiceResult
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

func (t Tray) Code() string { return "tray" }

// Pot implements CropContainerType
type Pot struct{}

func (p Pot) Code() string { return "pot" }

type CropContainer struct {
	Quantity int
	Type     CropContainerType
}

type InitialArea struct {
	AreaUID         uuid.UUID `json:"area_id"`
	InitialQuantity int       `json:"initial_quantity"`
	CurrentQuantity int       `json:"current_quantity"`
}

type MovedArea struct {
	AreaUID         uuid.UUID
	SourceAreaUID   uuid.UUID
	InitialQuantity int
	CurrentQuantity int
	CreatedDate     time.Time
	LastUpdated     time.Time
}

type HarvestedStorage struct {
	Quantity      int
	SourceAreaUID uuid.UUID
	CreatedDate   time.Time
	LastUpdated   time.Time
}

type Trash struct {
	Quantity      int
	SourceAreaUID uuid.UUID
	CreatedDate   time.Time
	LastUpdated   time.Time
}

type CropNote struct {
	UID         uuid.UUID `json:"uid"`
	Content     string    `json:"content"`
	CreatedDate time.Time `json:"created_date"`
}

type CropInventory struct {
	UID           uuid.UUID `json:"uid"`
	PlantTypeCode string    `json:"plant_type"`
	Variety       string    `json:"variety"`
}

type CropArea struct {
	UID      uuid.UUID    `json:"uid"`
	Name     string       `json:"name"`
	Size     CropAreaUnit `json:"size"`
	Type     string       `json:"type"`
	Location string       `json:"location"`
	FarmUID  uuid.UUID    `json:"farm_uid"`
}

type CropAreaUnit struct {
	Value  float32 `json:"value"`
	Symbol string  `json:"symbol"`
}

type CropFarm struct {
	UID  uuid.UUID
	Name string
}

type CropPhoto struct {
	Filename string `json:"filename"`
	MimeType string `json:"mime_type"`
	Size     int    `json:"size"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
}

func CreateCropBatch(
	cropService CropService,
	areaUID uuid.UUID,
	cropType string,
	inventoryUID uuid.UUID,
	quantity int, containerType CropContainerType) (Crop, error) {

	serviceResult := cropService.FindAreaByID(areaUID)
	if serviceResult.Error != nil {
		return Crop{}, serviceResult.Error
	}

	area := serviceResult.Result.(CropArea)

	ct := GetCropType(cropType)
	if ct == (CropType{}) {
		return Crop{}, CropError{Code: CropErrorInvalidCropType}
	}

	serviceResult = cropService.FindInventoryMaterialByID(inventoryUID)
	if serviceResult.Error != nil {
		return Crop{}, serviceResult.Error
	}

	inv := serviceResult.Result.(CropInventory)

	err := validateContainer(quantity, containerType)
	if err != nil {
		return Crop{}, err
	}

	cropContainer := CropContainer{
		Quantity: quantity,
		Type:     containerType,
	}

	createdDate := time.Now()

	batchID, err := generateBatchID(cropService, inv, createdDate)
	if err != nil {
		return Crop{}, err
	}

	uid, err := uuid.NewV4()
	if err != nil {
		return Crop{}, err
	}

	return Crop{
		UID:          uid,
		BatchID:      batchID,
		Status:       GetCropStatus(CropActive),
		Type:         ct,
		Container:    cropContainer,
		InventoryUID: inv.UID,
		CreatedDate:  createdDate,
		InitialArea: InitialArea{
			AreaUID:         area.UID,
			InitialQuantity: quantity,
			CurrentQuantity: quantity,
		},
		FarmUID: area.FarmUID,
	}, nil
}

func (c *Crop) MoveToArea(cropService CropService, sourceAreaUID uuid.UUID, destinationAreaUID uuid.UUID, quantity int) error {
	// Validate //
	// Check if source area is exist in DB
	serviceResult := cropService.FindAreaByID(sourceAreaUID)
	if serviceResult.Error != nil {
		return CropError{Code: CropMoveToAreaErrorInvalidSourceArea}
	}

	srcArea, ok := serviceResult.Result.(CropArea)
	if !ok {
		return CropError{Code: CropMoveToAreaErrorInvalidSourceArea}
	}

	if srcArea.UID == (uuid.UUID{}) {
		return CropError{Code: CropMoveToAreaErrorSourceAreaNotFound}
	}

	// Check if destination area is exist in DB
	serviceResult = cropService.FindAreaByID(destinationAreaUID)
	dstArea, ok := serviceResult.Result.(CropArea)
	if !ok {
		return CropError{Code: CropMoveToAreaErrorInvalidDestinationArea}
	}

	if dstArea.UID == (uuid.UUID{}) {
		return CropError{Code: CropMoveToAreaErrorDestinationAreaNotFound}
	}

	// Check if movement rules for area type is valid
	isValidMoveRules := false
	if srcArea.Type == "seeding" && dstArea.Type == "growing" {
		isValidMoveRules = true
	} else if srcArea.Type == "seeding" && dstArea.Type == "seeding" {
		isValidMoveRules = true
	} else if srcArea.Type == "growing" && dstArea.Type == "growing" {
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

	// Check existance of the destination area input to the existing crop destination area.
	isDstExist := false
	for _, v := range c.MovedArea {
		if v.AreaUID == dstArea.UID {
			isDstExist = true
		}
	}

	// Process //
	if c.InitialArea.AreaUID == srcArea.UID {
		c.InitialArea.CurrentQuantity -= quantity
	}

	for i, v := range c.MovedArea {
		if v.AreaUID == srcArea.UID {
			c.MovedArea[i].CurrentQuantity -= quantity
		}
	}

	if isDstExist {
		for i, v := range c.MovedArea {
			if v.AreaUID == dstArea.UID {
				c.MovedArea[i].CurrentQuantity += quantity
				c.MovedArea[i].LastUpdated = time.Now()
			}
		}
	} else {
		c.MovedArea = append(c.MovedArea, MovedArea{
			AreaUID:         dstArea.UID,
			SourceAreaUID:   srcArea.UID,
			InitialQuantity: quantity,
			CurrentQuantity: quantity,
			CreatedDate:     time.Now(),
			LastUpdated:     time.Now(),
		})
	}

	return nil
}

func (c *Crop) Harvest(cropService CropService, sourceAreaUID uuid.UUID, quantity int) error {
	// Validate //
	// Check if source area is exist in DB
	serviceResult := cropService.FindAreaByID(sourceAreaUID)
	srcArea, ok := serviceResult.Result.(CropArea)
	if !ok {
		return CropError{Code: CropHarvestErrorInvalidSourceArea}
	}

	if srcArea == (CropArea{}) {
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

	if quantity <= 0 {
		return CropError{Code: CropHarvestErrorInvalidQuantity}
	}

	// Process //
	// Check source area existance. If already exist, then just update it
	isExist := false
	for i, v := range c.HarvestedStorage {
		if v.SourceAreaUID == srcArea.UID {
			c.HarvestedStorage[i].Quantity += quantity
			c.HarvestedStorage[i].LastUpdated = time.Now()
			isExist = true
		}
	}

	if !isExist {
		hs := HarvestedStorage{
			Quantity:      quantity,
			SourceAreaUID: srcArea.UID,
			CreatedDate:   time.Now(),
			LastUpdated:   time.Now(),
		}
		c.HarvestedStorage = append(c.HarvestedStorage, hs)
	}

	// Reduce the quantity in the area because it has been harvested
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

func (c *Crop) Dump(cropService CropService, sourceAreaUID uuid.UUID, quantity int) error {
	// Validate //
	// Check if source area is exist in DB
	serviceResult := cropService.FindAreaByID(sourceAreaUID)
	srcArea, ok := serviceResult.Result.(CropArea)
	if !ok {
		return CropError{Code: CropDumpErrorInvalidSourceArea}
	}

	if srcArea == (CropArea{}) {
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
	c.LastFertilized = time.Now()

	return nil
}

func (c *Crop) Prune() error {
	c.LastPruned = time.Now()

	return nil
}

func (c *Crop) Pesticide() error {
	c.LastPesticided = time.Now()

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
	serviceResult := cropService.FindInventoryMaterialByID(inventoryUID)

	if serviceResult.Error != nil {
		return serviceResult.Error
	}

	inventory := serviceResult.Result.(CropInventory)

	batchID, err := generateBatchID(cropService, inventory, c.CreatedDate)
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

	diff := now.Sub(c.CreatedDate)

	days := int(diff.Hours()) / 24

	return days
}

func generateBatchID(cropService CropService, inventory CropInventory, createdDate time.Time) (string, error) {
	// Generate Batch ID
	// Format the date to become daymonth format like 25jan
	dateFormat := strings.ToLower(createdDate.Format("2Jan"))

	// Get variety name and split it to slice
	varietySlice := strings.Fields(inventory.Variety)
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
