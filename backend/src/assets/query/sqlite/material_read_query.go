package sqlite

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/assets/domain"
	"github.com/usetania/tania-core/src/assets/query"
	"github.com/usetania/tania-core/src/assets/storage"
	"github.com/usetania/tania-core/src/helper/paginationhelper"
)

type MaterialReadQuerySqlite struct {
	DB *sql.DB
}

func NewMaterialReadQuerySqlite(db *sql.DB) query.MaterialRead {
	return MaterialReadQuerySqlite{DB: db}
}

type materialReadResult struct {
	UID            string
	Name           string
	PricePerUnit   string
	CurrencyCode   string
	Type           string
	TypeData       string
	Quantity       float32
	QuantityUnit   string
	ExpirationDate sql.NullString
	Notes          sql.NullString
	ProducedBy     sql.NullString
	CreatedDate    string
}

func (q MaterialReadQuerySqlite) FindAll(materialType, materialTypeDetail string, page, limit int) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		materialReads := []storage.MaterialRead{}

		var params []interface{}

		sql := "SELECT * FROM MATERIAL_READ WHERE 1 = 1"

		if materialType != "" {
			t := strings.Split(materialType, ",")

			sql += " AND TYPE = ?"

			params = append(params, t[0])

			for _, v := range t[1:] {
				sql += " OR TYPE = ?"

				params = append(params, v)
			}
		}

		if materialTypeDetail != "" {
			t := strings.Split(materialTypeDetail, ",")

			sql += " AND TYPE_DATA = ?"

			params = append(params, t[0])

			for _, v := range t[1:] {
				sql += " OR TYPE_DATA = ?"

				params = append(params, v)
			}
		}

		sql += " ORDER BY CREATED_DATE DESC"

		if page != 0 && limit != 0 {
			sql += " LIMIT ? OFFSET ?"
			offset := paginationhelper.CalculatePageToOffset(page, limit)
			params = append(params, limit, offset)
		}

		rows, err := q.DB.Query(sql, params...)
		if err != nil {
			result <- query.Result{Error: err}
		}

		for rows.Next() {
			rowsData := materialReadResult{}

			err = rows.Scan(
				&rowsData.UID,
				&rowsData.Name,
				&rowsData.PricePerUnit,
				&rowsData.CurrencyCode,
				&rowsData.Type,
				&rowsData.TypeData,
				&rowsData.Quantity,
				&rowsData.QuantityUnit,
				&rowsData.ExpirationDate,
				&rowsData.Notes,
				&rowsData.ProducedBy,
				&rowsData.CreatedDate,
			)

			if err != nil {
				result <- query.Result{Error: err}
			}

			materialUID, err := uuid.FromString(rowsData.UID)
			if err != nil {
				result <- query.Result{Error: err}
			}

			var mExpDate *time.Time

			if rowsData.ExpirationDate.Valid && rowsData.ExpirationDate.String != "" {
				date, err := time.Parse(time.RFC3339, rowsData.ExpirationDate.String)
				if err != nil {
					result <- query.Result{Error: err}
				}

				mExpDate = &date
			}

			mCreatedDate, err := time.Parse(time.RFC3339, rowsData.CreatedDate)
			if err != nil {
				result <- query.Result{Error: err}
			}

			pricePerUnit, err := domain.CreatePricePerUnit(rowsData.PricePerUnit, rowsData.CurrencyCode)
			if err != nil {
				result <- query.Result{Error: err}
			}

			var materialType storage.MaterialType

			switch rowsData.Type {
			case domain.MaterialTypePlantCode:
				materialType, err = domain.CreateMaterialTypePlant(rowsData.TypeData)
				if err != nil {
					result <- query.Result{Error: err}
				}
			case domain.MaterialTypeSeedCode:
				materialType, err = domain.CreateMaterialTypeSeed(rowsData.TypeData)
				if err != nil {
					result <- query.Result{Error: err}
				}
			case domain.MaterialTypeGrowingMediumCode:
				materialType = domain.MaterialTypeGrowingMedium{}
			case domain.MaterialTypeAgrochemicalCode:
				materialType, err = domain.CreateMaterialTypeAgrochemical(rowsData.TypeData)
				if err != nil {
					result <- query.Result{Error: err}
				}
			case domain.MaterialTypeLabelAndCropSupportCode:
				materialType = domain.MaterialTypeLabelAndCropSupport{}
			case domain.MaterialTypeSeedingContainerCode:
				materialType, err = domain.CreateMaterialTypeSeedingContainer(rowsData.TypeData)
				if err != nil {
					result <- query.Result{Error: err}
				}
			case domain.MaterialTypePostHarvestSupplyCode:
				materialType = domain.MaterialTypePostHarvestSupply{}
			case domain.MaterialTypeOtherCode:
				materialType = domain.MaterialTypeOther{}
			default:
				result <- query.Result{Error: errors.New("invalid material type")}
			}

			qtyUnit := domain.GetMaterialQuantityUnit(rowsData.Type, rowsData.QuantityUnit)
			if qtyUnit == (domain.MaterialQuantityUnit{}) {
				result <- query.Result{Error: errors.New("invalid quantity unit")}
			}

			var notes *string
			if rowsData.Notes.Valid {
				notes = &rowsData.Notes.String
			}

			var producedBy *string
			if rowsData.ProducedBy.Valid {
				producedBy = &rowsData.ProducedBy.String
			}

			materialReads = append(materialReads, storage.MaterialRead{
				UID:          materialUID,
				Name:         rowsData.Name,
				PricePerUnit: storage.PricePerUnit(pricePerUnit),
				Type:         materialType,
				Quantity: storage.MaterialQuantity{
					Unit:  qtyUnit,
					Value: rowsData.Quantity,
				},
				ExpirationDate: mExpDate,
				Notes:          notes,
				ProducedBy:     producedBy,
				CreatedDate:    mCreatedDate,
			})
		}

		result <- query.Result{Result: materialReads}
		close(result)
	}()

	return result
}

