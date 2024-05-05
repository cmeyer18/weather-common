package internal

import (
	"database/sql"
)

var _ IAlertV2ReferencesTable = (*PostgresAlertV2ReferencesTable)(nil)

type IAlertV2ReferencesTable interface {
	Insert(tx *sql.Tx, alertId string, referencedAlertIds []string) error

	SelectByAlertId(alertId string) ([]string, error)

	Delete(tx *sql.Tx, alertId string) error
}

type PostgresAlertV2ReferencesTable struct {
	db *sql.DB
}

func NewPostgresAlertV2ReferencesTable(db *sql.DB) PostgresAlertV2ReferencesTable {
	return PostgresAlertV2ReferencesTable{
		db: db,
	}
}

func (p *PostgresAlertV2ReferencesTable) Insert(tx *sql.Tx, alertId string, referencedAlertIds []string) error {
	for _, referencedAlertId := range referencedAlertIds {
		query := `INSERT INTO alertV2_References (alertId, referenceId) VALUES ($1, $2)`

		_, err := tx.Exec(query, alertId, referencedAlertId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *PostgresAlertV2ReferencesTable) SelectByAlertId(alertId string) ([]string, error) {
	query := `SELECT referenceId FROM alertV2_References WHERE alertId = $1`

	row, err := p.db.Query(query, alertId)
	if err != nil {
		return nil, err
	}

	var referenceIds []string
	for row.Next() {
		var referenceId string

		err := row.Scan(&referenceId)
		if err != nil {
			return nil, err
		}

		referenceIds = append(referenceIds, referenceId)
	}

	return referenceIds, nil
}

func (p *PostgresAlertV2ReferencesTable) Delete(tx *sql.Tx, alertId string) error {
	query := `DELETE FROM alertV2_References WHERE alertId = $1`

	_, err := tx.Exec(query, alertId)
	if err != nil {
		return err
	}

	return nil
}
