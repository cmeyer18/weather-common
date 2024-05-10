package sql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/lib/pq"

	"github.com/cmeyer18/weather-common/v4/data_structures"
	"github.com/cmeyer18/weather-common/v4/data_structures/geojson_v2"
	"github.com/cmeyer18/weather-common/v4/sql/internal"
	"github.com/cmeyer18/weather-common/v4/sql/internal/common_tables"
)

var _ IAlertV2Table = (*PostgresAlertV2Table)(nil)

type IAlertV2Table interface {
	common_tables.IIdTable[data_structures.AlertV2]

	SelectByLocation(codes []string, point geojson_v2.Point) ([]data_structures.AlertV2, error)

	Exists(id string) (bool, error)
}

type PostgresAlertV2Table struct {
	db              *sql.DB
	sameTable       internal.IAlertV2SAMECodesTable
	ugcTable        internal.IAlertV2UGCCodesTable
	referencesTable internal.IAlertV2ReferencesTable
}

func NewPostgresAlertV2Table(db *sql.DB) PostgresAlertV2Table {
	sameTable := internal.NewPostgresAlertV2SAMECodesTable(db)
	ugcTable := internal.NewPostgresAlertV2UGCCodesTable(db)
	referencesTable := internal.NewPostgresAlertV2ReferencesTable(db)

	return PostgresAlertV2Table{
		db:              db,
		sameTable:       &sameTable,
		ugcTable:        &ugcTable,
		referencesTable: &referencesTable,
	}
}

