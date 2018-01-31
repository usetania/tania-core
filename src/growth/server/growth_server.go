package server

import (
	"net/http"
	"strconv"

	"github.com/Tanibox/tania-server/config"
	"github.com/Tanibox/tania-server/src/growth/domain"
	"github.com/Tanibox/tania-server/src/growth/domain/service"
	"github.com/Tanibox/tania-server/src/growth/query/inmemory"
	"github.com/Tanibox/tania-server/src/helper/imagehelper"
	"github.com/Tanibox/tania-server/src/helper/stringhelper"

	assetsstorage "github.com/Tanibox/tania-server/src/assets/storage"
	"github.com/Tanibox/tania-server/src/growth/query"
	"github.com/Tanibox/tania-server/src/growth/repository"
	"github.com/Tanibox/tania-server/src/growth/storage"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
)

// GrowthServer ties the routes and handlers with injected dependencies
type GrowthServer struct {
	CropRepo               repository.CropRepository
	CropQuery              query.CropQuery
	CropService            domain.CropService
	AreaQuery              query.AreaQuery
	InventoryMaterialQuery query.InventoryMaterialQuery
	FarmQuery              query.FarmQuery
	File                   File
}

// NewGrowthServer initializes GrowthServer's dependencies and create new GrowthServer struct
func NewGrowthServer(
	cropStorage *storage.CropStorage,
	areaStorage *assetsstorage.AreaStorage,
	materialStorage *assetsstorage.MaterialStorage,
	farmStorage *assetsstorage.FarmStorage,
) (*GrowthServer, error) {
	cropRepo := repository.NewCropRepositoryInMemory(cropStorage)
	cropQuery := inmemory.NewCropQueryInMemory(cropStorage)

	areaQuery := inmemory.NewAreaQueryInMemory(areaStorage)
	inventoryMaterialQuery := inmemory.NewInventoryMaterialQueryInMemory(materialStorage)
	farmQuery := inmemory.NewFarmQueryInMemory(farmStorage)

	cropService := service.CropServiceInMemory{
		InventoryMaterialQuery: inventoryMaterialQuery,
		CropQuery:              cropQuery,
		AreaQuery:              areaQuery,
	}

	return &GrowthServer{
		CropRepo:               cropRepo,
		CropQuery:              cropQuery,
		CropService:            cropService,
		AreaQuery:              areaQuery,
		InventoryMaterialQuery: inventoryMaterialQuery,
		FarmQuery:              farmQuery,
		File:                   LocalFile{},
	}, nil
}

// Mount defines the GrowthServer's endpoints with its handlers
func (s *GrowthServer) Mount(g *echo.Group) {
	g.GET("/:id/crops", s.FindAllCrops)
	g.GET("/areas/:id/crops", s.FindAllCropsByArea)
	g.POST("/areas/:id/crops", s.SaveAreaCropBatch)
	g.GET("/crops/:id", s.FindCropByID)
	g.POST("/crops/:id/move", s.MoveCrop)
	g.POST("/crops/:id/harvest", s.HarvestCrop)
	g.POST("/crops/:id/dump", s.DumpCrop)
	g.POST("/crops/:id/notes", s.SaveCropNotes)
	g.DELETE("/crops/:crop_id/notes/:note_id", s.RemoveCropNotes)
	g.POST("/crops/:id/photos", s.UploadCropPhotos)
	g.GET("/crops/:id/photos", s.GetCropPhotos)

}

