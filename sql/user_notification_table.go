package sql

import (
	"database/sql"
	"errors"
	"github.com/cmeyer18/weather-common/v3/data_structures"
	"github.com/cmeyer18/weather-common/v3/sql/internal"
	"github.com/cmeyer18/weather-common/v3/sql/internal/common_tables"
	"strconv"
)

var _ IUserNotificationTable = (*PostgresUserNotificationTable)(nil)

type IUserNotificationTable interface {
	common_tables.IIdTable[data_structures.UserNotification]

	// Deprecated: use Insert
	Create(userNotification data_structures.UserNotification) error

	// Deprecated: use Select
	Get(notificationId string) (*data_structures.UserNotification, error)

	// Deprecated: use SelectAll
	GetAll() ([]data_structures.UserNotification, error)

	SelectAll() ([]data_structures.UserNotification, error)

	// Deprecated: use SelectByUserId
	GetByUserId(userId string) ([]data_structures.UserNotification, error)

	SelectByUserId(userId string) ([]data_structures.UserNotification, error)

	// Deprecated: use SelectByCodes
	GetByCodes(codes []string) ([]data_structures.UserNotification, error)

	SelectByCodes(codes []string) ([]data_structures.UserNotification, error)

	// Deprecated: use SelectNotificationsWithConvectiveOutlook
	GetNotificationsWithConvectiveOutlookOptions() ([]data_structures.UserNotification, error)

	SelectNotificationsWithConvectiveOutlook() ([]data_structures.UserNotification, error)
}

type PostgresUserNotificationTable struct {
	db                            *sql.DB
	alertOptionsTable             internal.IUserNotificationAlertOptionTable
	convectiveOutlookOptionsTable internal.IUserNotificationConvectiveOutlookOptionTable
}

func NewPostgresUserNotificationTable(db *sql.DB) PostgresUserNotificationTable {
	alertOptionsTable := internal.NewPostgresUserNotificationsAlertOptionsTable(db)
	convectiveOutlookOptionsTable := internal.NewPostgresUserNotificationConvectiveOutlookOptionTable(db)

	return PostgresUserNotificationTable{
		db:                            db,
		alertOptionsTable:             &alertOptionsTable,
		convectiveOutlookOptionsTable: &convectiveOutlookOptionsTable,
	}
}

// Deprecated: use the weather-db-migrator docker images for setup
func (p *PostgresUserNotificationTable) Init() error {
	return nil
}

