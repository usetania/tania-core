package server

import (
	"net/http"
	"strconv"

	"github.com/Tanibox/tania-server/config"
	"github.com/Tanibox/tania-server/src/assets/domain"
	"github.com/Tanibox/tania-server/src/assets/query"
	"github.com/Tanibox/tania-server/src/assets/repository"
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

	inventoryMaterialStorage := storage.InventoryMaterialStorage{InventoryMaterialMap: make(map[uuid.UUID]domain.InventoryMaterial)}
	inventoryMaterialRepo := repository.NewInventoryMaterialRepositoryInMemory(&inventoryMaterialStorage)
	inventoryMaterialQuery := query.NewInventoryMaterialQueryInMemory(&inventoryMaterialStorage)

	return &FarmServer{
		FarmRepo:               farmRepo,
		ReservoirRepo:          reservoirRepo,
		AreaRepo:               areaRepo,
		AreaQuery:              areaQuery,
		InventoryMaterialRepo:  inventoryMaterialRepo,
		InventoryMaterialQuery: inventoryMaterialQuery,
		File: LocalFile{},
	}, nil
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
	g.GET("/:id/reservoirs", s.GetFarmReservoirs)
	g.GET("/:farm_id/reservoirs/:reservoir_id", s.GetReservoirsByID)
	g.POST("/:id/areas", s.SaveArea)
	g.GET("/:id/areas", s.GetFarmAreas)
	g.POST("/areas/:id/crops", s.SaveAreaCropBatch)
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

	result := <-s.FarmRepo.Save(&farm)

	if result.Error != nil {
		return Error(c, result.Error)
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
	reservoirResult := <-s.ReservoirRepo.Save(&r)
	if reservoirResult.Error != nil {
		return Error(c, reservoirResult.Error)
	}

	farmResult := <-s.FarmRepo.Save(&farm)
	if farmResult.Error != nil {
		return Error(c, farmResult.Error)
	}

	detailReservoir, err := MapToDetailReservoir(s, r)
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
	reservoirResult := <-s.ReservoirRepo.Save(&reservoir)
	if reservoirResult.Error != nil {
		return Error(c, reservoirResult.Error)
	}

	areaResult := <-s.AreaRepo.Save(&area)
	if areaResult.Error != nil {
		return Error(c, areaResult.Error)
	}

	farmResult := <-s.FarmRepo.Save(&farm)
	if farmResult.Error != nil {
		return Error(c, farmResult.Error)
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
	result := <-s.InventoryMaterialRepo.Save(&inventoryMaterial)
	if result.Error != nil {
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
	repoResult := <-s.AreaRepo.FindByID(areaID)
	if repoResult.Error != nil {
		return Error(c, repoResult.Error)
	}

	area, ok := repoResult.Result.(domain.Area)
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

	err = cropBatch.ChangeInventory(inventoryMaterial)
	if err != nil {
		return Error(c, err)
	}

	err = cropBatch.ChangeContainer(cropContainer)
	if err != nil {
		return Error(c, err)
	}

	return c.JSON(http.StatusOK, cropBatch)
}
