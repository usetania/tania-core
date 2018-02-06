package domain

const (
	TaskDomainArea      = "AREA"
	TaskDomainCrop      = "CROP"
	TaskDomainFinance   = "FINANCE"
	TaskDomainGeneral   = "GENERAL"
	TaskDomainInventory = "INVENTORY"
	TaskDomainReservoir = "RESERVOIR"
)

type TaskDomain struct {
	Code string
	Name string
}

func FindAllTaskDomains() []TaskDomain {
	return []TaskDomain{
		TaskDomain{Code: TaskDomainArea, Name: "Area"},
		TaskDomain{Code: TaskDomainCrop, Name: "Crop"},
		TaskDomain{Code: TaskDomainFinance, Name: "Finance"},
		TaskDomain{Code: TaskDomainGeneral, Name: "General"},
		TaskDomain{Code: TaskDomainInventory, Name: "Inventory"},
		TaskDomain{Code: TaskDomainReservoir, Name: "Reservoir"},
	}
}

func FindTaskDomainByCode(code string) (TaskDomain, error) {
	items := FindAllTaskDomains()

	for _, item := range items {
		if item.Code == code {
			return item, nil
		}
	}

	return TaskDomain{}, TaskError{TaskErrorInvalidDomainCode}
}
