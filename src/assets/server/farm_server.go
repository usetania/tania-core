package server

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Tanibox/tania-server/config"
	"github.com/Tanibox/tania-server/src/assets/domain"
	"github.com/Tanibox/tania-server/src/assets/domain/service"
	"github.com/Tanibox/tania-server/src/assets/query"
	"github.com/Tanibox/tania-server/src/assets/query/inmemory"
	"github.com/Tanibox/tania-server/src/assets/repository"
	"github.com/Tanibox/tania-server/src/assets/storage"
	growthstorage "github.com/Tanibox/tania-server/src/growth/storage"
	"github.com/Tanibox/tania-server/src/helper/imagehelper"
	"github.com/Tanibox/tania-server/src/helper/stringhelper"
	"github.com/Tanibox/tania-server/src/helper/structhelper"
	"github.com/asaskevich/EventBus"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
)

// FarmServer ties the routes and handlers with injected dependencies
type FarmServer struct {
	FarmRepo            repository.FarmRepository
	FarmEventRepo       repository.FarmEventRepository
	FarmReadRepo        repository.FarmReadRepository
	FarmReadQuery       query.FarmReadQuery
	ReservoirRepo       repository.ReservoirRepository
	ReservoirEventRepo  repository.ReservoirEventRepository
	ReservoirEventQuery query.ReservoirEventQuery
	ReservoirReadRepo   repository.ReservoirReadRepository
	ReservoirReadQuery  query.ReservoirReadQuery
	ReservoirService    domain.ReservoirService
	AreaRepo            repository.AreaRepository
	AreaQuery           query.AreaQuery
	MaterialRepo        repository.MaterialRepository
	MaterialQuery       query.MaterialQuery
	CropQuery           query.CropQuery
	File                File
	EventBus            EventBus.Bus
}

// NewFarmServer initializes FarmServer's dependencies and create new FarmServer struct
func NewFarmServer(
	farmStorage *storage.FarmStorage,
	farmEventStorage *storage.FarmEventStorage,
	farmReadStorage *storage.FarmReadStorage,
	areaStorage *storage.AreaStorage,
	reservoirStorage *storage.ReservoirStorage,
	reservoirEventStorage *storage.ReservoirEventStorage,
	reservoirReadStorage *storage.ReservoirReadStorage,
	materialStorage *storage.MaterialStorage,
	cropStorage *growthstorage.CropStorage,
	eventBus EventBus.Bus,
) (*FarmServer, error) {
	farmRepo := repository.NewFarmRepositoryInMemory(farmStorage)
	farmEventRepo := repository.NewFarmEventRepositoryInMemory(farmEventStorage)
	farmReadRepo := repository.NewFarmReadRepositoryInMemory(farmReadStorage)
	farmReadQuery := inmemory.NewFarmReadQueryInMemory(farmReadStorage)

	areaRepo := repository.NewAreaRepositoryInMemory(areaStorage)
	areaQuery := inmemory.NewAreaQueryInMemory(areaStorage)

	reservoirRepo := repository.NewReservoirRepositoryInMemory(reservoirStorage)
	reservoirEventRepo := repository.NewReservoirEventRepositoryInMemory(reservoirEventStorage)
	reservoirEventQuery := inmemory.NewReservoirEventQueryInMemory(reservoirEventStorage)
	reservoirReadRepo := repository.NewReservoirReadRepositoryInMemory(reservoirReadStorage)
	reservoirReadQuery := inmemory.NewReservoirReadQueryInMemory(reservoirReadStorage)
	reservoirService := service.ReservoirServiceInMemory{FarmReadQuery: farmReadQuery}

	materialRepo := repository.NewMaterialRepositoryInMemory(materialStorage)
	materialQuery := inmemory.NewMaterialQueryInMemory(materialStorage)

	cropQuery := inmemory.NewCropQueryInMemory(cropStorage)

	farmServer := FarmServer{
		FarmRepo:            farmRepo,
		FarmEventRepo:       farmEventRepo,
		FarmReadRepo:        farmReadRepo,
		FarmReadQuery:       farmReadQuery,
		ReservoirRepo:       reservoirRepo,
		ReservoirEventRepo:  reservoirEventRepo,
		ReservoirEventQuery: reservoirEventQuery,
		ReservoirReadRepo:   reservoirReadRepo,
		ReservoirReadQuery:  reservoirReadQuery,
		ReservoirService:    reservoirService,
		AreaRepo:            areaRepo,
		AreaQuery:           areaQuery,
		MaterialRepo:        materialRepo,
		MaterialQuery:       materialQuery,
		CropQuery:           cropQuery,
		File:                LocalFile{},
		EventBus:            eventBus,
	}

	farmServer.InitSubscriber()

	return &farmServer, nil
}

