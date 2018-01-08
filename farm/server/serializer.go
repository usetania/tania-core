package server

import (
	"encoding/json"

	"github.com/Tanibox/tania-server/farm/entity"
)

type SimpleFarm entity.Farm

type DetailArea entity.Area

type AreaSquareMeter struct {
	entity.SquareMeter
}
type AreaHectare struct {
	entity.Hectare
}

func MapToSimpleFarm(farms []entity.Farm) []SimpleFarm {
	farmList := make([]SimpleFarm, len(farms))

	for i, farm := range farms {
		farmList[i] = SimpleFarm(farm)
	}

	return farmList
}

func MapToArea(areas []entity.Area) []entity.Area {
	areaList := make([]entity.Area, len(areas))

	for i, area := range areas {
		areaList[i] = area

		switch v := area.Size.(type) {
		case entity.SquareMeter:
			areaList[i].Size = AreaSquareMeter{SquareMeter: v}
		case entity.Hectare:
			areaList[i].Size = AreaHectare{Hectare: v}
		}
	}

	return areaList
}

func MapToDetailArea(area entity.Area) DetailArea {
	switch v := area.Size.(type) {
	case entity.SquareMeter:
		area.Size = AreaSquareMeter{SquareMeter: v}
	case entity.Hectare:
		area.Size = AreaHectare{Hectare: v}
	}

	return DetailArea(area)
}

func (sf SimpleFarm) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		UID  string `json:"uid"`
		Name string `json:"name"`
		Type string `json:"type"`
	}{
		UID:  sf.UID.String(),
		Name: sf.Name,
		Type: sf.Type,
	})
}

func (da DetailArea) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		UID       string           `json:"uid"`
		Name      string           `json:"name"`
		Size      entity.AreaUnit  `json:"size"`
		Type      string           `json:"type"`
		Location  string           `json:"location"`
		Photo     entity.AreaPhoto `json:"photo"`
		Reservoir entity.Reservoir `json:"reservoir"`
	}{
		UID:       da.UID.String(),
		Name:      da.Name,
		Size:      da.Size,
		Type:      da.Type,
		Location:  da.Location,
		Photo:     da.Photo,
		Reservoir: da.Reservoir,
	})
}

func (sm AreaSquareMeter) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Value  float32 `json:"value"`
		Symbol string  `json:"symbol"`
	}{
		Value:  sm.Value,
		Symbol: sm.Symbol(),
	})
}

func (h AreaHectare) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Value  float32 `json:"value"`
		Symbol string  `json:"symbol"`
	}{
		Value:  h.Value,
		Symbol: h.Symbol(),
	})
}
