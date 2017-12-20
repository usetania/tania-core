// Package farm provides the operation that farm holder can do
// to their farm
package farm

import "regexp"

type Farm struct {
	Name        string
	Description string
	Latitude    string
	Longitude   string
	Type        string
	CountryCode string
	CityCode    string
}

// DisplayAll dispalys all existing farms
func DisplayAll() []Farm {

}

// CreateNew registers a new farm to Tania
func CreateNew(name string, description string, latitude string, longitude string, farmType string, countryCode string, cityCode string) (Farm, error) {

	err := validateName(name)
	if err != nil {
		return Farm{}, err
	}

	err = validateGeoLocation(latitude, longitude)
	if err != nil {
		return Farm{}, err
	}

	err = validateFarmType(farmType)
	if err != nil {
		return Farm{}, err
	}

	// validate contry code

	// validate city code

	return Farm{
		Name:        name,
		Description: description,
		Latitude:    latitude,
		Longitude:   longitude,
		Type:        farmType,
		CountryCode: countryCode,
		CityCode:    cityCode,
	}, nil
}

func validateName(name string) error {
	if name == "" {
		return FarmError{FarmErrorEmptyNameCode}
	}
	if len(name) < 5 {
		return FarmError{FarmErrorNotEnoughCharacterCode}
	}

	return nil
}

func validateGeoLocation(latitude string, longitude string) error {
	patternLatitude := "^[-+]?([1-8]?\\d(\\.\\d+)?|90(\\.0+)?)$"
	patternLongitude := "^[-+]?(180(\\.0+)?|((1[0-7]\\d)|([1-9]?\\d))(\\.\\d+)?)$"

	rxLatitude := regexp.MustCompile(patternLatitude)
	rxLongitude := regexp.MustCompile(patternLatitude)

	ok := rxLatitude.MatchString(latitude)
	if !ok {
		return FarmError{FarmErrorInvalidLatitudeValueCode}
	}

	ok = rxLongitude.MatchString(longitude)
	if !ok {
		return FarmError{FarmErrorInvalidLongitudeValueCode}
	}

	return nil
}

func validateFarmType(code string) error {
	farm, err := FindFarmTypeByCode(code)

	if err != nil {
		return err
	}
	return nil
}

// ShowInformation shows information of a farm
func ShowInformation(uid string) *Farm {

}

// UpdateInformation updates the existing farm information in Tania
func UpdateInformation() {
}

// DestroyFarm destroys the farm and its properties. This is dangerous.
func Destroy() {

}
