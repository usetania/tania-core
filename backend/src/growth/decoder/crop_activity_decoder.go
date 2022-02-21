package decoder

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/usetania/tania-core/src/growth/storage"
)

type CropActivityTypeWrapper InterfaceWrapper

func (w *CropActivityTypeWrapper) UnmarshalJSON(b []byte) error {
	wrapper := InterfaceWrapper{}

	err := json.Unmarshal(b, &wrapper)
	if err != nil {
		return err
	}

	mapped, ok := wrapper.Data.(map[string]interface{})
	if !ok {
		return errors.New("error type assertion")
	}

	f := mapstructure.ComposeDecodeHookFunc(
		UIDHook(),
		TimeHook(time.RFC3339),
		CropContainerHook(),
	)

	switch wrapper.Name {
	case storage.SeedActivityCode:
		a := storage.SeedActivity{}

		_, err := Decode(f, &mapped, &a)
		if err != nil {
			return err
		}

		w.Data = a

	case storage.MoveActivityCode:
		a := storage.MoveActivity{}

		_, err := Decode(f, &mapped, &a)
		if err != nil {
			return err
		}

		w.Data = a

	case storage.HarvestActivityCode:
		a := storage.HarvestActivity{}

		_, err := Decode(f, &mapped, &a)
		if err != nil {
			return err
		}

		w.Data = a

	case storage.DumpActivityCode:
		a := storage.DumpActivity{}

		_, err := Decode(f, &mapped, &a)
		if err != nil {
			return err
		}

		w.Data = a

	case storage.WaterActivityCode:
		a := storage.WaterActivity{}

		_, err := Decode(f, &mapped, &a)
		if err != nil {
			return err
		}

		w.Data = a

	case storage.PhotoActivityCode:
		a := storage.PhotoActivity{}

		_, err := Decode(f, &mapped, &a)
		if err != nil {
			return err
		}

		w.Data = a

	case storage.TaskCropActivityCode:
		a := storage.TaskCropActivity{}

		_, err := Decode(f, &mapped, &a)
		if err != nil {
			return err
		}

		w.Data = a

	case storage.TaskNutrientActivityCode:
		a := storage.TaskNutrientActivity{}

		_, err := Decode(f, &mapped, &a)
		if err != nil {
			return err
		}

		w.Data = a

	case storage.TaskPestControlActivityCode:
		a := storage.TaskPestControlActivity{}

		_, err := Decode(f, &mapped, &a)
		if err != nil {
			return err
		}

		w.Data = a

	case storage.TaskSafetyActivityCode:
		a := storage.TaskSafetyActivity{}

		_, err := Decode(f, &mapped, &a)
		if err != nil {
			return err
		}

		w.Data = a

	case storage.TaskSanitationActivityCode:
		a := storage.TaskSanitationActivity{}

		_, err := Decode(f, &mapped, &a)
		if err != nil {
			return err
		}

		w.Data = a
	}

	return nil
}
