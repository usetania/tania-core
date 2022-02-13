package decoder

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/usetania/tania-core/src/assets/domain"
)

type FarmEventWrapper EventWrapper

func (w *FarmEventWrapper) UnmarshalJSON(b []byte) error {
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

	switch wrapper.EventName {
	case "FarmCreated":
		e := domain.FarmCreated{}

		_, err := Decode(f, &mapped, &e)
		if err != nil {
			return err
		}

		w.EventData = e

	case "FarmNameChanged":
		e := domain.FarmNameChanged{}

		_, err := Decode(f, &mapped, &e)
		if err != nil {
			return err
		}

		w.EventData = e

	case "FarmTypeChanged":
		e := domain.FarmTypeChanged{}

		_, err := Decode(f, &mapped, &e)
		if err != nil {
			return err
		}

		w.EventData = e

	case "FarmGeolocationChanged":
		e := domain.FarmGeolocationChanged{}

		_, err := Decode(f, &mapped, &e)
		if err != nil {
			return err
		}

		w.EventData = e

	case "FarmRegionChanged":
		e := domain.FarmRegionChanged{}

		_, err := Decode(f, &mapped, &e)
		if err != nil {
			return err
		}

		w.EventData = e
	}

	return nil
}
