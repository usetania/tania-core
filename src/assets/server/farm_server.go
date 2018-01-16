package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Tanibox/tania-server/config"
	"github.com/Tanibox/tania-server/src/assets/domain"
	"github.com/Tanibox/tania-server/src/assets/query"
	"github.com/Tanibox/tania-server/src/assets/repository"
	"github.com/Tanibox/tania-server/src/assets/service"
	"github.com/Tanibox/tania-server/src/assets/storage"
	"github.com/Tanibox/tania-server/src/helper/imagehelper"
	"github.com/Tanibox/tania-server/src/helper/stringhelper"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
)

// FarmServer ties the routes and handlers with injected dependencies
type FarmServer struct {
	FarmRepo               repository.FarmRepository
	ReservoirRepo          repository.ReservoirRepository
	AreaRepo               repository.AreaRepository
	AreaQuery              query.AreaQuery
	CropRepo               repository.CropRepository
	CropQuery              query.CropQuery
	CropService            service.CropService
	InventoryMaterialRepo  repository.InventoryMaterialRepository
	InventoryMaterialQuery query.InventoryMaterialQuery
	File                   File
}

// NewFarmServer initializes FarmServer's dependencies and create new FarmServer struct
func NewFarmServer() (*FarmServer, error) {
	farmStorage := storage.FarmStorage{FarmMap: make(map[uuid.UUID]domain.Farm)}
	farmRepo := repository.NewFarmRepositoryInMemory(&farmStorage)

	areaStorage := storage.AreaStorage{AreaMap: make(map[uuid.UUID]domain.Area)}
	areaRepo := repository.NewAreaRepositoryInMemory(&areaStorage)
	areaQuery := query.NewAreaQueryInMemory(&areaStorage)

	reservoirStorage := storage.ReservoirStorage{ReservoirMap: make(map[uuid.UUID]domain.Reservoir)}
	reservoirRepo := repository.NewReservoirRepositoryInMemory(&reservoirStorage)

	cropStorage := storage.CropStorage{CropMap: make(map[uuid.UUID]domain.Crop)}
	cropRepo := repository.NewCropRepositoryInMemory(&cropStorage)
	cropQuery := query.NewCropQueryInMemory(&cropStorage)
	cropService := service.CropService{CropQuery: cropQuery}

	inventoryMaterialStorage := storage.InventoryMaterialStorage{InventoryMaterialMap: make(map[uuid.UUID]domain.InventoryMaterial)}
	inventoryMaterialRepo := repository.NewInventoryMaterialRepositoryInMemory(&inventoryMaterialStorage)
	inventoryMaterialQuery := query.NewInventoryMaterialQueryInMemory(&inventoryMaterialStorage)

	farmServer := FarmServer{
		FarmRepo:               farmRepo,
		ReservoirRepo:          reservoirRepo,
		AreaRepo:               areaRepo,
		AreaQuery:              areaQuery,
		CropRepo:               cropRepo,
		CropQuery:              cropQuery,
		CropService:            cropService,
		InventoryMaterialRepo:  inventoryMaterialRepo,
		InventoryMaterialQuery: inventoryMaterialQuery,
		File: LocalFile{},
	}

	if *config.Config.DemoMode {
		initDataDemo(
			&farmServer,
			&farmStorage,
			&areaStorage,
			&reservoirStorage,
			&cropStorage,
			&inventoryMaterialStorage,
		)
	}

	return &farmServer, nil
}

// Mount defines the FarmServer's endpoints with its handlers
func (s *FarmServer) Mount(g *echo.Group) {
	g.GET("/types", s.GetTypes)
	g.GET("/inventories/plant_types", s.GetInventoryPlantTypes)
	g.GET("/inventories", s.GetAvailableInventories)
	g.POST("/inventories", s.SaveInventory)

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
	g.GET("/:id/crops", s.FindAllCrops)
	g.GET("/areas/:id/crops", s.FindAllCropsByArea)
	g.POST("/areas/:id/crops", s.SaveAreaCropBatch)
	g.POST("/crops/:id/notes", s.SaveCropNotes)
	g.DELETE("/crops/:crop_id/notes/:note_id", s.RemoveCropNotes)
	g.GET("/:farm_id/areas/:area_id", s.GetAreasByID)
	g.GET("/:farm_id/areas/:area_id/photos", s.GetAreaPhotos)
}