func (s *GrowthServer) SaveAreaCropBatch(c echo.Context) error {
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
	areaUID, err := uuid.FromString(areaID)
	if err != nil {
		return Error(c, NewRequestValidationError(NOT_FOUND, "id"))
	}

	areaResult := <-s.AreaQuery.FindByID(areaUID)
	if areaResult.Error != nil {
		return Error(c, areaResult.Error)
	}

	area, ok := areaResult.Result.(query.CropAreaQueryResult)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusBadRequest, "Internal server error"))
	}

	queryResult := <-s.InventoryMaterialQuery.FindInventoryByPlantTypeCodeAndVariety(plantType, variety)
	if queryResult.Error != nil {
		return Error(c, queryResult.Error)
	}

	inventoryMaterial, ok := queryResult.Result.(query.CropInventoryQueryResult)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusBadRequest, "Internal server error"))
	}

	var containerT domain.CropContainerType
	switch containerType {
	case domain.Tray{}.Code():
		containerT = domain.Tray{Cell: containerCell}
	case domain.Pot{}.Code():
		containerT = domain.Pot{}
	default:
		return Error(c, NewRequestValidationError(NOT_FOUND, "container_type"))
	}

	// Process //
	cropBatch, err := domain.CreateCropBatch(
		s.CropService,
		area.UID,
		cropType,
		inventoryMaterial.UID,
		containerQuantity,
		containerT,
	)
	if err != nil {
		return Error(c, err)
	}

	// Persists //
	err = <-s.CropRepo.Save(&cropBatch)
	if err != nil {
		return Error(c, err)
	}

	data := make(map[string]CropBatch)
	cb, err := MapToCropBatch(s, cropBatch)
	if err != nil {
		return Error(c, err)
	}

	data["data"] = cb

	return c.JSON(http.StatusOK, data)
}

func (s *GrowthServer) FindCropByID(c echo.Context) error {
	// Validate //
	result := <-s.CropRepo.FindByID(c.Param("id"))
	if result.Error != nil {
		return Error(c, result.Error)
	}

	crop, ok := result.Result.(domain.Crop)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusBadRequest, "Internal server error"))
	}

	if crop.UID == (uuid.UUID{}) {
		return Error(c, NewRequestValidationError(NOT_FOUND, "id"))
	}

	data := make(map[string]CropBatch)
	cropBatch, err := MapToCropBatch(s, crop)
	if err != nil {
		return Error(c, err)
	}
	data["data"] = cropBatch

	return c.JSON(http.StatusOK, data)
}

func (s *GrowthServer) MoveCrop(c echo.Context) error {
	cropID := c.Param("id")
	srcAreaID := c.FormValue("source_area_id")
	dstAreaID := c.FormValue("destination_area_id")
	quantity := c.FormValue("quantity")

	// VALIDATE //
	result := <-s.CropRepo.FindByID(cropID)
	if result.Error != nil {
		return Error(c, result.Error)
	}

	crop, ok := result.Result.(domain.Crop)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusBadRequest, "Internal server error"))
	}

	if crop.UID == (uuid.UUID{}) {
		return Error(c, NewRequestValidationError(NOT_FOUND, "id"))
	}

	srcAreaUID, err := uuid.FromString(srcAreaID)
	if err != nil {
		return Error(c, err)
	}

	dstAreaUID, err := uuid.FromString(dstAreaID)
	if err != nil {
		return Error(c, err)
	}

	qty, err := strconv.Atoi(quantity)
	if err != nil {
		return Error(c, err)
	}

	// PROCESS //
	err = crop.MoveToArea(s.CropService, srcAreaUID, dstAreaUID, qty)
	if err != nil {
		return Error(c, err)
	}

	// PERSIST //
	err = <-s.CropRepo.Save(&crop)
	if err != nil {
		return Error(c, err)
	}

	data := make(map[string]CropBatch)
	cropBatch, err := MapToCropBatch(s, crop)
	if err != nil {
		return Error(c, err)
	}
	data["data"] = cropBatch

	return c.JSON(http.StatusOK, data)
}