// InitSubscriber defines the mapping of which event this domain listen with their handler
func (s *FarmServer) InitSubscriber() {
	s.EventBus.Subscribe("FarmCreated", s.SaveToFarmReadModel)
	s.EventBus.Subscribe("ReservoirCreated", s.SaveToReservoirReadModel)
}

// Mount defines the FarmServer's endpoints with its handlers
func (s *FarmServer) Mount(g *echo.Group) {
	g.GET("/types", s.GetTypes)
	g.GET("/inventories/materials", s.GetMaterials)
	g.GET("/inventories/plant_types", s.GetInventoryPlantTypes)
	g.GET("/inventories/materials/available_plant_type", s.GetAvailableMaterialPlantType)
	g.POST("/inventories/materials/:type", s.SaveMaterial)

	g.POST("", s.SaveFarm)
	g.GET("", s.FindAllFarm)
	g.GET("/:id", s.FindFarmByID)

	g.POST("/:id/reservoirs", s.SaveReservoir)
	g.POST("/reservoirs/:id/notes", s.SaveReservoirNotes)
	g.DELETE("/reservoirs/:reservoir_id/notes/:note_id", s.RemoveReservoirNotes)
	g.GET("/:id/reservoirs", s.GetFarmReservoirs)
	g.GET("/:farm_id/reservoirs/:reservoir_id", s.GetReservoirsByID)

	g.POST("/:id/areas", s.SaveArea)
	g.POST("/areas/:id/notes", s.SaveAreaNotes)
	g.DELETE("/areas/:area_id/notes/:note_id", s.RemoveAreaNotes)
	g.GET("/:id/areas", s.GetFarmAreas)
	g.GET("/:farm_id/areas/:area_id", s.GetAreasByID)
	g.GET("/:farm_id/areas/:area_id/photos", s.GetAreaPhotos)
}

// GetTypes is a FarmServer's handle to get farm types
func (s *FarmServer) GetTypes(c echo.Context) error {
	types := domain.FindAllFarmTypes()

	return c.JSON(http.StatusOK, types)
}

func (s FarmServer) FindAllFarm(c echo.Context) error {
	result := <-s.FarmReadQuery.FindAll()
	if result.Error != nil {
		return result.Error
	}

	farms, ok := result.Result.([]storage.FarmRead)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	data := make(map[string][]storage.FarmRead)
	for _, v := range farms {
		data["data"] = append(data["data"], v)
	}

	return c.JSON(http.StatusOK, data)
}

// SaveFarm is a FarmServer's handler to save new Farm
func (s *FarmServer) SaveFarm(c echo.Context) error {
	farm, err := domain.CreateFarm(
		c.FormValue("name"),
		c.FormValue("farm_type"),
		c.FormValue("latitude"),
		c.FormValue("longitude"),
		c.FormValue("country_code"),
		c.FormValue("city_code"),
	)
	if err != nil {
		return Error(c, err)
	}

	err = <-s.FarmEventRepo.Save(farm.UID, farm.Version, farm.UncommittedChanges)
	if err != nil {
		return Error(c, err)
	}

	s.publishUncommittedEvents(farm)

	data := make(map[string]*storage.FarmRead)
	data["data"] = MapToFarmRead(farm)

	return c.JSON(http.StatusOK, data)
}

