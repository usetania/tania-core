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
	FarmEventRepo       repository.FarmEventRepository
	FarmReadRepo        repository.FarmReadRepository
	FarmReadQuery       query.FarmReadQuery
	ReservoirEventRepo  repository.ReservoirEventRepository
	ReservoirEventQuery query.ReservoirEventQuery
	ReservoirReadRepo   repository.ReservoirReadRepository
	ReservoirReadQuery  query.ReservoirReadQuery
	ReservoirService    domain.ReservoirService
	AreaEventRepo       repository.AreaEventRepository
	AreaReadRepo        repository.AreaReadRepository
	AreaEventQuery      query.AreaEventQuery
	AreaReadQuery       query.AreaReadQuery
	AreaService         domain.AreaService
	MaterialEventRepo   repository.MaterialEventRepository
	MaterialEventQuery  query.MaterialEventQuery
	MaterialReadRepo    repository.MaterialReadRepository
	MaterialReadQuery   query.MaterialReadQuery
	CropReadQuery       query.CropReadQuery
	File                File
	EventBus            EventBus.Bus
}

// NewFarmServer initializes FarmServer's dependencies and create new FarmServer struct
func NewFarmServer(
	farmEventStorage *storage.FarmEventStorage,
	farmReadStorage *storage.FarmReadStorage,
	areaEventStorage *storage.AreaEventStorage,
	areaReadStorage *storage.AreaReadStorage,
	reservoirEventStorage *storage.ReservoirEventStorage,
	reservoirReadStorage *storage.ReservoirReadStorage,
	materialEventStorage *storage.MaterialEventStorage,
	materialReadStorage *storage.MaterialReadStorage,
	cropReadStorage *growthstorage.CropReadStorage,
	eventBus EventBus.Bus,
) (*FarmServer, error) {
	farmEventRepo := repository.NewFarmEventRepositoryInMemory(farmEventStorage)
	farmReadRepo := repository.NewFarmReadRepositoryInMemory(farmReadStorage)
	farmReadQuery := inmemory.NewFarmReadQueryInMemory(farmReadStorage)

	areaEventRepo := repository.NewAreaEventRepositoryInMemory(areaEventStorage)
	areaReadRepo := repository.NewAreaReadRepositoryInMemory(areaReadStorage)
	areaEventQuery := inmemory.NewAreaEventQueryInMemory(areaEventStorage)
	areaReadQuery := inmemory.NewAreaReadQueryInMemory(areaReadStorage)

	reservoirEventRepo := repository.NewReservoirEventRepositoryInMemory(reservoirEventStorage)
	reservoirEventQuery := inmemory.NewReservoirEventQueryInMemory(reservoirEventStorage)
	reservoirReadRepo := repository.NewReservoirReadRepositoryInMemory(reservoirReadStorage)
	reservoirReadQuery := inmemory.NewReservoirReadQueryInMemory(reservoirReadStorage)

	materialEventRepo := repository.NewMaterialEventRepositoryInMemory(materialEventStorage)
	materialEventQuery := inmemory.NewMaterialEventQueryInMemory(materialEventStorage)
	materialReadRepo := repository.NewMaterialReadRepositoryInMemory(materialReadStorage)
	materialReadQuery := inmemory.NewMaterialReadQueryInMemory(materialReadStorage)

	cropReadQuery := inmemory.NewCropReadQueryInMemory(cropReadStorage)

	areaService := service.AreaServiceInMemory{
		FarmReadQuery:      farmReadQuery,
		ReservoirReadQuery: reservoirReadQuery,
	}
	reservoirService := service.ReservoirServiceInMemory{FarmReadQuery: farmReadQuery}

	farmServer := FarmServer{
		FarmEventRepo:       farmEventRepo,
		FarmReadRepo:        farmReadRepo,
		FarmReadQuery:       farmReadQuery,
		ReservoirEventRepo:  reservoirEventRepo,
		ReservoirEventQuery: reservoirEventQuery,
		ReservoirReadRepo:   reservoirReadRepo,
		ReservoirReadQuery:  reservoirReadQuery,
		ReservoirService:    reservoirService,
		AreaEventRepo:       areaEventRepo,
		AreaReadRepo:        areaReadRepo,
		AreaEventQuery:      areaEventQuery,
		AreaReadQuery:       areaReadQuery,
		AreaService:         areaService,
		MaterialEventRepo:   materialEventRepo,
		MaterialEventQuery:  materialEventQuery,
		MaterialReadRepo:    materialReadRepo,
		MaterialReadQuery:   materialReadQuery,
		CropReadQuery:       cropReadQuery,
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
	s.EventBus.Subscribe("ReservoirNoteAdded", s.SaveToReservoirReadModel)
	s.EventBus.Subscribe("ReservoirNoteRemoved", s.SaveToReservoirReadModel)

	s.EventBus.Subscribe("AreaCreated", s.SaveToAreaReadModel)
	s.EventBus.Subscribe("AreaPhotoAdded", s.SaveToAreaReadModel)
	s.EventBus.Subscribe("AreaNoteAdded", s.SaveToAreaReadModel)
	s.EventBus.Subscribe("AreaNoteRemoved", s.SaveToAreaReadModel)

	s.EventBus.Subscribe("MaterialCreated", s.SaveToMaterialReadModel)
	s.EventBus.Subscribe("MaterialNameChanged", s.SaveToMaterialReadModel)
	s.EventBus.Subscribe("MaterialPriceChanged", s.SaveToMaterialReadModel)
	s.EventBus.Subscribe("MaterialQuantityChanged", s.SaveToMaterialReadModel)
	s.EventBus.Subscribe("MaterialTypeChanged", s.SaveToMaterialReadModel)
	s.EventBus.Subscribe("MaterialExpirationDateChanged", s.SaveToMaterialReadModel)
	s.EventBus.Subscribe("MaterialNotesChanged", s.SaveToMaterialReadModel)
	s.EventBus.Subscribe("MaterialProducedByChanged", s.SaveToMaterialReadModel)

}

// Mount defines the FarmServer's endpoints with its handlers
func (s *FarmServer) Mount(g *echo.Group) {
	g.GET("/types", s.GetTypes)
	g.GET("/inventories/materials", s.GetMaterials)
	g.GET("/inventories/plant_types", s.GetInventoryPlantTypes)
	g.GET("/inventories/materials/available_plant_type", s.GetAvailableMaterialPlantType)
	g.POST("/inventories/materials/:type", s.SaveMaterial)
	g.PUT("/inventories/materials/:type/:id", s.UpdateMaterial)

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
	g.GET("/:id/areas/total", s.GetTotalAreas)
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
	data["data"] = farms
	if len(farms) == 0 {
		data["data"] = []storage.FarmRead{}
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

	resRead, err := MapToReservoirRead(s, *r)
	if err != nil {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	data := make(map[string]storage.ReservoirRead)
	data["data"] = resRead

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

	reservoirRead, ok := queryResult.Result.(storage.ReservoirRead)
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

	resRead, err := MapToReservoirRead(s, *reservoir)
	if err != nil {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	data := make(map[string]storage.ReservoirRead)
	data["data"] = resRead

	return c.JSON(http.StatusOK, data)
}

func (s *FarmServer) RemoveReservoirNotes(c echo.Context) error {
	reservoirUID, err := uuid.FromString(c.Param("reservoir_id"))
	if err != nil {
		return Error(c, err)
	}

	noteUID, err := uuid.FromString(c.Param("note_id"))
	if err != nil {
		return Error(c, err)
	}

	// Validate //
	queryResult := <-s.ReservoirReadQuery.FindByID(reservoirUID)
	if queryResult.Error != nil {
		return Error(c, queryResult.Error)
	}

	reservoirRead, ok := queryResult.Result.(storage.ReservoirRead)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	if reservoirRead.UID == (uuid.UUID{}) {
		return Error(c, NewRequestValidationError(NOT_FOUND, "reservoir_id"))
	}

	noteFound := false
	for _, v := range reservoirRead.Notes {
		if v.UID == noteUID {
			noteFound = true
		}
	}

	if !noteFound {
		return Error(c, NewRequestValidationError(NOT_FOUND, "note_id"))
	}

	// Process //
	eventQueryResult := <-s.ReservoirEventQuery.FindAllByID(reservoirRead.UID)
	if eventQueryResult.Error != nil {
		return Error(c, eventQueryResult.Error)
	}

	events := eventQueryResult.Result.([]storage.ReservoirEvent)
	reservoir := repository.NewReservoirFromHistory(events)

	err = reservoir.RemoveNote(noteUID)
	if err != nil {
		return Error(c, err)
	}

	// Persists //
	err = <-s.ReservoirEventRepo.Save(reservoir.UID, reservoir.Version, reservoir.UncommittedChanges)
	if err != nil {
		return Error(c, err)
	}

	resRead, err := MapToReservoirRead(s, *reservoir)
	if err != nil {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	data := make(map[string]storage.ReservoirRead)
	data["data"] = resRead

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

	reservoirs, ok := result.Result.([]storage.ReservoirRead)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	data := make(map[string][]storage.ReservoirRead)
	for _, v := range reservoirs {
		data["data"] = append(data["data"], MapToReservoirReadFromRead(v))
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

	reservoir, ok := result.Result.(storage.ReservoirRead)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusBadRequest, "Internal server error"))
	}

	if reservoir.UID == (uuid.UUID{}) {
		return Error(c, NewRequestValidationError(REQUIRED, "reservoir_id"))
	}

	data := make(map[string]storage.ReservoirRead)
	data["data"] = MapToReservoirReadFromRead(reservoir)

	return c.JSON(http.StatusOK, data)
}

func (s *FarmServer) SaveArea(c echo.Context) error {
	validation := RequestValidation{}

	farmUID, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return Error(c, err)
	}

	reservoirUID, err := uuid.FromString(c.FormValue("reservoir_id"))
	if err != nil {
		return Error(c, err)
	}

	// Validation //
	reservoir, err := validation.ValidateReservoir(*s, reservoirUID)
	if err != nil {
		return Error(c, err)
	}

	farm, err := validation.ValidateFarm(*s, farmUID)
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
	area, err := domain.CreateArea(s.AreaService, farm.UID, reservoir.UID, c.FormValue("name"), c.FormValue("type"), size, location)
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

		area.ChangePhoto(areaPhoto)
	}

	// Persists //
	err = <-s.AreaEventRepo.Save(area.UID, area.Version, area.UncommittedChanges)
	if err != nil {
		return Error(c, err)
	}

	// Publish //
	s.publishUncommittedEvents(area)

	data := make(map[string]DetailArea)
	detailArea, err := MapToDetailArea(s, *area)
	if err != nil {
		return Error(c, err)
	}

	data["data"] = detailArea

	return c.JSON(http.StatusOK, data)
}

func (s *FarmServer) SaveAreaNotes(c echo.Context) error {
	areaUID, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return Error(c, err)
	}
	content := c.FormValue("content")

	// Validate //
	queryResult := <-s.AreaReadQuery.FindByID(areaUID)
	if queryResult.Error != nil {
		return Error(c, queryResult.Error)
	}

	areaRead, ok := queryResult.Result.(storage.AreaRead)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	if areaRead.UID == (uuid.UUID{}) {
		return Error(c, NewRequestValidationError(NOT_FOUND, "id"))
	}

	if content == "" {
		return Error(c, NewRequestValidationError(REQUIRED, "content"))
	}

	// Process //
	eventQueryResult := <-s.AreaEventQuery.FindAllByID(areaRead.UID)
	if eventQueryResult.Error != nil {
		return Error(c, eventQueryResult.Error)
	}

	events := eventQueryResult.Result.([]storage.AreaEvent)
	area := repository.NewAreaFromHistory(events)

	err = area.AddNewNote(content)
	if err != nil {
		return Error(c, err)
	}

	// Persists //
	err = <-s.AreaEventRepo.Save(area.UID, area.Version, area.UncommittedChanges)
	if err != nil {
		return Error(c, err)
	}

	// Publish //
	s.publishUncommittedEvents(area)

	detailArea, err := MapToDetailArea(s, *area)
	if err != nil {
		return Error(c, err)
	}

	data := make(map[string]DetailArea)
	data["data"] = detailArea

	return c.JSON(http.StatusOK, data)
}

