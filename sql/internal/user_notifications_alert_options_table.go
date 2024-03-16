package internal

import (
	"database/sql"

	"github.com/cmeyer18/weather-common/v3/generative/golang"
)

var _ IUserNotificationAlertOptionTable = (*PostgresUserNotificationAlertOptionTable)(nil)

type IUserNotificationAlertOptionTable interface {
	Insert(tx *sql.Tx, notificationId string, alertOptions []golang.AlertType) error

	Update(tx *sql.Tx, notificationId string, alertOptions []golang.AlertType) error

	SelectByNotificationId(notificationId string) ([]golang.AlertType, error)

	Delete(tx *sql.Tx, notificationId string) error
}

type PostgresUserNotificationAlertOptionTable struct {
	db *sql.DB
}

func NewPostgresUserNotificationsAlertOptionsTable(db *sql.DB) PostgresUserNotificationAlertOptionTable {
	return PostgresUserNotificationAlertOptionTable{
		db: db,
	}
}

func (p *PostgresUserNotificationAlertOptionTable) insert(tx *sql.Tx, notificationId string, alertOption golang.AlertType) error {
	//language=SQL
	query := `INSERT INTO userNotificationAlertOption (notificationId, alertOption) VALUES ($1, $2)`

	convertedAlertOption := string(alertOption)

	_, err := tx.Exec(query, notificationId, convertedAlertOption)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresUserNotificationAlertOptionTable) Insert(tx *sql.Tx, notificationId string, alertOptions []golang.AlertType) error {
	for _, alertOption := range alertOptions {
		err := p.insert(tx, notificationId, alertOption)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *PostgresUserNotificationAlertOptionTable) Update(tx *sql.Tx, notificationId string, alertOptions []golang.AlertType) error {
	err := p.Delete(tx, notificationId)
	if err != nil {
		return err
	}

	err = p.Insert(tx, notificationId, alertOptions)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresUserNotificationAlertOptionTable) SelectByNotificationId(notificationId string) ([]golang.AlertType, error) {
	query := `SELECT alertOption FROM userNotificationAlertOption WHERE notificationId = $1`

	row, err := p.db.Query(query, notificationId)
	if err != nil {
		return nil, err
	}

	var alertTypes []golang.AlertType
	for row.Next() {
		var alertTypeString string

		err := row.Scan(&alertTypeString)
		if err != nil {
			return nil, err
		}

		alertType := golang.AlertType(alertTypeString)

		alertTypes = append(alertTypes, alertType)
	}

	return alertTypes, nil
}

func (p *PostgresUserNotificationAlertOptionTable) Delete(tx *sql.Tx, notificationId string) error {
	query := `DELETE FROM userNotificationAlertOption WHERE notificationId = $1`

	_, err := tx.Exec(query, notificationId)
	if err != nil {
		return err
	}

	return nil
}
