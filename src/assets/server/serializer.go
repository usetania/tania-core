package server

import (
	"encoding/json"
	"net/http"
	"sort"
	"time"

	"github.com/Tanibox/tania-server/src/assets/domain"
	"github.com/Tanibox/tania-server/src/assets/query"
	"github.com/Tanibox/tania-server/src/assets/storage"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
)

type SimpleFarm domain.Farm
type SimpleArea struct {
	UID  uuid.UUID
	Name string
	Type string
}
type AreaList struct {
	UID            uuid.UUID       `json:"uid"`
	Name           string          `json:"name"`
	Type           string          `json:"type"`
	Size           domain.AreaSize `json:"size"`
	TotalCropBatch int             `json:"total_crop_batch"`
	PlantQuantity  int             `json:"plant_quantity"`
}
type DetailArea struct {
	UID            uuid.UUID                   `json:"uid"`
	Name           string                      `json:"name"`
	Size           domain.AreaSize             `json:"size"`
	Type           string                      `json:"type"`
	Location       string                      `json:"location"`
	Photo          domain.AreaPhoto            `json:"photo"`
	Reservoir      domain.Reservoir            `json:"reservoir"`
	TotalCropBatch int                         `json:"total_crop_batch"`
	TotalVariety   int                         `json:"total_variety"`
	Crops          []query.AreaCropQueryResult `json:"crops"`
	Notes          SortedAreaNotes             `json:"notes"`
}

type DetailReservoir struct {
	UID              uuid.UUID            `json:"uid"`
	Name             string               `json:"name"`
	WaterSource      WaterSource          `json:"water_source"`
	Farm             SimpleFarm           `json:"farm"`
	CreatedDate      time.Time            `json:"created_date"`
	Notes            SortedReservoirNotes `json:"notes"`
	InstalledToAreas []SimpleArea         `json:"installed_to_areas"`
}

type WaterSource struct {
	Type     string
	Capacity float32
}

type ReservoirNote struct {
	UID         uuid.UUID
	Content     string
	CreatedDate time.Time
}

type SortedReservoirNotes []ReservoirNote

// Len is part of sort.Interface.
func (sn SortedReservoirNotes) Len() int { return len(sn) }

// Swap is part of sort.Interface.
func (sn SortedReservoirNotes) Swap(i, j int) { sn[i], sn[j] = sn[j], sn[i] }

// Less is part of sort.Interface.
func (sn SortedReservoirNotes) Less(i, j int) bool { return sn[i].CreatedDate.After(sn[j].CreatedDate) }

type ReservoirBucket struct{ domain.Bucket }
type ReservoirTap struct{ domain.Tap }

type Material struct {
	UID            uuid.UUID        `json:"uid"`
	Name           string           `json:"name"`
	PricePerUnit   Money            `json:"price_per_unit"`
	Type           MaterialType     `json:"type"`
	Quantity       MaterialQuantity `json:"quantity"`
	ExpirationDate *time.Time       `json:"expiration_date,omitempty"`
	Notes          *string          `json:"notes"`
	IsExpense      *bool            `json:"is_expense"`
	ProducedBy     *string          `json:"produced_by"`
}

type Money struct {
	Code   string `json:"code"`
	Symbol string `json:"symbol"`
	Amount string `json:"amount"`
}

type MaterialQuantity struct {
	Value float32 `json:"value"`
	Unit  string  `json:"unit"`
}

type MaterialType struct {
	Code               string      `json:"code"`
	MaterialTypeDetail interface{} `json:"type_detail"`
}

type MaterialTypeSeed struct {
	PlantType domain.PlantType `json:"plant_type"`
}

type MaterialTypePlant struct {
	PlantType domain.PlantType `json:"plant_type"`
}

type MaterialTypeAgrochemical struct {
	ChemicalType domain.ChemicalType `json:"chemical_type"`
}

type MaterialTypeSeedingContainer struct {
	ContainerType domain.ContainerType `json:"container_type"`
}

type AvailableMaterialPlantType struct {
	PlantType string   `json:"plant_type"`
	Names     []string `json:"names"`
}

type SortedAreaNotes []domain.AreaNote

// Len is part of sort.Interface.
func (sn SortedAreaNotes) Len() int { return len(sn) }

// Swap is part of sort.Interface.
func (sn SortedAreaNotes) Swap(i, j int) { sn[i], sn[j] = sn[j], sn[i] }

