package server

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/usetania/tania-core/config"
	"github.com/usetania/tania-core/src/assets/domain"
	"github.com/usetania/tania-core/src/assets/domain/service"
	"github.com/usetania/tania-core/src/assets/query"
	queryInMem "github.com/usetania/tania-core/src/assets/query/inmemory"
	queryMysql "github.com/usetania/tania-core/src/assets/query/mysql"
	querySqlite "github.com/usetania/tania-core/src/assets/query/sqlite"
	"github.com/usetania/tania-core/src/assets/repository"
	repoInMem "github.com/usetania/tania-core/src/assets/repository/inmemory"
	repoMysql "github.com/usetania/tania-core/src/assets/repository/mysql"
	repoSqlite "github.com/usetania/tania-core/src/assets/repository/sqlite"
	"github.com/usetania/tania-core/src/assets/storage"
	"github.com/usetania/tania-core/src/eventbus"
	growthstorage "github.com/usetania/tania-core/src/growth/storage"
	"github.com/usetania/tania-core/src/helper/imagehelper"
	"github.com/usetania/tania-core/src/helper/paginationhelper"
	"github.com/usetania/tania-core/src/helper/stringhelper"
	"github.com/usetania/tania-core/src/helper/structhelper"
)

// FarmServer ties the routes and handlers with injected dependencies
type FarmServer struct {
	FarmEventRepo       repository.FarmEvent
	FarmEventQuery      query.FarmEvent
	FarmReadRepo        repository.FarmRead
	FarmReadQuery       query.FarmRead
	ReservoirEventRepo  repository.ReservoirEvent
	ReservoirEventQuery query.ReservoirEvent
	ReservoirReadRepo   repository.ReservoirRead
	ReservoirReadQuery  query.ReservoirRead
	ReservoirService    domain.ReservoirService
	AreaEventRepo       repository.AreaEvent
	AreaReadRepo        repository.AreaRead
	AreaEventQuery      query.AreaEvent
	AreaReadQuery       query.AreaRead
	AreaService         domain.AreaService
	MaterialEventRepo   repository.MaterialEvent
	MaterialEventQuery  query.MaterialEvent
	MaterialReadRepo    repository.MaterialRead
	MaterialReadQuery   query.MaterialRead
	CropReadQuery       query.CropRead
	File                File
	EventBus            eventbus.TaniaEventBus
}

