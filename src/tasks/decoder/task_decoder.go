package decoder

import (
	"encoding/json"
	"time"

	"github.com/Tanibox/tania-server/src/tasks/domain"
	"github.com/mitchellh/mapstructure"
)

type TaskEventWrapper InterfaceWrapper

func (w *TaskEventWrapper) UnmarshalJSON(b []byte) error {
	wrapper := InterfaceWrapper{}

	err := json.Unmarshal(b, &wrapper)
	if err != nil {
		return err
	}

	mapped := wrapper.Data.(map[string]interface{})

	f := mapstructure.ComposeDecodeHookFunc(
		UIDHook(),
		TimeHook(time.RFC3339),
		TaskDomainDetailHook(),
	)

	switch wrapper.Name {
	case "TaskCreated":
		e := domain.TaskCreated{}

		_, err := Decode(f, &mapped, &e)
		if err != nil {
			return err
		}

		w.Data = e

	case "TaskModified":
		e := domain.TaskModified{}

		_, err := Decode(f, &mapped, &e)
		if err != nil {
			return err
		}

		w.Data = e

	case "TaskCompleted":
		e := domain.TaskCompleted{}

		_, err := Decode(f, &mapped, &e)
		if err != nil {
			return err
		}

		w.Data = e

	case "TaskCancelled":
		e := domain.TaskCancelled{}

		_, err := Decode(f, &mapped, &e)
		if err != nil {
			return err
		}

		w.Data = e

	case "TaskDue":
		e := domain.TaskDue{}

		_, err := Decode(f, &mapped, &e)
		if err != nil {
			return err
		}

		w.Data = e

	}

	return nil
}
