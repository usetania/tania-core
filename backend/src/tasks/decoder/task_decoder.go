package decoder

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gofrs/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/usetania/tania-core/src/tasks/domain"
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
	case domain.TaskCreatedCode:
		e := domain.TaskCreated{}

		_, err := Decode(f, &mapped, &e)
		if err != nil {
			return err
		}

		if v, ok := mapped["domain_details"]; ok {
			domainCode, ok := mapped["domain"].(string)
			if !ok {
				return errors.New("error type assertion")
			}

			e.DomainDetails, err = makeDomainDetails(v, domainCode)
			if err != nil {
				return err
			}
		}

		w.Data = e

	case domain.TaskTitleChangedCode:
		e := domain.TaskTitleChanged{}

		_, err := Decode(f, &mapped, &e)
		if err != nil {
			return err
		}

		w.Data = e

	case domain.TaskDescriptionChangedCode:
		e := domain.TaskDescriptionChanged{}

		_, err := Decode(f, &mapped, &e)
		if err != nil {
			return err
		}

		w.Data = e

	case domain.TaskPriorityChangedCode:
		e := domain.TaskPriorityChanged{}

		_, err := Decode(f, &mapped, &e)
		if err != nil {
			return err
		}

		w.Data = e

	case domain.TaskDueDateChangedCode:
		e := domain.TaskDueDateChanged{}

		_, err := Decode(f, &mapped, &e)
		if err != nil {
			return err
		}

		w.Data = e

	case domain.TaskCategoryChangedCode:
		e := domain.TaskCategoryChanged{}

		_, err := Decode(f, &mapped, &e)
		if err != nil {
			return err
		}

		w.Data = e

	case domain.TaskDetailsChangedCode:
		e := domain.TaskDetailsChanged{}

		_, err := Decode(f, &mapped, &e)
		if err != nil {
			return err
		}

		if v, ok := mapped["domain_details"]; ok {
			domainCode, ok := mapped["domain"].(string)
			if !ok {
				return errors.New("error type assertion")
			}

			e.DomainDetails, err = makeDomainDetails(v, domainCode)
			if err != nil {
				return err
			}
		}

		w.Data = e

	case domain.TaskAssetIDChangedCode:
		e := domain.TaskAssetIDChanged{}

		_, err := Decode(f, &mapped, &e)
		if err != nil {
			return err
		}

		w.Data = e

	case domain.TaskCompletedCode:
		e := domain.TaskCompleted{}

		_, err := Decode(f, &mapped, &e)
		if err != nil {
			return err
		}

		w.Data = e

	case domain.TaskCancelledCode:
		e := domain.TaskCancelled{}

		_, err := Decode(f, &mapped, &e)
		if err != nil {
			return err
		}

		w.Data = e

	case domain.TaskDueCode:
		e := domain.TaskDue{}

		_, err := Decode(f, &mapped, &e)
		if err != nil {
			return err
		}

		w.Data = e
	}

	return nil
}

func makeDomainDetails(v interface{}, domainCode string) (domain.TaskDomain, error) {
	mapped := v.(map[string]interface{})

	var domainDetails domain.TaskDomain

	switch domainCode {
	case domain.TaskDomainAreaCode:
		taskDomainArea := domain.TaskDomainArea{}

		if v2, ok2 := mapped["material_id"]; ok2 {
			val, ok2 := v2.(string)
			if !ok2 {
				return domain.TaskDomainArea{}, nil
			}

			uid, err := uuid.FromString(val)
			if err != nil {
				return domain.TaskDomainArea{}, err
			}

			taskDomainArea.MaterialID = &uid
		}

		domainDetails = taskDomainArea

	case domain.TaskDomainCropCode:
		taskDomainCrop := domain.TaskDomainCrop{}

		if v2, ok2 := mapped["material_id"]; ok2 {
			val, ok2 := v2.(string)
			if !ok2 {
				return domain.TaskDomainCrop{}, nil
			}

			uid, err := uuid.FromString(val)
			if err != nil {
				return domain.TaskDomainCrop{}, err
			}

			taskDomainCrop.MaterialID = &uid
		}

		if v2, ok2 := mapped["area_id"]; ok2 {
			val, ok := v2.(string)
			if !ok {
				return domain.TaskDomainCrop{}, nil
			}

			uid, err := uuid.FromString(val)
			if err != nil {
				return domain.TaskDomainCrop{}, err
			}

			taskDomainCrop.AreaID = &uid
		}

		domainDetails = taskDomainCrop
	case domain.TaskDomainFinanceCode:
		domainDetails = domain.TaskDomainFinance{}
	case domain.TaskDomainGeneralCode:
		domainDetails = domain.TaskDomainGeneral{}
	case domain.TaskDomainInventoryCode:
		domainDetails = domain.TaskDomainInventory{}
	case domain.TaskDomainReservoirCode:
		taskDomainReservoir := domain.TaskDomainReservoir{}

		if v2, ok2 := mapped["material_id"]; ok2 {
			val, ok2 := v2.(string)
			if !ok2 {
				return domain.TaskDomainReservoir{}, nil
			}

			uid, err := uuid.FromString(val)
			if err != nil {
				return domain.TaskDomainReservoir{}, err
			}

			taskDomainReservoir.MaterialID = &uid
		}

		domainDetails = taskDomainReservoir
	}

	return domainDetails, nil
}