// NewFarmServer initializes FarmServer's dependencies and create new FarmServer struct
func NewFarmServer(
	db *sql.DB,
	farmEventStorage *storage.FarmEventStorage,
	farmReadStorage *storage.FarmReadStorage,
	areaEventStorage *storage.AreaEventStorage,
	areaReadStorage *storage.AreaReadStorage,
	reservoirEventStorage *storage.ReservoirEventStorage,
	reservoirReadStorage *storage.ReservoirReadStorage,
	materialEventStorage *storage.MaterialEventStorage,
	materialReadStorage *storage.MaterialReadStorage,
	cropReadStorage *growthstorage.CropReadStorage,
	eventBus eventbus.TaniaEventBus,
) (*FarmServer, error) {
	farmServer := &FarmServer{
		File:     LocalFile{},
		EventBus: eventBus,
	}

	switch *config.Config.TaniaPersistenceEngine {
	case config.DBInmemory:
		farmServer.FarmEventRepo = repoInMem.NewFarmEventRepositoryInMemory(farmEventStorage)
		farmServer.FarmEventQuery = queryInMem.NewFarmEventQueryInMemory(farmEventStorage)
		farmServer.FarmReadRepo = repoInMem.NewFarmReadRepositoryInMemory(farmReadStorage)
		farmServer.FarmReadQuery = queryInMem.NewFarmReadQueryInMemory(farmReadStorage)

		farmServer.AreaEventRepo = repoInMem.NewAreaEventRepositoryInMemory(areaEventStorage)
		farmServer.AreaEventQuery = queryInMem.NewAreaEventQueryInMemory(areaEventStorage)
		farmServer.AreaReadRepo = repoInMem.NewAreaReadRepositoryInMemory(areaReadStorage)
		farmServer.AreaReadQuery = queryInMem.NewAreaReadQueryInMemory(areaReadStorage)

		farmServer.ReservoirEventRepo = repoInMem.NewReservoirEventRepositoryInMemory(reservoirEventStorage)
		farmServer.ReservoirEventQuery = queryInMem.NewReservoirEventQueryInMemory(reservoirEventStorage)
		farmServer.ReservoirReadRepo = repoInMem.NewReservoirReadRepositoryInMemory(reservoirReadStorage)
		farmServer.ReservoirReadQuery = queryInMem.NewReservoirReadQueryInMemory(reservoirReadStorage)

		farmServer.MaterialEventRepo = repoInMem.NewMaterialEventRepositoryInMemory(materialEventStorage)
		farmServer.MaterialEventQuery = queryInMem.NewMaterialEventQueryInMemory(materialEventStorage)
		farmServer.MaterialReadRepo = repoInMem.NewMaterialReadRepositoryInMemory(materialReadStorage)
		farmServer.MaterialReadQuery = queryInMem.NewMaterialReadQueryInMemory(materialReadStorage)

		farmServer.CropReadQuery = queryInMem.NewCropReadQueryInMemory(cropReadStorage)

		// TODO: AreaServiceInMemory should be renamed. It doesn't need InMemory name
		farmServer.AreaService = service.AreaServiceInMemory{
			FarmReadQuery:      farmServer.FarmReadQuery,
			ReservoirReadQuery: farmServer.ReservoirReadQuery,
			CropReadQuery:      farmServer.CropReadQuery,
		}
		// TODO: ReservoirServiceInMemory should be renamed. It doesn't need InMemory name
		farmServer.ReservoirService = service.ReservoirServiceInMemory{
			FarmReadQuery: farmServer.FarmReadQuery,
		}

	case config.DBSqlite:
		farmServer.FarmEventRepo = repoSqlite.NewFarmEventRepositorySqlite(db)
		farmServer.FarmEventQuery = querySqlite.NewFarmEventQuerySqlite(db)
		farmServer.FarmReadRepo = repoSqlite.NewFarmReadRepositorySqlite(db)
		farmServer.FarmReadQuery = querySqlite.NewFarmReadQuerySqlite(db)

		farmServer.AreaEventRepo = repoSqlite.NewAreaEventRepositorySqlite(db)
		farmServer.AreaEventQuery = querySqlite.NewAreaEventQuerySqlite(db)
		farmServer.AreaReadRepo = repoSqlite.NewAreaReadRepositorySqlite(db)
		farmServer.AreaReadQuery = querySqlite.NewAreaReadQuerySqlite(db)

		farmServer.ReservoirEventRepo = repoSqlite.NewReservoirEventRepositorySqlite(db)
		farmServer.ReservoirEventQuery = querySqlite.NewReservoirEventQuerySqlite(db)
		farmServer.ReservoirReadRepo = repoSqlite.NewReservoirReadRepositorySqlite(db)
		farmServer.ReservoirReadQuery = querySqlite.NewReservoirReadQuerySqlite(db)

		farmServer.MaterialEventRepo = repoSqlite.NewMaterialEventRepositorySqlite(db)
		farmServer.MaterialEventQuery = querySqlite.NewMaterialEventQuerySqlite(db)
		farmServer.MaterialReadRepo = repoSqlite.NewMaterialReadRepositorySqlite(db)
		farmServer.MaterialReadQuery = querySqlite.NewMaterialReadQuerySqlite(db)

		farmServer.CropReadQuery = querySqlite.NewCropReadQuerySqlite(db)

		// TODO: AreaServiceInMemory should be renamed. It doesn't need InMemory name
		farmServer.AreaService = service.AreaServiceInMemory{
			FarmReadQuery:      farmServer.FarmReadQuery,
			ReservoirReadQuery: farmServer.ReservoirReadQuery,
			CropReadQuery:      farmServer.CropReadQuery,
		}
		// TODO: ReservoirServiceInMemory should be renamed. It doesn't need InMemory name
		farmServer.ReservoirService = service.ReservoirServiceInMemory{
			FarmReadQuery: farmServer.FarmReadQuery,
		}

	case config.DBMysql:
		farmServer.FarmEventRepo = repoMysql.NewFarmEventRepositoryMysql(db)
		farmServer.FarmEventQuery = queryMysql.NewFarmEventQueryMysql(db)
		farmServer.FarmReadRepo = repoMysql.NewFarmReadRepositoryMysql(db)
		farmServer.FarmReadQuery = queryMysql.NewFarmReadQueryMysql(db)

		farmServer.AreaEventRepo = repoMysql.NewAreaEventRepositoryMysql(db)
		farmServer.AreaEventQuery = queryMysql.NewAreaEventQueryMysql(db)
		farmServer.AreaReadRepo = repoMysql.NewAreaReadRepositoryMysql(db)
		farmServer.AreaReadQuery = queryMysql.NewAreaReadQueryMysql(db)

		farmServer.ReservoirEventRepo = repoMysql.NewReservoirEventRepositoryMysql(db)
		farmServer.ReservoirEventQuery = queryMysql.NewReservoirEventQueryMysql(db)
		farmServer.ReservoirReadRepo = repoMysql.NewReservoirReadRepositoryMysql(db)
		farmServer.ReservoirReadQuery = queryMysql.NewReservoirReadQueryMysql(db)

		farmServer.MaterialEventRepo = repoMysql.NewMaterialEventRepositoryMysql(db)
		farmServer.MaterialEventQuery = queryMysql.NewMaterialEventQueryMysql(db)
		farmServer.MaterialReadRepo = repoMysql.NewMaterialReadRepositoryMysql(db)
		farmServer.MaterialReadQuery = queryMysql.NewMaterialReadQueryMysql(db)

		farmServer.CropReadQuery = queryMysql.NewCropReadQueryMysql(db)

		// TODO: AreaServiceInMemory should be renamed. It doesn't need InMemory name
		farmServer.AreaService = service.AreaServiceInMemory{
			FarmReadQuery:      farmServer.FarmReadQuery,
			ReservoirReadQuery: farmServer.ReservoirReadQuery,
			CropReadQuery:      farmServer.CropReadQuery,
		}
		// TODO: ReservoirServiceInMemory should be renamed. It doesn't need InMemory name
		farmServer.ReservoirService = service.ReservoirServiceInMemory{
			FarmReadQuery: farmServer.FarmReadQuery,
		}
	}

	farmServer.InitSubscriber()

	return farmServer, nil
}

