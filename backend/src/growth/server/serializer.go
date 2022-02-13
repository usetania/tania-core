package server

import (
	"encoding/json"
	"net/http"
	"sort"
	"time"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/usetania/tania-core/src/growth/domain"
	"github.com/usetania/tania-core/src/growth/query"
	"github.com/usetania/tania-core/src/growth/storage"
)

type CropListInArea struct {
	UID              uuid.UUID       `json:"uid"`
	BatchID          string          `json:"batch_id"`
	DaysSinceSeeding int             `json:"days_since_seeding"`
	InitialQuantity  int             `json:"initial_quantity"`
	CurrentQuantity  int             `json:"current_quantity"`
	Container        query.Container `json:"container"`
	LastWatered      *time.Time      `json:"last_watered"`
	SeedingDate      time.Time       `json:"seeding_date"`
	MovingDate       time.Time       `json:"moving_date"`
	InitialArea      InitialArea     `json:"initial_area"`
	Inventory        query.Inventory `json:"inventory"`
}

type InitialArea struct {
	AreaUID uuid.UUID `json:"area_id"`
	Name    string    `json:"name"`
}

type SortedCropNotes []domain.CropNote

// Len is part of sort.Interface.
func (sn SortedCropNotes) Len() int { return len(sn) }

// Swap is part of sort.Interface.
func (sn SortedCropNotes) Swap(i, j int) { sn[i], sn[j] = sn[j], sn[i] }

// Less is part of sort.Interface.
func (sn SortedCropNotes) Less(i, j int) bool { return sn[i].CreatedDate.After(sn[j].CreatedDate) }

type (
	CropActivity            storage.CropActivity
	SeedActivity            struct{ *storage.SeedActivity }
	MoveActivity            struct{ *storage.MoveActivity }
	HarvestActivity         struct{ *storage.HarvestActivity }
	DumpActivity            struct{ *storage.DumpActivity }
	PhotoActivity           struct{ *storage.PhotoActivity }
	WaterActivity           struct{ *storage.WaterActivity }
	TaskCropActivity        struct{ *storage.TaskCropActivity }
	TaskNutrientActivity    struct{ *storage.TaskNutrientActivity }
	TaskPestControlActivity struct {
		*storage.TaskPestControlActivity
	}
)

type (
	TaskSafetyActivity     struct{ *storage.TaskSafetyActivity }
	TaskSanitationActivity struct {
		*storage.TaskSanitationActivity
	}
)

func MapToCropActivity(activity storage.CropActivity) CropActivity {
	ca := CropActivity(activity)

	switch v := ca.ActivityType.(type) {
	case storage.SeedActivity:
		ca.ActivityType = SeedActivity{&v}
	case storage.MoveActivity:
		ca.ActivityType = MoveActivity{&v}
	case storage.HarvestActivity:
		ca.ActivityType = HarvestActivity{&v}
	case storage.DumpActivity:
		ca.ActivityType = DumpActivity{&v}
	case storage.PhotoActivity:
		ca.ActivityType = PhotoActivity{&v}
	case storage.WaterActivity:
		ca.ActivityType = WaterActivity{&v}
	case storage.TaskCropActivity:
		ca.ActivityType = TaskCropActivity{&v}
	case storage.TaskNutrientActivity:
		ca.ActivityType = TaskNutrientActivity{&v}
	case storage.TaskPestControlActivity:
		ca.ActivityType = TaskPestControlActivity{&v}
	case storage.TaskSanitationActivity:
		ca.ActivityType = TaskSanitationActivity{&v}
	case storage.TaskSafetyActivity:
		ca.ActivityType = TaskSafetyActivity{&v}
	}

	return ca
}

