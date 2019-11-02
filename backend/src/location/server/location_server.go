package location

import (
	"net/http"
	"sort"

	"github.com/labstack/echo/v4"
	"github.com/pariz/gountries"
)

// Country is country
type Country struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type LocationServer struct{}

func NewLocationServer() (*LocationServer, error) {
	return &LocationServer{}, nil
}

func (s *LocationServer) Mount(g *echo.Group) {
	g.GET("/countries", LocationsGetCountries)
}

// LocationsGetCountries displays all available location in Tania
func LocationsGetCountries(c echo.Context) error {
	var countries []Country
	query := gountries.New()
	items := query.FindAllCountries()

	for _, item := range items {
		countries = append(countries, Country{
			ID:   item.Codes.Alpha2,
			Name: item.Name.Common,
		})
	}

	sort.Slice(countries, func(i, j int) bool {
		return countries[i].Name < countries[j].Name
	})

	return c.JSON(http.StatusOK, countries)
}