// GetTypes is a FarmServer's handle to get farm types
func (s *FarmServer) GetTypes(c echo.Context) error {
	types := domain.FindAllFarmTypes()

	return c.JSON(http.StatusOK, types)
}

func (s FarmServer) FindAllFarm(c echo.Context) error {
	data := make(map[string][]SimpleFarm)

	result := <-s.FarmRepo.FindAll()
	if result.Error != nil {
		return result.Error
	}

	farms, ok := result.Result.([]domain.Farm)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	data["data"] = MapToSimpleFarm(farms)

	return c.JSON(http.StatusOK, data)
}

// SaveFarm is a FarmServer's handler to save new Farm
func (s *FarmServer) SaveFarm(c echo.Context) error {
	data := make(map[string]domain.Farm)

	farm, err := domain.CreateFarm(c.FormValue("name"), c.FormValue("farm_type"))
	if err != nil {
		return Error(c, err)
	}

	err = farm.ChangeGeoLocation(c.FormValue("latitude"), c.FormValue("longitude"))
	if err != nil {
		return Error(c, err)
	}

	err = farm.ChangeRegion(c.FormValue("country_code"), c.FormValue("city_code"))
	if err != nil {
		return Error(c, err)
	}

	err = <-s.FarmRepo.Save(&farm)
	if err != nil {
		return Error(c, err)
	}

	data["data"] = farm

	return c.JSON(http.StatusOK, data)
}

func (s *FarmServer) FindFarmByID(c echo.Context) error {
	data := make(map[string]domain.Farm)

	result := <-s.FarmRepo.FindByID(c.Param("id"))
	if result.Error != nil {
		return result.Error
	}

	farm, ok := result.Result.(domain.Farm)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	data["data"] = farm

	return c.JSON(http.StatusOK, data)
}

