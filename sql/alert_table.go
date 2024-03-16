package sql

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"strconv"

	"github.com/lib/pq"

	"github.com/cmeyer18/weather-common/v4/data_structures"
	"github.com/cmeyer18/weather-common/v4/sql/internal/common_tables"
)

var _ IAlertTable = (*PostgresAlertTable)(nil)

type IAlertTable interface {
	common_tables.IIdTable[data_structures.Alert]

	SelectAlertsByCode(codes []string) ([]data_structures.Alert, error)

	DeleteExpiredAlerts(id string) error

	Exists(id string) (bool, error)
}

type PostgresAlertTable struct {
	db *sql.DB
}

func NewPostgresAlertTable(db *sql.DB) PostgresAlertTable {
	return PostgresAlertTable{
		db: db,
	}
}

func (p *PostgresAlertTable) Insert(alert data_structures.Alert) error {
	//language=SQL
	query := `INSERT INTO alerts (id, payload) VALUES ($1, $2)`

	marshalledAlert, err := json.Marshal(alert)
	if err != nil {
		return err
	}

	_, err = p.db.Exec(query, alert.ID, marshalledAlert)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresAlertTable) Select(id string) (*data_structures.Alert, error) {
	query := `SELECT payload FROM alerts WHERE id = $1`

	row := p.db.QueryRow(query, id)
	var rawAlert []byte
	err := row.Scan(rawAlert)
	if err != nil {
		return nil, err
	}

	var alert data_structures.Alert
	err = json.Unmarshal(rawAlert, &alert)
	if err != nil {
		return nil, err
	}

	return &alert, nil
}

func (p *PostgresAlertTable) SelectAlertsByCode(codes []string) ([]data_structures.Alert, error) {
	query := `
		SELECT payload
		FROM alerts
		WHERE 
		    (alerts.payload->'properties' ->> 'ends')::timestamptz >= NOW() AND
		    alerts.payload->'properties'->'geocode'->'UGC' IS NOT NULL
			AND EXISTS (
				SELECT 1
				FROM jsonb_array_elements_text(alerts.payload->'properties'->'geocode'->'UGC') AS ugc_element
				WHERE ugc_element.value = ANY ($1)
			)`

	// Execute the SQL query
	rows, err := p.db.Query(query, pq.Array(codes))
	if err != nil {
		log.Fatal(err)
	}

	var alerts []data_structures.Alert
	for rows.Next() {
		var alert data_structures.Alert

		var rawAlert []byte
		err := rows.Scan(&rawAlert)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(rawAlert, &alert)
		if err != nil {
			return nil, err
		}

		alerts = append(alerts, alert)
	}

	return alerts, nil
}

func (p *PostgresAlertTable) Exists(id string) (bool, error) {
	query := `SELECT count(id) FROM alerts WHERE id = $1`

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

func (p *PostgresAlertTable) Delete(id string) error {
	query := `DELETE FROM alerts WHERE id = $1`

	exec, err := p.db.Exec(query, id)
	if err != nil {
		return err
	}

	affected, err := exec.RowsAffected()
	if err != nil {
		return err
	}

	if affected != 1 {
		return errors.New("unexpected number of rows deleted, expected: 1 got:" + strconv.FormatInt(affected, 10))
	}

	return nil
}

func (p *PostgresAlertTable) DeleteExpiredAlerts(id string) error {
	query := `DELETE FROM alerts WHERE id = $1`

	exec, err := p.db.Exec(query, id)
	if err != nil {
		return err
	}

	affected, err := exec.RowsAffected()
	if err != nil {
		return err
	}

	if affected != 1 {
		return errors.New("unexpected number of rows deleted, expected: 1 got:" + strconv.FormatInt(affected, 10))
	}

	return nil
}
