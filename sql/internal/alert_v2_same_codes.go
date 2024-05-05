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
		query := `INSERT INTO alertV2_SAMECodes (alertId, code) VALUES ($1, $2)`

		_, err := tx.Exec(query, alertId, code)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *PostgresAlertV2SAMECodesTable) SelectByAlertId(alertId string) ([]string, error) {
	query := `SELECT code FROM alertV2_SAMECodes WHERE alertId = $1`

	row, err := p.db.Query(query, alertId)
	if err != nil {
		return nil, err
	}

	var sameCodes []string
	for row.Next() {
		var sameCode string

		err := row.Scan(&sameCode)
		if err != nil {
			return nil, err
		}

		sameCodes = append(sameCodes, sameCode)
	}

	return sameCodes, nil
}

func (p *PostgresAlertV2SAMECodesTable) Delete(tx *sql.Tx, alertId string) error {
	query := `DELETE FROM alertV2_SAMECodes WHERE alertId = $1`

	_, err := tx.Exec(query, alertId)
	if err != nil {
		return err
	}

	return nil
}