func (s *FarmServer) FindFarmByID(c echo.Context) error {
	farmUID, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return Error(c, err)
	}

	result := <-s.FarmReadQuery.FindByID(farmUID)
	if result.Error != nil {
		return result.Error
	}

	farm, ok := result.Result.(storage.FarmRead)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	data := make(map[string]storage.FarmRead)
	data["data"] = farm

	return c.JSON(http.StatusOK, data)
}

// SaveReservoir is a FarmServer's handler to save new Reservoir and place it to a Farm
func (s *FarmServer) SaveReservoir(c echo.Context) error {
	validation := RequestValidation{}

	// Validate requests //
	name, err := validation.ValidateReservoirName(c.FormValue("name"))
	if err != nil {
		return Error(c, err)
	}

	waterSourceType, err := validation.ValidateType(c.FormValue("type"))
	if err != nil {
		return Error(c, err)
	}

	capacity, err := validation.ValidateCapacity(waterSourceType, c.FormValue("capacity"))
	if err != nil {
		return Error(c, err)
	}

	farmUID, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return err
	}

	// Process //
	r, err := domain.CreateReservoir(s.ReservoirService, farmUID, name, waterSourceType, capacity)
	if err != nil {
		return Error(c, err)
	}

	// Persists //
	err = <-s.ReservoirEventRepo.Save(r.UID, r.Version, r.UncommittedChanges)
	if err != nil {
		return Error(c, err)
	}

	// Publish //
	s.publishUncommittedEvents(r)

	detailReservoir, err := MapToDetailReservoir(s, *r)
	if err != nil {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	data := make(map[string]DetailReservoir)
	data["data"] = detailReservoir

	return c.JSON(http.StatusOK, data)
}

