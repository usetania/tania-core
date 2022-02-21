package decoder

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/usetania/tania-core/src/assets/domain"
)

type AreaEventWrapper EventWrapper

func (w *AreaEventWrapper) UnmarshalJSON(b []byte) error {
	wrapper := EventWrapper{}

	err := json.Unmarshal(b, &wrapper)
	if err != nil {
		return err
	}

	mapped, ok := wrapper.EventData.(map[string]interface{})
	if !ok {
		return errors.New("error type assertion")
	}

	f := mapstructure.ComposeDecodeHookFunc(
		UIDHook(),
		TimeHook(time.RFC3339),
	)

	var e interface{}

	switch wrapper.EventName {
	case "AreaCreated":
		e = domain.AreaCreated{}
	case "AreaNameChanged":
		e = domain.AreaNameChanged{}
	case "AreaSizeChanged":
		e = domain.AreaSizeChanged{}
	case "AreaTypeChanged":
		e = domain.AreaTypeChanged{}
	case "AreaLocationChanged":
		e = domain.AreaLocationChanged{}
	case "AreaReservoirChanged":
		e = domain.AreaReservoirChanged{}
	case "AreaPhotoAdded":
		e = domain.AreaPhotoAdded{}
	case "AreaNoteAdded":
		e = domain.AreaNoteAdded{}
	case "AreaNoteRemoved":
		e = domain.AreaNoteRemoved{}
	}

	_, err = Decode(f, &mapped, &e)
	if err != nil {
		return err
	}

	w.EventData = e

	return nil
}
