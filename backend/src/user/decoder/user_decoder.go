package decoder

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/usetania/tania-core/src/user/domain"
)

type UserEventWrapper EventWrapper

func (w *UserEventWrapper) UnmarshalJSON(b []byte) error {
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
		PasswordHook(),
	)

	switch wrapper.EventName {
	case "UserCreated":
		e := domain.UserCreated{}

		_, err := Decode(f, &mapped, &e)
		if err != nil {
			return err
		}

		w.EventData = e

	case "PasswordChanged":
		e := domain.PasswordChanged{}

		_, err := Decode(f, &mapped, &e)
		if err != nil {
			return err
		}

		w.EventData = e
	}

	return nil
}
