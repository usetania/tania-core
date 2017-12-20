// Package routing provides the list of functions for each HTTP routing
package routing

import (
	"net/http"

	"github.com/Tanibox/tania-server/location"
	"github.com/labstack/echo"
)

func LocationsRouter(g *echo.Group) {
	g.GET("/countries", LocationsGetCountries)
	g.GET("/cities", LocationsGetCities)
}

// LocationsGetCountries displays all available location in Tania
func LocationsGetCountries(c echo.Context) error {
	countries := location.FindAllCountries()
	return c.JSON(http.StatusOK, countries)
}

// LocationsGetCities displays all available location in Tania
func LocationsGetCities(c echo.Context) error {
	country := c.QueryParam("country_id")
	cities, _ := location.FindAllCitiesByCountryCode(country)
	return c.JSON(http.StatusOK, cities)
}