func (s *FarmServer) SaveReservoirNotes(c echo.Context) error {
	reservoirUID, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return Error(c, err)
	}

	content := c.FormValue("content")

	// Validate //
	queryResult := <-s.ReservoirReadQuery.FindByID(reservoirUID)
	if queryResult.Error != nil {
		return Error(c, queryResult.Error)
	}

	reservoirRead, ok := queryResult.Result.(query.ReservoirReadQueryResult)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	if reservoirRead.UID == (uuid.UUID{}) {
		return Error(c, NewRequestValidationError(REQUIRED, "id"))
	}

	if content == "" {
		return Error(c, NewRequestValidationError(REQUIRED, "content"))
	}

	// Process //
	eventQueryResult := <-s.ReservoirEventQuery.FindAllByID(reservoirRead.UID)
	if eventQueryResult.Error != nil {
		return Error(c, eventQueryResult.Error)
	}

	events := eventQueryResult.Result.([]storage.ReservoirEvent)
	reservoir := repository.NewReservoirFromHistory(events)

	err = reservoir.AddNewNote(content)
	if err != nil {
		return Error(c, err)
	}

	// Persists //
	resultSave := <-s.ReservoirEventRepo.Save(reservoir.UID, reservoir.Version, reservoir.UncommittedChanges)
	if resultSave != nil {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	// Publish //
	s.publishUncommittedEvents(reservoir)

	detailReservoir, err := MapToDetailReservoir(s, *reservoir)
	if err != nil {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	data := make(map[string]DetailReservoir)
	data["data"] = detailReservoir

	return c.JSON(http.StatusOK, data)
}

func (s *FarmServer) RemoveReservoirNotes(c echo.Context) error {
	reservoirUID, err := uuid.FromString(c.Param("reservoir_id"))
	noteID := c.Param("note_id")

	// Validate //
	queryResult := <-s.ReservoirReadQuery.FindByID(reservoirUID)
	if queryResult.Error != nil {
		return Error(c, queryResult.Error)
	}

	reservoirRead, ok := queryResult.Result.(query.ReservoirReadQueryResult)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	// Process //
	eventQueryResult := <-s.ReservoirEventQuery.FindAllByID(reservoirRead.UID)
	if eventQueryResult.Error != nil {
		return Error(c, eventQueryResult.Error)
	}

	events := eventQueryResult.Result.([]storage.ReservoirEvent)
	reservoir := repository.NewReservoirFromHistory(events)

	err = reservoir.RemoveNote(noteID)
	if err != nil {
		return Error(c, err)
	}

	// Persists //
	err = <-s.ReservoirEventRepo.Save(reservoir.UID, reservoir.Version, reservoir.UncommittedChanges)
	if err != nil {
		return Error(c, err)
	}

	detailReservoir, err := MapToDetailReservoir(s, *reservoir)
	if err != nil {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	data := make(map[string]DetailReservoir)
	data["data"] = detailReservoir

	return c.JSON(http.StatusOK, data)
}

func (s *FarmServer) GetFarmReservoirs(c echo.Context) error {
	farmUID, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return Error(c, err)
	}

	result := <-s.ReservoirReadQuery.FindAllByFarm(farmUID)
	if result.Error != nil {
		return result.Error
	}

	reservoirReads, ok := result.Result.([]query.ReservoirReadQueryResult)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	reservoirs, err := MapToReservoirListFromQuery(s, reservoirReads)
	if err != nil {
		return Error(c, err)
	}

	data := make(map[string][]DetailReservoir)
	data["data"] = reservoirs
	if len(reservoirs) == 0 {
		data["data"] = []DetailReservoir{}
	}

	return c.JSON(http.StatusOK, data)
}

func (s *FarmServer) GetReservoirsByID(c echo.Context) error {
	reservoirUID, err := uuid.FromString(c.Param("reservoir_id"))
	if err != nil {
		return Error(c, err)
	}

	// Validate //
	result := <-s.ReservoirReadQuery.FindByID(reservoirUID)
	if result.Error != nil {
		return Error(c, result.Error)
	}

	reservoir, ok := result.Result.(query.ReservoirReadQueryResult)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusBadRequest, "Internal server error"))
	}

	if reservoir.UID == (uuid.UUID{}) {
		return Error(c, NewRequestValidationError(REQUIRED, "reservoir_id"))
	}

	data := make(map[string]DetailReservoir)
	detailReservoir, err := MapToReservoirFromQuery(s, reservoir)
	if err != nil {
		return Error(c, err)
	}

	data["data"] = detailReservoir

	return c.JSON(http.StatusOK, data)
}

func (s *FarmServer) SaveArea(c echo.Context) error {
	data := make(map[string]DetailArea)
	validation := RequestValidation{}

	// Validation //
	farmUID, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return Error(c, err)
	}

	// TODO: NEED TO BE FIXED
	farm := domain.Farm{UID: farmUID}

	reservoir, err := validation.ValidateReservoir(*s, c.FormValue("reservoir_id"))
	if err != nil {
		return Error(c, err)
	}

	size, err := validation.ValidateAreaSize(c.FormValue("size"), c.FormValue("size_unit"))
	if err != nil {
		return Error(c, err)
	}

	location, err := validation.ValidateAreaLocation(c.FormValue("location"))
	if err != nil {
		return Error(c, err)
	}

	// Process //
	area, err := domain.CreateArea(farm, c.FormValue("name"), c.FormValue("type"))
	if err != nil {
		return Error(c, err)
	}

	err = area.ChangeSize(size)
	if err != nil {
		return Error(c, err)
	}

	err = area.ChangeLocation(location)
	if err != nil {
		return Error(c, err)
	}

	photo, err := c.FormFile("photo")
	if err == nil {
		destPath := stringhelper.Join(*config.Config.UploadPathArea, "/", photo.Filename)
		err = s.File.Upload(photo, destPath)

		if err != nil {
			return Error(c, err)
		}

		width, height, err := imagehelper.GetImageDimension(destPath)
		if err != nil {
			return Error(c, err)
		}

		areaPhoto := domain.AreaPhoto{
			Filename: photo.Filename,
			MimeType: photo.Header["Content-Type"][0],
			Size:     int(photo.Size),
			Width:    width,
			Height:   height,
		}

		area.Photo = areaPhoto
	}

	area.Farm = farm
	area.Reservoir = reservoir

	err = farm.AddArea(&area)
	if err != nil {
		return Error(c, err)
	}

	// Persists //
	err = <-s.ReservoirRepo.Save(&reservoir)
	if err != nil {
		return Error(c, err)
	}

	err = <-s.AreaRepo.Save(&area)
	if err != nil {
		return Error(c, err)
	}

	err = <-s.FarmRepo.Save(&farm)
	if err != nil {
		return Error(c, err)
	}

	detailArea, err := MapToDetailArea(s, area)
	if err != nil {
		return Error(c, err)
	}

	data["data"] = detailArea

	return c.JSON(http.StatusOK, data)
}

