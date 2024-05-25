package sql

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/cmeyer18/weather-common/v5/data_structures"
	"github.com/cmeyer18/weather-common/v5/sql/internal"
	"github.com/cmeyer18/weather-common/v5/sql/internal/common_tables"
)

var _ IUserNotificationTable = (*PostgresUserNotificationTable)(nil)

type IUserNotificationTable interface {
	common_tables.IIdTable[data_structures.UserNotification]

	SelectAll() ([]data_structures.UserNotification, error)

	SelectByUserId(userId string) ([]data_structures.UserNotification, error)

	SelectByCodes(codes []string) ([]data_structures.UserNotification, error)

	SelectNotificationsWithMDNotifications() ([]data_structures.UserNotification, error)

	SelectNotificationsWithConvectiveOutlook() ([]data_structures.UserNotification, error)

	Update(id string, userNotification data_structures.UserNotification) error
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

func (p *PostgresUserNotificationTable) Insert(userNotification data_structures.UserNotification) error {
	//language=SQL
	query := `
	INSERT INTO userNotification (
	  	notificationId, 
	  	userId, 
	  	zoneCode, 
	  	countyCode, 
	  	creationTime, 
	  	lat, 
		lng, 
	  	formattedAddress, 
		apnKey, 
	  	locationName, 
	  	mesoscaleDiscussionNotifications,
	    liveActivities,
	    isDeviceLocation
	) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	ctx := context.Background()
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(
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
		userNotification.MesoscaleDiscussionNotifications,
		userNotification.LiveActivities,
		userNotification.IsDeviceLocation,
	)
	if err != nil {
		return err
	}

	err = p.alertOptionsTable.Insert(tx, userNotification.NotificationId, userNotification.AlertOptions)
	if err != nil {
		return err
	}

	err = p.convectiveOutlookOptionsTable.Insert(tx, userNotification.NotificationId, userNotification.ConvectiveOutlookOptions)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (p *PostgresUserNotificationTable) Select(id string) (*data_structures.UserNotification, error) {
	query := `
	SELECT 
		notificationId, 
		userId, 
		zoneCode,
		countyCode, 
		creationTime, 
		lat, 
		lng, 
		formattedAddress, 
		apnKey, 
		locationName,
		mesoscaleDiscussionNotifications,
		liveActivities,
		isDeviceLocation		
	FROM userNotification 
	WHERE notificationId = $1`

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
		&userNotification.MesoscaleDiscussionNotifications,
		&userNotification.LiveActivities,
		&userNotification.IsDeviceLocation,
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

func (p *PostgresUserNotificationTable) SelectAll() ([]data_structures.UserNotification, error) {
	query := `
	SELECT 
	    notificationId, 
	    userId, 
	    zoneCode,
	    countyCode,
	    creationTime, 
	    lat,
	    lng, 
	    formattedAddress,
	    apnKey,
	    locationName,
		mesoscaleDiscussionNotifications,
		liveActivities,
		isDeviceLocation
	FROM userNotification
`

	rows, err := p.db.Query(query)
	if err != nil {
		return nil, err
	}

	return p.scanRows(rows, false)
}

func (p *PostgresUserNotificationTable) SelectByUserId(userId string) ([]data_structures.UserNotification, error) {
	query := `
	SELECT 
	    notificationId, 
	    userId,
	    zoneCode, 
	    countyCode,
	    creationTime,
	    lat, 
	    lng,
	    formattedAddress,
	    apnKey, 
	    locationName,
		mesoscaleDiscussionNotifications,
		liveActivities,
		isDeviceLocation
	FROM userNotification WHERE userId = $1`

	rows, err := p.db.Query(query, userId)
	if err != nil {
		return nil, err
	}

	return p.scanRows(rows, false)
}

func (p *PostgresUserNotificationTable) SelectByCodes(codes []string) ([]data_structures.UserNotification, error) {
	var userNotifications []data_structures.UserNotification
	for _, code := range codes {
		query := `
		SELECT 
		    notificationId, 
		    userId, 
		    zoneCode, 
		    countyCode, 
		    creationTime,
		    lat, 
		    lng, 
		    formattedAddress,
		    apnKey, 
		    locationName,
			mesoscaleDiscussionNotifications,
			liveActivities,
			isDeviceLocation
		FROM userNotification WHERE zoneCode = $1 OR countyCode = $1`

		rows, err := p.db.Query(query, code)
		if err != nil {
			return nil, err
		}

		userNotificationsToAppend, err := p.scanRows(rows, false)
		if err != nil {
			return nil, err
		}

		userNotifications = append(userNotifications, userNotificationsToAppend...)
	}

	return userNotifications, nil
}

func (p *PostgresUserNotificationTable) SelectNotificationsWithConvectiveOutlook() ([]data_structures.UserNotification, error) {
	query := `
	SELECT 
	    notificationId, 
	    userId, 
	    zoneCode, 
	    countyCode, 
	    creationTime, 
	    lat, 
	    lng, 
	    formattedAddress, 
	    apnKey, 
	    locationName,
		mesoscaleDiscussionNotifications,
		liveActivities,
		isDeviceLocation
	FROM userNotification WHERE notificationId = $1`

	rows, err := p.db.Query(query)
	if err != nil {
		return nil, err
	}

	userNotifications, err := p.scanRows(rows, true)
	if err != nil {
		return nil, err
	}

	return userNotifications, nil
}

// SelectNotificationsWithMDNotifications Selects all of the notifications that want mesoscale discussion notifications.
// Note this does not fill out AlertOptions or SPCOptions in the returned UserNotifications struct
func (p *PostgresUserNotificationTable) SelectNotificationsWithMDNotifications() ([]data_structures.UserNotification, error) {
	query := `
	SELECT 
	    notificationId, 
	    userId, 
	    zoneCode,
	    countyCode, 
	    creationTime, 
	    lat, 
	    lng, 
	    formattedAddress, 
	    apnKey,
	    locationName,
		mesoscaleDiscussionNotifications,
		liveActivities,
		isDeviceLocation
	FROM userNotification WHERE mesoscaleDiscussionNotifications = TRUE`

	rows, err := p.db.Query(query)
	if err != nil {
		return nil, err
	}

	userNotifications, err := p.scanRows(rows, false)
	if err != nil {
		return nil, err
	}

	return userNotifications, nil
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

func (p *PostgresUserNotificationTable) scanRows(rows *sql.Rows, includeOnlyWithConvectiveOptions bool) ([]data_structures.UserNotification, error) {
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
			&userNotification.MesoscaleDiscussionNotifications,
			&userNotification.LiveActivities,
			&userNotification.IsDeviceLocation,
		)
		if err != nil {
			return nil, err
		}

		convectiveOptions, err := p.convectiveOutlookOptionsTable.SelectByNotificationId(userNotification.NotificationId)
		if err != nil {
			return nil, err
		}
		if includeOnlyWithConvectiveOptions && len(convectiveOptions) == 0 {
			continue
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

func (p *PostgresUserNotificationTable) Update(id string, userNotification data_structures.UserNotification) error {
	//language=SQL
	query := `
	UPDATE userNotification SET (
	  	userId, 
	  	zoneCode, 
	  	countyCode, 
	  	creationTime, 
	  	lat, 
		lng, 
	  	formattedAddress, 
		apnKey, 
	  	locationName, 
	  	mesoscaleDiscussionNotifications,
		liveActivities,
		isDeviceLocation
	) =  ($2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	WHERE notificationid = ($1)
`
	ctx := context.Background()
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		query,
		id,
		userNotification.UserID,
		userNotification.ZoneCode,
		userNotification.CountyCode,
		userNotification.CreationTime,
		userNotification.Lat,
		userNotification.Lng,
		userNotification.FormattedAddress,
		userNotification.APNKey,
		userNotification.LocationName,
		userNotification.MesoscaleDiscussionNotifications,
		userNotification.LiveActivities,
		userNotification.IsDeviceLocation,
	)
	if err != nil {
		return err
	}

	err = p.alertOptionsTable.Update(tx, id, userNotification.AlertOptions)
	if err != nil {
		return err
	}

	err = p.convectiveOutlookOptionsTable.Update(tx, id, userNotification.ConvectiveOutlookOptions)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