// SaveReservoir is a FarmServer's handler to save new Reservoir and place it to a Farm
func (s *FarmServer) SaveReservoir(c echo.Context) error {
	data := make(map[string]DetailReservoir)
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

	farm, err := validation.ValidateFarm(*s, c.Param("id"))
	if err != nil {
		return Error(c, err)
	}

	// Process //
	r, err := domain.CreateReservoir(farm, name)
	if err != nil {
		return Error(c, err)
	}

	if waterSourceType == "bucket" {
		b, err := domain.CreateBucket(capacity, 0)
		if err != nil {
			return Error(c, err)
		}

		r.AttachBucket(b)
	} else if waterSourceType == "tap" {
		t, err := domain.CreateTap()
		if err != nil {
			return Error(c, err)
		}

		r.AttachTap(t)
	}

	err = farm.AddReservoir(&r)
	if err != nil {
		return Error(c, err)
	}

	// Persists //
	err = <-s.ReservoirRepo.Save(&r)
	if err != nil {
		return Error(c, err)
	}

	err = <-s.FarmRepo.Save(&farm)
	if err != nil {
		return Error(c, err)
	}

	detailReservoir, err := MapToDetailReservoir(s, r)
	if err != nil {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	data["data"] = detailReservoir

	return c.JSON(http.StatusOK, data)
}

func (s *FarmServer) SaveReservoirNotes(c echo.Context) error {
	data := make(map[string]DetailReservoir)

	reservoirID := c.Param("id")
	content := c.FormValue("content")

	// Validate //
	result := <-s.ReservoirRepo.FindByID(reservoirID)
	if result.Error != nil {
		return Error(c, result.Error)
	}

	reservoir, ok := result.Result.(domain.Reservoir)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	result = <-s.FarmRepo.FindByID(reservoir.Farm.UID.String())
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
	reservoir.AddNewNote(content)
	farm.ChangeReservoirInformation(reservoir)

	// Persists //
	resultSave := <-s.ReservoirRepo.Save(&reservoir)
	if resultSave != nil {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	resultSave = <-s.FarmRepo.Save(&farm)
	if resultSave != nil {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	detailReservoir, err := MapToDetailReservoir(s, reservoir)
	if err != nil {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	data["data"] = detailReservoir

	return c.JSON(http.StatusOK, data)
}

func (s *FarmServer) RemoveReservoirNotes(c echo.Context) error {
	data := make(map[string]DetailReservoir)

	reservoirID := c.Param("reservoir_id")
	noteID := c.Param("note_id")

	// Validate //
	result := <-s.ReservoirRepo.FindByID(reservoirID)
	if result.Error != nil {
		return Error(c, result.Error)
	}

	reservoir, ok := result.Result.(domain.Reservoir)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	result = <-s.FarmRepo.FindByID(reservoir.Farm.UID.String())
	if result.Error != nil {
		return Error(c, result.Error)
	}

	farm, ok := result.Result.(domain.Farm)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	// Process //
	err := reservoir.RemoveNote(noteID)
	if err != nil {
		return Error(c, err)
	}

	err = farm.ChangeReservoirInformation(reservoir)
	if err != nil {
		return Error(c, err)
	}

	// Persists //
	resultSave := <-s.ReservoirRepo.Save(&reservoir)
	if resultSave != nil {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	resultSave = <-s.FarmRepo.Save(&farm)
	if resultSave != nil {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	detailReservoir, err := MapToDetailReservoir(s, reservoir)
	if err != nil {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	data["data"] = detailReservoir

	return c.JSON(http.StatusOK, data)
}

func (s *FarmServer) GetFarmReservoirs(c echo.Context) error {
	data := make(map[string][]DetailReservoir)

	result := <-s.FarmRepo.FindByID(c.Param("id"))
	if result.Error != nil {
		return result.Error
	}

	farm, ok := result.Result.(domain.Farm)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	reservoirs, err := MapToReservoir(s, farm.Reservoirs)
	if err != nil {
		return Error(c, err)
	}

	data["data"] = reservoirs
	if len(farm.Reservoirs) == 0 {
		data["data"] = []DetailReservoir{}
	}

	return c.JSON(http.StatusOK, data)
}

func (s *FarmServer) GetReservoirsByID(c echo.Context) error {
	data := make(map[string]DetailReservoir)

	// Validate //
	result := <-s.FarmRepo.FindByID(c.Param("farm_id"))
	if result.Error != nil {
		return Error(c, result.Error)
	}

	_, ok := result.Result.(domain.Farm)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusBadRequest, "Internal server error"))
	}

	result = <-s.ReservoirRepo.FindByID(c.Param("reservoir_id"))
	if result.Error != nil {
		return Error(c, result.Error)
	}

	reservoir, ok := result.Result.(domain.Reservoir)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusBadRequest, "Internal server error"))
	}

	detailReservoir, err := MapToDetailReservoir(s, reservoir)
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
	farm, err := validation.ValidateFarm(*s, c.Param("id"))
	if err != nil {
		return Error(c, err)
	}

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

	data["data"] = MapToDetailArea(area)

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

	data["data"] = MapToDetailArea(area)

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

	data["data"] = MapToDetailArea(area)

	return c.JSON(http.StatusOK, data)
}

func (s *FarmServer) GetFarmAreas(c echo.Context) error {
	data := make(map[string][]domain.Area)

	result := <-s.FarmRepo.FindByID(c.Param("id"))
	if result.Error != nil {
		return Error(c, result.Error)
	}

	farm, ok := result.Result.(domain.Farm)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusBadRequest, "Internal server error"))
	}

	data["data"] = MapToArea(farm.Areas)
	if len(farm.Areas) == 0 {
		data["data"] = []domain.Area{}
	}

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

	data["data"] = MapToDetailArea(area)

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
	plantTypes := MapToPlantType(domain.GetPlantTypes())

	return c.JSON(http.StatusOK, plantTypes)
}

func (s *FarmServer) SaveInventory(c echo.Context) error {
	data := make(map[string]string)

	pType := c.FormValue("plant_type")
	variety := c.FormValue("variety")

	// Validate //
	var plantType domain.PlantType
	switch pType {
	case "vegetable":
		plantType = domain.Vegetable{}
	case "fruit":
		plantType = domain.Fruit{}
	case "herb":
		plantType = domain.Herb{}
	case "flower":
		plantType = domain.Flower{}
	case "tree":
		plantType = domain.Tree{}
	default:
		return Error(c, NewRequestValidationError(NOT_FOUND, "plant_type"))
	}

	// Process //
	inventoryMaterial, err := domain.CreateInventoryMaterial(plantType, variety)
	if err != nil {
		return Error(c, err)
	}

	// Persist //
	err = <-s.InventoryMaterialRepo.Save(&inventoryMaterial)
	if err != nil {
		return Error(c, err)
	}

	data["data"] = inventoryMaterial.UID.String()

	return c.JSON(http.StatusOK, data)
}

