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
	return nil
}
