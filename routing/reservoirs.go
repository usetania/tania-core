// Package routing provides the list of functions for each HTTP routing
package routing

import (
	"net/http"

	"github.com/Tanibox/tania-server/reservoir"
	"github.com/labstack/echo"
)

type ReservoirsData struct {
	Data []reservoir.Reservoir
}

func ReservoirsRouter(g *echo.Group) {
	g.GET("/", ReservoirsIndex)
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