func (s *FarmServer) RemoveAreaNotes(c echo.Context) error {
	data := make(map[string]DetailArea)

	areaUID, err := uuid.FromString(c.Param("area_id"))
	if err != nil {
		return Error(c, err)
	}

	noteUID, err := uuid.FromString(c.Param("note_id"))
	if err != nil {
		return Error(c, err)
	}

	// Validate //
	queryResult := <-s.AreaReadQuery.FindByID(areaUID)
	if queryResult.Error != nil {
		return Error(c, queryResult.Error)
	}

	areaRead, ok := queryResult.Result.(storage.AreaRead)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	if areaRead.UID == (uuid.UUID{}) {
		return Error(c, NewRequestValidationError(NOT_FOUND, "area_id"))
	}

	found := false
	for _, v := range areaRead.Notes {
		if v.UID == noteUID {
			found = true
		}
	}

	if !found {
		return Error(c, NewRequestValidationError(NOT_FOUND, "note_id"))
	}

	// // Process //
	eventQueryResult := <-s.AreaEventQuery.FindAllByID(areaRead.UID)
	if eventQueryResult.Error != nil {
		return Error(c, eventQueryResult.Error)
	}

	events := eventQueryResult.Result.([]storage.AreaEvent)
	area := repository.NewAreaFromHistory(events)

	err = area.RemoveNote(noteUID)
	if err != nil {
		return Error(c, err)
	}

	// Persists //
	resultSave := <-s.AreaEventRepo.Save(area.UID, area.Version, area.UncommittedChanges)
	if resultSave != nil {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	// Publish //
	s.publishUncommittedEvents(area)

	detailArea, err := MapToDetailArea(s, *area)
	if err != nil {
		return Error(c, err)
	}

	data["data"] = detailArea

	return c.JSON(http.StatusOK, data)
}

func (s *FarmServer) GetFarmAreas(c echo.Context) error {
	farmUID, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return Error(c, err)
	}

	queryResult := <-s.AreaReadQuery.FindAllByFarm(farmUID)
	if queryResult.Error != nil {
		return Error(c, queryResult.Error)
	}

	areas, ok := queryResult.Result.([]storage.AreaRead)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusBadRequest, "Internal server error"))
	}

	areaList, err := MapToAreaList(s, areas)
	if err != nil {
		return Error(c, err)
	}

	data := make(map[string][]AreaList)
	data["data"] = areaList

	return c.JSON(http.StatusOK, data)
}

