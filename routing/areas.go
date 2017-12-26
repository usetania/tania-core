// Package routing provides the list of functions for each HTTP routing
package routing

import (
	"net/http"

	"github.com/Tanibox/tania-server/area"
	"github.com/labstack/echo"
)

type AreasData struct {
	Data []area.Area
}

func AreasRouter(g *echo.Group) {
	g.GET("/", AreasIndex)
	g.POST("/", AreasCreate)
	g.GET("/:uid", AreasHarvest)
	g.PUT("/:uid", AreasDispose)
	g.DELETE("/:uid", AreasDestroy)
}

// AreasIndex displays all available areas in a farm
func AreasIndex(c echo.Context) error {
	farmID := c.QueryParam("farm_id")

	if farmID == "" {
		return c.JSON(http.StatusOK, RequestError{ErrorFarmID})
	}

	areasJSON := area.DisplayAll()

	return c.JSON(http.StatusOK, AreasData{
		Data: areasJSON,
	})
}

// AreasCreate registers a new area inside a farm
func AreasCreate(c echo.Context) error {
	farmID := c.QueryParam("farm_id")

	if farmID == "" {
		return c.JSON(http.StatusOK, RequestError{ErrorFarmID})
	}

	return c.String(http.StatusOK, "not implemented yet")
}

// AreasHarvest harvests all planted crops from a particular area
func AreasHarvest(c echo.Context) error {
	farmID := c.QueryParam("farm_id")

	if farmID == "" {
		return c.JSON(http.StatusOK, RequestError{ErrorFarmID})
	}

	return c.String(http.StatusOK, "not implemented yet")
}

// AreasDispose disposes all planted crops from a particular area
func AreasDispose(c echo.Context) error {
	farmID := c.QueryParam("farm_id")

	if farmID == "" {
		return c.JSON(http.StatusOK, RequestError{ErrorFarmID})
	}

	return c.String(http.StatusOK, "not implemented yet")
}

// AreasDestroy destroys the particular area from a farm
func AreasDestroy(c echo.Context) error {
	farmID := c.QueryParam("farm_id")

	if farmID == "" {
		return c.JSON(http.StatusOK, RequestError{ErrorFarmID})
	}

	return c.String(http.StatusOK, "not implemented yet")
}
