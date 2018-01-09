package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Tanibox/tania-server/farm/entity"
	"github.com/labstack/echo"
)

type SimpleFarm entity.Farm
type SimpleArea entity.Area
type DetailArea entity.Area
type DetailReservoir struct {
	entity.Reservoir
	InstalledToAreas []SimpleArea
}

type AreaSquareMeter struct {
	entity.SquareMeter
}
type AreaHectare struct {
	entity.Hectare
}
type ReservoirBucket struct {
	entity.Bucket
}
type ReservoirTap struct {
	entity.Tap
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

func MapToSimpleArea(areas []entity.Area) []SimpleArea {
	installedAreaList := make([]SimpleArea, len(areas))

	for i, area := range areas {
		installedAreaList[i] = SimpleArea(area)
	}

	return installedAreaList
}

func MapToReservoir(s *FarmServer, reservoirs []entity.Reservoir) ([]DetailReservoir, error) {
	reservoirList := make([]DetailReservoir, len(reservoirs))

	for i, reservoir := range reservoirs {
		reservoirList[i] = DetailReservoir{Reservoir: reservoir}

		switch v := reservoir.WaterSource.(type) {
		case entity.Bucket:
			reservoirList[i].WaterSource = ReservoirBucket{Bucket: v}
		case entity.Tap:
			reservoirList[i].WaterSource = ReservoirTap{Tap: v}
		}

		queryResult := <-s.AreaQuery.FindAreasByReservoirID(reservoir.UID.String())
		if queryResult.Error != nil {
			return nil, echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
		}

		areas, ok := queryResult.Result.([]entity.Area)
		if !ok {
			return nil, echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
		}

		reservoirList[i].InstalledToAreas = MapToSimpleArea(areas)
	}

	return reservoirList, nil
}

func MapToDetailReservoir(s *FarmServer, reservoir entity.Reservoir) (DetailReservoir, error) {
	detailReservoir := DetailReservoir{Reservoir: reservoir}

	switch v := detailReservoir.WaterSource.(type) {
	case entity.Bucket:
		detailReservoir.WaterSource = ReservoirBucket{Bucket: v}
	case entity.Tap:
		detailReservoir.WaterSource = ReservoirTap{Tap: v}
	}

	queryResult := <-s.AreaQuery.FindAreasByReservoirID(detailReservoir.UID.String())
	if queryResult.Error != nil {
		return DetailReservoir{}, echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	areas, ok := queryResult.Result.([]entity.Area)
	if !ok {
		return DetailReservoir{}, echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	detailReservoir.InstalledToAreas = MapToSimpleArea(areas)

	return detailReservoir, nil
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

func (sa SimpleArea) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		UID  string `json:"uid"`
		Name string `json:"name"`
		Type string `json:"type"`
	}{
		UID:  sa.UID.String(),
		Name: sa.Name,
		Type: sa.Type,
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

func (dr DetailReservoir) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		UID              string             `json:"uid"`
		Name             string             `json:"name"`
		PH               float32            `json:"ph"`
		EC               float32            `json:"ec"`
		Temperature      float32            `json:"temperature"`
		WaterSource      entity.WaterSource `json:"water_source"`
		CreatedDate      time.Time          `json:"created_date"`
		InstalledToAreas []SimpleArea       `json:"installed_to_areas"`
	}{
		UID:              dr.UID.String(),
		Name:             dr.Name,
		PH:               dr.PH,
		EC:               dr.EC,
		Temperature:      dr.Temperature,
		WaterSource:      dr.WaterSource,
		CreatedDate:      dr.CreatedDate,
		InstalledToAreas: dr.InstalledToAreas,
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

func (rb ReservoirBucket) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type     string  `json:"type"`
		Capacity float32 `json:"capacity"`
		Volume   float32 `json:"volume"`
	}{
		Type:     rb.Type(),
		Capacity: rb.Capacity,
		Volume:   rb.Volume,
	})
}

func (rt ReservoirTap) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string `json:"type"`
	}{
		Type: rt.Type(),
	})
}
