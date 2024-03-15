package internal

import (
	"database/sql"
	"errors"
	"github.com/cmeyer18/weather-common/v3/generative/golang"
	"strconv"
)

var _ IUserNotificationConvectiveOutlookOptionTable = (*PostgresUserNotificationConvectiveOutlookOptionTable)(nil)

type IUserNotificationConvectiveOutlookOptionTable interface {
	Insert(notificationId string, convectiveOutlookOptions []golang.ConvectiveOutlookType) error

	Update(notificationId string, convectiveOutlookOptions []golang.ConvectiveOutlookType) error

	SelectByNotificationId(notificationId string) ([]golang.ConvectiveOutlookType, error)

	Delete(notificationId string) error
}

type PostgresUserNotificationConvectiveOutlookOptionTable struct {
	db *sql.DB
}

func NewPostgresUserNotificationConvectiveOutlookOptionTable(db *sql.DB) PostgresUserNotificationConvectiveOutlookOptionTable {
	return PostgresUserNotificationConvectiveOutlookOptionTable{
		db: db,
	}
}

func (p *PostgresUserNotificationConvectiveOutlookOptionTable) insert(notificationId string, convectiveOutlookOption golang.ConvectiveOutlookType) error {
	//language=SQL
	query := `INSERT INTO userNotificationConvectiveOutlookOption (notificationId, convectiveOutlookOption) VALUES ($1, $2)`

	convertedConvectiveOutlookOption := string(convectiveOutlookOption)

	_, err := p.db.Exec(query, notificationId, convertedConvectiveOutlookOption)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresUserNotificationConvectiveOutlookOptionTable) Insert(notificationId string, convectiveOutlookOptions []golang.ConvectiveOutlookType) error {
	for _, convectiveOutlookOption := range convectiveOutlookOptions {
		err := p.insert(notificationId, convectiveOutlookOption)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *PostgresUserNotificationConvectiveOutlookOptionTable) Update(notificationId string, convectiveOutlookOptions []golang.ConvectiveOutlookType) error {
	err := p.Delete(notificationId)
	if err != nil {
		return err
	}

	err = p.Insert(notificationId, convectiveOutlookOptions)
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

func (p *PostgresUserNotificationConvectiveOutlookOptionTable) Delete(notificationId string) error {
	query := `DELETE FROM userNotificationConvectiveOutlookOption WHERE notificationId = $1`

	exec, err := p.db.Exec(query, notificationId)
	if err != nil {
		return err
	}

	affected, err := exec.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return errors.New("unexpected number of rows deleted, expected: 1 got:" + strconv.FormatInt(affected, 10))
	}

	return nil
}
