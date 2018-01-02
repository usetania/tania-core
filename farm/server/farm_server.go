package server

import (
	"net/http"

	"github.com/Tanibox/tania-server/farm/entity"
	"github.com/Tanibox/tania-server/farm/repository"
	"github.com/Tanibox/tania-server/helper/stringhelper"
	"github.com/labstack/echo"
)

// FarmServer ties the routes and handlers with injected dependencies
type FarmServer struct {
	FarmRepo      repository.FarmRepository
	ReservoirRepo repository.ReservoirRepository
	AreaRepo      repository.AreaRepository
	File          File
}

// NewFarmServer initializes FarmServer's dependencies and create new FarmServer struct
func NewFarmServer() (*FarmServer, error) {
	return &FarmServer{
		FarmRepo:      repository.NewFarmRepositoryInMemory(),
		ReservoirRepo: repository.NewReservoirRepositoryInMemory(),
		AreaRepo:      repository.NewAreaRepositoryInMemory(),
		File:          InitLocalFile(),
	}, nil
}

// Mount defines the FarmServer's endpoints with its handlers
func (s *FarmServer) Mount(g *echo.Group) {
	g.GET("/types", s.GetTypes)

	g.POST("", s.SaveFarm)
	g.GET("", s.FindAllFarm)
	g.GET("/:id", s.FindFarmByID)
	g.POST("/:id/reservoirs", s.SaveReservoir)
	g.GET("/:id/reservoirs", s.GetFarmReservoirs)
	g.POST("/:id/areas", s.SaveArea)
	g.GET("/:id/areas", s.GetFarmAreas)
	g.GET("/:farm_id/areas/:area_id/files/:file_id", s.GetAreaPhoto)
}

// GetTypes is a FarmServer's handle to get farm types
func (s *FarmServer) GetTypes(c echo.Context) error {
	types := entity.FindAllFarmTypes()

	return c.JSON(http.StatusOK, types)
}

func (s FarmServer) FindAllFarm(c echo.Context) error {
	data := make(map[string][]SimpleFarm)

	result := <-s.FarmRepo.FindAll()
	if result.Error != nil {
		return result.Error
	}

	farms, ok := result.Result.([]entity.Farm)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	data["data"] = MapToSimpleFarm(farms)

	return c.JSON(http.StatusOK, data)
}

// SaveFarm is a FarmServer's handler to save new Farm
func (s *FarmServer) SaveFarm(c echo.Context) error {
	data := make(map[string]string)

	farm, err := entity.CreateFarm(c.FormValue("name"), c.FormValue("farm_type"))
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

	farm.UID = repository.GetRandomUID()

	result := <-s.FarmRepo.Save(&farm)

	if result.Error != nil {
		return result.Error
	}

	uid, ok := result.Result.(string)

	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	data["data"] = uid

	return c.JSON(http.StatusOK, data)
}

func (s *FarmServer) FindFarmByID(c echo.Context) error {
	data := make(map[string]entity.Farm)

	result := <-s.FarmRepo.FindByID(c.Param("id"))
	if result.Error != nil {
		return result.Error
	}

	farm, ok := result.Result.(entity.Farm)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	data["data"] = farm

	return c.JSON(http.StatusOK, data)
}

// SaveReservoir is a FarmServer's handler to save new Reservoir and place it to a Farm
func (s *FarmServer) SaveReservoir(c echo.Context) error {
	data := make(map[string]string)
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
	r, err := entity.CreateReservoir(farm, name)
	if err != nil {
		return Error(c, err)
	}

	if waterSourceType == "bucket" {
		b, err := entity.CreateBucket(capacity, 0)
		if err != nil {
			return Error(c, err)
		}

		r.AttachBucket(&b)
	} else if waterSourceType == "tap" {
		t, err := entity.CreateTap()
		if err != nil {
			return Error(c, err)
		}

		r.AttachTap(&t)
	}

	r.UID = repository.GetRandomUID()

	err = farm.AddReservoir(&r)
	if err != nil {
		return Error(c, err)
	}

	// Persists //
	reservoirResult := <-s.ReservoirRepo.Save(&r)
	if reservoirResult.Error != nil {
		return reservoirResult.Error
	}

	uid, ok := reservoirResult.Result.(string)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	farmResult := <-s.FarmRepo.Save(&farm)
	if farmResult.Error != nil {
		return farmResult.Error
	}

	data["data"] = uid

	return c.JSON(http.StatusOK, data)
}