func (s *FarmServer) SaveAreaNotes(c echo.Context) error {
	data := make(map[string]DetailArea)

	areaID := c.Param("id")
	content := c.FormValue("content")

	// Validate //
	result := <-s.AreaRepo.FindByID(areaID)
	if result.Error != nil {
		return Error(c, result.Error)
	}

	area, ok := result.Result.(domain.Area)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	result = <-s.FarmRepo.FindByID(area.Farm.UID.String())
	if result.Error != nil {
		return Error(c, result.Error)
	}

	farm, ok := result.Result.(domain.Farm)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	if content == "" {
		return Error(c, NewRequestValidationError(REQUIRED, "content"))
	}

	// Process //
	area.AddNewNote(content)
	farm.ChangeAreaInformation(area)

	// Persists //
	resultSave := <-s.AreaRepo.Save(&area)
	if resultSave != nil {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	resultSave = <-s.FarmRepo.Save(&farm)
	if resultSave != nil {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	detailArea, err := MapToDetailArea(s, area)
	if err != nil {
		return Error(c, err)
	}

	data["data"] = detailArea

	return c.JSON(http.StatusOK, data)
}

func (s *FarmServer) RemoveAreaNotes(c echo.Context) error {
	data := make(map[string]DetailArea)

	areaID := c.Param("area_id")
	noteID := c.Param("note_id")

	// Validate //
	result := <-s.AreaRepo.FindByID(areaID)
	if result.Error != nil {
		return Error(c, result.Error)
	}

	area, ok := result.Result.(domain.Area)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	result = <-s.FarmRepo.FindByID(area.Farm.UID.String())
	if result.Error != nil {
		return Error(c, result.Error)
	}

	farm, ok := result.Result.(domain.Farm)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	// Process //
	err := area.RemoveNote(noteID)
	if err != nil {
		return Error(c, err)
	}

	err = farm.ChangeAreaInformation(area)
	if err != nil {
		return Error(c, err)
	}

	// Persists //
	resultSave := <-s.AreaRepo.Save(&area)
	if resultSave != nil {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	resultSave = <-s.FarmRepo.Save(&farm)
	if resultSave != nil {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	detailArea, err := MapToDetailArea(s, area)
	if err != nil {
		return Error(c, err)
	}

	data["data"] = detailArea

	return c.JSON(http.StatusOK, data)
}

func (s *FarmServer) GetFarmAreas(c echo.Context) error {
	data := make(map[string][]AreaList)

	result := <-s.FarmRepo.FindByID(c.Param("id"))
	if result.Error != nil {
		return Error(c, result.Error)
	}

	farm, ok := result.Result.(domain.Farm)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusBadRequest, "Internal server error"))
	}

	areaList, err := MapToAreaList(s, farm.Areas)
	if err != nil {
		return Error(c, err)
	}

	data["data"] = areaList

	return c.JSON(http.StatusOK, data)
}

func (s *FarmServer) GetAreasByID(c echo.Context) error {
	data := make(map[string]DetailArea)

	// Validate //
	result := <-s.FarmRepo.FindByID(c.Param("farm_id"))
	if result.Error != nil {
		return Error(c, result.Error)
	}

	_, ok := result.Result.(domain.Farm)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusBadRequest, "Internal server error"))
	}

	result = <-s.AreaRepo.FindByID(c.Param("area_id"))
	if result.Error != nil {
		return Error(c, result.Error)
	}

	area, ok := result.Result.(domain.Area)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusBadRequest, "Internal server error"))
	}

	detailArea, err := MapToDetailArea(s, area)
	if err != nil {
		return Error(c, err)
	}

	data["data"] = detailArea

	return c.JSON(http.StatusOK, data)
}

