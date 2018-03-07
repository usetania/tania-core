package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	assetsdomain "github.com/Tanibox/tania-server/src/assets/domain"
	assetsstorage "github.com/Tanibox/tania-server/src/assets/storage"
	"github.com/Tanibox/tania-server/src/tasks/query"
	uuid "github.com/satori/go.uuid"
)

type MaterialQuerySqlite struct {
	DB *sql.DB
}

func NewMaterialQuerySqlite(db *sql.DB) query.MaterialQuery {
	return MaterialQuerySqlite{DB: db}
}

func (s MaterialQuerySqlite) FindMaterialByID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		rowsData := struct {
			UID      string
			Name     string
			Type     string
			TypeData string
		}{}
		material := query.TaskMaterialQueryResult{}

		err := s.DB.QueryRow(`SELECT UID, NAME
			FROM MATERIAL_READ WHERE UID = ?`, uid).Scan(&rowsData.UID, &rowsData.Name)

		fmt.Println("UID", rowsData.UID)
		materialUID, err := uuid.FromString(rowsData.UID)
		fmt.Println("MUID", materialUID)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		material.UID = materialUID
		material.Name = rowsData.Name

		var materialType assetsstorage.MaterialType
		switch rowsData.Type {
		case assetsdomain.MaterialTypePlantCode:
			materialType, err = assetsdomain.CreateMaterialTypePlant(rowsData.TypeData)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}
			material.Type = materialType
		case assetsdomain.MaterialTypeSeedCode:
			materialType, err = assetsdomain.CreateMaterialTypeSeed(rowsData.TypeData)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}
			material.Type = materialType
		case assetsdomain.MaterialTypeGrowingMediumCode:
			material.Type = assetsdomain.MaterialTypeGrowingMedium{}
		case assetsdomain.MaterialTypeAgrochemicalCode:
			materialType, err = assetsdomain.CreateMaterialTypeAgrochemical(rowsData.TypeData)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}
			material.Type = materialType
		case assetsdomain.MaterialTypeLabelAndCropSupportCode:
			material.Type = assetsdomain.MaterialTypeLabelAndCropSupport{}
		case assetsdomain.MaterialTypeSeedingContainerCode:
			materialType, err = assetsdomain.CreateMaterialTypeSeedingContainer(rowsData.TypeData)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}
			material.Type = materialType
		case assetsdomain.MaterialTypePostHarvestSupplyCode:
			material.Type = assetsdomain.MaterialTypePostHarvestSupply{}
		case assetsdomain.MaterialTypeOtherCode:
			material.Type = assetsdomain.MaterialTypeOther{}
		default:
			result <- query.QueryResult{Error: errors.New("Invalid material type")}
		}
		result <- query.QueryResult{Result: material}

		close(result)
	}()

	return result
}
