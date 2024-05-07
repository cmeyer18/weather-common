package sql

import (
	"database/sql"

	"github.com/cmeyer18/weather-common/v4/data_structures"
)

var _ INotificationQueries = (*PostgresNotificationQueries)(nil)

type INotificationQueries interface {
	GetDevicesForAlertId(alertId string) ([]data_structures.Device, map[string][]string, error)
}

type PostgresNotificationQueries struct {
	db *sql.DB
}

func NewNotificationQueries(db *sql.DB) PostgresNotificationQueries {
	return PostgresNotificationQueries{
		db: db,
	}
}

func (n *PostgresNotificationQueries) GetDevicesForAlertId(alertId string) ([]data_structures.Device, map[string][]string, error) {
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
		where a.id = $1 and (au.code = un.zonecode or au.code = un.countycode)
	`)
	if err != nil {
		return nil, nil, err
	}
	defer statement.Close()

	rows, err := statement.Query(alertId)
	if err != nil {
		return nil, nil, err
	}

	var devices []data_structures.Device
	deviceIdToNotificationNames := make(map[string][]string)
	for rows.Next() {
		var device data_structures.Device
		var locationName string

		err := rows.Scan(&device.DeviceId, &device.UserId, &device.APNSToken, &locationName)
		if err != nil {
			return nil, nil, err
		}

		if len(deviceIdToNotificationNames[device.DeviceId]) == 0 {
			devices = append(devices, device)

		}
		deviceIdToNotificationNames[device.DeviceId] = append(deviceIdToNotificationNames[device.DeviceId], locationName)

		devices = append(devices, device)
	}

	return devices, deviceIdToNotificationNames, nil
}

func (n *PostgresNotificationQueries) GetDevicesForConvectiveOutlookId(convectiveOutlookId string) ([]data_structures.Device, map[string][]string, error) {
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
		where a.id = $1 and (au.code = un.zonecode or au.code = un.countycode)
	`)
	if err != nil {
		return nil, nil, err
	}
	defer statement.Close()

	rows, err := statement.Query(convectiveOutlookId)
	if err != nil {
		return nil, nil, err
	}

	var devices []data_structures.Device
	deviceIdToNotificationNames := make(map[string][]string)
	for rows.Next() {
		var device data_structures.Device
		var locationName string

		err := rows.Scan(&device.DeviceId, &device.UserId, &device.APNSToken, &locationName)
		if err != nil {
			return nil, nil, err
		}

		if len(deviceIdToNotificationNames[device.DeviceId]) == 0 {
			devices = append(devices, device)

		}
		deviceIdToNotificationNames[device.DeviceId] = append(deviceIdToNotificationNames[device.DeviceId], locationName)

		devices = append(devices, device)
	}

	return devices, deviceIdToNotificationNames, nil
}

func (n *PostgresNotificationQueries) GetDevicesForMesoscaleDiscussion(mesoscaleDiscussionId string) ([]data_structures.Device, map[string][]string, error) {
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
		where a.id = $1 and (au.code = un.zonecode or au.code = un.countycode)
	`)
	if err != nil {
		return nil, nil, err
	}
	defer statement.Close()

	rows, err := statement.Query(mesoscaleDiscussionId)
	if err != nil {
		return nil, nil, err
	}

	var devices []data_structures.Device
	deviceIdToNotificationNames := make(map[string][]string)
	for rows.Next() {
		var device data_structures.Device
		var locationName string

		err := rows.Scan(&device.DeviceId, &device.UserId, &device.APNSToken, &locationName)
		if err != nil {
			return nil, nil, err
		}

		if len(deviceIdToNotificationNames[device.DeviceId]) == 0 {
			devices = append(devices, device)

		}
		deviceIdToNotificationNames[device.DeviceId] = append(deviceIdToNotificationNames[device.DeviceId], locationName)

		devices = append(devices, device)
	}

	return devices, deviceIdToNotificationNames, nil
}