func (s *FarmServer) GetAreaPhotos(c echo.Context) error {
	// Validate //
	result := <-s.FarmRepo.FindByID(c.Param("farm_id"))
	if result.Error != nil {
		return Error(c, result.Error)
	}

	_, ok := result.Result.(domain.Farm)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusBadRequest, "Internal server error"))
	}

	result = <-s.AreaRepo.FindByID(c.Param("area_id"))
	if result.Error != nil {
		return Error(c, result.Error)
	}

	area, ok := result.Result.(domain.Area)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusBadRequest, "Internal server error"))
	}

	if area.Photo.Filename == "" {
		return Error(c, NewRequestValidationError(NOT_FOUND, "photo"))
	}

	// Process //
	srcPath := stringhelper.Join(*config.Config.UploadPathArea, "/", area.Photo.Filename)

	return c.File(srcPath)
}

func (s *FarmServer) GetInventoryPlantTypes(c echo.Context) error {
	data := make(map[string][]string)

	plantTypes := MapToPlantType(domain.PlantTypes())

	data["data"] = plantTypes

	return c.JSON(http.StatusOK, data)
}

func (s *FarmServer) GetMaterials(c echo.Context) error {
	data := make(map[string][]Material)

	queryResult := <-s.MaterialQuery.FindAll()
	if queryResult.Error != nil {
		return Error(c, queryResult.Error)
	}

	results, ok := queryResult.Result.([]domain.Material)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	materials := []Material{}
	for _, v := range results {
		materials = append(materials, MapToMaterial(v))
	}

	data["data"] = materials

	return c.JSON(http.StatusOK, data)
}