func (s *FarmServer) GetAvailableInventories(c echo.Context) error {
	data := make(map[string][]AvailableInventory)

	// Process //
	result := <-s.InventoryMaterialRepo.FindAll()

	inventories, ok := result.Result.([]domain.InventoryMaterial)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusBadRequest, "Internal server error"))
	}

	data["data"] = MapToAvailableInventories(inventories)

	return c.JSON(http.StatusOK, data)
}

func (s *FarmServer) SaveAreaCropBatch(c echo.Context) error {
	// Form Value
	areaID := c.Param("id")
	cropType := c.FormValue("crop_type")
	plantType := c.FormValue("plant_type")
	variety := c.FormValue("variety")

	containerQuantity, err := strconv.Atoi(c.FormValue("container_quantity"))
	if err != nil {
		return Error(c, err)
	}

	containerType := c.FormValue("container_type")
	containerCell, err := strconv.Atoi(c.FormValue("container_cell"))
	if err != nil {
		return Error(c, err)
	}

	// Validate //
	areaResult := <-s.AreaRepo.FindByID(areaID)
	if areaResult.Error != nil {
		return Error(c, areaResult.Error)
	}

	area, ok := areaResult.Result.(domain.Area)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusBadRequest, "Internal server error"))
	}

	var cropT domain.CropType
	switch cropType {
	case "nursery":
		cropT = domain.Nursery{}
	case "growing":
		cropT = domain.Growing{}
	default:
		return Error(c, NewRequestValidationError(INVALID_OPTION, "crop_type"))
	}

	var plantT domain.PlantType
	switch plantType {
	case "vegetable":
		plantT = domain.Vegetable{}
	case "fruit":
		plantT = domain.Fruit{}
	case "herb":
		plantT = domain.Herb{}
	case "flower":
		plantT = domain.Flower{}
	case "tree":
		plantT = domain.Tree{}
	default:
		return Error(c, NewRequestValidationError(NOT_FOUND, "plant_type"))
	}

	queryResult := <-s.InventoryMaterialQuery.FindInventoryByPlantTypeAndVariety(plantT, variety)
	if queryResult.Error != nil {
		return Error(c, queryResult.Error)
	}

	inventoryMaterial, ok := queryResult.Result.(domain.InventoryMaterial)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusBadRequest, "Internal server error"))
	}

	var containerT domain.CropContainerType
	switch containerType {
	case "tray":
		containerT = domain.Tray{Cell: containerCell}
	case "pot":
		containerT = domain.Pot{}
	default:
		return Error(c, NewRequestValidationError(NOT_FOUND, "container_type"))
	}

	cropContainer := domain.CropContainer{
		Quantity: containerQuantity,
		Type:     containerT,
	}

	// Process //
	cropBatch, err := domain.CreateCropBatch(area)
	if err != nil {
		return Error(c, err)
	}

	err = cropBatch.ChangeCropType(cropT)
	if err != nil {
		return Error(c, err)
	}

	err = s.CropService.ChangeInventory(&cropBatch, inventoryMaterial)
	if err != nil {
		return Error(c, err)
	}

	err = cropBatch.ChangeContainer(cropContainer)
	if err != nil {
		return Error(c, err)
	}

	// Persists //
	err = <-s.CropRepo.Save(&cropBatch)
	if err != nil {
		return Error(c, err)
	}

	data := make(map[string]CropBatch)
	data["data"] = MapToCropBatch(cropBatch)

	return c.JSON(http.StatusOK, data)
}

