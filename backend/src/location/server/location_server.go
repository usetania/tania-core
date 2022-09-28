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

type Server struct{}

func NewServer() (*Server, error) {
	return &Server{}, nil
}

func (*Server) Mount(g *echo.Group) {
	g.GET("/countries", LocationsGetCountries)
}

// LocationsGetCountries displays all available location in Tania
func LocationsGetCountries(c echo.Context) error {
	query := gountries.New()
	items := query.FindAllCountries()
	countries := make([]Country, 0, len(items))

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
