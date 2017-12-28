package server

import (
	"net/http"

	"github.com/Tanibox/tania-server/farm/entity"
	"github.com/Tanibox/tania-server/farm/repository"
	"github.com/labstack/echo"
)

// FarmServer ties the routes and handlers with injected dependencies
type FarmServer struct {
	FarmRepo      repository.FarmRepository
	ReservoirRepo repository.ReservoirRepository
}

// NewFarmServer initializes FarmServer's dependencies and create new FarmServer struct
func NewFarmServer() (*FarmServer, error) {
	return &FarmServer{
		FarmRepo:      repository.NewFarmRepositoryInMemory(),
		ReservoirRepo: repository.NewReservoirRepositoryInMemory(),
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
}

// GetTypes is a FarmServer's handle to get farm types
func (s *FarmServer) GetTypes(c echo.Context) error {
	types := entity.FindAllFarmTypes()

	return c.JSON(http.StatusOK, types)
}

func (s FarmServer) FindAllFarm(c echo.Context) error {
	data := make(map[string][]entity.Farm)

	result := <-s.FarmRepo.FindAll()
	if result.Error != nil {
		return result.Error
	}

	farms, ok := result.Result.([]entity.Farm)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	data["data"] = farms

	return c.JSON(http.StatusOK, data)
}

// SaveFarm is a FarmServer's handler to save new Farm
func (s *FarmServer) SaveFarm(c echo.Context) error {
	data := make(map[string]string)

	r, err := entity.CreateFarm(
		c.FormValue("name"),
		c.FormValue("description"),
		c.FormValue("latitude"),
		c.FormValue("longitude"),
		c.FormValue("farm_type"),
		c.FormValue("country_code"),
		c.FormValue("city_code"),
	)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	result := <-s.FarmRepo.Save(&r)

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
	return nil
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

	err = farm.AddReservoir(r)
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

	farmResult := <-s.FarmRepo.Update(&farm)
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

	return c.JSON(http.StatusOK, data)
}
