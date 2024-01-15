package sql

import (
	"errors"
	"github.com/cmeyer18/weather-common/v3/data_structures"
	"strconv"
)

var _ PostgresTable[data_structures.UserNotification] = (*PostgresUserNotificationTable)(nil)

type PostgresUserNotificationTable struct {
	conn                          Connection
	alertOptionsTable             PostgresUserNotificationAlertOptionTable
	convectiveOutlookOptionsTable PostgresUserNotificationConvectiveOutlookOptionTable
}

func NewPostgresUserNotificationTable(conn Connection) PostgresUserNotificationTable {
	return PostgresUserNotificationTable{
		conn: conn,
	}
}

func (p *PostgresUserNotificationTable) Init() error {
	p.alertOptionsTable = NewPostgresUserNotificationsAlertOptionsTable(p.conn)
	p.convectiveOutlookOptionsTable = NewPostgresUserNotificationConvectiveOutlookOptionTable(p.conn)

	err := p.alertOptionsTable.Init()
	if err != nil {
		return err
	}

	err = p.convectiveOutlookOptionsTable.Init()
	if err != nil {
		return err
	}

	//language=SQL
	query := `CREATE TABLE IF NOT EXISTS userNotification(
		notificationId       varchar(255) primary key,
		userId				 varchar(255),
		zoneCode			 varchar(255),
		countyCode			 varchar(255),
		creationTime		 timestamptz,
		lat					 float,
		lng					 float,
		formattedAddress	 varchar(500),
		apnKey				 varchar(255),
		locationName		 varchar(255)
	)`

	err = p.conn.AddTable(query)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresUserNotificationTable) Create(userNotification data_structures.UserNotification) error {
	//language=SQL
	query := `INSERT INTO userNotification (notificationId, userId, zoneCode, countyCode, creationTime, lat, lng, formattedAddress, apnKey, locationName) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := p.conn.db.Exec(
		query,
		userNotification.NotificationId,
		userNotification.UserID,
		userNotification.ZoneCode,
		userNotification.CountyCode,
		userNotification.CreationTime,
		userNotification.Lat,
		userNotification.Lng,
		userNotification.FormattedAddress,
		userNotification.APNKey,
		userNotification.LocationName,
	)
	if err != nil {
		return err
	}

	err = p.alertOptionsTable.CreateMany(userNotification.NotificationId, userNotification.AlertOptions)
	if err != nil {
		return err
	}

	err = p.convectiveOutlookOptionsTable.CreateMany(userNotification.NotificationId, userNotification.ConvectiveOutlookOptions)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresUserNotificationTable) Find(notificationId string) (*data_structures.UserNotification, error) {
	query := `SELECT notificationId, userId, zoneCode, countyCode, creationTime, lat, lng, formattedAddress, apnKey, locationName FROM userNotification WHERE notificationId = $1`

	userNotification := data_structures.UserNotification{}

	row := p.conn.db.QueryRow(query, notificationId)

	err := row.Scan(&userNotification.NotificationId,
		&userNotification.UserID,
		&userNotification.ZoneCode,
		&userNotification.CountyCode,
		&userNotification.CreationTime,
		&userNotification.Lat,
		&userNotification.Lng,
		&userNotification.FormattedAddress,
		&userNotification.APNKey,
		&userNotification.LocationName,
	)

	if err != nil {
		return nil, err
	}

	return &userNotification, nil
}

func (p *PostgresUserNotificationTable) FindByUserId(userId string) ([]data_structures.UserNotification, error) {
	query := `SELECT notificationId, userId, zoneCode, countyCode, creationTime, lat, lng, formattedAddress, apnKey, locationName FROM userNotification WHERE userId = $1`

	row, err := p.conn.db.Query(query, userId)
	if err != nil {
		return nil, err
	}

	var userNotifications []data_structures.UserNotification
	for row.Next() {
		userNotification := data_structures.UserNotification{}

		err := row.Scan(&userNotification.NotificationId,
			&userNotification.UserID,
			&userNotification.ZoneCode,
			&userNotification.CountyCode,
			&userNotification.CreationTime,
			&userNotification.Lat,
			&userNotification.Lng,
			&userNotification.FormattedAddress,
			&userNotification.APNKey,
			&userNotification.LocationName,
		)
		if err != nil {
			return nil, err
		}

		userNotifications = append(userNotifications, userNotification)
	}

	return userNotifications, nil
}

func (p *PostgresUserNotificationTable) Delete(notificationId string) error {
	query := `DELETE FROM userNotification WHERE notificationId = $1`

	exec, err := p.conn.db.Exec(query, notificationId)
	if err != nil {
		return err
	}

	affected, err := exec.RowsAffected()
	if err != nil {
		return err
	}

	if affected != 1 {
		return errors.New("unexpected number of rows deleted, expected: 1 got:" + strconv.FormatInt(affected, 10))
	}

	return nil
}