func (s *FarmServer) GetFarmReservoirs(c echo.Context) error {
	data := make(map[string][]entity.Reservoir)

	result := <-s.FarmRepo.FindByID(c.Param("id"))
	if result.Error != nil {
		return result.Error
	}

	farm, ok := result.Result.(entity.Farm)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	data["data"] = farm.Reservoirs
	if len(farm.Reservoirs) == 0 {
		data["data"] = []entity.Reservoir{}
	}

	return c.JSON(http.StatusOK, data)
}

func (s *FarmServer) SaveArea(c echo.Context) error {
	data := make(map[string]string)
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

	size, err := validation.ValidateAreaSize(c.FormValue("size"))
	if err != nil {
		return Error(c, err)
	}

	sizeUnit, err := validation.ValidateAreaSizeUnit(c.FormValue("size_unit"))
	if err != nil {
		return Error(c, err)
	}

	location, err := validation.ValidateAreaLocation(c.FormValue("location"))
	if err != nil {
		return Error(c, err)
	}

	// Process //
	area, err := entity.CreateArea(farm, c.FormValue("name"), c.FormValue("type"))
	if err != nil {
		return Error(c, err)
	}

	area.UID = repository.GetRandomUID()
	area.SizeUnit = sizeUnit
	area.Location = location

	err = area.ChangeSize(size)
	if err != nil {
		return Error(c, err)
	}

	// TODO: validate location first
	err = area.ChangeLocation(c.FormValue("location"))
	if err != nil {
		return Error(c, err)
	}

	photo, err := c.FormFile("photo")
	if err == nil {
		destPath := stringhelper.Join(s.File.GetAreaFilepath(), photo.Filename)
		err = s.File.Upload(photo, destPath)

		if err != nil {
			return Error(c, err)
		}

		areaPhoto := entity.AreaPhoto{
			Filename: photo.Filename,
			MimeType: photo.Header["Content-Type"][0],
			Size:     int(photo.Size),
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
	areaResult := <-s.AreaRepo.Save(&area)
	if areaResult.Error != nil {
		return areaResult.Error
	}

	uid, ok := areaResult.Result.(string)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	farmResult := <-s.FarmRepo.Save(&farm)
	if farmResult.Error != nil {
		return farmResult.Error
	}

	data["data"] = uid

	return c.JSON(http.StatusOK, data)
}

func (s *FarmServer) GetFarmAreas(c echo.Context) error {
	data := make(map[string][]entity.Area)

	result := <-s.FarmRepo.FindByID(c.Param("id"))
	if result.Error != nil {
		return result.Error
	}

	farm, ok := result.Result.(entity.Farm)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	data["data"] = farm.Areas
	if len(farm.Areas) == 0 {
		data["data"] = []entity.Area{}
	}

	return c.JSON(http.StatusOK, data)
}

func (s *FarmServer) GetAreaPhoto(c echo.Context) error {
	// Validate //
	result := <-s.FarmRepo.FindByID(c.Param("farm_id"))
	if result.Error != nil {
		return result.Error
	}

	_, ok := result.Result.(entity.Farm)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	result = <-s.AreaRepo.FindByID(c.Param("area_id"))
	if result.Error != nil {
		return result.Error
	}

	area, ok := result.Result.(entity.Area)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	if area.Photo.Filename == "" {
		return Error(c, NewRequestValidationError(NOT_FOUND, "photo"))
	}

	// Process //
	srcPath := stringhelper.Join(s.File.GetAreaFilepath(), area.Photo.Filename)

	return c.File(srcPath)
}
