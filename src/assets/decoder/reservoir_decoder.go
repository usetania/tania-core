package decoder

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/Tanibox/tania-core/src/assets/domain"
	"github.com/mitchellh/mapstructure"
)

type ReservoirEventWrapper EventWrapper

func (w *ReservoirEventWrapper) UnmarshalJSON(b []byte) error {
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
		WaterSourceHook(),
	)

	var e interface{}

	switch wrapper.EventName {
	case "ReservoirCreated":
		e = domain.ReservoirCreated{}
	case "ReservoirWaterSourceChanged":
		e = domain.ReservoirWaterSourceChanged{}
	case "ReservoirNameChanged":
		e = domain.ReservoirNameChanged{}
	case "ReservoirNoteAdded":
		e = domain.ReservoirNoteAdded{}
	case "ReservoirNoteRemoved":
		e = domain.ReservoirNoteRemoved{}
	}

	_, err = Decode(f, &mapped, &e)
	if err != nil {
		return err
	}

	w.EventData = e

	return nil
}
