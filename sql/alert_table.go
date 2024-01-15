package sql

import (
	"encoding/json"
	"errors"
	"github.com/cmeyer18/weather-common/v3/data_structures"
	"strconv"
)

var _ PostgresTable[data_structures.Alert] = (*PostgresAlertTable)(nil)

type PostgresAlertTable struct {
	conn Connection
}

func NewPostgresAlertTable(conn Connection) PostgresAlertTable {
	return PostgresAlertTable{
		conn: conn,
	}
}

func (p *PostgresAlertTable) Init() error {
	//language=SQL
	query := `CREATE TABLE IF NOT EXISTS alerts(
		id       varchar(255) primary key, 
		payload  jsonb not null
	)`

	err := p.conn.AddTable(query)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresAlertTable) Create(alert *data_structures.Alert) error {
	//language=SQL
	query := `INSERT INTO alerts (id, payload) VALUES ($1, $2)`

	marshalledAlert, err := json.Marshal(alert)
	if err != nil {
		return err
	}

	_, err = p.conn.db.Exec(query, alert.ID, marshalledAlert)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresAlertTable) Find(id string) (*data_structures.Alert, error) {
	query := `SELECT payload FROM alerts WHERE id = $1`

	row := p.conn.db.QueryRow(query, id)
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

func (p *PostgresAlertTable) Exists(id string) (bool, error) {
	query := `SELECT count(id) FROM alerts WHERE id = $1`

	row := p.conn.db.QueryRow(query, id)

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

	exec, err := p.conn.db.Exec(query, id)
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

	exec, err := p.conn.db.Exec(query, id)
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