func (p *PostgresAlertV2Table) Insert(alert data_structures.AlertV2) error {
	ctx := context.Background()
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	statement, err := tx.Prepare(`INSERT INTO alertV2 (
		id, type, geometry, areaDesc, sent, 
		effective, onset, expires, ends, status, 
		messageType, category, severity, certainty, urgency, 
		event, sender, senderName, headline, description,
	 	instruction, response, parameters
	) 
	VALUES (
		$1, $2, 
		CASE 
			WHEN $3::TEXT IS NULL OR $3::TEXT = '' OR jsonb_typeof($3::JSONB) = 'null' THEN NULL 
			ELSE ST_GeomFromGeoJSON($3::JSONB) 
    	END,
		$4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23
	);`)
	if err != nil {
		return err
	}
	defer statement.Close()

	marshalledParameters, err := json.Marshal(alert.Parameters)
	if err != nil {
		return err
	}

	marshalledGeometryBytes, err := json.Marshal(alert.Geometry)
	if err != nil {
		return err
	}

	_, err = statement.Exec(
		alert.ID, alert.Type, marshalledGeometryBytes, alert.AreaDesc, alert.Sent,
		alert.Effective, alert.Onset, alert.Expires, alert.Ends, alert.Status,
		alert.MessageType, alert.Category, alert.Severity, alert.Certainty, alert.Urgency,
		alert.Event, alert.Sender, alert.SenderName, alert.Headline, alert.Description,
		alert.Instruction, alert.Response, marshalledParameters,
	)
	if err != nil {
		return err
	}

	if alert.Geocode != nil {
		err := p.sameTable.Insert(tx, alert.ID, alert.Geocode.SAME)
		if err != nil {
			return err
		}

		err = p.ugcTable.Insert(tx, alert.ID, alert.Geocode.UGC)
		if err != nil {
			return err
		}
	}
	if len(alert.References) != 0 {
		err := p.referencesTable.Insert(tx, alert.ID, alert.References)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresAlertV2Table) Select(id string) (*data_structures.AlertV2, error) {
	statement, err := p.db.Prepare(`
        SELECT 
            id, type, geometry::JSONB, areaDesc, sent, effective, onset, 
            expires, ends, status, messageType, category, severity, 
            certainty, urgency, event, sender, senderName, headline, 
            description, instruction, response, parameters
        FROM alertV2
        WHERE id = $1
    `)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	rows, err := statement.Query(id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	alerts, err := p.processAlertRows(rows)
	if err != nil {
		return nil, err
	}

	if len(alerts) == 0 {
		return nil, nil
	}

	return &alerts[0], nil
}

func (p *PostgresAlertV2Table) SelectByLocation(codes []string, point geojson_v2.Point) ([]data_structures.AlertV2, error) {
	statement, err := p.db.Prepare(`
	SELECT DISTINCT a.id, a.type, a.geometry::JSONB, a.areaDesc, a.sent, a.effective, a.onset, 
            a.expires, a.ends, a.status, a.messageType, a.category, a.severity, 
            a.certainty, a.urgency, a.event, a.sender, a.senderName, a.headline, 
            a.description, a.instruction, a.response, a.parameters 
	FROM alertV2 a
		LEFT JOIN alertV2_SAMECodes same ON a.id = same.alertId
		LEFT JOIN alertV2_UGCCodes ugc ON a.id = ugc.alertId
	WHERE 
	    a.geometry IS NULL AND
		a.ends >= NOW() AND (
		same.code = ANY($1::VARCHAR[]) OR 
		ugc.code = ANY($1::VARCHAR[]));`)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	rows, err := statement.Query(pq.Array(codes))
	if err != nil {
		return nil, err
	}

	geocodeActiveAlerts, err := p.processAlertRows(rows)

	statement, err = p.db.Prepare(`
	SELECT DISTINCT a.id, a.type, a.geometry::JSONB, a.areaDesc, a.sent, a.effective, a.onset, 
		a.expires, a.ends, a.status, a.messageType, a.category, a.severity, 
		a.certainty, a.urgency, a.event, a.sender, a.senderName, a.headline, 
		a.description, a.instruction, a.response, a.parameters 
	FROM alertV2 a
	WHERE 
	    ST_Contains(a.geometry, ST_GeomFromText($1, 4326)) AND
	    a.ends >= NOW()
	`)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	pointString := fmt.Sprintf("POINT (%f %f)", point.Longitude, point.Latitude)

	rows, err = statement.Query(pointString)
	if err != nil {
		return nil, err
	}

	geometryAlerts, err := p.processAlertRows(rows)

	activeAlerts := append(geocodeActiveAlerts, geometryAlerts...)

	return activeAlerts, nil
}

func (p *PostgresAlertV2Table) Exists(id string) (bool, error) {
	query := `SELECT count(id) FROM alertV2 WHERE id = $1`

	row := p.db.QueryRow(query, id)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func (p *PostgresAlertV2Table) Delete(id string) error {
	query := `DELETE FROM alertV2 WHERE id = $1`

	_, err := p.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresAlertV2Table) processAlertRows(rows *sql.Rows) ([]data_structures.AlertV2, error) {
	var alerts []data_structures.AlertV2
	for rows.Next() {
		var alert data_structures.AlertV2
		var marshalledParameters []byte
		var marshalledGeometry []byte

		err := rows.Scan(
			&alert.ID, &alert.Type, &marshalledGeometry, &alert.AreaDesc, &alert.Sent, &alert.Effective,
			&alert.Onset, &alert.Expires, &alert.Ends, &alert.Status, &alert.MessageType, &alert.Category,
			&alert.Severity, &alert.Certainty, &alert.Urgency, &alert.Event, &alert.Sender, &alert.SenderName,
			&alert.Headline, &alert.Description, &alert.Instruction, &alert.Response, &marshalledParameters,
		)
		if err != nil {
			return nil, err
		}

		if !(string(marshalledParameters) == "" || string(marshalledParameters) == `""` || string(marshalledParameters) == "null") {
			err = json.Unmarshal(marshalledParameters, &alert.Parameters)
			if err != nil {
				return nil, err
			}
		}

		if !(string(marshalledGeometry) == "" || string(marshalledGeometry) == `""` || string(marshalledGeometry) == "null") {
			err = json.Unmarshal(marshalledGeometry, &alert.Geometry)
			if err != nil {
				return nil, err
			}
		}

		sameIds, err := p.sameTable.SelectByAlertId(alert.ID)
		if err != nil {
			return nil, err
		}

		ugcIds, err := p.ugcTable.SelectByAlertId(alert.ID)
		if err != nil {
			return nil, err
		}

		if len(sameIds) > 0 || len(ugcIds) > 0 {
			alert.Geocode = &data_structures.AlertPropertiesGeocodeV2{
				SAME: sameIds,
				UGC:  ugcIds,
			}
		}

		alert.References, err = p.ugcTable.SelectByAlertId(alert.ID)
		if err != nil {
			return nil, err
		}

		alerts = append(alerts, alert)
	}

	return alerts, nil
}
