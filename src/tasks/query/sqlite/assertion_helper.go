package sqlite

import (
	"time"

	"github.com/Tanibox/tania-core/src/tasks/domain"
	uuid "github.com/satori/go.uuid"
)

func makeUUID(v interface{}) (uuid.UUID, error) {
	val := v.(string)
	uid, err := uuid.FromString(val)
	if err != nil {
		return uuid.UUID{}, err
	}

	return uid, nil
}

func makeTime(v interface{}) (time.Time, error) {
	val := v.(string)
	date, err := time.Parse(time.RFC3339, val)
	if err != nil {
		return time.Time{}, err
	}

	return date, nil
}

func makeTimePointer(v interface{}) (*time.Time, error) {
	val, ok := v.(string)
	var date *time.Time

	if ok {
		d, err := time.Parse(time.RFC3339, val)
		if err != nil {
			return nil, err
		}

		date = &d
	}

	return date, nil
}

func makeDomainDetails(v interface{}, domainCode string) (domain.TaskDomain, error) {
	mapped := v.(map[string]interface{})

	var domainDetails domain.TaskDomain
	switch domainCode {
	case domain.TaskDomainAreaCode:
		domainDetails = domain.TaskDomainArea{}
	case domain.TaskDomainCropCode:
		if v, ok := mapped["InventoryUID"]; ok {
			val, ok2 := v.(string)
			if ok2 {
				uid, err := makeUUID(val)
				if err != nil {
					return nil, err
				}

				domainDetails = domain.TaskDomainCrop{MaterialID: &uid}
			} else {
				domainDetails = domain.TaskDomainCrop{}
			}
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
