package internal

import (
	"database/sql"
	"errors"
	"github.com/cmeyer18/weather-common/v3/generative/golang"
	"strconv"
)

var _ IUserNotificationAlertOptionTable = (*PostgresUserNotificationAlertOptionTable)(nil)

type IUserNotificationAlertOptionTable interface {
	Insert(notificationId string, alertOptions []golang.AlertType) error

	Update(notificationId string, alertOptions []golang.AlertType) error

	SelectByNotificationId(notificationId string) ([]golang.AlertType, error)

	Delete(notificationId string) error
}

type PostgresUserNotificationAlertOptionTable struct {
	db *sql.DB
}

func NewPostgresUserNotificationsAlertOptionsTable(db *sql.DB) PostgresUserNotificationAlertOptionTable {
	return PostgresUserNotificationAlertOptionTable{
		db: db,
	}
}

func (p *PostgresUserNotificationAlertOptionTable) insert(notificationId string, alertOption golang.AlertType) error {
	//language=SQL
	query := `INSERT INTO userNotificationAlertOption (notificationId, alertOption) VALUES ($1, $2)`

	convertedAlertOption := string(alertOption)

	_, err := p.db.Exec(query, notificationId, convertedAlertOption)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresUserNotificationAlertOptionTable) Insert(notificationId string, alertOptions []golang.AlertType) error {
	for _, alertOption := range alertOptions {
		err := p.insert(notificationId, alertOption)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *PostgresUserNotificationAlertOptionTable) Update(notificationId string, alertOptions []golang.AlertType) error {
	err := p.Delete(notificationId)
	if err != nil {
		return err
	}

	err = p.Insert(notificationId, alertOptions)
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

func (p *PostgresUserNotificationAlertOptionTable) Delete(notificationId string) error {
	query := `DELETE FROM userNotificationAlertOption WHERE notificationId = $1`

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
