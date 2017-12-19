package location

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

type LocationStoreMock struct {
	mock.Mock
	countryMap map[string]Country
}

func (m *LocationStoreMock) initPackageStore() error {
	m.countryMap["ID"] = Country{ID: "ID", Name: "Indonesia"}
	m.countryMap["SG"] = Country{ID: "SG", Name: "Singapore"}
	m.countryMap["MY"] = Country{ID: "MY", Name: "Malaysia"}
	return nil
}

// func (m *LocationStoreMock) FindAllCountries(t *testing.T) ([]Country, error) {

// }

func TestFindAllCountries(t *testing.T) {
	obj := new(LocationStoreMock)
	obj.On("FindAllCountries")

	countries, err := obj.On("FindAllCountries").Return([]{})
}