// InitSubscriber defines the mapping of which event this domain listen with their handler
func (s *FarmServer) InitSubscriber() {
	s.EventBus.Subscribe("FarmCreated", s.SaveToFarmReadModel)
	s.EventBus.Subscribe("FarmNameChanged", s.SaveToFarmReadModel)
	s.EventBus.Subscribe("FarmTypeChanged", s.SaveToFarmReadModel)
	s.EventBus.Subscribe("FarmGeolocationChanged", s.SaveToFarmReadModel)
	s.EventBus.Subscribe("FarmRegionChanged", s.SaveToFarmReadModel)

	s.EventBus.Subscribe("ReservoirCreated", s.SaveToReservoirReadModel)
	s.EventBus.Subscribe("ReservoirNameChanged", s.SaveToReservoirReadModel)
	s.EventBus.Subscribe("ReservoirWaterSourceChanged", s.SaveToReservoirReadModel)
	s.EventBus.Subscribe("ReservoirNoteAdded", s.SaveToReservoirReadModel)
	s.EventBus.Subscribe("ReservoirNoteRemoved", s.SaveToReservoirReadModel)

	s.EventBus.Subscribe("AreaCreated", s.SaveToAreaReadModel)
	s.EventBus.Subscribe("AreaNameChanged", s.SaveToAreaReadModel)
	s.EventBus.Subscribe("AreaSizeChanged", s.SaveToAreaReadModel)
	s.EventBus.Subscribe("AreaTypeChanged", s.SaveToAreaReadModel)
	s.EventBus.Subscribe("AreaLocationChanged", s.SaveToAreaReadModel)
	s.EventBus.Subscribe("AreaReservoirChanged", s.SaveToAreaReadModel)
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
	g.GET("/inventories/materials/simple", s.GetMaterialsSimple)
	g.GET("/inventories/plant_types", s.GetInventoryPlantTypes)
	g.GET("/inventories/materials/available_plant_type", s.GetAvailableMaterialPlantType)
	g.POST("/inventories/materials/:type", s.SaveMaterial)
	g.PUT("/inventories/materials/:type/:id", s.UpdateMaterial)
	g.GET("/inventories/materials/:id", s.GetMaterialByID)

	g.POST("", s.SaveFarm)
	g.PUT("/:id", s.UpdateFarm)
	g.GET("", s.FindAllFarm)
	g.GET("/:id", s.FindFarmByID)

	g.POST("/:id/reservoirs", s.SaveReservoir)
	g.PUT("/reservoirs/:id", s.UpdateReservoir)
	g.POST("/reservoirs/:id/notes", s.SaveReservoirNotes)
	g.DELETE("/reservoirs/:reservoir_id/notes/:note_id", s.RemoveReservoirNotes)
	g.GET("/:id/reservoirs", s.GetFarmReservoirs)
	g.GET("/:farm_id/reservoirs/:reservoir_id", s.GetReservoirsByID)

	g.POST("/:id/areas", s.SaveArea)
	g.PUT("/areas/:id", s.UpdateArea)
	g.POST("/areas/:id/notes", s.SaveAreaNotes)
	g.DELETE("/areas/:area_id/notes/:note_id", s.RemoveAreaNotes)
	g.GET("/:id/areas/total", s.GetTotalAreas)
	g.GET("/:id/areas", s.GetFarmAreas)
	g.GET("/:farm_id/areas/:area_id", s.GetAreasByID)
	g.GET("/:farm_id/areas/:area_id/photos", s.GetAreaPhotos)
}

