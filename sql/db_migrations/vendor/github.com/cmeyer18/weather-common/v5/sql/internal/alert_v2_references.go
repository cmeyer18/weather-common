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
		statement, err := tx.Prepare(`INSERT INTO alertV2_References (alertId, referenceId) VALUES ($1, $2)`)
		if err != nil {
			return err
		}
		defer statement.Close()

		_, err = statement.Exec(alertId, referencedAlertId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *PostgresAlertV2ReferencesTable) SelectByAlertId(alertId string) ([]string, error) {
	statement, err := p.db.Prepare(`SELECT referenceId FROM alertV2_References WHERE alertId = $1`)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	rows, err := statement.Query(alertId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var referenceIds []string
	for rows.Next() {
		var referenceId string
		err := rows.Scan(&referenceId)
		if err != nil {
			return nil, err
		}

		referenceIds = append(referenceIds, referenceId)
	}

	return referenceIds, nil
}

func (p *PostgresAlertV2ReferencesTable) Delete(tx *sql.Tx, alertId string) error {
	statement, err := p.db.Prepare(`DELETE FROM alertV2_References WHERE alertId = $1`)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = tx.Exec(alertId)
	if err != nil {
		return err
	}

	return nil
}
