// Package routing provides the list of functions for each HTTP routing
package routing

import (
	"../area"
	"github.com/labstack/echo"
	"net/http"
)

// AreasIndex displays all available areas in a farm
func AreasIndex(c echo.Context) error {
	farmID := c.QueryParam("farm_id")

	if farmID == "" {
		return c.JSON(http.StatusOK, RequestError{ErrorFarmID})
	}

	areasJSON := area.DisplayAll()

	return c.JSON(http.StatusOK, areasJSON)
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