func (s *FarmServer) SaveCropNotes(c echo.Context) error {
	data := make(map[string]CropBatch)

	cropID := c.Param("id")
	content := c.FormValue("content")

	// Validate //
	result := <-s.CropRepo.FindByID(cropID)
	if result.Error != nil {
		return Error(c, result.Error)
	}

	crop, ok := result.Result.(domain.Crop)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	if content == "" {
		return Error(c, NewRequestValidationError(REQUIRED, "content"))
	}

	// Process //
	crop.AddNewNote(content)

	// Persists //
	resultSave := <-s.CropRepo.Save(&crop)
	if resultSave != nil {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	data["data"] = MapToCropBatch(crop)

	return c.JSON(http.StatusOK, data)
}

func (s *FarmServer) RemoveCropNotes(c echo.Context) error {
	data := make(map[string]CropBatch)

	cropID := c.Param("crop_id")
	noteID := c.Param("note_id")

	// Validate //
	result := <-s.CropRepo.FindByID(cropID)
	if result.Error != nil {
		return Error(c, result.Error)
	}

	crop, ok := result.Result.(domain.Crop)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	// Process //
	err := crop.RemoveNote(noteID)
	if err != nil {
		return Error(c, err)
	}

	// Persists //
	resultSave := <-s.CropRepo.Save(&crop)
	if resultSave != nil {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	data["data"] = MapToCropBatch(crop)

	return c.JSON(http.StatusOK, data)
}

// TODO: The crops should be found by its Farm
func (s *FarmServer) FindAllCrops(c echo.Context) error {
	data := make(map[string][]CropBatch)

	// Params //
	farmID := c.Param("id")

	// Validate //
	result := <-s.FarmRepo.FindByID(farmID)
	if result.Error != nil {
		return Error(c, result.Error)
	}

	_, ok := result.Result.(domain.Farm)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	// Process //
	result = <-s.CropRepo.FindAll()
	if result.Error != nil {
		return Error(c, result.Error)
	}

	crops, ok := result.Result.([]domain.Crop)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	data["data"] = []CropBatch{}
	for _, v := range crops {
		data["data"] = append(data["data"], MapToCropBatch(v))
	}

	return c.JSON(http.StatusOK, data)
}

func (s *FarmServer) FindAllCropsByArea(c echo.Context) error {
	data := make(map[string][]CropBatch)

	// Params //
	areaID := c.Param("id")

	// Validate //
	result := <-s.AreaRepo.FindByID(areaID)
	if result.Error != nil {
		return Error(c, result.Error)
	}

	area, ok := result.Result.(domain.Area)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	// Process //
	resultQuery := <-s.CropQuery.FindAllCropsByArea(area)
	if resultQuery.Error != nil {
		return Error(c, resultQuery.Error)
	}

	crops, ok := resultQuery.Result.([]domain.Crop)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	data["data"] = []CropBatch{}
	for _, v := range crops {
		data["data"] = append(data["data"], MapToCropBatch(v))
	}

	return c.JSON(http.StatusOK, data)
}

func initDataDemo(
	server *FarmServer,
	farmStorage *storage.FarmStorage,
	areaStorage *storage.AreaStorage,
	reservoirStorage *storage.ReservoirStorage,
	cropStorage *storage.CropStorage,
	inventoryMaterialStorage *storage.InventoryMaterialStorage,
) {
	farmUID, _ := uuid.NewV4()
	farm1 := domain.Farm{
		UID:         farmUID,
		Name:        "MyFarm",
		Type:        "organic",
		Latitude:    "10.00",
		Longitude:   "11.00",
		CountryCode: "ID",
		CityCode:    "JK",
		IsActive:    true,
	}

	farmStorage.FarmMap[farmUID] = farm1

	uid, _ := uuid.NewV4()

	noteUID, _ := uuid.NewV4()
	reservoirNotes := make(map[uuid.UUID]domain.ReservoirNote, 0)
	reservoirNotes[noteUID] = domain.ReservoirNote{
		UID:         noteUID,
		Content:     "Don't forget to close the bucket after using",
		CreatedDate: time.Now(),
	}

	reservoir1 := domain.Reservoir{
		UID:         uid,
		Name:        "MyBucketReservoir",
		PH:          8,
		EC:          12.5,
		Temperature: 29,
		WaterSource: domain.Bucket{Capacity: 100, Volume: 10},
		Farm:        farm1,
		Notes:       reservoirNotes,
		CreatedDate: time.Now(),
	}

	farm1.AddReservoir(&reservoir1)
	farmStorage.FarmMap[farmUID] = farm1
	reservoirStorage.ReservoirMap[uid] = reservoir1

	uid, _ = uuid.NewV4()
	reservoir2 := domain.Reservoir{
		UID:         uid,
		Name:        "MyTapReservoir",
		PH:          8,
		EC:          12.5,
		Temperature: 29,
		WaterSource: domain.Tap{},
		Farm:        farm1,
		Notes:       make(map[uuid.UUID]domain.ReservoirNote),
		CreatedDate: time.Now(),
	}

	farm1.AddReservoir(&reservoir2)
	farmStorage.FarmMap[farmUID] = farm1
	reservoirStorage.ReservoirMap[uid] = reservoir2

	uid, _ = uuid.NewV4()

	noteUID, _ = uuid.NewV4()
	areaNotes := make(map[uuid.UUID]domain.AreaNote, 0)
	areaNotes[noteUID] = domain.AreaNote{
		UID:         noteUID,
		Content:     "This area should only be used for seeding.",
		CreatedDate: time.Now(),
	}

	area1 := domain.Area{
		UID:       uid,
		Name:      "MySeedingArea",
		Size:      domain.SquareMeter{Value: 10},
		Type:      "nursery",
		Location:  "indoor",
		Photo:     domain.AreaPhoto{},
		Notes:     areaNotes,
		Reservoir: reservoir2,
		Farm:      farm1,
	}

	farm1.AddArea(&area1)
	farmStorage.FarmMap[farmUID] = farm1
	areaStorage.AreaMap[uid] = area1

	uid, _ = uuid.NewV4()
	area2 := domain.Area{
		UID:       uid,
		Name:      "MyGrowingArea",
		Size:      domain.SquareMeter{Value: 100},
		Type:      "growing",
		Location:  "outdoor",
		Photo:     domain.AreaPhoto{},
		Notes:     make(map[uuid.UUID]domain.AreaNote),
		Reservoir: reservoir1,
		Farm:      farm1,
	}

	farm1.AddArea(&area2)
	farmStorage.FarmMap[farmUID] = farm1
	areaStorage.AreaMap[uid] = area2

	uid, _ = uuid.NewV4()
	inventory1 := domain.InventoryMaterial{
		UID:       uid,
		PlantType: domain.Vegetable{},
		Variety:   "Bayam Lu Hsieh",
	}

	inventoryMaterialStorage.InventoryMaterialMap[uid] = inventory1

	uid, _ = uuid.NewV4()
	inventory2 := domain.InventoryMaterial{
		UID:       uid,
		PlantType: domain.Vegetable{},
		Variety:   "Tomat Super One",
	}

	inventoryMaterialStorage.InventoryMaterialMap[uid] = inventory2

	uid, _ = uuid.NewV4()
	inventory3 := domain.InventoryMaterial{
		UID:       uid,
		PlantType: domain.Fruit{},
		Variety:   "Apple Rome Beauty",
	}

	inventoryMaterialStorage.InventoryMaterialMap[uid] = inventory3

	uid, _ = uuid.NewV4()
	inventory4 := domain.InventoryMaterial{
		UID:       uid,
		PlantType: domain.Fruit{},
		Variety:   "Orange Sweet Mandarin",
	}

	inventoryMaterialStorage.InventoryMaterialMap[uid] = inventory4

	uid, _ = uuid.NewV4()

	noteUID, _ = uuid.NewV4()
	cropNotes := make(map[uuid.UUID]domain.CropNote, 0)
	cropNotes[noteUID] = domain.CropNote{
		UID:         noteUID,
		Content:     "This crop must be intensely watched because its expensive af",
		CreatedDate: time.Now(),
	}

	crop1 := domain.Crop{
		UID:          uid,
		InitialArea:  area1,
		CurrentAreas: []domain.Area{area1},
		Type:         domain.Nursery{},
		Container:    domain.CropContainer{Quantity: 5, Type: domain.Tray{Cell: 10}},
		Notes:        cropNotes,
		CreatedDate:  time.Now(),
	}

	server.CropService.ChangeInventory(&crop1, inventory1)
	cropStorage.CropMap[uid] = crop1

	uid, _ = uuid.NewV4()
	crop2 := domain.Crop{
		UID:          uid,
		InitialArea:  area1,
		CurrentAreas: []domain.Area{area1},
		Type:         domain.Nursery{},
		Container:    domain.CropContainer{Quantity: 10, Type: domain.Pot{}},
		CreatedDate:  time.Now(),
	}

	server.CropService.ChangeInventory(&crop2, inventory2)
	cropStorage.CropMap[uid] = crop2

	uid, _ = uuid.NewV4()
	crop3 := domain.Crop{
		UID:          uid,
		InitialArea:  area1,
		CurrentAreas: []domain.Area{area2},
		Type:         domain.Growing{},
		CreatedDate:  time.Now(),
	}

	server.CropService.ChangeInventory(&crop3, inventory3)
	cropStorage.CropMap[uid] = crop3
}
