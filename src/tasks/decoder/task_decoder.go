package decoder

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/Tanibox/tania-server/src/tasks/domain"
	"github.com/mitchellh/mapstructure"
	uuid "github.com/satori/go.uuid"
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

		if v, ok := mapped["domain_details"]; ok {
			domainCode, ok := mapped["domain"].(string)
			if !ok {
				return errors.New("Error type assertion")
			}

			e.DomainDetails, err = makeDomainDetails(v, domainCode)
			if err != nil {
				return err
			}
		}

		w.Data = e

	case "TaskModified":
		e := domain.TaskModified{}

		_, err := Decode(f, &mapped, &e)
		if err != nil {
			return err
		}

		if v, ok := mapped["domain_details"]; ok {
			domainCode, ok := mapped["domain"].(string)
			if !ok {
				return errors.New("Error type assertion")
			}

			e.DomainDetails, err = makeDomainDetails(v, domainCode)
			if err != nil {
				return err
			}
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

func makeDomainDetails(v interface{}, domainCode string) (domain.TaskDomain, error) {
	mapped := v.(map[string]interface{})

	var domainDetails domain.TaskDomain
	switch domainCode {
	case domain.TaskDomainAreaCode:
		domainDetails = domain.TaskDomainArea{}
	case domain.TaskDomainCropCode:
		if v2, ok2 := mapped["InventoryUID"]; ok2 {
			val, ok := v2.(string)
			if !ok {
				return domain.TaskDomainCrop{}, nil
			}

			uid, err := uuid.FromString(val)
			if err != nil {
				return domain.TaskDomainCrop{}, err
			}

			domainDetails = domain.TaskDomainCrop{InventoryUID: &uid}
		}
	case domain.TaskDomainFinanceCode:
		domainDetails = domain.TaskDomainFinance{}
	case domain.TaskDomainGeneralCode:
		domainDetails = domain.TaskDomainGeneral{}
	case domain.TaskDomainInventoryCode:
		domainDetails = domain.TaskDomainInventory{}
	case domain.TaskDomainReservoirCode:
		domainDetails = domain.TaskDomainReservoir{}
	}

	return domainDetails, nil
}
