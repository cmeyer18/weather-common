package sql

import (
	"database/sql"

	"github.com/cmeyer18/weather-common/v5/data_structures"
)

var _ INotificationQueries = (*PostgresNotificationQueries)(nil)

type INotificationQueries interface {
	GetDevicesForAlertId(alertId string) ([]AlertNotification, error)

	GetDevicesForConvectiveOutlookId(convectiveOutlookId string) ([]ConvectiveOutlookNotification, error)

	GetDevicesForMesoscaleDiscussion(mesoscaleDiscussionId string) ([]MesoscaleDiscussionNotification, error)
}

type PostgresNotificationQueries struct {
	db *sql.DB
}

func NewNotificationQueries(db *sql.DB) PostgresNotificationQueries {
	return PostgresNotificationQueries{
		db: db,
	}
}

type AlertNotification struct {
	Device       data_structures.Device
	LocationName string
}

func (n *PostgresNotificationQueries) GetDevicesForAlertId(alertId string) ([]AlertNotification, error) {
	statement, err := n.db.Prepare(`
		select DISTINCT
	    d.id, 
	    d.userId, 
	    d.apnsToken,
	    un.locationname
		from alertv2 a 
		inner join alertv2_ugccodes au on a.id = au.alertid
		inner join usernotificationalertoption unao on a.event = unao.alertoption 
		inner join usernotification un on un.notificationid = unao.notificationid 
		inner join device d on un.userid = d.userid 
		where a.id = $1 AND 
	  	CASE
    		WHEN (a.geometry IS NOT NULL)
				THEN ST_Contains(a.geometry, ST_SetSRID(ST_MakePoint(un.lng , un.lat), 4326))
    		ELSE
	  	    	au.code = un.zonecode OR au.code = un.countycode
		END
`)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	rows, err := statement.Query(alertId)
	if err != nil {
		return nil, err
	}

	var alertNotifications []AlertNotification
	for rows.Next() {
		var device data_structures.Device
		var locationName string

		err := rows.Scan(&device.DeviceId, &device.UserId, &device.APNSToken, &locationName)
		if err != nil {
			return nil, err
		}

		alertNotifications = append(alertNotifications, AlertNotification{
			Device:       device,
			LocationName: locationName,
		})
	}

	return alertNotifications, nil
}

type ConvectiveOutlookNotification struct {
	Device                 data_structures.Device
	LocationName           string
	ConvectiveOutlookLabel string
}

func (n *PostgresNotificationQueries) GetDevicesForConvectiveOutlookId(convectiveOutlookId string) ([]ConvectiveOutlookNotification, error) {
	statement, err := n.db.Prepare(`
		select DISTINCT
			d.id, 
			d.userId, 
			d.apnsToken,
			un.locationname,
			c.label
		from convectiveoutlookv2 c 
		inner join usernotificationconvectiveoutlookoption co ON c.outlooktype = co.convectiveoutlookoption 
		inner join usernotification un on un.notificationid = co.notificationid
		inner join device d on un.userid = d.userid 
		where c.id = $1 and ST_Contains(c.geometry, ST_SetSRID(ST_MakePoint(un.lng , un.lat), 4326)) 
	`)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	rows, err := statement.Query(convectiveOutlookId)
	if err != nil {
		return nil, err
	}

	var convectiveOutlookNotification []ConvectiveOutlookNotification
	for rows.Next() {
		var device data_structures.Device
		var locationName string
		var convectiveOutlookLabel string

		err := rows.Scan(&device.DeviceId, &device.UserId, &device.APNSToken, &locationName, &convectiveOutlookLabel)
		if err != nil {
			return nil, err
		}

		convectiveOutlookNotification = append(convectiveOutlookNotification, ConvectiveOutlookNotification{
			Device:                 device,
			LocationName:           locationName,
			ConvectiveOutlookLabel: convectiveOutlookLabel,
		})
	}

	return convectiveOutlookNotification, nil
}

type MesoscaleDiscussionNotification struct {
	Device       data_structures.Device
	LocationName string
}

func (n *PostgresNotificationQueries) GetDevicesForMesoscaleDiscussion(mesoscaleDiscussionId string) ([]MesoscaleDiscussionNotification, error) {
	statement, err := n.db.Prepare(`
		select DISTINCT
			d.id, 
			d.userId, 
			d.apnsToken,
			un.locationname
		from mesoscaleDiscussionV2 m  
		join usernotification un on un.mesoscalediscussionnotifications = TRUE
		inner join device d on un.userid = d.userid 
		where m.id = $1 and ST_Contains(m.geometry, ST_SetSRID(ST_MakePoint(un.lng , un.lat), 4326)) 
	`)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	rows, err := statement.Query(mesoscaleDiscussionId)
	if err != nil {
		return nil, err
	}

	var mesoscaleDiscussionNotifications []MesoscaleDiscussionNotification
	for rows.Next() {
		var device data_structures.Device
		var locationName string

		err := rows.Scan(&device.DeviceId, &device.UserId, &device.APNSToken, &locationName)
		if err != nil {
			return nil, err
		}

		mesoscaleDiscussionNotifications = append(mesoscaleDiscussionNotifications, MesoscaleDiscussionNotification{
			Device:       device,
			LocationName: locationName,
		})
	}

	return mesoscaleDiscussionNotifications, nil
}