func (s *GrowthServer) HarvestCrop(c echo.Context) error {
	cropID := c.Param("id")
	srcAreaID := c.FormValue("source_area_id")
	harvestType := c.FormValue("harvest_type")
	producedQuantity := c.FormValue("produced_quantity")
	producedUnit := c.FormValue("produced_unit")

	// VALIDATE //
	result := <-s.CropRepo.FindByID(cropID)
	if result.Error != nil {
		return Error(c, result.Error)
	}

	crop, ok := result.Result.(domain.Crop)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusBadRequest, "Internal server error"))
	}

	if crop.UID == (uuid.UUID{}) {
		return Error(c, NewRequestValidationError(NOT_FOUND, "id"))
	}

	srcAreaUID, err := uuid.FromString(srcAreaID)
	if err != nil {
		return Error(c, err)
	}

	ht := domain.GetHarvestType(harvestType)
	if ht == (domain.HarvestType{}) {
		return Error(c, NewRequestValidationError(INVALID_OPTION, "harvest_type"))
	}

	if producedQuantity == "" {
		return Error(c, NewRequestValidationError(REQUIRED, "produced_quantity"))
	}

	prodQty, err := strconv.ParseFloat(producedQuantity, 32)
	if err != nil {
		return Error(c, err)
	}

	prodUnit := domain.GetProducedUnit(producedUnit)
	if prodUnit == (domain.ProducedUnit{}) {
		return Error(c, NewRequestValidationError(INVALID_OPTION, "produced_unit"))
	}

	// PROCESS //
	err = crop.Harvest(s.CropService, srcAreaUID, harvestType, float32(prodQty), prodUnit)
	if err != nil {
		return Error(c, err)
	}

	// PERSIST //
	err = <-s.CropRepo.Save(&crop)
	if err != nil {
		return Error(c, err)
	}

	data := make(map[string]CropBatch)
	cropBatch, err := MapToCropBatch(s, crop)
	if err != nil {
		return Error(c, err)
	}
	data["data"] = cropBatch

	return c.JSON(http.StatusOK, data)
}

func (s *GrowthServer) DumpCrop(c echo.Context) error {
	cropID := c.Param("id")
	srcAreaID := c.FormValue("source_area_id")
	quantity := c.FormValue("quantity")

	// VALIDATE //
	result := <-s.CropRepo.FindByID(cropID)
	if result.Error != nil {
		return Error(c, result.Error)
	}

	crop, ok := result.Result.(domain.Crop)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusBadRequest, "Internal server error"))
	}

	if crop.UID == (uuid.UUID{}) {
		return Error(c, NewRequestValidationError(NOT_FOUND, "id"))
	}

	srcAreaUID, err := uuid.FromString(srcAreaID)
	if err != nil {
		return Error(c, err)
	}

	qty, err := strconv.Atoi(quantity)
	if err != nil {
		return Error(c, err)
	}

	// PROCESS //
	err = crop.Dump(s.CropService, srcAreaUID, qty)
	if err != nil {
		return Error(c, err)
	}

	// PERSIST //
	err = <-s.CropRepo.Save(&crop)
	if err != nil {
		return Error(c, err)
	}

	data := make(map[string]CropBatch)
	cropBatch, err := MapToCropBatch(s, crop)
	if err != nil {
		return Error(c, err)
	}
	data["data"] = cropBatch

	return c.JSON(http.StatusOK, data)
}

func (s *GrowthServer) SaveCropNotes(c echo.Context) error {
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

	data := make(map[string]CropBatch)
	cropBatch, err := MapToCropBatch(s, crop)
	if err != nil {
		return Error(c, err)
	}
	data["data"] = cropBatch

	return c.JSON(http.StatusOK, data)
}

