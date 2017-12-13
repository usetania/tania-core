// Package routing provides the list of functions for each HTTP routing
package routing

import (
	"net/http"

	"github.com/Tanibox/tania-server/farm"
	"github.com/labstack/echo"
)

func FarmsRouter(g *echo.Group) {
	g.GET("/", FarmsIndex)
	g.POST("/", FarmsCreate)
	g.GET("/:uid", FarmsShow)
	g.PUT("/:uid", FarmsUpdate)
	g.DELETE("/:uid", FarmsDestroy)
}

// FarmsIndex displays all available farms in Tania
func FarmsIndex(c echo.Context) error {
	farmJSON := farm.DisplayAll()

	return c.JSON(http.StatusOK, farmJSON)
}

// FarmsCreate registers a new farm into Tania
func FarmsCreate(c echo.Context) error {
	return c.String(http.StatusOK, "not implemented yet")
}

// FarmsShow shows the detail of a particular farm.
func FarmsShow(c echo.Context) error {
	uid := c.Param("uid")
	farmJSON := farm.ShowInformation(uid)

	return c.JSON(http.StatusOK, farmJSON)
}

// FarmsUpdate updates the information of a particular farm.
func FarmsUpdate(c echo.Context) error {
	//uid := c.Param("uid")
	return c.String(http.StatusOK, "not implemented yet")
}

// FarmsDestroy deletes the farm from Tania
func FarmsDestroy(c echo.Context) error {
	//uid := c.Param("uid")
	return c.String(http.StatusOK, "not implemented yet")
}
