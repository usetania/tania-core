// Package server contains routes, handlers and dependency injection of repository.
package server

import (
	"net/http"
	"strconv"

	"github.com/Tanibox/tania-server/reservoir"
	"github.com/Tanibox/tania-server/reservoir/repository"
	"github.com/Tanibox/tania-server/reservoir/validation"
	"github.com/labstack/echo"
)

// ReservoirServer ties the routes and handlers with injected dependencies
type ReservoirServer struct {
	ReservoirRepo repository.ReservoirRepository
}

// NewServer initializes ReservoirServer's dependencies and create new ReservoirServer struct
func NewServer() (*ReservoirServer, error) {
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
	data := make(map[string][]reservoir.Reservoir)

	result := <-s.ReservoirRepo.FindAll()

	if result.Error != nil {
		return result.Error
	}

	reservoirs, ok := result.Result.([]reservoir.Reservoir)

	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	data["data"] = reservoirs

	return c.JSON(http.StatusOK, data)
}

// Save is a ReservoirServer's handler to save new Reservoir
func (s *ReservoirServer) Save(c echo.Context) error {
	data := make(map[string]string)
	validation := validation.RequestValidation{}

	// Validate requests //
	name, err := validation.ValidateName(c.FormValue("name"))
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

	// Process //
	r, err := reservoir.CreateReservoir(name)
	if err != nil {
		return Error(c, err)
	}

	if waterSourceType == "bucket" {
		b, err := reservoir.CreateBucket(capacity, 0)
		if err != nil {
			return Error(c, err)
		}

		r.AttachBucket(&b)
	} else if waterSourceType == "tap" {
		t, err := reservoir.CreateTap()
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

func Error(c echo.Context, err error) error {
	if re, ok := err.(reservoir.ReservoirError); ok {
		errMap := map[string]string{
			"field_name":    "",
			"error_code":    strconv.Itoa(re.Code),
			"error_message": re.Error(),
		}

		return c.JSON(http.StatusBadRequest, errMap)
	} else if rve, ok := err.(validation.RequestValidationError); ok {
		return c.JSON(http.StatusBadRequest, rve)
	}

	return c.JSON(http.StatusInternalServerError, err)
}
