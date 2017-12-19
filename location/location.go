package location

import "github.com/pariz/gountries"

// Country is a representation of a Country
type Country struct {
	ID   string
	Name string
}

// City is a representation of a city
type City struct {
	ID          string
	Name        string
	CountryCode string
}

// CountryRepo is a interface how country can be accessed
type CountryRepo interface {
	FindAllCountries() []Country
	FindCountryByCode(code string) (*Country, error)
}

// The LocationStore struct
type LocationStore struct {
	countryMap map[string]Country
	cityMap    map[string]City
}

// initPackageStore
func (cm *LocationStore) initPackageStore() error {
	query := gountries.New()
	items := query.FindAllCountries()

	for _, country := range items {
		cm.countryMap[country.Codes.Alpha2] = Country{
			ID:   country.Codes.Alpha2,
			Name: country.Name.Common,
		}

		for _, city := range country.SubDivisions() {
			cm.cityMap[city.Code] = City{
				ID:          city.Code,
				Name:        city.Name,
				CountryCode: city.CountryAlpha2,
			}
		}
	}

	return nil
}

// FindAllCountries find all available countries
func (cm *LocationStore) FindAllCountries() []Country {
	var countries []Country

	for _, item := range cm.countryMap {
		countries = append(countries, Country{
			ID:   item.ID,
			Name: item.Name,
		})
	}

	return countries
}

// FindCountryByCode find country by code
func (cm *LocationStore) FindCountryByCode(code string) (Country, error) {
	country, ok := cm.countryMap[code]

	if !ok {
		return Country{}, LocationError{LocationErrorInvalidCountryCode}
	}

	return Country{
		ID:   country.ID,
		Name: country.Name,
	}, nil
}

// FindAllCitiesByCountryCode find cities by country code
func (cm *LocationStore) FindAllCitiesByCountryCode(code string) ([]City, error) {
	var cities []City
	country := cm.countryMap[code]

	for _, city := range cm.cityMap {
		if city.CountryCode == code {
			cities = append(cities, City{
				ID:          city.ID,
				Name:        city.Name,
				CountryCode: country.ID,
			})
		}
	}

	return cities, nil
}

// FindCityByCode find city by city code
func (cm *LocationStore) FindCityByCode(code string) (City, error) {
	city, ok := cm.cityMap[code]

	if !ok {
		return City{}, LocationError{LocationErrorInvalidCityCode}
	}

	return City{
		ID:          city.ID,
		Name:        city.Name,
		CountryCode: city.CountryCode,
	}, nil
}
