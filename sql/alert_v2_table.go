package sql

import (
	"context"
	"database/sql"
	"encoding/json"
	"regexp"

	"github.com/cmeyer18/weather-common/v4/data_structures"
	"github.com/cmeyer18/weather-common/v4/sql/internal"
	"github.com/cmeyer18/weather-common/v4/sql/internal/common_tables"
)

var _ IAlertV2Table = (*PostgresAlertV2Table)(nil)

type IAlertV2Table interface {
	common_tables.IIdTable[data_structures.AlertV2]

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
		id, type, geometry, areaDesc, sent, effective, onset, 
		expires, ends, status, messageType, category, severity, 
		certainty, urgency, event, sender, senderName, headline, 
		description, instruction, response, parameters
	) 
	VALUES (
		$1, $2, 
		CASE 
			WHEN $3::TEXT IS NULL OR $3::TEXT = '' OR  jsonb_typeof($3::JSONB) = 'null' THEN NULL 
			ELSE ST_GeomFromGeoJSON($3::JSONB) 
    	END,
		$4, $5, $6, $7, $8, $9, $10, $11, $12, $13, 
		$14, $15, $16, $17, $18, $19, $20, $21, $22, $23
	);`)
	if err != nil {
		return err
	}

	marshalledParameters, err := json.Marshal(alert.Parameters)
	if err != nil {
		return err
	}

	marshalledGeometryBytes, err := json.Marshal(alert.Geometry)
	if err != nil {
		return err
	}

	// Clean up the json, SQL doesn't like these escape characters.
	pattern := regexp.MustCompile(`\\+`)
	unescapedString := pattern.ReplaceAllString(string(marshalledGeometryBytes), "")

	_, err = statement.Exec(
		alert.ID, alert.Type, unescapedString, alert.AreaDesc, alert.Sent, alert.Effective,
		alert.Onset, alert.Expires, alert.Ends, alert.Status, alert.MessageType, alert.Category,
		alert.Severity, alert.Certainty, alert.Urgency, alert.Event, alert.Sender, alert.SenderName,
		alert.Headline, alert.Description, alert.Instruction, alert.Response, marshalledParameters,
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
	stmt, err := p.db.Prepare(`
        SELECT 
            id, type, geometry, areaDesc, sent, effective, onset, 
            expires, ends, status, messageType, category, severity, 
            certainty, urgency, event, sender, senderName, headline, 
            description, instruction, response, parameters
        FROM alertV2
        WHERE id = $1
    `)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []data_structures.AlertV2
	for rows.Next() {
		var alert data_structures.AlertV2
		var marshalledParameters string

		err = rows.Scan(
			&alert.ID, &alert.Type, &alert.Geometry, &alert.AreaDesc, &alert.Sent, &alert.Effective,
			&alert.Onset, &alert.Expires, &alert.Ends, &alert.Status, &alert.MessageType, &alert.Category,
			&alert.Severity, &alert.Certainty, &alert.Urgency, &alert.Event, &alert.Sender, &alert.SenderName,
			&alert.Headline, &alert.Description, &alert.Instruction, &alert.Response, &marshalledParameters,
		)
		if err != nil {
			return nil, err
		}

		var parameters map[string]interface{}
		err := json.Unmarshal([]byte(marshalledParameters), &parameters)
		if err != nil {
			return nil, err
		}

		sameIds, err := p.sameTable.SelectByAlertId(id)
		if err != nil {
			return nil, err
		}

		ugcIds, err := p.ugcTable.SelectByAlertId(id)
		if err != nil {
			return nil, err
		}

		if len(sameIds) > 0 || len(ugcIds) > 0 {
			alert.Geocode = &data_structures.AlertPropertiesGeocodeV2{
				SAME: sameIds,
				UGC:  ugcIds,
			}
		}

		alert.References, err = p.ugcTable.SelectByAlertId(id)
		if err != nil {
			return nil, err
		}
	}

	if len(alerts) == 0 {
		return nil, nil
	}

	return &alerts[0], nil
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