func MapToCropRead(s *GrowthServer, crop domain.Crop) (storage.CropRead, error) {
	queryResult := <-s.MaterialReadQuery.FindByID(crop.InventoryUID)
	if queryResult.Error != nil {
		return storage.CropRead{}, queryResult.Error
	}

	inv, ok := queryResult.Result.(query.CropMaterialQueryResult)
	if !ok {
		return storage.CropRead{}, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	queryResult = <-s.AreaReadQuery.FindByID(crop.InitialArea.AreaUID)
	if queryResult.Error != nil {
		return storage.CropRead{}, queryResult.Error
	}

	totalSeeding := 0
	totalGrowing := 0
	totalDumped := 0

	initialArea, ok := queryResult.Result.(query.CropAreaQueryResult)
	if !ok {
		return storage.CropRead{}, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	if initialArea.Type == "SEEDING" {
		totalSeeding += crop.InitialArea.CurrentQuantity
	} else if initialArea.Type == "GROWING" {
		totalGrowing += crop.InitialArea.CurrentQuantity
	}

	movedAreas := []storage.MovedArea{}

	for _, v := range crop.MovedArea {
		queryResult = <-s.AreaReadQuery.FindByID(v.AreaUID)
		if queryResult.Error != nil {
			return storage.CropRead{}, queryResult.Error
		}

		area, ok := queryResult.Result.(query.CropAreaQueryResult)
		if !ok {
			return storage.CropRead{}, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}

		if area.Type == "SEEDING" {
			totalSeeding += v.CurrentQuantity
		} else if area.Type == "GROWING" {
			totalGrowing += v.CurrentQuantity
		}

		var lastWatered *time.Time
		if !v.LastWatered.IsZero() {
			lastWatered = &v.LastWatered
		}

		var lastFertilized *time.Time
		if !v.LastFertilized.IsZero() {
			lastFertilized = &v.LastFertilized
		}

		var lastPesticided *time.Time
		if !v.LastPesticided.IsZero() {
			lastPesticided = &v.LastPesticided
		}

		var lastPruned *time.Time
		if !v.LastPruned.IsZero() {
			lastPruned = &v.LastPruned
		}

		movedAreas = append(movedAreas, storage.MovedArea{
			AreaUID:         area.UID,
			Name:            area.Name,
			InitialQuantity: v.InitialQuantity,
			CurrentQuantity: v.CurrentQuantity,
			LastWatered:     lastWatered,
			LastFertilized:  lastFertilized,
			LastPesticided:  lastPesticided,
			LastPruned:      lastPruned,
			CreatedDate:     v.CreatedDate,
			LastUpdated:     v.LastUpdated,
		})
	}

	harvestedStorage := []storage.HarvestedStorage{}

	for _, v := range crop.HarvestedStorage {
		queryResult = <-s.AreaReadQuery.FindByID(v.SourceAreaUID)
		if queryResult.Error != nil {
			return storage.CropRead{}, queryResult.Error
		}

		area, ok := queryResult.Result.(query.CropAreaQueryResult)
		if !ok {
			return storage.CropRead{}, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}

		harvestedStorage = append(harvestedStorage, storage.HarvestedStorage{
			Quantity:             v.Quantity,
			ProducedGramQuantity: v.ProducedGramQuantity,
			SourceAreaUID:        area.UID,
			SourceAreaName:       area.Name,
			CreatedDate:          v.CreatedDate,
			LastUpdated:          v.LastUpdated,
		})
	}

	trash := []storage.Trash{}

	for _, v := range crop.Trash {
		queryResult = <-s.AreaReadQuery.FindByID(v.SourceAreaUID)
		if queryResult.Error != nil {
			return storage.CropRead{}, queryResult.Error
		}

		area, ok := queryResult.Result.(query.CropAreaQueryResult)
		if !ok {
			return storage.CropRead{}, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}

		totalDumped += v.Quantity

		trash = append(trash, storage.Trash{
			Quantity:       v.Quantity,
			SourceAreaUID:  area.UID,
			SourceAreaName: area.Name,
			CreatedDate:    v.CreatedDate,
			LastUpdated:    v.LastUpdated,
		})
	}

	cropRead := storage.CropRead{}
	cropRead.UID = crop.UID
	cropRead.BatchID = crop.BatchID
	cropRead.Status = crop.Status.Code
	cropRead.Type = crop.Type.Code

	containerCell := 0
	switch v := crop.Container.Type.(type) {
	case domain.Tray:
		containerCell = v.Cell
	}

	cropRead.Container = storage.Container{
		Type:     crop.Container.Type.Code(),
		Quantity: crop.Container.Quantity,
		Cell:     containerCell,
	}

	cropRead.Inventory = storage.Inventory{
		UID:       inv.UID,
		PlantType: inv.PlantTypeCode,
		Type:      inv.TypeCode,
		Name:      inv.Name,
	}

	cropRead.AreaStatus = storage.AreaStatus{
		Seeding: totalSeeding,
		Growing: totalGrowing,
		Dumped:  totalDumped,
	}

	for _, v := range crop.Photos {
		cropRead.Photos = append(cropRead.Photos, storage.CropPhoto{
			UID:         v.UID,
			Filename:    v.Filename,
			MimeType:    v.MimeType,
			Size:        v.Size,
			Width:       v.Width,
			Height:      v.Height,
			Description: v.Description,
		})
	}

	cropRead.FarmUID = crop.FarmUID

	var lastWatered *time.Time
	if !crop.InitialArea.LastWatered.IsZero() {
		lastWatered = &crop.InitialArea.LastWatered
	}

	var lastFertilized *time.Time
	if !crop.InitialArea.LastFertilized.IsZero() {
		lastFertilized = &crop.InitialArea.LastFertilized
	}

	var lastPesticided *time.Time
	if !crop.InitialArea.LastPesticided.IsZero() {
		lastPesticided = &crop.InitialArea.LastPesticided
	}

	var lastPruned *time.Time
	if !crop.InitialArea.LastPruned.IsZero() {
		lastPruned = &crop.InitialArea.LastPruned
	}

	cropRead.InitialArea = storage.InitialArea{
		AreaUID:         initialArea.UID,
		Name:            initialArea.Name,
		InitialQuantity: crop.InitialArea.InitialQuantity,
		CurrentQuantity: crop.InitialArea.CurrentQuantity,
		LastWatered:     lastWatered,
		LastFertilized:  lastFertilized,
		LastPesticided:  lastPesticided,
		LastPruned:      lastPruned,
		CreatedDate:     crop.InitialArea.CreatedDate,
		LastUpdated:     crop.InitialArea.LastUpdated,
	}

	cropRead.MovedArea = movedAreas
	cropRead.HarvestedStorage = harvestedStorage
	cropRead.Trash = trash

	for _, v := range crop.Notes {
		cropRead.Notes = append(cropRead.Notes, v)
	}

	sort.Slice(cropRead.Notes, func(i, j int) bool {
		return cropRead.Notes[i].CreatedDate.After(cropRead.Notes[j].CreatedDate)
	})

	return cropRead, nil
}

func MapToCropListInArea(crop query.CropAreaByAreaQueryResult) (CropListInArea, error) {
	cl := CropListInArea{}

	cl.UID = crop.UID
	cl.BatchID = crop.BatchID
	cl.SeedingDate = crop.Area.InitialArea.CreatedDate

	now := time.Now()
	diff := now.Sub(cl.SeedingDate)
	cl.DaysSinceSeeding = int(diff.Hours()) / 24

	cl.InitialQuantity = crop.Area.InitialQuantity
	cl.CurrentQuantity = crop.Area.CurrentQuantity
	cl.Container = crop.Container
	cl.Inventory = crop.Inventory

	if crop.Area.LastWatered != nil && !crop.Area.LastWatered.IsZero() {
		cl.LastWatered = crop.Area.LastWatered
	}

	cl.MovingDate = crop.Area.MovingDate
	cl.InitialArea = InitialArea{
		AreaUID: crop.Area.InitialArea.UID,
		Name:    crop.Area.InitialArea.Name,
	}

	return cl, nil
}

func (a SeedActivity) MarshalJSON() ([]byte, error) {
	type Alias SeedActivity

	return json.Marshal(struct {
		*Alias
		Code string `json:"code"`
	}{
		Alias: (*Alias)(&a),
		Code:  a.Code(),
	})
}

func (a MoveActivity) MarshalJSON() ([]byte, error) {
	type Alias MoveActivity

	return json.Marshal(struct {
		*Alias
		Code string `json:"code"`
	}{
		Alias: (*Alias)(&a),
		Code:  a.Code(),
	})
}

func (a HarvestActivity) MarshalJSON() ([]byte, error) {
	type Alias HarvestActivity

	return json.Marshal(struct {
		*Alias
		Code string `json:"code"`
	}{
		Alias: (*Alias)(&a),
		Code:  a.Code(),
	})
}

func (a DumpActivity) MarshalJSON() ([]byte, error) {
	type Alias DumpActivity

	return json.Marshal(struct {
		*Alias
		Code string `json:"code"`
	}{
		Alias: (*Alias)(&a),
		Code:  a.Code(),
	})
}

func (a PhotoActivity) MarshalJSON() ([]byte, error) {
	type Alias PhotoActivity

	return json.Marshal(struct {
		*Alias
		Code string `json:"code"`
	}{
		Alias: (*Alias)(&a),
		Code:  a.Code(),
	})
}

func (a WaterActivity) MarshalJSON() ([]byte, error) {
	type Alias WaterActivity

	return json.Marshal(struct {
		*Alias
		Code string `json:"code"`
	}{
		Alias: (*Alias)(&a),
		Code:  a.Code(),
	})
}

func (a TaskCropActivity) MarshalJSON() ([]byte, error) {
	type Alias TaskCropActivity

	return json.Marshal(struct {
		*Alias
		Code string `json:"code"`
	}{
		Alias: (*Alias)(&a),
		Code:  a.Code(),
	})
}

func (a TaskNutrientActivity) MarshalJSON() ([]byte, error) {
	type Alias TaskNutrientActivity

	return json.Marshal(struct {
		*Alias
		Code string `json:"code"`
	}{
		Alias: (*Alias)(&a),
		Code:  a.Code(),
	})
}

func (a TaskPestControlActivity) MarshalJSON() ([]byte, error) {
	type Alias TaskPestControlActivity

	return json.Marshal(struct {
		*Alias
		Code string `json:"code"`
	}{
		Alias: (*Alias)(&a),
		Code:  a.Code(),
	})
}

func (a TaskSafetyActivity) MarshalJSON() ([]byte, error) {
	type Alias TaskSafetyActivity

	return json.Marshal(struct {
		*Alias
		Code string `json:"code"`
	}{
		Alias: (*Alias)(&a),
		Code:  a.Code(),
	})
}

func (a TaskSanitationActivity) MarshalJSON() ([]byte, error) {
	type Alias TaskSanitationActivity

	return json.Marshal(struct {
		*Alias
		Code string `json:"code"`
	}{
		Alias: (*Alias)(&a),
		Code:  a.Code(),
	})
}
