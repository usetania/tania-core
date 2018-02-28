package sqlite

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

func makeUUID(v interface{}) (uuid.UUID, error) {
	val := v.(string)
	uid, err := uuid.FromString(val)
	if err != nil {
		return uuid.UUID{}, err
	}

	return uid, nil
}

func makeTime(v interface{}) (time.Time, error) {
	val := v.(string)

	date, err := time.Parse(time.RFC3339, val)
	if err != nil {
		return time.Time{}, err
	}

	return date, nil
}
