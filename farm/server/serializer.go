package server

import (
	"encoding/json"

	"github.com/Tanibox/tania-server/farm/entity"
)

type SimpleFarm entity.Farm

func MapToSimpleFarm(farms []entity.Farm) []SimpleFarm {
	farmList := make([]SimpleFarm, len(farms))

	for i, farm := range farms {
		farmList[i] = SimpleFarm(farm)
	}

	return farmList
}

func (sf SimpleFarm) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		UID  string `json:"uid"`
		Name string `json:"name"`
		Type string `json:"type"`
	}{
		UID:  sf.UID,
		Name: sf.Name,
		Type: sf.Type,
	})
}
