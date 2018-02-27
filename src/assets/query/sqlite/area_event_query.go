package sqlite

import (
	"database/sql"

	"github.com/Tanibox/tania-server/src/assets/query"
	uuid "github.com/satori/go.uuid"
)

type AreaEventQuerySqlite struct {
	DB *sql.DB
}

func NewAreaEventQuerySqlite(db *sql.DB) query.AreaEventQuery {
	return &AreaEventQuerySqlite{DB: db}
}

func (f *AreaEventQuerySqlite) FindAllByID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	// go func() {
	// 	events := []storage.FarmEvent{}

	// 	rows, err := f.DB.Query("SELECT * FROM FARM_EVENT WHERE FARM_UID = ? ORDER BY VERSION ASC", uid)
	// 	if err != nil {
	// 		result <- query.QueryResult{Error: err}
	// 	}

	// 	rowsData := struct {
	// 		ID          int
	// 		FarmUID     string
	// 		Version     int
	// 		CreatedDate string
	// 		Event       []byte
	// 	}{}

	// 	for rows.Next() {
	// 		rows.Scan(&rowsData.ID, &rowsData.FarmUID, &rowsData.Version, &rowsData.CreatedDate, &rowsData.Event)

	// 		wrapper := query.EventWrapper{}
	// 		json.Unmarshal(rowsData.Event, &wrapper)

	// 		event, err := assertFarmEvent(wrapper)
	// 		if err != nil {
	// 			result <- query.QueryResult{Error: err}
	// 		}

	// 		farmUID, err := uuid.FromString(rowsData.FarmUID)
	// 		if err != nil {
	// 			result <- query.QueryResult{Error: err}
	// 		}

	// 		createdDate, err := time.Parse(time.RFC3339, rowsData.CreatedDate)
	// 		if err != nil {
	// 			result <- query.QueryResult{Error: err}
	// 		}

	// 		events = append(events, storage.FarmEvent{
	// 			FarmUID:     farmUID,
	// 			Version:     rowsData.Version,
	// 			CreatedDate: createdDate,
	// 			Event:       event,
	// 		})
	// 	}

	// 	result <- query.QueryResult{Result: events}
	// 	close(result)
	// }()

	return result
}