func (s *FarmServer) GetAreasByID(c echo.Context) error {
	// Validate //
	farmUID, err := uuid.FromString(c.Param("farm_id"))
	if err != nil {
		return Error(c, err)
	}

	areaUID, err := uuid.FromString(c.Param("area_id"))
	if err != nil {
		return Error(c, err)
	}

	queryResult := <-s.FarmReadQuery.FindByID(farmUID)
	if queryResult.Error != nil {
		return Error(c, queryResult.Error)
	}

	farmRead, ok := queryResult.Result.(storage.FarmRead)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusBadRequest, "Internal server error"))
	}

	if farmRead.UID == (uuid.UUID{}) {
		return Error(c, NewRequestValidationError(NOT_FOUND, "farm_id"))
	}

	queryResult = <-s.AreaReadQuery.FindByIDAndArea(areaUID, farmUID)
	if queryResult.Error != nil {
		return Error(c, queryResult.Error)
	}

	areaRead, ok := queryResult.Result.(storage.AreaRead)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusBadRequest, "Internal server error"))
	}

	if areaRead.UID == (uuid.UUID{}) {
		return Error(c, NewRequestValidationError(NOT_FOUND, "area_id"))
	}

	detailArea, err := MapToDetailAreaFromStorage(s, areaRead)
	if err != nil {
		return Error(c, err)
	}

	data := make(map[string]DetailArea)
	data["data"] = detailArea

	return c.JSON(http.StatusOK, data)
}

