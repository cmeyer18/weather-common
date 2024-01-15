package sql

import (
	"errors"
	"github.com/cmeyer18/weather-common/v3/generative/golang"
	"strconv"
)

var _ PostgresTable[string] = (*PostgresUserNotificationAlertOptionTable)(nil)

type PostgresUserNotificationAlertOptionTable struct {
	conn Connection
}

func NewPostgresUserNotificationsAlertOptionsTable(conn Connection) PostgresUserNotificationAlertOptionTable {
	return PostgresUserNotificationAlertOptionTable{
		conn: conn,
	}
}

func (p *PostgresUserNotificationAlertOptionTable) Init() error {
	//language=SQL
	query := `CREATE TABLE IF NOT EXISTS userNotificationAlertOption(
		notificationId       varchar(255) primary key,
		alertOption		 varchar(255)
	)`

	err := p.conn.AddTable(query)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresUserNotificationAlertOptionTable) Create(notificationId string, alertOption golang.AlertType) error {
	//language=SQL
	query := `INSERT INTO userNotificationAlertOption (notificationId, alertOption) VALUES ($1, $2)`

	convertedAlertOption := string(alertOption)

	_, err := p.conn.db.Exec(query, notificationId, convertedAlertOption)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresUserNotificationAlertOptionTable) CreateMany(notificationId string, alertOptions []golang.AlertType) error {
	for _, alertOption := range alertOptions {
		err := p.Create(notificationId, alertOption)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *PostgresUserNotificationAlertOptionTable) GetAlertOptionsForNotificationId(notificationId string) ([]golang.AlertType, error) {
	query := `SELECT alertOption FROM userNotificationAlertOption WHERE notificationId = $1`

	row, err := p.conn.db.Query(query, notificationId)
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

	exec, err := p.conn.db.Exec(query, notificationId)
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
