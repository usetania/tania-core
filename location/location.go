package location

import "github.com/pariz/gountries"

// Country is country
type Country struct {
	ID   string
	Name string
}

// City is city
type City struct {
	ID          string
	Name        string
	CountryCode string
}

// FindAllCountries Find All Country
func FindAllCountries() []Country {
	var countries []Country
	query := gountries.New()
	items := query.FindAllCountries()

	for _, item := range items {
		countries = append(countries, Country{
			ID:   item.Codes.Alpha2,
			Name: item.Name.Common,
		})
	}
	return countries
}

// FindCountryByCountryCode find country by code
func FindCountryByCountryCode(code string) (Country, error) {
	query := gountries.New()
	country, err := query.FindCountryByAlpha(code)

	if err != nil {
		return Country{}, LocationError{LocationErrorInvalidCountryCode}
	}

	return Country{
		ID:   country.Codes.Alpha2,
		Name: country.Name.Common,
	}, nil
}

// FindAllCitiesByCountryCode find all cities
func FindAllCitiesByCountryCode(code string) ([]City, error) {
	var cities []City
	query := gountries.New()
	country, err := query.FindCountryByAlpha(code)

	if err != nil {
		return cities, LocationError{LocationErrorInvalidCountryCode}
	}

	items := country.SubDivisions()

	for _, item := range items {
		cities = append(cities, City{
			ID:          item.Code,
			Name:        item.Name,
			CountryCode: item.CountryAlpha2,
		})
	}
	return cities, nil
}

// FindCityByCityCode find city by city code
func FindCityByCityCode(countryCode string, code string) (City, error) {
	items, err := FindAllCitiesByCountryCode(countryCode)

	if err != nil {
		return City{}, LocationError{LocationErrorInvalidCountryCode}
	}

	for _, item := range items {
		if item.ID == code {
			return City{
				ID:          item.ID,
				Name:        item.Name,
				CountryCode: item.CountryCode,
			}, nil
		}
	}

	return City{}, LocationError{LocationErrorInvalidCityCode}
}
