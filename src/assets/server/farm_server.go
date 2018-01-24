package server

import (
	"net/http"

	"github.com/Tanibox/tania-server/config"
	"github.com/Tanibox/tania-server/src/assets/domain"
	"github.com/Tanibox/tania-server/src/assets/query"
	"github.com/Tanibox/tania-server/src/assets/query/inmemory"
	"github.com/Tanibox/tania-server/src/assets/repository"
	"github.com/Tanibox/tania-server/src/assets/storage"
	growthstorage "github.com/Tanibox/tania-server/src/growth/storage"
	"github.com/Tanibox/tania-server/src/helper/imagehelper"
	"github.com/Tanibox/tania-server/src/helper/stringhelper"
	"github.com/labstack/echo"
)

// FarmServer ties the routes and handlers with injected dependencies
type FarmServer struct {
	FarmRepo               repository.FarmRepository
	ReservoirRepo          repository.ReservoirRepository
	AreaRepo               repository.AreaRepository
	AreaQuery              query.AreaQuery
	InventoryMaterialRepo  repository.InventoryMaterialRepository
	InventoryMaterialQuery query.InventoryMaterialQuery
	CropQuery              query.CropQuery
	File                   File
}

// NewFarmServer initializes FarmServer's dependencies and create new FarmServer struct
func NewFarmServer(
	farmStorage *storage.FarmStorage,
	areaStorage *storage.AreaStorage,
	reservoirStorage *storage.ReservoirStorage,
	inventoryMaterialStorage *storage.InventoryMaterialStorage,
	cropStorage *growthstorage.CropStorage,
) (*FarmServer, error) {
	farmRepo := repository.NewFarmRepositoryInMemory(farmStorage)

	areaRepo := repository.NewAreaRepositoryInMemory(areaStorage)
	areaQuery := inmemory.NewAreaQueryInMemory(areaStorage)

	reservoirRepo := repository.NewReservoirRepositoryInMemory(reservoirStorage)

	inventoryMaterialRepo := repository.NewInventoryMaterialRepositoryInMemory(inventoryMaterialStorage)
	inventoryMaterialQuery := inmemory.NewInventoryMaterialQueryInMemory(inventoryMaterialStorage)

	cropQuery := inmemory.NewCropQueryInMemory(cropStorage)

	farmServer := FarmServer{
		FarmRepo:               farmRepo,
		ReservoirRepo:          reservoirRepo,
		AreaRepo:               areaRepo,
		AreaQuery:              areaQuery,
		InventoryMaterialRepo:  inventoryMaterialRepo,
		InventoryMaterialQuery: inventoryMaterialQuery,
		CropQuery:              cropQuery,
		File:                   LocalFile{},
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
	plantTypes := MapToPlantType(domain.GetPlantTypes())

	return c.JSON(http.StatusOK, plantTypes)
}

func (s *FarmServer) SaveInventory(c echo.Context) error {
	data := make(map[string]InventoryMaterial)

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

	data["data"] = MapToInventoryMaterial(inventoryMaterial)

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