// GetTypes is a FarmServer's handle to get farm types
func (*FarmServer) GetTypes(c echo.Context) error {
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
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
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
		c.FormValue("country"),
		c.FormValue("city"),
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

func (s *FarmServer) UpdateFarm(c echo.Context) error {
	farmUID, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return Error(c, err)
	}

	name := c.FormValue("name")
	farmType := c.FormValue("farm_type")
	latitude := c.FormValue("latitude")
	longitude := c.FormValue("longitude")
	country := c.FormValue("country")
	city := c.FormValue("city")

	// Validate //
	queryResult := <-s.FarmReadQuery.FindByID(farmUID)
	if queryResult.Error != nil {
		return Error(c, queryResult.Error)
	}

	farmRead, ok := queryResult.Result.(storage.FarmRead)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	if farmRead.UID == (uuid.UUID{}) {
		return Error(c, NewRequestValidationError(NotFound, "id"))
	}

	if latitude != "" && longitude == "" {
		return Error(c, NewRequestValidationError(Required, "longitude"))
	}

	if longitude != "" && latitude == "" {
		return Error(c, NewRequestValidationError(Required, "latitude"))
	}

	if country != "" && city == "" {
		return Error(c, NewRequestValidationError(Required, "city"))
	}

	if city != "" && country == "" {
		return Error(c, NewRequestValidationError(Required, "country"))
	}

	// Process //
	queryResult = <-s.FarmEventQuery.FindAllByID(farmUID)
	if queryResult.Error != nil {
		return Error(c, queryResult.Error)
	}

	events, ok := queryResult.Result.([]storage.FarmEvent)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	farm := repository.NewFarmFromHistory(events)

	if name != "" {
		err = farm.ChangeName(name)
		if err != nil {
			return Error(c, err)
		}
	}

	if farmType != "" {
		err = farm.ChangeType(farmType)
		if err != nil {
			return Error(c, err)
		}
	}

	if latitude != "" && longitude != "" {
		err = farm.ChangeGeoLocation(latitude, longitude)
		if err != nil {
			return Error(c, err)
		}
	}

	if country != "" && city != "" {
		err = farm.ChangeRegion(country, city)
		if err != nil {
			return Error(c, err)
		}
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

func (s *FarmServer) UpdateReservoir(c echo.Context) error {
	validation := RequestValidation{}

	reservoirUID, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return Error(c, err)
	}

	name := c.FormValue("name")
	resType := c.FormValue("type")
	capacity := c.FormValue("capacity")

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
		return Error(c, NewRequestValidationError(Required, "id"))
	}

	if resType == domain.BucketType && capacity == "" {
		return Error(c, NewRequestValidationError(Required, "capacity"))
	}

	if resType != "" {
		_, err := validation.ValidateType(resType)
		if err != nil {
			return Error(c, err)
		}
	}

	capacityFloat := float32(0)
	if capacity != "" {
		capacityFloat, err = validation.ValidateCapacity(resType, capacity)
		if err != nil {
			return Error(c, err)
		}
	}

	// Process //
	eventQueryResult := <-s.ReservoirEventQuery.FindAllByID(reservoirRead.UID)
	if eventQueryResult.Error != nil {
		return Error(c, eventQueryResult.Error)
	}

	events := eventQueryResult.Result.([]storage.ReservoirEvent)
	reservoir := repository.NewReservoirFromHistory(events)

	if name != "" {
		err = reservoir.ChangeName(name)
		if err != nil {
			return Error(c, err)
		}
	}

	if resType != "" && capacity != "" {
		err = reservoir.ChangeWaterSource(resType, capacityFloat)
		if err != nil {
			return Error(c, err)
		}
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
		return Error(c, NewRequestValidationError(Required, "id"))
	}

	if content == "" {
		return Error(c, NewRequestValidationError(Required, "content"))
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
		return Error(c, NewRequestValidationError(NotFound, "reservoir_id"))
	}

	noteFound := false

	for _, v := range reservoirRead.Notes {
		if v.UID == noteUID {
			noteFound = true
		}
	}

	if !noteFound {
		return Error(c, NewRequestValidationError(NotFound, "note_id"))
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

func (s *FarmServer) GetFarmReservoirs(c echo.Context) error {
	farmUID, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return Error(c, err)
	}

	result := <-s.ReservoirReadQuery.FindAllByFarm(farmUID)
	if result.Error != nil {
		return Error(c, result.Error)
	}

	reservoirs, ok := result.Result.([]storage.ReservoirRead)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	data := make(map[string][]storage.ReservoirRead)

	for _, v := range reservoirs {
		r, err := MapToReservoirReadFromRead(s, v)
		if err != nil {
			return Error(c, err)
		}

		data["data"] = append(data["data"], r)
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
		return Error(c, NewRequestValidationError(Required, "reservoir_id"))
	}

	data := make(map[string]storage.ReservoirRead)

	data["data"], err = MapToReservoirReadFromRead(s, reservoir)
	if err != nil {
		Error(c, err)
	}

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

func (s *FarmServer) UpdateArea(c echo.Context) error {
	validation := RequestValidation{}

	areaUID, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return Error(c, err)
	}

	name := c.FormValue("name")
	size := c.FormValue("size")
	sizeUnit := c.FormValue("size_unit")
	areaType := c.FormValue("type")
	location := c.FormValue("location")
	reservoirID := c.FormValue("reservoir_id")
	photo, photoErr := c.FormFile("photo")

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
		return Error(c, NewRequestValidationError(NotFound, "id"))
	}

	if size != "" && sizeUnit == "" {
		return Error(c, NewRequestValidationError(Required, "size_unit"))
	}

	if sizeUnit != "" && size == "" {
		return Error(c, NewRequestValidationError(Required, "size"))
	}

	// Process //
	eventQueryResult := <-s.AreaEventQuery.FindAllByID(areaRead.UID)
	if eventQueryResult.Error != nil {
		return Error(c, eventQueryResult.Error)
	}

	events := eventQueryResult.Result.([]storage.AreaEvent)

	area := repository.NewAreaFromHistory(events)

	if name != "" {
		err := area.ChangeName(name)
		if err != nil {
			return Error(c, err)
		}
	}

	if size != "" && sizeUnit != "" {
		areaSize, err := validation.ValidateAreaSize(size, sizeUnit)
		if err != nil {
			return Error(c, err)
		}

		err = area.ChangeSize(areaSize)
		if err != nil {
			return Error(c, err)
		}
	}

	if areaType != "" {
		err = area.ChangeType(s.AreaService, areaType)
		if err != nil {
			return Error(c, err)
		}
	}

	if location != "" {
		err = area.ChangeLocation(location)
		if err != nil {
			return Error(c, err)
		}
	}

	if reservoirID != "" {
		resUID, err := uuid.FromString(reservoirID)
		if err != nil {
			return Error(c, err)
		}

		err = area.ChangeReservoir(resUID)
		if err != nil {
			return Error(c, err)
		}
	}

	if photoErr == nil {
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

	detailArea, err := MapToDetailArea(s, *area)
	if err != nil {
		return Error(c, err)
	}

	data := make(map[string]DetailArea)
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
		return Error(c, NewRequestValidationError(NotFound, "id"))
	}

	if content == "" {
		return Error(c, NewRequestValidationError(Required, "content"))
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
		return Error(c, NewRequestValidationError(NotFound, "area_id"))
	}

	found := false

	for _, v := range areaRead.Notes {
		if v.UID == noteUID {
			found = true
		}
	}

	if !found {
		return Error(c, NewRequestValidationError(NotFound, "note_id"))
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
		return Error(c, NewRequestValidationError(NotFound, "farm_id"))
	}

	queryResult = <-s.AreaReadQuery.FindByIDAndFarm(areaUID, farmUID)
	if queryResult.Error != nil {
		return Error(c, queryResult.Error)
	}

	areaRead, ok := queryResult.Result.(storage.AreaRead)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusBadRequest, "Internal server error"))
	}

	if areaRead.UID == (uuid.UUID{}) {
		return Error(c, NewRequestValidationError(NotFound, "area_id"))
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
		return Error(c, NewRequestValidationError(NotFound, "farm_id"))
	}

	queryResult = <-s.AreaReadQuery.FindByIDAndFarm(areaUID, farmUID)
	if queryResult.Error != nil {
		return Error(c, queryResult.Error)
	}

	areaRead, ok := queryResult.Result.(storage.AreaRead)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusBadRequest, "Internal server error"))
	}

	if areaRead.UID == (uuid.UUID{}) {
		return Error(c, NewRequestValidationError(NotFound, "area_id"))
	}

	if areaRead.Photo.Filename == "" {
		return Error(c, NewRequestValidationError(NotFound, "photo"))
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

func (*FarmServer) GetInventoryPlantTypes(c echo.Context) error {
	data := make(map[string][]string)

	plantTypes := MapToPlantType(domain.PlantTypes())

	data["data"] = plantTypes

	return c.JSON(http.StatusOK, data)
}

func (s *FarmServer) GetMaterials(c echo.Context) error {
	materialType := c.QueryParam("type")
	materialTypeDetail := c.QueryParam("type_detail")
	page := c.QueryParam("page")
	limit := c.QueryParam("limit")

	pageInt, limitInt, err := paginationhelper.ParsePagination(page, limit)
	if err != nil {
		return Error(c, err)
	}

	queryResult := <-s.MaterialReadQuery.FindAll(materialType, materialTypeDetail, pageInt, limitInt)
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

	queryResult = <-s.MaterialReadQuery.CountAll(materialType, materialTypeDetail)
	if queryResult.Error != nil {
		return Error(c, queryResult.Error)
	}

	total, ok := queryResult.Result.(int)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	data := make(map[string]interface{})
	data["data"] = materials
	data["total"] = total
	data["page"] = pageInt

	return c.JSON(http.StatusOK, data)
}

func (s *FarmServer) GetMaterialsSimple(c echo.Context) error {
	materialType := c.QueryParam("type")
	materialTypeDetail := c.QueryParam("type_detail")

	queryResult := <-s.MaterialReadQuery.FindAll(materialType, materialTypeDetail, 0, 0)
	if queryResult.Error != nil {
		return Error(c, queryResult.Error)
	}

	results, ok := queryResult.Result.([]storage.MaterialRead)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	data := make(map[string][]MaterialSimple)
	data["data"] = MapToMaterialSimpleFromRead(results)

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
		return Error(c, NewRequestValidationError(InvalidOption, "quantity"))
	}

	var expDate *time.Time

	if expirationDate != "" {
		tp, err := time.Parse("2006-01-02", expirationDate)
		if err != nil {
			return Error(c, NewRequestValidationError(ParseFailed, "expiration_date"))
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
			return Error(c, NewRequestValidationError(InvalidOption, "plant_type"))
		}

		mt, err = domain.CreateMaterialTypeSeed(pt.Code)
		if err != nil {
			return Error(c, NewRequestValidationError(InvalidOption, "type"))
		}
	case strings.ToLower(domain.MaterialTypeAgrochemicalCode):
		ct := domain.GetChemicalType(chemicalType)
		if ct == (domain.ChemicalType{}) {
			return Error(c, NewRequestValidationError(InvalidOption, "chemical_type"))
		}

		mt, err = domain.CreateMaterialTypeAgrochemical(ct.Code)
		if err != nil {
			return Error(c, NewRequestValidationError(InvalidOption, "type"))
		}
	case strings.ToLower(domain.MaterialTypeGrowingMediumCode):
		mt = domain.MaterialTypeGrowingMedium{}
	case strings.ToLower(domain.MaterialTypeLabelAndCropSupportCode):
		mt = domain.MaterialTypeLabelAndCropSupport{}
	case strings.ToLower(domain.MaterialTypeSeedingContainerCode):
		ct := domain.GetContainerType(containerType)
		if ct == (domain.ContainerType{}) {
			return Error(c, NewRequestValidationError(InvalidOption, "container_type"))
		}

		mt, err = domain.CreateMaterialTypeSeedingContainer(ct.Code)
		if err != nil {
			return Error(c, NewRequestValidationError(InvalidOption, "type"))
		}
	case strings.ToLower(domain.MaterialTypePostHarvestSupplyCode):
		mt = domain.MaterialTypePostHarvestSupply{}
	case strings.ToLower(domain.MaterialTypeOtherCode):
		mt = domain.MaterialTypeOther{}
	case strings.ToLower(domain.MaterialTypePlantCode):
		pt := domain.GetPlantType(plantType)
		if pt == (domain.PlantType{}) {
			return Error(c, NewRequestValidationError(InvalidOption, "plant_type"))
		}

		mt, err = domain.CreateMaterialTypePlant(pt.Code)
		if err != nil {
			return Error(c, NewRequestValidationError(InvalidOption, "type"))
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
		return Error(c, NewRequestValidationError(Required, "currency_code"))
	}

	if currencyCode != "" && pricePerUnit == "" {
		return Error(c, NewRequestValidationError(Required, "price_per_unit"))
	}

	if quantity != "" && quantityUnit == "" {
		return Error(c, NewRequestValidationError(Required, "quantity_unit"))
	}

	if quantityUnit != "" && quantity == "" {
		return Error(c, NewRequestValidationError(Required, "quantity"))
	}

	var expDate *time.Time

	if expirationDate != "" {
		tp, err := time.Parse("2006-01-02", expirationDate)
		if err != nil {
			return Error(c, NewRequestValidationError(ParseFailed, "expiration_date"))
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

	if materialRead.UID == (uuid.UUID{}) {
		return Error(c, NewRequestValidationError(NotFound, "id"))
	}

	// Process //
	var mt domain.MaterialType

	switch materialTypeParam {
	case strings.ToLower(domain.MaterialTypeSeedCode):
		if plantType != "" {
			pt := domain.GetPlantType(plantType)
			if pt == (domain.PlantType{}) {
				return Error(c, NewRequestValidationError(InvalidOption, "plant_type"))
			}

			mt, err = domain.CreateMaterialTypeSeed(pt.Code)
			if err != nil {
				return Error(c, NewRequestValidationError(InvalidOption, "type"))
			}

			materialRead.Type = mt
		}
	case strings.ToLower(domain.MaterialTypeAgrochemicalCode):
		if chemicalType != "" {
			ct := domain.GetChemicalType(chemicalType)
			if ct == (domain.ChemicalType{}) {
				return Error(c, NewRequestValidationError(InvalidOption, "chemical_type"))
			}

			mt, err = domain.CreateMaterialTypeAgrochemical(ct.Code)
			if err != nil {
				return Error(c, NewRequestValidationError(InvalidOption, "type"))
			}

			materialRead.Type = mt
		}
	case strings.ToLower(domain.MaterialTypeSeedingContainerCode):
		if containerType != "" {
			ct := domain.GetContainerType(containerType)
			if ct == (domain.ContainerType{}) {
				return Error(c, NewRequestValidationError(InvalidOption, "container_type"))
			}

			mt, err = domain.CreateMaterialTypeSeedingContainer(ct.Code)
			if err != nil {
				return Error(c, NewRequestValidationError(InvalidOption, "type"))
			}

			materialRead.Type = mt
		}
	case strings.ToLower(domain.MaterialTypePlantCode):
		if plantType != "" {
			pt := domain.GetPlantType(plantType)
			if pt == (domain.PlantType{}) {
				return Error(c, NewRequestValidationError(InvalidOption, "plant_type"))
			}

			mt, err = domain.CreateMaterialTypePlant(pt.Code)
			if err != nil {
				return Error(c, NewRequestValidationError(InvalidOption, "type"))
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

func (s *FarmServer) GetMaterialByID(c echo.Context) error {
	materialUID, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return Error(c, err)
	}

	queryResult := <-s.MaterialReadQuery.FindByID(materialUID)
	if queryResult.Error != nil {
		return Error(c, err)
	}

	materialRead := queryResult.Result.(storage.MaterialRead)

	if materialRead.UID == (uuid.UUID{}) {
		return Error(c, NewRequestValidationError(NotFound, "id"))
	}

	data := make(map[string]Material)
	data["data"] = MapToMaterialFromRead(materialRead)

	return c.JSON(http.StatusOK, data)
}

func (s *FarmServer) GetAvailableMaterialPlantType(c echo.Context) error {
	data := make(map[string][]AvailableMaterialPlantType)

	params := domain.MaterialTypeSeedCode + "," + domain.MaterialTypePlantCode

	// Process //
	// TODO: Refactor this query to only get material by plant type
	result := <-s.MaterialReadQuery.FindAll(params, "", 0, 100)

	materials, ok := result.Result.([]storage.MaterialRead)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusBadRequest, "Internal server error"))
	}

	data["data"] = MapToAvailableMaterialPlantType(materials)

	return c.JSON(http.StatusOK, data)
}

func (s *FarmServer) publishUncommittedEvents(entity interface{}) {
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
}
