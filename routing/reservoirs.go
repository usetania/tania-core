// Package routing provides the list of functions for each HTTP routing
package routing

import (
	"../reservoir"
	"github.com/labstack/echo"
	"net/http"
)

type ReservoirsData struct {
	Data []reservoir.Reservoir
}

// ReservoirsIndex displays all available reservoirs in a farm
func ReservoirsIndex(c echo.Context) error {
	farmID := c.QueryParam("farm_id")

	if farmID == "" {
		return c.JSON(http.StatusOK, RequestError{ErrorFarmID})
	}

	reservoirsJSON := reservoir.DisplayAll()

	return c.JSON(http.StatusOK, ReservoirsData{
		Data: reservoirsJSON,
	})
}
