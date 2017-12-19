// Package routing provides the list of functions for each HTTP routing
package routing

import (
	"net/http"

	"github.com/labstack/echo"
)

func LocationsRouter(g *echo.Group) {
	// Bootstap countries data
	// location.CountryRepoInMemory{ countryMap: [] }

	g.GET("/countries", LocationsGetCountries)
	g.GET("/cities", LocationsGetCities)
}

// LocationsGetCountries displays all available location in Tania
func LocationsGetCountries(c echo.Context) error {
	return c.JSON(http.StatusOK, "")
}

// LocationsGetCities displays all available location in Tania
func LocationsGetCities(c echo.Context) error {
	// country := c.QueryParam("country_id")
	return c.JSON(http.StatusOK, "")
}
