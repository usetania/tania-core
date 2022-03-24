package domain

import (
	"regexp"

	"github.com/usetania/tania-core/src/helper/validationhelper"
)

func validateFarmName(name string) error {
	if name == "" {
		return FarmError{FarmErrorNameEmptyCode}
	}

	if !validationhelper.IsAlphanumSpaceHyphenUnderscore(name) {
		return FarmError{FarmErrorNameAlphanumericOnlyCode}
	}

	if len(name) < 5 {
		return FarmError{FarmErrorNameNotEnoughCharacterCode}
	}

	if len(name) > 100 {
		return FarmError{FarmErrorNameExceedMaximunCharacterCode}
	}

	return nil
}

func validateGeoLocation(latitude, longitude string) error {
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

func validateCountry(country string) error {
	if country == "" {
		return FarmError{FarmErrorInvalidCountry}
	}

	return nil
}

func validateCity(city string) error {
	if city == "" {
		return FarmError{FarmErrorInvalidCity}
	}

	return nil
}