func (p *PostgresUserNotificationTable) Insert(userNotification data_structures.UserNotification) error {
	//language=SQL
	query := `INSERT INTO userNotification (notificationId, userId, zoneCode, countyCode, creationTime, lat, lng, formattedAddress, apnKey, locationName) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := p.db.Exec(
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

	err = p.alertOptionsTable.Insert(userNotification.NotificationId, userNotification.AlertOptions)
	if err != nil {
		return err
	}

	err = p.convectiveOutlookOptionsTable.Insert(userNotification.NotificationId, userNotification.ConvectiveOutlookOptions)
	if err != nil {
		return err
	}

	return nil
}

// Deprecated: use Insert
func (p *PostgresUserNotificationTable) Create(userNotification data_structures.UserNotification) error {
	return p.Insert(userNotification)
}

func (p *PostgresUserNotificationTable) Select(id string) (*data_structures.UserNotification, error) {
	query := `SELECT notificationId, userId, zoneCode, countyCode, creationTime, lat, lng, formattedAddress, apnKey, locationName FROM userNotification WHERE notificationId = $1`

	userNotification := data_structures.UserNotification{}

	row := p.db.QueryRow(query, id)

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

	convectiveOptions, err := p.convectiveOutlookOptionsTable.SelectByNotificationId(id)
	if err != nil {
		return nil, err
	}

	alertOptions, err := p.alertOptionsTable.SelectByNotificationId(id)
	if err != nil {
		return nil, err
	}

	userNotification.ConvectiveOutlookOptions = convectiveOptions
	userNotification.AlertOptions = alertOptions

	if err != nil {
		return nil, err
	}

	return &userNotification, nil
}

// Deprecated: use Select
func (p *PostgresUserNotificationTable) Get(notificationId string) (*data_structures.UserNotification, error) {
	return p.Select(notificationId)
}

func (p *PostgresUserNotificationTable) SelectAll() ([]data_structures.UserNotification, error) {
	query := `SELECT notificationId, userId, zoneCode, countyCode, creationTime, lat, lng, formattedAddress, apnKey, locationName FROM userNotification`

	row, err := p.db.Query(query)
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

		convectiveOptions, err := p.convectiveOutlookOptionsTable.SelectByNotificationId(userNotification.NotificationId)
		if err != nil {
			return nil, err
		}

		alertOptions, err := p.alertOptionsTable.SelectByNotificationId(userNotification.NotificationId)
		if err != nil {
			return nil, err
		}

		userNotification.ConvectiveOutlookOptions = convectiveOptions
		userNotification.AlertOptions = alertOptions

		userNotifications = append(userNotifications, userNotification)
	}

	return userNotifications, nil
}

// Deprecated: use SelectAll
func (p *PostgresUserNotificationTable) GetAll() ([]data_structures.UserNotification, error) {
	return p.SelectAll()
}

func (p *PostgresUserNotificationTable) SelectByUserId(userId string) ([]data_structures.UserNotification, error) {
	query := `SELECT notificationId, userId, zoneCode, countyCode, creationTime, lat, lng, formattedAddress, apnKey, locationName FROM userNotification WHERE userId = $1`

	row, err := p.db.Query(query, userId)
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

		convectiveOptions, err := p.convectiveOutlookOptionsTable.SelectByNotificationId(userNotification.NotificationId)
		if err != nil {
			return nil, err
		}

		alertOptions, err := p.alertOptionsTable.SelectByNotificationId(userNotification.NotificationId)
		if err != nil {
			return nil, err
		}

		userNotification.ConvectiveOutlookOptions = convectiveOptions
		userNotification.AlertOptions = alertOptions

		userNotifications = append(userNotifications, userNotification)
	}

	return userNotifications, nil
}

// Deprecated: use SelectByUserId
func (p *PostgresUserNotificationTable) GetByUserId(userId string) ([]data_structures.UserNotification, error) {
	return p.SelectByUserId(userId)
}

func (p *PostgresUserNotificationTable) SelectByCodes(codes []string) ([]data_structures.UserNotification, error) {
	var userNotifications []data_structures.UserNotification
	for _, code := range codes {
		query := `SELECT notificationId, userId, zoneCode, countyCode, creationTime, lat, lng, formattedAddress, apnKey, locationName FROM userNotification WHERE zoneCode = $1 OR countyCode = $1`

		userNotification := data_structures.UserNotification{}

		row := p.db.QueryRow(query, code)

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

		convectiveOptions, err := p.convectiveOutlookOptionsTable.SelectByNotificationId(userNotification.NotificationId)
		if err != nil {
			return nil, err
		}

		alertOptions, err := p.alertOptionsTable.SelectByNotificationId(userNotification.NotificationId)
		if err != nil {
			return nil, err
		}

		userNotification.ConvectiveOutlookOptions = convectiveOptions
		userNotification.AlertOptions = alertOptions

		userNotifications = append(userNotifications, userNotification)
	}

	return userNotifications, nil
}

// Deprecated: use SelectByCodes
func (p *PostgresUserNotificationTable) GetByCodes(codes []string) ([]data_structures.UserNotification, error) {
	return p.SelectByCodes(codes)
}

func (p *PostgresUserNotificationTable) SelectNotificationsWithConvectiveOutlook() ([]data_structures.UserNotification, error) {
	query := `SELECT notificationId, userId, zoneCode, countyCode, creationTime, lat, lng, formattedAddress, apnKey, locationName FROM userNotification WHERE notificationId = $1`

	rows, err := p.db.Query(query)
	if err != nil {
		return nil, err
	}

	var userNotifications []data_structures.UserNotification
	for rows.Next() {
		userNotification := data_structures.UserNotification{}

		err := rows.Scan(&userNotification.NotificationId,
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

		convectiveOptions, err := p.convectiveOutlookOptionsTable.SelectByNotificationId(userNotification.NotificationId)
		if err != nil {
			return nil, err
		}

		alertOptions, err := p.alertOptionsTable.SelectByNotificationId(userNotification.NotificationId)
		if err != nil {
			return nil, err
		}

		userNotification.ConvectiveOutlookOptions = convectiveOptions
		userNotification.AlertOptions = alertOptions

		userNotifications = append(userNotifications, userNotification)
	}

	return userNotifications, nil
}

// Deprecated: use SelectNotificationsWithConvectiveOutlook
func (p *PostgresUserNotificationTable) GetNotificationsWithConvectiveOutlookOptions() ([]data_structures.UserNotification, error) {
	return p.SelectNotificationsWithConvectiveOutlook()
}

func (p *PostgresUserNotificationTable) Delete(notificationId string) error {
	query := `DELETE FROM userNotification WHERE notificationId = $1`

	exec, err := p.db.Exec(query, notificationId)
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