func (s *FarmServer) GetAreaPhotos(c echo.Context) error {
	// Validate //
	farmUID, err := uuid.FromString(c.Param("farm_id"))
	if err != nil {
		return Error(c, err)
	}

	areaUID, err := uuid.FromString(c.Param("area_id"))
	if err != nil {
		return Error(c, err)
	}

	queryResult := <-s.FarmReadQuery.FindByID(farmUID)
	if queryResult.Error != nil {
		return Error(c, queryResult.Error)
	}

	farmRead, ok := queryResult.Result.(storage.FarmRead)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusBadRequest, "Internal server error"))
	}

	if farmRead.UID == (uuid.UUID{}) {
		return Error(c, NewRequestValidationError(NOT_FOUND, "farm_id"))
	}

	queryResult = <-s.AreaReadQuery.FindByIDAndArea(areaUID, farmUID)
	if queryResult.Error != nil {
		return Error(c, queryResult.Error)
	}

	areaRead, ok := queryResult.Result.(storage.AreaRead)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusBadRequest, "Internal server error"))
	}

	if areaRead.UID == (uuid.UUID{}) {
		return Error(c, NewRequestValidationError(NOT_FOUND, "area_id"))
	}

	if areaRead.Photo.Filename == "" {
		return Error(c, NewRequestValidationError(NOT_FOUND, "photo"))
	}

	// Process //
	srcPath := stringhelper.Join(*config.Config.UploadPathArea, "/", areaRead.Photo.Filename)

	return c.File(srcPath)
}

