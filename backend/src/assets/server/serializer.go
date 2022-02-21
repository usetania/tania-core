package server

import (
	"encoding/json"
	"net/http"
	"sort"
	"time"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/usetania/tania-core/src/assets/domain"
	"github.com/usetania/tania-core/src/assets/query"
	"github.com/usetania/tania-core/src/assets/storage"
)

type (
	SimpleFarm domain.Farm
	SimpleArea struct {
		UID  uuid.UUID
		Name string
		Type string
	}
)

type AreaList struct {
	UID            uuid.UUID        `json:"uid"`
	Name           string           `json:"name"`
	Type           string           `json:"type"`
	Size           storage.AreaSize `json:"size"`
	TotalCropBatch int              `json:"total_crop_batch"`
	PlantQuantity  int              `json:"plant_quantity"`
}

type DetailArea struct {
	storage.AreaRead
	TotalCropBatch int `json:"total_crop_batch"`
	TotalVariety   int `json:"total_variety"`
	PlantQuantity  int `json:"plant_quantity"`
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

type SortedReservoirNotes []domain.ReservoirNote

// Len is part of sort.Interface.
func (sn SortedReservoirNotes) Len() int { return len(sn) }

// Swap is part of sort.Interface.
func (sn SortedReservoirNotes) Swap(i, j int) { sn[i], sn[j] = sn[j], sn[i] }

// Less is part of sort.Interface.
func (sn SortedReservoirNotes) Less(i, j int) bool { return sn[i].CreatedDate.After(sn[j].CreatedDate) }

type (
	ReservoirBucket struct{ domain.Bucket }
	ReservoirTap    struct{ domain.Tap }
)

type MaterialSimple struct {
	UID  uuid.UUID    `json:"uid"`
	Name string       `json:"name"`
	Type MaterialType `json:"type"`
}

type Material struct {
	UID            uuid.UUID        `json:"uid"`
	Name           string           `json:"name"`
	PricePerUnit   PricePerUnit     `json:"price_per_unit"`
	Type           MaterialType     `json:"type"`
	Quantity       MaterialQuantity `json:"quantity"`
	ExpirationDate *time.Time       `json:"expiration_date,omitempty"`
	Notes          *string          `json:"notes"`
	ProducedBy     *string          `json:"produced_by"`
	CreatedDate    time.Time        `json:"created_date"`
}

type PricePerUnit struct {
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
	farmRead.Country = farm.Country
	farmRead.City = farm.City
	farmRead.CreatedDate = farm.CreatedDate
	farmRead.IsActive = farm.IsActive

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

func MapToAreaList(s *FarmServer, areas []storage.AreaRead) ([]AreaList, error) {
	areaList := make([]AreaList, len(areas))

	for i, area := range areas {
		queryResult := <-s.CropReadQuery.CountCropsByArea(area.UID)
		if queryResult.Error != nil {
			return []AreaList{}, queryResult.Error
		}

		cropCount, ok := queryResult.Result.(query.CountAreaCropResult)
		if !ok {
			return []AreaList{}, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}

		areaList[i] = AreaList{
			UID:            area.UID,
			Name:           area.Name,
			Type:           area.Type,
			Size:           area.Size,
			TotalCropBatch: cropCount.TotalCropBatch,
			PlantQuantity:  cropCount.PlantQuantity,
		}
	}

	return areaList, nil
}

func MapToReservoirRead(s *FarmServer, reservoir domain.Reservoir) (storage.ReservoirRead, error) {
	resRead := storage.ReservoirRead{}

	resRead.UID = reservoir.UID
	resRead.Name = reservoir.Name
	resRead.CreatedDate = reservoir.CreatedDate

	switch v := reservoir.WaterSource.(type) {
	case domain.Bucket:
		resRead.WaterSource = storage.WaterSource{
			Type:     v.Type(),
			Capacity: v.Capacity,
		}
	case domain.Tap:
		resRead.WaterSource = storage.WaterSource{
			Type: v.Type(),
		}
	}

	for _, v := range reservoir.Notes {
		resRead.Notes = append(resRead.Notes, storage.ReservoirNote{
			UID:         v.UID,
			Content:     v.Content,
			CreatedDate: v.CreatedDate,
		})
	}

	sort.Slice(resRead.Notes, func(i, j int) bool {
		return resRead.Notes[i].CreatedDate.After(resRead.Notes[j].CreatedDate)
	})

	queryResult := <-s.FarmReadQuery.FindByID(reservoir.FarmUID)
	if queryResult.Error != nil {
		return storage.ReservoirRead{}, echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	farm, ok := queryResult.Result.(storage.FarmRead)
	if !ok {
		return storage.ReservoirRead{}, echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	resRead.Farm = storage.ReservoirFarm{
		UID:  farm.UID,
		Name: farm.Name,
	}

	queryResult = <-s.AreaReadQuery.FindAreasByReservoirID(reservoir.UID)
	if queryResult.Error != nil {
		return storage.ReservoirRead{}, echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	areas, ok := queryResult.Result.([]storage.AreaRead)
	if !ok {
		return storage.ReservoirRead{}, echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	for _, v := range areas {
		resRead.InstalledToArea = append(resRead.InstalledToArea, storage.AreaInstalled{
			UID:  v.UID,
			Name: v.Name,
		})
	}

	return resRead, nil
}

func MapToReservoirReadFromRead(s *FarmServer, reservoir storage.ReservoirRead) (storage.ReservoirRead, error) {
	queryResult := <-s.AreaReadQuery.FindAreasByReservoirID(reservoir.UID)
	if queryResult.Error != nil {
		return storage.ReservoirRead{}, echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	areas, ok := queryResult.Result.([]storage.AreaRead)
	if !ok {
		return storage.ReservoirRead{}, echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	for _, v := range areas {
		reservoir.InstalledToArea = append(reservoir.InstalledToArea, storage.AreaInstalled{
			UID:  v.UID,
			Name: v.Name,
		})
	}

	sort.Slice(reservoir.Notes, func(i, j int) bool {
		return reservoir.Notes[i].CreatedDate.After(reservoir.Notes[j].CreatedDate)
	})

	return reservoir, nil
}

func MapToDetailAreaFromStorage(s *FarmServer, areaRead storage.AreaRead) (DetailArea, error) {
	detailArea := DetailArea{}

	detailArea.UID = areaRead.UID
	detailArea.Name = areaRead.Name
	detailArea.Type = areaRead.Type
	detailArea.Location = areaRead.Location
	detailArea.Photo = areaRead.Photo
	detailArea.Size = areaRead.Size
	detailArea.CreatedDate = areaRead.CreatedDate
	detailArea.Reservoir = areaRead.Reservoir
	detailArea.Farm = areaRead.Farm

	queryResult := <-s.CropReadQuery.CountCropsByArea(areaRead.UID)
	if queryResult.Error != nil {
		return DetailArea{}, queryResult.Error
	}

	cropCount, ok := queryResult.Result.(query.CountAreaCropResult)
	if !ok {
		return DetailArea{}, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	detailArea.TotalCropBatch = cropCount.TotalCropBatch
	detailArea.PlantQuantity = cropCount.PlantQuantity

	queryResult = <-s.CropReadQuery.FindAllCropByArea(areaRead.UID)
	if queryResult.Error != nil {
		return DetailArea{}, queryResult.Error
	}

	crops, ok := queryResult.Result.([]query.AreaCropResult)
	if !ok {
		return DetailArea{}, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	uniqueInventories := make(map[uuid.UUID]bool)
	for _, v := range crops {
		if _, ok := uniqueInventories[v.Inventory.UID]; !ok {
			uniqueInventories[v.Inventory.UID] = true
		}
	}

	detailArea.TotalVariety = len(uniqueInventories)

	detailArea.Notes = areaRead.Notes
	sort.Slice(detailArea.Notes, func(i, j int) bool {
		return detailArea.Notes[i].CreatedDate.After(detailArea.Notes[j].CreatedDate)
	})

	return detailArea, nil
}

func MapToDetailArea(s *FarmServer, area domain.Area) (DetailArea, error) {
	areaRead := DetailArea{}

	areaRead.UID = area.UID
	areaRead.Name = area.Name
	areaRead.Type = area.Type.Code
	areaRead.Location = storage.AreaLocation(area.Location)
	areaRead.Photo = storage.AreaPhoto(area.Photo)
	areaRead.Size = storage.AreaSize(area.Size)
	areaRead.CreatedDate = area.CreatedDate

	queryResult := <-s.ReservoirReadQuery.FindByID(area.ReservoirUID)
	if queryResult.Error != nil {
		return DetailArea{}, echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	reservoir, ok := queryResult.Result.(storage.ReservoirRead)
	if !ok {
		return DetailArea{}, echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	areaRead.Reservoir = storage.AreaReservoir{
		UID:  reservoir.UID,
		Name: reservoir.Name,
	}

	queryResult = <-s.FarmReadQuery.FindByID(area.FarmUID)
	if queryResult.Error != nil {
		return DetailArea{}, echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	farm, ok := queryResult.Result.(storage.FarmRead)
	if !ok {
		return DetailArea{}, echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	areaRead.Farm = storage.AreaFarm{
		UID:  farm.UID,
		Name: farm.Name,
	}

	queryResult = <-s.CropReadQuery.CountCropsByArea(area.UID)
	if queryResult.Error != nil {
		return DetailArea{}, queryResult.Error
	}

	cropCount, ok := queryResult.Result.(query.CountAreaCropResult)
	if !ok {
		return DetailArea{}, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	areaRead.TotalCropBatch = cropCount.TotalCropBatch
	areaRead.PlantQuantity = cropCount.PlantQuantity

	queryResult = <-s.CropReadQuery.FindAllCropByArea(area.UID)
	if queryResult.Error != nil {
		return DetailArea{}, queryResult.Error
	}

	crops, ok := queryResult.Result.([]query.AreaCropResult)
	if !ok {
		return DetailArea{}, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	uniqueInventories := make(map[uuid.UUID]bool)
	for _, v := range crops {
		if _, ok := uniqueInventories[v.Inventory.UID]; !ok {
			uniqueInventories[v.Inventory.UID] = true
		}
	}

	areaRead.TotalVariety = len(uniqueInventories)

	for _, v := range area.Notes {
		areaRead.Notes = append(areaRead.Notes, storage.AreaNote(v))
	}

	sort.Slice(areaRead.Notes, func(i, j int) bool {
		return areaRead.Notes[i].CreatedDate.After(areaRead.Notes[j].CreatedDate)
	})

	return areaRead, nil
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
	m.PricePerUnit = PricePerUnit{
		Code:   material.PricePerUnit.CurrencyCode,
		Symbol: material.PricePerUnit.Symbol(),
		Amount: material.PricePerUnit.Amount,
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

	m.ProducedBy = nil
	if material.ProducedBy != nil {
		m.ProducedBy = material.ProducedBy
	}

	m.CreatedDate = material.CreatedDate

	return m
}

func MapToMaterialFromRead(material storage.MaterialRead) Material {
	m := Material{}

	m.UID = material.UID
	m.Name = material.Name

	ppu := domain.PricePerUnit(material.PricePerUnit)
	m.PricePerUnit = PricePerUnit{
		Amount: ppu.Amount,
		Code:   ppu.CurrencyCode,
		Symbol: ppu.Symbol(),
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

	m.ProducedBy = nil
	if material.ProducedBy != nil {
		m.ProducedBy = material.ProducedBy
	}

	m.CreatedDate = material.CreatedDate

	return m
}

func MapToMaterialSimpleFromRead(materials []storage.MaterialRead) []MaterialSimple {
	m := []MaterialSimple{}

	for _, val := range materials {
		mSimple := MaterialSimple{
			UID:  val.UID,
			Name: val.Name,
		}

		switch v := val.Type.(type) {
		case domain.MaterialTypeSeed:
			mSimple.Type = MaterialType{
				Code: v.Code(),
				MaterialTypeDetail: MaterialTypeSeed{
					PlantType: v.PlantType,
				},
			}
		case domain.MaterialTypePlant:
			mSimple.Type = MaterialType{
				Code: v.Code(),
				MaterialTypeDetail: MaterialTypePlant{
					PlantType: v.PlantType,
				},
			}
		case domain.MaterialTypeAgrochemical:
			mSimple.Type = MaterialType{
				Code: v.Code(),
				MaterialTypeDetail: MaterialTypeAgrochemical{
					ChemicalType: v.ChemicalType,
				},
			}
		case domain.MaterialTypeGrowingMedium:
			mSimple.Type = MaterialType{Code: v.Code()}
		case domain.MaterialTypeLabelAndCropSupport:
			mSimple.Type = MaterialType{Code: v.Code()}
		case domain.MaterialTypeSeedingContainer:
			mSimple.Type = MaterialType{
				Code: v.Code(),
				MaterialTypeDetail: MaterialTypeSeedingContainer{
					ContainerType: v.ContainerType,
				},
			}
		case domain.MaterialTypePostHarvestSupply:
			mSimple.Type = MaterialType{Code: v.Code()}
		case domain.MaterialTypeOther:
			mSimple.Type = MaterialType{Code: v.Code()}
		}

		m = append(m, mSimple)
	}

	return m
}

func MapToAvailableMaterialPlantType(materials []storage.MaterialRead) []AvailableMaterialPlantType {
	ai := make(map[string]AvailableMaterialPlantType)

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