// Less is part of sort.Interface.
func (sn SortedAreaNotes) Less(i, j int) bool { return sn[i].CreatedDate.After(sn[j].CreatedDate) }

func MapToFarmRead(farm *domain.Farm) *storage.FarmRead {
	farmRead := &storage.FarmRead{}
	farmRead.UID = farm.UID
	farmRead.Name = farm.Name
	farmRead.Type = farm.Type
	farmRead.Latitude = farm.Latitude
	farmRead.Longitude = farm.Longitude
	farmRead.CountryCode = farm.CountryCode
	farmRead.CityCode = farm.CityCode
	farmRead.CreatedDate = farm.CreatedDate

	return farmRead
}

func MapToSimpleFarm(farms []domain.Farm) []SimpleFarm {
	farmList := make([]SimpleFarm, len(farms))

	for i, farm := range farms {
		farmList[i] = SimpleFarm(farm)
	}

	return farmList
}

func MapToSimpleArea(areas []domain.Area) []SimpleArea {
	simpleAreaList := make([]SimpleArea, len(areas))

	for i, area := range areas {
		simpleAreaList[i] = SimpleArea{
			UID:  area.UID,
			Name: area.Name,
			Type: area.Type.Code,
		}
	}

	return simpleAreaList
}

func MapToAreaList(s *FarmServer, areas []domain.Area) ([]AreaList, error) {
	areaList := make([]AreaList, len(areas))

	for i, area := range areas {
		queryResult := <-s.CropQuery.CountCropsByArea(area.UID)
		if queryResult.Error != nil {
			return []AreaList{}, queryResult.Error
		}

		cropCount, ok := queryResult.Result.(query.CountAreaCropQueryResult)
		if !ok {
			return []AreaList{}, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}

		areaList[i] = AreaList{
			UID:            area.UID,
			Name:           area.Name,
			Type:           area.Type.Code,
			Size:           area.Size,
			TotalCropBatch: cropCount.TotalCropBatch,
			PlantQuantity:  cropCount.PlantQuantity,
		}
	}

	return areaList, nil
}

func MapToReservoirListFromQuery(s *FarmServer, reservoirs []query.ReservoirReadQueryResult) ([]DetailReservoir, error) {
	reservoirList := make([]DetailReservoir, len(reservoirs))

	for _, reservoir := range reservoirs {
		detailReservoir, err := MapToReservoirFromQuery(s, reservoir)
		if err != nil {
			return nil, err
		}

		reservoirList = append(reservoirList, detailReservoir)
	}

	return reservoirList, nil
}