func (s *FarmServer) GetTotalAreas(c echo.Context) error {
	farmUID, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return Error(c, err)
	}

	queryResult := <-s.AreaReadQuery.CountAreas(farmUID)
	if queryResult.Error != nil {
		return Error(c, queryResult.Error)
	}

	totalArea, ok := queryResult.Result.(int)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusBadRequest, "Internal server error"))
	}

	data := make(map[string]int)
	data["data"] = totalArea

	return c.JSON(http.StatusOK, data)
}

func (s *FarmServer) GetInventoryPlantTypes(c echo.Context) error {
	data := make(map[string][]string)

	plantTypes := MapToPlantType(domain.PlantTypes())

	data["data"] = plantTypes

	return c.JSON(http.StatusOK, data)
}

func (s *FarmServer) GetMaterials(c echo.Context) error {
	queryResult := <-s.MaterialReadQuery.FindAll()
	if queryResult.Error != nil {
		return Error(c, queryResult.Error)
	}

	results, ok := queryResult.Result.([]storage.MaterialRead)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	materials := []Material{}
	for _, v := range results {
		materials = append(materials, MapToMaterialFromRead(v))
	}

	data := make(map[string][]Material)
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

	material, err := domain.CreateMaterial(
		name, pricePerUnit, currencyCode, mt, float32(q), quantityUnit,
		expDate, n, pb)
	if err != nil {
		return Error(c, err)
	}

	// Persist //
	err = <-s.MaterialEventRepo.Save(material.UID, material.Version, material.UncommittedChanges)
	if err != nil {
		return Error(c, err)
	}

	// Publish //
	s.publishUncommittedEvents(material)

	data["data"] = MapToMaterial(*material)

	return c.JSON(http.StatusOK, data)
}

