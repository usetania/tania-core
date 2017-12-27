// Package server contains routes, handlers and dependency injection of repository.
package server

import (
	"net/http"

	"github.com/Tanibox/tania-server/farm/entity"
	"github.com/Tanibox/tania-server/farm/repository"
	"github.com/labstack/echo"
)

// ReservoirServer ties the routes and handlers with injected dependencies
type ReservoirServer struct {
	ReservoirRepo repository.ReservoirRepository
}

// NewReservoirServer initializes ReservoirServer's dependencies and create new ReservoirServer struct
func NewReservoirServer() (*ReservoirServer, error) {
	reservoirRepo := repository.NewReservoirRepositoryInMemory()

	return &ReservoirServer{
		ReservoirRepo: reservoirRepo,
	}, nil
}

// Mount defines the ReservoirServer's endpoints with its handlers
func (s *ReservoirServer) Mount(g *echo.Group) {
	g.GET("", s.FindAll)
	g.POST("", s.Save)
}

// FindAll is a ResevoirServer's handler to get all Reservoir
func (s *ReservoirServer) FindAll(c echo.Context) error {
	data := make(map[string][]entity.Reservoir)

	result := <-s.ReservoirRepo.FindAll()

	if result.Error != nil {
		return result.Error
	}

	reservoirs, ok := result.Result.([]entity.Reservoir)

	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	data["data"] = reservoirs

	return c.JSON(http.StatusOK, data)
}

// Save is a ReservoirServer's handler to save new Reservoir
func (s *ReservoirServer) Save(c echo.Context) error {
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

	farm, err := validation.ValidateFarm(c.FormValue("farm_id"))
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

	// Persists //
	result := <-s.ReservoirRepo.Save(&r)

	if result.Error != nil {
		return result.Error
	}

	uid, ok := result.Result.(string)

	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	data["data"] = uid

	return c.JSON(http.StatusOK, data)
}
