package mysql

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/assets/decoder"
	"github.com/usetania/tania-core/src/assets/domain"
	"github.com/usetania/tania-core/src/assets/repository"
	"github.com/usetania/tania-core/src/helper/structhelper"
)

type MaterialEventRepositoryMysql struct {
	DB *sql.DB
}

func NewMaterialEventRepositoryMysql(db *sql.DB) repository.MaterialEvent {
	return &MaterialEventRepositoryMysql{DB: db}
}

func (f *MaterialEventRepositoryMysql) Save(uid uuid.UUID, latestVersion int, events []interface{}) <-chan error {
	result := make(chan error)

	go func() {
		for _, v := range events {
			stmt, err := f.DB.Prepare(`INSERT INTO MATERIAL_EVENT
				(MATERIAL_UID, VERSION, CREATED_DATE, EVENT) VALUES (?, ?, ?, ?)`)
			if err != nil {
				result <- err
			}

			latestVersion++

			var eTemp interface{}

			switch val := v.(type) {
			case domain.MaterialCreated:
				val.Type = repository.MaterialEventTypeWrapper{
					Type: val.Type.Code(),
					Data: val.Type,
				}

				eTemp = val

			case domain.MaterialTypeChanged:
				val.MaterialType = repository.MaterialEventTypeWrapper{
					Type: val.MaterialType.Code(),
					Data: val.MaterialType,
				}

				eTemp = val

			default:
				eTemp = val
			}

			e, err := json.Marshal(decoder.EventWrapper{
				EventName: structhelper.GetName(eTemp),
				EventData: eTemp,
			})
			if err != nil {
				result <- err
			}

			_, err = stmt.Exec(uid.Bytes(), latestVersion, time.Now(), e)
			if err != nil {
				result <- err
			}
		}

		result <- nil
		close(result)
	}()

	return result
}
