package sqlite

import (
	"database/sql"
	"time"

	"github.com/usetania/tania-core/src/assets/domain"
	"github.com/usetania/tania-core/src/assets/repository"
	"github.com/usetania/tania-core/src/assets/storage"
)

type MaterialReadRepositorySqlite struct {
	DB *sql.DB
}

func NewMaterialReadRepositorySqlite(db *sql.DB) repository.MaterialRead {
	return &MaterialReadRepositorySqlite{DB: db}
}

func (f *MaterialReadRepositorySqlite) Save(materialRead *storage.MaterialRead) <-chan error {
	result := make(chan error)

	go func() {
		count := 0

		err := f.DB.QueryRow(`SELECT COUNT(*) FROM MATERIAL_READ WHERE UID = ?`, materialRead.UID).Scan(&count)
		if err != nil {
			result <- err
		}

		var typeData string
		switch t := materialRead.Type.(type) {
		case domain.MaterialTypeSeed:
			typeData = t.PlantType.Code
		case domain.MaterialTypePlant:
			typeData = t.PlantType.Code
		case domain.MaterialTypeAgrochemical:
			typeData = t.ChemicalType.Code
		case domain.MaterialTypeSeedingContainer:
			typeData = t.ContainerType.Code
		}

		expirationDate := ""
		if materialRead.ExpirationDate != nil {
			expirationDate = materialRead.ExpirationDate.Format(time.RFC3339)
		}

		if count > 0 {
			_, err = f.DB.Exec(`UPDATE MATERIAL_READ SET
				NAME = ?, PRICE_PER_UNIT = ?, CURRENCY_CODE = ?, TYPE = ?, TYPE_DATA = ?,
				QUANTITY = ?, QUANTITY_UNIT = ?, EXPIRATION_DATE = ?, NOTES = ?,
				PRODUCED_BY = ?, CREATED_DATE = ?
				WHERE UID = ?`,
				materialRead.Name,
				materialRead.PricePerUnit.Amount,
				materialRead.PricePerUnit.CurrencyCode,
				materialRead.Type.Code(),
				typeData,
				materialRead.Quantity.Value,
				materialRead.Quantity.Unit.Code,
				expirationDate,
				materialRead.Notes,
				materialRead.ProducedBy,
				materialRead.CreatedDate.Format(time.RFC3339),
				materialRead.UID)

			if err != nil {
				result <- err
			}
		} else {
			_, err = f.DB.Exec(`INSERT INTO MATERIAL_READ
				(UID, NAME, PRICE_PER_UNIT, CURRENCY_CODE, TYPE, TYPE_DATA, QUANTITY,
				QUANTITY_UNIT, EXPIRATION_DATE, NOTES, PRODUCED_BY, CREATED_DATE)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				materialRead.UID,
				materialRead.Name,
				materialRead.PricePerUnit.Amount,
				materialRead.PricePerUnit.CurrencyCode,
				materialRead.Type.Code(),
				typeData,
				materialRead.Quantity.Value,
				materialRead.Quantity.Unit.Code,
				expirationDate,
				materialRead.Notes,
				materialRead.ProducedBy,
				materialRead.CreatedDate.Format(time.RFC3339))

			if err != nil {
				result <- err
			}
		}

		result <- nil
		close(result)
	}()

	return result
}
