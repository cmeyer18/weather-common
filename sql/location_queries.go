package sql

import (
	"database/sql"

	"github.com/cmeyer18/weather-common/v5/data_structures"
)

var _ ILocationQueries = (*PostgresLocationQueries)(nil)

type ILocationQueries interface {
	GetDevicesForAlertID(alertID string) (map[data_structures.Device][]string, error)

	GetDevicesForConvectiveOutlookID(convectiveOutlookID string) (map[string]map[data_structures.Device][]string, error)

	GetDevicesForMesoscaleDiscussionID(mesoscaleDiscussionID string) (map[data_structures.Device][]string, error)
}

type PostgresLocationQueries struct {
	db *sql.DB
}

func NewLocationQueries(db *sql.DB) PostgresLocationQueries {
	return PostgresLocationQueries{
		db: db,
	}
}

// GetDevicesForAlertID returns a mapping from device to list of LocationNames
func (n *PostgresLocationQueries) GetDevicesForAlertID(alertId string) (map[data_structures.Device][]string, error) {
	statement, err := n.db.Prepare(`
		SELECT DISTINCT
			device.id, 
			device.userId, 
			device.apnsToken,
			location.locationname
		FROM 
		    alertv2 a 
		INNER JOIN 
			alertv2_ugccodes au 
			ON a.id = au.alertid
		INNER JOIN 
			locationoptions 
			ON locationoptions.optiontype = 0 AND a.event = locationoptions.option
		INNER JOIN 
			location 
			ON location.locationID = locationOptions.locationID 
		INNER JOIN 
			device 
		    ON (
		        location.locationType = 0 AND location.locationReferenceID = device.id 
			) OR (
			    location.locationType = 1 AND location.locationReferenceID = device.userid 
			)
		WHERE a.id = $1 
			AND CASE
				WHEN (a.geometry IS NOT NULL)
					THEN ST_Contains(a.geometry, ST_SetSRID(ST_MakePoint(location.lng , location.lat), 4326))
				ELSE
					au.code = location.zonecode OR au.code = location.countycode
			END`)

	if err != nil {
		return nil, err
	}
	defer statement.Close()

	rows, err := statement.Query(alertId)
	if err != nil {
		return nil, err
	}

	deviceToLocationNames := make(map[data_structures.Device][]string)
	for rows.Next() {
		var device data_structures.Device
		var locationName string

		err := rows.Scan(&device.DeviceId, &device.UserId, &device.APNSToken, &locationName)
		if err != nil {
			return nil, err
		}

		deviceToLocationNames[device] = append(deviceToLocationNames[device], locationName)
	}

	return deviceToLocationNames, nil
}

// GetDevicesForConvectiveOutlookID returns a mapping from level to device to list of LocationNames
func (n *PostgresLocationQueries) GetDevicesForConvectiveOutlookID(convectiveOutlookId string) (map[string]map[data_structures.Device][]string, error) {
	statement, err := n.db.Prepare(`
		SELECT DISTINCT
			device.id, 
			device.userId, 
			device.apnsToken,
			location.locationname,
			convectiveoutlookv2.label
		FROM convectiveoutlookv2 
		INNER JOIN 
			locationoptions 
			ON locationoptions.optiontype = 1 AND convectiveoutlookv2.outlooktype = locationoptions.option
		INNER JOIN 
			location 
			ON location.locationID = locationOptions.locationID 
		INNER JOIN 
			device 
		    ON (
		        location.locationType = 0 AND location.locationReferenceID = device.id 
			) OR (
			    location.locationType = 1 AND location.locationReferenceID = device.userid 
			)
		WHERE
		    convectiveoutlookv2.id = $1
		  	AND ST_Contains(convectiveoutlookv2.geometry, ST_SetSRID(ST_MakePoint(location.lng , location.lat), 4326)) 
	`)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	rows, err := statement.Query(convectiveOutlookId)
	if err != nil {
		return nil, err
	}

	levelToDeviceToLocationNames := make(map[string]map[data_structures.Device][]string)

	for rows.Next() {
		var device data_structures.Device
		var locationName string
		var convectiveOutlookLabel string

		err := rows.Scan(&device.DeviceId, &device.UserId, &device.APNSToken, &locationName, &convectiveOutlookLabel)
		if err != nil {
			return nil, err
		}

		if levelToDeviceToLocationNames[convectiveOutlookLabel] == nil {
			levelToDeviceToLocationNames[convectiveOutlookLabel] = make(map[data_structures.Device][]string)
		}

		levelToDeviceToLocationNames[convectiveOutlookLabel][device] = append(
			levelToDeviceToLocationNames[convectiveOutlookLabel][device],
			locationName,
		)
	}

	return levelToDeviceToLocationNames, nil
}

// GetDevicesForMesoscaleDiscussionID returns a mapping from device to list of LocationNames
func (n *PostgresLocationQueries) GetDevicesForMesoscaleDiscussionID(mesoscaleDiscussionID string) (map[data_structures.Device][]string, error) {
	statement, err := n.db.Prepare(`
		SELECT DISTINCT
			device.id, 
			device.userId, 
			device.apnsToken,
			location.locationname
		FROM mesoscaleDiscussionV2 m  
		INNER JOIN 
			locationoptions 
			ON locationoptions.optiontype = 2 AND locationoptions.option = 'true'
		INNER JOIN 
			location 
			ON location.locationID = locationOptions.locationID 
		INNER JOIN 
			device 
		    ON (
		        location.locationType = 0 AND location.locationReferenceID = device.id 
			) OR (
			    location.locationType = 1 AND location.locationReferenceID = device.userid 
			)
		WHERE m.id = $1 AND ST_Contains(m.geometry, ST_SetSRID(ST_MakePoint(location.lng , location.lat), 4326)) 
	`)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	rows, err := statement.Query(mesoscaleDiscussionID)
	if err != nil {
		return nil, err
	}

	deviceToLocationNames := make(map[data_structures.Device][]string)
	for rows.Next() {
		var device data_structures.Device
		var locationName string

		err := rows.Scan(&device.DeviceId, &device.UserId, &device.APNSToken, &locationName)
		if err != nil {
			return nil, err
		}

		deviceToLocationNames[device] = append(deviceToLocationNames[device], locationName)
	}

	return deviceToLocationNames, nil
}
