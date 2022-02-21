package decoder

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/usetania/tania-core/src/assets/domain"
)

type MaterialEventWrapper EventWrapper

func (w *MaterialEventWrapper) UnmarshalJSON(b []byte) error {
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
		MaterialTypeHook(),
	)

	switch wrapper.EventName {
	case "MaterialCreated":
		e := domain.MaterialCreated{}

		_, err := Decode(f, &mapped, &e)
		if err != nil {
			return err
		}

		w.EventData = e

	case "MaterialNameChanged":
		e := domain.MaterialNameChanged{}

		_, err := Decode(f, &mapped, &e)
		if err != nil {
			return err
		}

		w.EventData = e

	case "MaterialPriceChanged":
		e := domain.MaterialPriceChanged{}

		_, err := Decode(f, &mapped, &e)
		if err != nil {
			return err
		}

		w.EventData = e

	case "MaterialQuantityChanged":
		e := domain.MaterialQuantityChanged{}

		_, err := Decode(f, &mapped, &e)
		if err != nil {
			return err
		}

		w.EventData = e

	case "MaterialTypeChanged":
		e := domain.MaterialTypeChanged{}

		_, err := Decode(f, &mapped, &e)
		if err != nil {
			return err
		}

		w.EventData = e

	case "MaterialExpirationDateChanged":
		e := domain.MaterialExpirationDateChanged{}

		_, err := Decode(f, &mapped, &e)
		if err != nil {
			return err
		}

		w.EventData = e

	case "MaterialNotesChanged":
		e := domain.MaterialNotesChanged{}

		_, err := Decode(f, &mapped, &e)
		if err != nil {
			return err
		}

		w.EventData = e

	case "MaterialProducedByChanged":
		e := domain.MaterialProducedByChanged{}

		_, err := Decode(f, &mapped, &e)
		if err != nil {
			return err
		}

		w.EventData = e
	}

	return nil
}
