package internal

import (
	"database/sql"
)

var _ IAlertV2UGCCodesTable = (*PostgresAlertV2UGCCodesTable)(nil)

type IAlertV2UGCCodesTable interface {
	Insert(tx *sql.Tx, alertId string, codes []string) error

	SelectByAlertId(alertId string) ([]string, error)

	Delete(tx *sql.Tx, alertId string) error
}

type PostgresAlertV2UGCCodesTable struct {
	db *sql.DB
}

func NewPostgresAlertV2UGCCodesTable(db *sql.DB) PostgresAlertV2UGCCodesTable {
	return PostgresAlertV2UGCCodesTable{
		db: db,
	}
}

func (p *PostgresAlertV2UGCCodesTable) Insert(tx *sql.Tx, alertId string, codes []string) error {
	for _, code := range codes {
		statement, err := tx.Prepare(`INSERT INTO alertV2_UGCCodes (alertId, code) VALUES ($1, $2)`)
		if err != nil {
			return err
		}
		defer statement.Close()

		_, err = statement.Exec(alertId, code)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *PostgresAlertV2UGCCodesTable) SelectByAlertId(alertId string) ([]string, error) {
	statement, err := p.db.Prepare(`SELECT code FROM alertV2_UGCCodes WHERE alertId = $1`)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	row, err := statement.Query(alertId)
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

func (p *PostgresAlertV2UGCCodesTable) Delete(tx *sql.Tx, alertId string) error {
	statement, err := p.db.Prepare(`DELETE FROM alertV2_UGCCodes WHERE alertId = $1`)
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