func MapToReservoirFromQuery(s *FarmServer, reservoir query.ReservoirReadQueryResult) (DetailReservoir, error) {
	detailReservoir := DetailReservoir{}
	detailReservoir.UID = reservoir.UID
	detailReservoir.Name = reservoir.Name
	detailReservoir.WaterSource = WaterSource{
		Type:     reservoir.WaterSource.Type,
		Capacity: reservoir.WaterSource.Capacity,
	}
	detailReservoir.CreatedDate = reservoir.CreatedDate

	for _, v := range reservoir.Notes {
		detailReservoir.Notes = append(detailReservoir.Notes, ReservoirNote{
			UID:         v.UID,
			Content:     v.Content,
			CreatedDate: v.CreatedDate,
		})
	}
	sort.Sort(detailReservoir.Notes)

	queryResult := <-s.FarmReadQuery.FindByID(reservoir.FarmUID)
	if queryResult.Error != nil {
		return DetailReservoir{}, echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	farm, ok := queryResult.Result.(query.FarmReadQueryResult)
	if !ok {
		return DetailReservoir{}, echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	detailReservoir.Farm = SimpleFarm{
		UID:  farm.UID,
		Name: farm.Name,
		Type: farm.Type,
	}

	queryResult = <-s.AreaQuery.FindAreasByReservoirID(reservoir.UID.String())
	if queryResult.Error != nil {
		return DetailReservoir{}, echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	areas, ok := queryResult.Result.([]domain.Area)
	if !ok {
		return DetailReservoir{}, echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	detailReservoir.InstalledToAreas = MapToSimpleArea(areas)

	return detailReservoir, nil
}

func MapToDetailReservoir(s *FarmServer, reservoir domain.Reservoir) (DetailReservoir, error) {
	detailReservoir := DetailReservoir{}

	detailReservoir.UID = reservoir.UID
	detailReservoir.Name = reservoir.Name
	detailReservoir.CreatedDate = reservoir.CreatedDate

	switch v := reservoir.WaterSource.(type) {
	case domain.Bucket:
		detailReservoir.WaterSource = WaterSource{
			Type:     v.Type(),
			Capacity: v.Capacity,
		}
	case domain.Tap:
		detailReservoir.WaterSource = WaterSource{
			Type: v.Type(),
		}
	}

	for _, v := range reservoir.Notes {
		detailReservoir.Notes = append(detailReservoir.Notes, ReservoirNote{
			UID:         v.UID,
			Content:     v.Content,
			CreatedDate: v.CreatedDate,
		})
	}
	sort.Sort(detailReservoir.Notes)

	queryResult := <-s.FarmReadQuery.FindByID(reservoir.FarmUID)
	if queryResult.Error != nil {
		return DetailReservoir{}, echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	farm, ok := queryResult.Result.(query.FarmReadQueryResult)
	if !ok {
		return DetailReservoir{}, echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	detailReservoir.Farm = SimpleFarm{
		UID:  farm.UID,
		Name: farm.Name,
		Type: farm.Type,
	}

	queryResult = <-s.AreaQuery.FindAreasByReservoirID(reservoir.UID.String())
	if queryResult.Error != nil {
		return DetailReservoir{}, echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	areas, ok := queryResult.Result.([]domain.Area)
	if !ok {
		return DetailReservoir{}, echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	detailReservoir.InstalledToAreas = MapToSimpleArea(areas)

	return detailReservoir, nil
}

func MapToDetailArea(s *FarmServer, area domain.Area) (DetailArea, error) {
	detailArea := DetailArea{}

	detailArea.UID = area.UID
	detailArea.Name = area.Name
	detailArea.Type = area.Type.Code
	detailArea.Location = area.Location.Code
	detailArea.Photo = area.Photo
	detailArea.Size = area.Size

	detailArea.Reservoir = area.Reservoir
	switch v := area.Reservoir.WaterSource.(type) {
	case domain.Bucket:
		detailArea.Reservoir.WaterSource = ReservoirBucket{Bucket: v}
	case domain.Tap:
		detailArea.Reservoir.WaterSource = ReservoirTap{Tap: v}
	}

	queryResult := <-s.CropQuery.CountCropsByArea(area.UID)
	if queryResult.Error != nil {
		return DetailArea{}, queryResult.Error
	}

	cropCount, ok := queryResult.Result.(query.CountAreaCropQueryResult)
	if !ok {
		return DetailArea{}, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	detailArea.TotalCropBatch = cropCount.TotalCropBatch

	queryResult = <-s.CropQuery.FindAllCropByArea(area.UID)
	if queryResult.Error != nil {
		return DetailArea{}, queryResult.Error
	}

	crops, ok := queryResult.Result.([]query.AreaCropQueryResult)
	if !ok {
		return DetailArea{}, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	for i, v := range crops {
		repoResult := <-s.AreaRepo.FindByID(v.InitialArea.AreaUID.String())

		if repoResult.Error != nil {
			return DetailArea{}, queryResult.Error
		}

		a, ok := repoResult.Result.(domain.Area)
		if !ok {
			return DetailArea{}, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}

		repoResult = <-s.MaterialRepo.FindByID(v.Inventory.UID.String())
		if repoResult.Error != nil {
			return DetailArea{}, repoResult.Error
		}

		mat, ok := repoResult.Result.(domain.Material)
		if !ok {
			return DetailArea{}, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}

		now := time.Now()
		diff := now.Sub(v.CreatedDate)
		days := int(diff.Hours()) / 24

		crops[i].DaysSinceSeeding = days
		crops[i].InitialArea.Name = a.Name
		crops[i].Inventory.UID = mat.UID
		crops[i].Inventory.Name = mat.Name
	}

	detailArea.Crops = crops

	uniqueInventories := make(map[uuid.UUID]bool)
	for _, v := range crops {
		if _, ok := uniqueInventories[v.Inventory.UID]; !ok {
			uniqueInventories[v.Inventory.UID] = true
		}
	}

	detailArea.TotalVariety = len(uniqueInventories)

	notes := make(SortedAreaNotes, 0, len(area.Notes))
	for _, v := range area.Notes {
		notes = append(notes, v)
	}

	sort.Sort(notes)

	detailArea.Notes = notes

	return detailArea, nil
}

func MapToPlantType(plantTypes []domain.PlantType) []string {
	pt := make([]string, len(plantTypes))

	for i, v := range plantTypes {
		pt[i] = v.Code
	}

	return pt
}

func MapToMaterial(material domain.Material) Material {
	m := Material{}

	m.UID = material.UID
	m.Name = material.Name
	m.PricePerUnit = Money{
		Code:   material.PricePerUnit.Code(),
		Symbol: material.PricePerUnit.Symbol(),
		Amount: material.PricePerUnit.Amount(),
	}

	m.ExpirationDate = nil
	if material.ExpirationDate != nil {
		m.ExpirationDate = material.ExpirationDate
	}

	switch v := material.Type.(type) {
	case domain.MaterialTypeSeed:
		m.Type = MaterialType{
			Code: v.Code(),
			MaterialTypeDetail: MaterialTypeSeed{
				PlantType: v.PlantType,
			},
		}
	case domain.MaterialTypePlant:
		m.Type = MaterialType{
			Code: v.Code(),
			MaterialTypeDetail: MaterialTypePlant{
				PlantType: v.PlantType,
			},
		}
	case domain.MaterialTypeAgrochemical:
		m.Type = MaterialType{
			Code: v.Code(),
			MaterialTypeDetail: MaterialTypeAgrochemical{
				ChemicalType: v.ChemicalType,
			},
		}
	case domain.MaterialTypeGrowingMedium:
		m.Type = MaterialType{Code: v.Code()}
		m.ExpirationDate = nil
	case domain.MaterialTypeLabelAndCropSupport:
		m.Type = MaterialType{Code: v.Code()}
	case domain.MaterialTypeSeedingContainer:
		m.Type = MaterialType{
			Code: v.Code(),
			MaterialTypeDetail: MaterialTypeSeedingContainer{
				ContainerType: v.ContainerType,
			},
		}
	case domain.MaterialTypePostHarvestSupply:
		m.Type = MaterialType{Code: v.Code()}
	case domain.MaterialTypeOther:
		m.Type = MaterialType{Code: v.Code()}
	}

	m.Quantity = MaterialQuantity{
		Value: material.Quantity.Value,
		Unit:  material.Quantity.Unit.Code,
	}

	m.Notes = nil
	if material.Notes != nil {
		m.Notes = material.Notes
	}

	m.IsExpense = nil
	if material.IsExpense != nil {
		m.IsExpense = material.IsExpense
	}

	m.ProducedBy = nil
	if material.ProducedBy != nil {
		m.ProducedBy = material.ProducedBy
	}

	return m
}

func MapToAvailableMaterialPlantType(materials []domain.Material) []AvailableMaterialPlantType {
	ai := make(map[string]AvailableMaterialPlantType, 0)

	// Convert domain.Material to AvailableMaterialPlantType first with Map
	for _, v := range materials {
		switch mt := v.Type.(type) {
		case domain.MaterialTypeSeed:
			asm := AvailableMaterialPlantType{
				PlantType: mt.PlantType.Code,
				Names:     append(ai[mt.PlantType.Code].Names, v.Name),
			}

			ai[mt.PlantType.Code] = asm
		case domain.MaterialTypePlant:
			asm := AvailableMaterialPlantType{
				PlantType: mt.PlantType.Code,
				Names:     append(ai[mt.PlantType.Code].Names, v.Name),
			}

			ai[mt.PlantType.Code] = asm
		}
	}

	// From Map, we need to change it to slice for the json response purpose
	aiSlice := []AvailableMaterialPlantType{}
	for _, v := range ai {
		aiSlice = append(aiSlice, v)
	}

	return aiSlice
}

func (sf SimpleFarm) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		UID  string `json:"uid"`
		Name string `json:"name"`
		Type string `json:"type"`
	}{
		UID:  sf.UID.String(),
		Name: sf.Name,
		Type: sf.Type,
	})
}

func (sa SimpleArea) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		UID  string `json:"uid"`
		Name string `json:"name"`
		Type string `json:"type"`
	}{
		UID:  sa.UID.String(),
		Name: sa.Name,
		Type: sa.Type,
	})
}

func (rb ReservoirBucket) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type     string  `json:"type"`
		Capacity float32 `json:"capacity"`
	}{
		Type:     rb.Type(),
		Capacity: rb.Capacity,
	})
}

func (rt ReservoirTap) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string `json:"type"`
	}{
		Type: rt.Type(),
	})
}