func (s *FarmServer) SaveMaterial(c echo.Context) error {
	data := make(map[string]Material)

	materialTypeParam := c.Param("type")
	name := c.FormValue("name")

	plantType := c.FormValue("plant_type")
	chemicalType := c.FormValue("chemical_type")
	containerType := c.FormValue("container_type")

	pricePerUnit := c.FormValue("price_per_unit")
	currencyCode := c.FormValue("currency_code")
	quantity := c.FormValue("quantity")
	quantityUnit := c.FormValue("quantity_unit")
	expirationDate := c.FormValue("expiration_date")
	notes := c.FormValue("notes")
	isExpense := c.FormValue("is_expense")
	producedBy := c.FormValue("produced_by")

	// Validate //
	q, err := strconv.ParseFloat(quantity, 32)
	if err != nil {
		return Error(c, NewRequestValidationError(INVALID_OPTION, "quantity"))
	}

	var expDate *time.Time
	if expirationDate != "" {
		tp, err := time.Parse("2006-01-02", expirationDate)
		if err != nil {
			return Error(c, NewRequestValidationError(PARSE_FAILED, "expiration_date"))
		}

		expDate = &tp
	}

	var n *string
	if notes != "" {
		n = &notes
	}

	var pb *string
	if producedBy != "" {
		pb = &producedBy
	}

	var isExpenseBool *bool
	if isExpense != "" && isExpense == "true" {
		b := true
		isExpenseBool = &b
	}

	// Process //
	var mt domain.MaterialType
	switch materialTypeParam {
	case strings.ToLower(domain.MaterialTypeSeedCode):
		pt := domain.GetPlantType(plantType)
		if pt == (domain.PlantType{}) {
			return Error(c, NewRequestValidationError(INVALID_OPTION, "plant_type"))
		}

		mt, err = domain.CreateMaterialTypeSeed(pt.Code)
		if err != nil {
			return Error(c, NewRequestValidationError(INVALID_OPTION, "type"))
		}
	case strings.ToLower(domain.MaterialTypeAgrochemicalCode):
		ct := domain.GetChemicalType(chemicalType)
		if ct == (domain.ChemicalType{}) {
			return Error(c, NewRequestValidationError(INVALID_OPTION, "chemical_type"))
		}

		mt, err = domain.CreateMaterialTypeAgrochemical(ct.Code)
		if err != nil {
			return Error(c, NewRequestValidationError(INVALID_OPTION, "type"))
		}
	case strings.ToLower(domain.MaterialTypeGrowingMediumCode):
		mt = domain.MaterialTypeGrowingMedium{}
	case strings.ToLower(domain.MaterialTypeLabelAndCropSupportCode):
		mt = domain.MaterialTypeLabelAndCropSupport{}
	case strings.ToLower(domain.MaterialTypeSeedingContainerCode):
		ct := domain.GetContainerType(containerType)
		if ct == (domain.ContainerType{}) {
			return Error(c, NewRequestValidationError(INVALID_OPTION, "container_type"))
		}

		mt, err = domain.CreateMaterialTypeSeedingContainer(ct.Code)
		if err != nil {
			return Error(c, NewRequestValidationError(INVALID_OPTION, "type"))
		}
	case strings.ToLower(domain.MaterialTypePostHarvestSupplyCode):
		mt = domain.MaterialTypePostHarvestSupply{}
	case strings.ToLower(domain.MaterialTypeOtherCode):
		mt = domain.MaterialTypeOther{}
	case strings.ToLower(domain.MaterialTypePlantCode):
		pt := domain.GetPlantType(plantType)
		if pt == (domain.PlantType{}) {
			return Error(c, NewRequestValidationError(INVALID_OPTION, "plant_type"))
		}

		mt, err = domain.CreateMaterialTypePlant(pt.Code)
		if err != nil {
			return Error(c, NewRequestValidationError(INVALID_OPTION, "type"))
		}
	}

	material, err := domain.CreateMaterial(name, pricePerUnit, currencyCode, mt, float32(q), quantityUnit)
	if err != nil {
		return Error(c, err)
	}

	material.ExpirationDate = expDate
	material.Notes = n
	material.ProducedBy = pb
	material.IsExpense = isExpenseBool

	// Persist //
	err = <-s.MaterialRepo.Save(&material)
	if err != nil {
		return Error(c, err)
	}

	data["data"] = MapToMaterial(material)

	return c.JSON(http.StatusOK, data)
}

func (s *FarmServer) GetAvailableMaterialPlantType(c echo.Context) error {
	data := make(map[string][]AvailableMaterialPlantType)

	// Process //
	result := <-s.MaterialRepo.FindAll()

	materials, ok := result.Result.([]domain.Material)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusBadRequest, "Internal server error"))
	}

	data["data"] = MapToAvailableMaterialPlantType(materials)

	return c.JSON(http.StatusOK, data)
}

func (s *FarmServer) publishUncommittedEvents(entity interface{}) error {
	switch e := entity.(type) {
	case *domain.Farm:
		for _, v := range e.UncommittedChanges {
			name := structhelper.GetName(v)
			s.EventBus.Publish(name, v)
		}
	case *domain.Reservoir:
		for _, v := range e.UncommittedChanges {
			name := structhelper.GetName(v)
			s.EventBus.Publish(name, v)
		}
	}

	return nil
}
