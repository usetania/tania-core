// Package server contains routes, handlers and dependency injection of repository.
package server

import (
	"net/http"
	"strconv"

	"github.com/Tanibox/tania-server/reservoir"
	"github.com/Tanibox/tania-server/reservoir/repository"
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

	r, err := reservoir.CreateReservoir(c.FormValue("name"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if c.FormValue("type") == "bucket" {
		capacity, err := strconv.ParseFloat(c.FormValue("capacity"), 32)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		b, err := reservoir.CreateBucket(float32(capacity), 0)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		r.AttachBucket(&b)
	} else if c.FormValue("type") == "tap" {
		t, err := reservoir.CreateTap()
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		r.AttachTap(&t)
	}

	result := <-s.ReservoirRepo.Save(&r)

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