func (q MaterialReadQuerySqlite) CountAll(materialType, materialTypeDetail string) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		total := 0

		var params []interface{}

		sql := "SELECT COUNT(UID) FROM MATERIAL_READ WHERE 1 = 1"

		if materialType != "" {
			t := strings.Split(materialType, ",")

			sql += " AND TYPE = ?"

			params = append(params, t[0])

			for _, v := range t[1:] {
				sql += " OR TYPE = ?"

				params = append(params, v)
			}
		}

		if materialTypeDetail != "" {
			t := strings.Split(materialTypeDetail, ",")

			sql += " AND TYPE_DATA = ?"

			params = append(params, t[0])

			for _, v := range t[1:] {
				sql += " OR TYPE_DATA = ?"

				params = append(params, v)
			}
		}

		err := q.DB.QueryRow(sql, params...).Scan(&total)
		if err != nil {
			result <- query.Result{Error: err}
		}

		result <- query.Result{Result: total}
		close(result)
	}()

	return result
}

func (q MaterialReadQuerySqlite) FindByID(materialUID uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		materialRead := storage.MaterialRead{}
		rowsData := materialReadResult{}

		err := q.DB.QueryRow(`SELECT * FROM MATERIAL_READ WHERE UID = ?`, materialUID).Scan(
			&rowsData.UID,
			&rowsData.Name,
			&rowsData.PricePerUnit,
			&rowsData.CurrencyCode,
			&rowsData.Type,
			&rowsData.TypeData,
			&rowsData.Quantity,
			&rowsData.QuantityUnit,
			&rowsData.ExpirationDate,
			&rowsData.Notes,
			&rowsData.ProducedBy,
			&rowsData.CreatedDate,
		)

		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			result <- query.Result{Error: err}
		}

		if errors.Is(err, sql.ErrNoRows) {
			result <- query.Result{Result: materialRead}
		}

		materialUID, err := uuid.FromString(rowsData.UID)
		if err != nil {
			result <- query.Result{Error: err}
		}

		var mExpDate *time.Time

		if rowsData.ExpirationDate.Valid && rowsData.ExpirationDate.String != "" {
			date, err := time.Parse(time.RFC3339, rowsData.ExpirationDate.String)
			if err != nil {
				result <- query.Result{Error: err}
			}

			mExpDate = &date
		}

		mCreatedDate, err := time.Parse(time.RFC3339, rowsData.CreatedDate)
		if err != nil {
			result <- query.Result{Error: err}
		}

		pricePerUnit, err := domain.CreatePricePerUnit(rowsData.PricePerUnit, rowsData.CurrencyCode)
		if err != nil {
			result <- query.Result{Error: err}
		}

		var materialType storage.MaterialType

		switch rowsData.Type {
		case domain.MaterialTypePlantCode:
			materialType, err = domain.CreateMaterialTypePlant(rowsData.TypeData)
			if err != nil {
				result <- query.Result{Error: err}
			}
		case domain.MaterialTypeSeedCode:
			materialType, err = domain.CreateMaterialTypeSeed(rowsData.TypeData)
			if err != nil {
				result <- query.Result{Error: err}
			}
		case domain.MaterialTypeGrowingMediumCode:
			materialType = domain.MaterialTypeGrowingMedium{}
		case domain.MaterialTypeAgrochemicalCode:
			materialType, err = domain.CreateMaterialTypeAgrochemical(rowsData.TypeData)
			if err != nil {
				result <- query.Result{Error: err}
			}
		case domain.MaterialTypeLabelAndCropSupportCode:
			materialType = domain.MaterialTypeLabelAndCropSupport{}
		case domain.MaterialTypeSeedingContainerCode:
			materialType, err = domain.CreateMaterialTypeSeedingContainer(rowsData.TypeData)
			if err != nil {
				result <- query.Result{Error: err}
			}
		case domain.MaterialTypePostHarvestSupplyCode:
			materialType = domain.MaterialTypePostHarvestSupply{}
		case domain.MaterialTypeOtherCode:
			materialType = domain.MaterialTypeOther{}
		default:
			result <- query.Result{Error: errors.New("invalid material type")}
		}

		qtyUnit := domain.GetMaterialQuantityUnit(rowsData.Type, rowsData.QuantityUnit)
		if qtyUnit == (domain.MaterialQuantityUnit{}) {
			result <- query.Result{Error: errors.New("invalid quantity unit")}
		}

		var notes *string
		if rowsData.Notes.Valid {
			notes = &rowsData.Notes.String
		}

		var producedBy *string
		if rowsData.ProducedBy.Valid {
			producedBy = &rowsData.ProducedBy.String
		}

		materialRead = storage.MaterialRead{
			UID:          materialUID,
			Name:         rowsData.Name,
			PricePerUnit: storage.PricePerUnit(pricePerUnit),
			Type:         materialType,
			Quantity: storage.MaterialQuantity{
				Unit:  qtyUnit,
				Value: rowsData.Quantity,
			},
			ExpirationDate: mExpDate,
			Notes:          notes,
			ProducedBy:     producedBy,
			CreatedDate:    mCreatedDate,
		}

		result <- query.Result{Result: materialRead}
		close(result)
	}()

	return result
}
