package server

import (
	"net/http"

	"github.com/Tanibox/tania-server/farm/entity"
	"github.com/Tanibox/tania-server/farm/repository"
	"github.com/labstack/echo"
)

// FarmServer ties the routes and handlers with injected dependencies
type FarmServer struct {
	FarmRepo repository.FarmRepository
}

// NewFarmServer initializes FarmServer's dependencies and create new FarmServer struct
func NewFarmServer() (*FarmServer, error) {
	farmRepo := repository.NewFarmRepositoryInMemory()

	return &FarmServer{
		FarmRepo: farmRepo,
	}, nil
}

// Mount defines the FarmServer's endpoints with its handlers
func (s *FarmServer) Mount(g *echo.Group) {
	g.GET("/types", s.GetTypes)

	g.POST("", s.Save)
}

// GetTypes is a FarmServer's handle to get farm types
func (s *FarmServer) GetTypes(c echo.Context) error {
	types := entity.FindAllFarmTypes()

	return c.JSON(http.StatusOK, types)
}

// Save is a FarmServer's handler to save new Reservoir
func (s *FarmServer) Save(c echo.Context) error {
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