func (s *GrowthServer) RemoveCropNotes(c echo.Context) error {
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
	noteUID, err := uuid.FromString(noteID)
	if err != nil {
		return Error(c, err)
	}

	err = crop.RemoveNote(noteUID)
	if err != nil {
		return Error(c, err)
	}

	// Persists //
	resultSave := <-s.CropRepo.Save(&crop)
	if resultSave != nil {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	data := make(map[string]CropBatch)
	cropBatch, err := MapToCropBatch(s, crop)
	if err != nil {
		return Error(c, err)
	}
	data["data"] = cropBatch

	return c.JSON(http.StatusOK, data)
}

// TODO: The crops should be found by its Farm
func (s *GrowthServer) FindAllCrops(c echo.Context) error {
	data := make(map[string][]CropBatch)

	// Params //
	farmID := c.Param("id")

	// Validate //
	farmUID, err := uuid.FromString(farmID)
	if err != nil {
		return Error(c, err)
	}

	result := <-s.FarmQuery.FindByID(farmUID)
	if result.Error != nil {
		return Error(c, result.Error)
	}

	farm, ok := result.Result.(query.CropFarmQueryResult)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	// Process //
	resultQuery := <-s.CropQuery.FindAllCropsByFarm(farm.UID)
	if resultQuery.Error != nil {
		return Error(c, resultQuery.Error)
	}

	crops, ok := resultQuery.Result.([]domain.Crop)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	data["data"] = []CropBatch{}
	for _, v := range crops {
		cropBatch, err := MapToCropBatch(s, v)
		if err != nil {
			return Error(c, err)
		}
		data["data"] = append(data["data"], cropBatch)
	}

	return c.JSON(http.StatusOK, data)
}

func (s *GrowthServer) FindAllCropsByArea(c echo.Context) error {
	data := make(map[string][]CropBatch)

	// Params //
	areaID := c.Param("id")

	// Validate //
	areaUID, err := uuid.FromString(areaID)
	if err != nil {
		return Error(c, err)
	}

	result := <-s.AreaQuery.FindByID(areaUID)
	if result.Error != nil {
		return Error(c, result.Error)
	}

	area, ok := result.Result.(query.CropAreaQueryResult)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	// Process //
	resultQuery := <-s.CropQuery.FindAllCropsByArea(area.UID)
	if resultQuery.Error != nil {
		return Error(c, resultQuery.Error)
	}

	crops, ok := resultQuery.Result.([]domain.Crop)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	data["data"] = []CropBatch{}
	for _, v := range crops {
		cropBatch, err := MapToCropBatch(s, v)
		if err != nil {
			return Error(c, err)
		}
		data["data"] = append(data["data"], cropBatch)
	}

	return c.JSON(http.StatusOK, data)
}

func (s *GrowthServer) UploadCropPhotos(c echo.Context) error {
	// Validate //
	photo, err := c.FormFile("photo")
	if err != nil {
		return Error(c, NewRequestValidationError(REQUIRED, "photo"))
	}

	result := <-s.CropRepo.FindByID(c.Param("id"))
	if result.Error != nil {
		return Error(c, result.Error)
	}

	crop, ok := result.Result.(domain.Crop)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusBadRequest, "Internal server error"))
	}

	if crop.UID == (uuid.UUID{}) {
		return Error(c, NewRequestValidationError(NOT_FOUND, "id"))
	}

	// Process
	destPath := stringhelper.Join(*config.Config.UploadPathCrop, "/", photo.Filename)
	err = s.File.Upload(photo, destPath)

	if err != nil {
		return Error(c, err)
	}

	width, height, err := imagehelper.GetImageDimension(destPath)
	if err != nil {
		return Error(c, err)
	}

	cropPhoto := domain.CropPhoto{
		Filename: photo.Filename,
		MimeType: photo.Header["Content-Type"][0],
		Size:     int(photo.Size),
		Width:    width,
		Height:   height,
	}

	crop.Photo = cropPhoto

	// Persists //
	resultSave := <-s.CropRepo.Save(&crop)
	if resultSave != nil {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	data := make(map[string]CropBatch)
	cropBatch, err := MapToCropBatch(s, crop)
	if err != nil {
		return Error(c, err)
	}
	data["data"] = cropBatch

	return c.JSON(http.StatusOK, data)
}

func (s *GrowthServer) GetCropPhotos(c echo.Context) error {
	// Validate //
	result := <-s.CropRepo.FindByID(c.Param("id"))
	if result.Error != nil {
		return Error(c, result.Error)
	}

	crop, ok := result.Result.(domain.Crop)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusBadRequest, "Internal server error"))
	}

	if crop.Photo.Filename == "" {
		return Error(c, NewRequestValidationError(NOT_FOUND, "photo"))
	}

	// Process //
	srcPath := stringhelper.Join(*config.Config.UploadPathCrop, "/", crop.Photo.Filename)

	return c.File(srcPath)
}
