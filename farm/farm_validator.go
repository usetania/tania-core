package farm

import (
	"regexp"

	"github.com/Tanibox/tania-server/location"
)

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
	rxLongitude := regexp.MustCompile(patternLongitude)

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
	_, err := FindFarmTypeByCode(code)

	if err != nil {
		return err
	}
	return nil
}

func validateCountryCode(code string) error {
	_, err := location.FindCountryByCountryCode(code)

	if err != nil {
		return FarmError{FarmErrorInvalidCountryCode}
	}
	return nil
}

func validateCityCode(countryCode string, cityCode string) error {
	_, err := location.FindCityByCityCode(countryCode, cityCode)

	if err != nil {
		return FarmError{FarmErrorInvalidCityCode}
	}
	return nil
}