func (s *FarmServer) UpdateMaterial(c echo.Context) error {
	data := make(map[string]Material)

	materialUID, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return Error(c, err)
	}

	materialTypeParam := c.Param("type")

	plantType := c.FormValue("plant_type")
	chemicalType := c.FormValue("chemical_type")
	containerType := c.FormValue("container_type")

	name := c.FormValue("name")
	pricePerUnit := c.FormValue("price_per_unit")
	currencyCode := c.FormValue("currency_code")
	quantity := c.FormValue("quantity")
	quantityUnit := c.FormValue("quantity_unit")
	expirationDate := c.FormValue("expiration_date")
	notes := c.FormValue("notes")
	producedBy := c.FormValue("produced_by")

	// Validate //
	if pricePerUnit != "" && currencyCode == "" {
		return Error(c, NewRequestValidationError(REQUIRED, "currency_code"))
	}

	if currencyCode != "" && pricePerUnit == "" {
		return Error(c, NewRequestValidationError(REQUIRED, "price_per_unit"))
	}

	if quantity != "" && quantityUnit == "" {
		return Error(c, NewRequestValidationError(REQUIRED, "quantity_unit"))
	}

	if quantityUnit != "" && quantity == "" {
		return Error(c, NewRequestValidationError(REQUIRED, "quantity"))
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

	queryResult := <-s.MaterialReadQuery.FindByID(materialUID)
	if queryResult.Error != nil {
		return Error(c, queryResult.Error)
	}

	materialRead, ok := queryResult.Result.(storage.MaterialRead)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	// Process //
	var mt domain.MaterialType
	switch materialTypeParam {
	case strings.ToLower(domain.MaterialTypeSeedCode):
		if plantType != "" {
			pt := domain.GetPlantType(plantType)
			if pt == (domain.PlantType{}) {
				return Error(c, NewRequestValidationError(INVALID_OPTION, "plant_type"))
			}

			mt, err = domain.CreateMaterialTypeSeed(pt.Code)
			if err != nil {
				return Error(c, NewRequestValidationError(INVALID_OPTION, "type"))
			}

			materialRead.Type = mt
		}
	case strings.ToLower(domain.MaterialTypeAgrochemicalCode):
		if chemicalType != "" {
			ct := domain.GetChemicalType(chemicalType)
			if ct == (domain.ChemicalType{}) {
				return Error(c, NewRequestValidationError(INVALID_OPTION, "chemical_type"))
			}

			mt, err = domain.CreateMaterialTypeAgrochemical(ct.Code)
			if err != nil {
				return Error(c, NewRequestValidationError(INVALID_OPTION, "type"))
			}

			materialRead.Type = mt
		}
	case strings.ToLower(domain.MaterialTypeSeedingContainerCode):
		if containerType != "" {
			ct := domain.GetContainerType(containerType)
			if ct == (domain.ContainerType{}) {
				return Error(c, NewRequestValidationError(INVALID_OPTION, "container_type"))
			}

			mt, err = domain.CreateMaterialTypeSeedingContainer(ct.Code)
			if err != nil {
				return Error(c, NewRequestValidationError(INVALID_OPTION, "type"))
			}

			materialRead.Type = mt
		}
	case strings.ToLower(domain.MaterialTypePlantCode):
		if plantType != "" {
			pt := domain.GetPlantType(plantType)
			if pt == (domain.PlantType{}) {
				return Error(c, NewRequestValidationError(INVALID_OPTION, "plant_type"))
			}

			mt, err = domain.CreateMaterialTypePlant(pt.Code)
			if err != nil {
				return Error(c, NewRequestValidationError(INVALID_OPTION, "type"))
			}

			materialRead.Type = mt
		}
	}

	eventQueryResult := <-s.MaterialEventQuery.FindAllByID(materialRead.UID)
	if eventQueryResult.Error != nil {
		return Error(c, eventQueryResult.Error)
	}

	events := eventQueryResult.Result.([]storage.MaterialEvent)
	material := repository.NewMaterialFromHistory(events)

	if name != "" {
		material.ChangeName(name)
	}

	if mt != nil {
		material.ChangeType(mt)
	}

	if pricePerUnit != "" && currencyCode != "" {
		material.ChangePricePerUnit(pricePerUnit, currencyCode)
	}

	if quantity != "" && quantityUnit != "" {
		q, err := strconv.ParseFloat(quantity, 32)
		if err != nil {
			return Error(c, err)
		}

		material.ChangeQuantityUnit(float32(q), quantityUnit, materialRead.Type)
	}

	if expDate != nil {
		material.ChangeExpirationDate(*expDate)
	}

	if n != nil {
		material.ChangeNotes(*n)
	}

	if pb != nil {
		material.ChangeProducedBy(*pb)
	}

	// Persist //
	err = <-s.MaterialEventRepo.Save(material.UID, material.Version, material.UncommittedChanges)
	if err != nil {
		return Error(c, err)
	}

	// Publish //
	s.publishUncommittedEvents(material)

	data["data"] = MapToMaterial(*material)

	return c.JSON(http.StatusOK, data)
}

func (s *FarmServer) GetAvailableMaterialPlantType(c echo.Context) error {
	data := make(map[string][]AvailableMaterialPlantType)

	// Process //
	result := <-s.MaterialReadQuery.FindAll()

	materials, ok := result.Result.([]storage.MaterialRead)
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
	case *domain.Area:
		for _, v := range e.UncommittedChanges {
			name := structhelper.GetName(v)
			s.EventBus.Publish(name, v)
		}
	case *domain.Material:
		for _, v := range e.UncommittedChanges {
			name := structhelper.GetName(v)
			s.EventBus.Publish(name, v)
		}
	}

	return nil
}
