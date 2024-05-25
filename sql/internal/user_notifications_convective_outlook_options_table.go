package internal

import (
	"database/sql"

	"github.com/cmeyer18/weather-common/v5/generative/golang"
)

var _ IUserNotificationConvectiveOutlookOptionTable = (*PostgresUserNotificationConvectiveOutlookOptionTable)(nil)

type IUserNotificationConvectiveOutlookOptionTable interface {
	Insert(tx *sql.Tx, notificationId string, convectiveOutlookOptions []golang.ConvectiveOutlookType) error

	Update(tx *sql.Tx, notificationId string, convectiveOutlookOptions []golang.ConvectiveOutlookType) error

	SelectByNotificationId(notificationId string) ([]golang.ConvectiveOutlookType, error)

	Delete(tx *sql.Tx, notificationId string) error
}

type PostgresUserNotificationConvectiveOutlookOptionTable struct {
	db *sql.DB
}

func NewPostgresUserNotificationConvectiveOutlookOptionTable(db *sql.DB) PostgresUserNotificationConvectiveOutlookOptionTable {
	return PostgresUserNotificationConvectiveOutlookOptionTable{
		db: db,
	}
}

func (p *PostgresUserNotificationConvectiveOutlookOptionTable) insert(tx *sql.Tx, notificationId string, convectiveOutlookOption golang.ConvectiveOutlookType) error {
	//language=SQL
	query := `INSERT INTO userNotificationConvectiveOutlookOption (notificationId, convectiveOutlookOption) VALUES ($1, $2)`

	convertedConvectiveOutlookOption := string(convectiveOutlookOption)

	_, err := tx.Exec(query, notificationId, convertedConvectiveOutlookOption)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresUserNotificationConvectiveOutlookOptionTable) Insert(tx *sql.Tx, notificationId string, convectiveOutlookOptions []golang.ConvectiveOutlookType) error {
	for _, convectiveOutlookOption := range convectiveOutlookOptions {
		err := p.insert(tx, notificationId, convectiveOutlookOption)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *PostgresUserNotificationConvectiveOutlookOptionTable) Update(tx *sql.Tx, notificationId string, convectiveOutlookOptions []golang.ConvectiveOutlookType) error {
	err := p.Delete(tx, notificationId)
	if err != nil {
		return err
	}

	err = p.Insert(tx, notificationId, convectiveOutlookOptions)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresUserNotificationConvectiveOutlookOptionTable) SelectByNotificationId(notificationId string) ([]golang.ConvectiveOutlookType, error) {
	query := `SELECT convectiveOutlookOption FROM userNotificationConvectiveOutlookOption WHERE notificationId = $1`

	row, err := p.db.Query(query, notificationId)
	if err != nil {
		return nil, err
	}

	var convectiveOutlookTypes []golang.ConvectiveOutlookType
	for row.Next() {
		var convectiveOutlookTypeString string

		err := row.Scan(&convectiveOutlookTypeString)
		if err != nil {
			return nil, err
		}

		convectiveOutlookType := golang.ConvectiveOutlookType(convectiveOutlookTypeString)

		convectiveOutlookTypes = append(convectiveOutlookTypes, convectiveOutlookType)
	}

	return convectiveOutlookTypes, nil
}

func (p *PostgresUserNotificationConvectiveOutlookOptionTable) Delete(tx *sql.Tx, notificationId string) error {
	query := `DELETE FROM userNotificationConvectiveOutlookOption WHERE notificationId = $1`

	_, err := tx.Exec(query, notificationId)
	if err != nil {
		return err
	}

	return nil
}
