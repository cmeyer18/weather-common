package internal

import (
	"database/sql"
)

var _ IAlertV2SAMECodesTable = (*PostgresAlertV2SAMECodesTable)(nil)

type IAlertV2SAMECodesTable interface {
	Insert(tx *sql.Tx, alertId string, codes []string) error

	SelectByAlertId(alertId string) ([]string, error)

	Delete(tx *sql.Tx, alertId string) error
}

type PostgresAlertV2SAMECodesTable struct {
	db *sql.DB
}

func NewPostgresAlertV2SAMECodesTable(db *sql.DB) PostgresAlertV2SAMECodesTable {
	return PostgresAlertV2SAMECodesTable{
		db: db,
	}
}

func (p *PostgresAlertV2SAMECodesTable) Insert(tx *sql.Tx, alertId string, codes []string) error {
	for _, code := range codes {
		statement, err := tx.Prepare(`INSERT INTO alertV2_SAMECodes (alertId, code) VALUES ($1, $2)`)
		if err != nil {
			return err
		}

		_, err = statement.Exec(alertId, code)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *PostgresAlertV2SAMECodesTable) SelectByAlertId(alertId string) ([]string, error) {
	statement, err := p.db.Prepare(`SELECT code FROM alertV2_SAMECodes WHERE alertId = $1`)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	rows, err := statement.Query(alertId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sameCodes []string
	for rows.Next() {
		var sameCode string
		err := rows.Scan(&sameCode)
		if err != nil {
			return nil, err
		}

		sameCodes = append(sameCodes, sameCode)
	}

	return sameCodes, nil
}

func (p *PostgresAlertV2SAMECodesTable) Delete(tx *sql.Tx, alertId string) error {
	statement, err := p.db.Prepare(`DELETE FROM alertV2_SAMECodes WHERE alertId = $1`)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(alertId)
	if err != nil {
		return err
	}

	return nil
}
