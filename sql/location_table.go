package sql

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/cmeyer18/weather-common/v6/data_structures"
	"github.com/cmeyer18/weather-common/v6/generative/golang"
	"github.com/cmeyer18/weather-common/v6/sql/internal/common_tables"
)

var _ ILocationTable = (*PostgresLocationTable)(nil)

type ILocationTable interface {
	common_tables.IIdTable[data_structures.Location]

	SelectByUserID(userID string) ([]data_structures.Location, error)

	SelectByDeviceID(deviceID string) ([]data_structures.Location, error)

	SelectByCodes(codes []string) ([]data_structures.Location, error)

	SelectNotificationsWithMDNotifications() ([]data_structures.Location, error)

	SelectNotificationsWithConvectiveOutlook() ([]data_structures.Location, error)

	Update(location data_structures.Location) error
}

type PostgresLocationTable struct {
	db *sql.DB
}

func NewPostgresLocationTable(db *sql.DB) PostgresLocationTable {
	return PostgresLocationTable{
		db: db,
	}
}

type LocationOptionType int8

const (
	LocationOptionType_AlertOption                      LocationOptionType = 0
	LocationOptionType_ConvectiveOutlookOption          LocationOptionType = 1
	LocationOptionType_MesoscaleDiscussionNotifications LocationOptionType = 2
)

func (p *PostgresLocationTable) Insert(location data_structures.Location) error {
	ctx := context.Background()
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = p.insert(tx, location)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (p *PostgresLocationTable) insert(transaction *sql.Tx, location data_structures.Location) error {
	//language=SQL
	locationOptionQuery := `
	INSERT INTO locationOptions (
		locationID,
		optionType,
		option
	) 
	VALUES ($1, $2, $3)`

	//language=SQL
	locationQuery := `
	INSERT INTO location (
		locationID,
		locationType,
		locationReferenceID,
		zoneCode,
		countyCode,
		latitude,
		longitude,
		locationName
	) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := transaction.Exec(
		locationQuery,
		location.LocationID,
		location.LocationType,
		location.LocationReferenceID,
		location.ZoneCode,
		location.CountyCode,
		location.Latitude,
		location.Longitude,
		location.LocationName,
	)
	if err != nil {
		return err
	}

	for _, alertOption := range location.AlertOptions {
		_, err = transaction.Exec(
			locationOptionQuery,
			location.LocationID,
			int8(LocationOptionType_AlertOption),
			string(alertOption),
		)
		if err != nil {
			return err
		}
	}

	for _, convectiveOutlookOption := range location.ConvectiveOutlookOptions {
		_, err = transaction.Exec(
			locationOptionQuery,
			location.LocationID,
			int8(LocationOptionType_ConvectiveOutlookOption),
			string(convectiveOutlookOption),
		)
		if err != nil {
			return err
		}
	}

	mesoscaleOptionToString := "true"
	if !location.MesoscaleDiscussionNotifications {
		mesoscaleOptionToString = "false"
	}

	_, err = transaction.Exec(
		locationOptionQuery,
		location.LocationID,
		int8(LocationOptionType_MesoscaleDiscussionNotifications),
		mesoscaleOptionToString,
	)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresLocationTable) Select(locationID string) (*data_structures.Location, error) {
	query := `
	SELECT 
		location.locationID,
		location.locationType,
		location.locationReferenceID,
		location.zoneCode,
		location.countyCode,
		location.latitude,
		location.lng,
		location.locationName,
		location.created,
		locationOptions.option,
		locationOptions.optionType
	FROM location 
	JOIN locationOptions ON location.locationID = locationOptions.locationID
	WHERE location.locationID = $1`

	rows, err := p.db.Query(query, locationID)
	if err != nil {
		return nil, err
	}

	locations, err := p.scanRows(rows)
	if err != nil {
		return nil, err
	}

	if len(locations) == 0 {
		return nil, nil
	}

	return &locations[0], nil
}

func (p *PostgresLocationTable) SelectByUserID(userID string) ([]data_structures.Location, error) {
	query := `
	SELECT 
		location.locationID,
		location.locationType,
		location.locationReferenceID,
		location.zoneCode,
		location.countyCode,
		location.latitude,
		location.longitude,
		location.locationName,
		location.created,
		locationOptions.option,
		locationOptions.optionType
	FROM location
	JOIN locationOptions ON location.locationID = locationOptions.locationID
	WHERE locationReferenceID = $1 AND locationType = 1`

	rows, err := p.db.Query(query, userID)
	if err != nil {
		return nil, err
	}

	return p.scanRows(rows)
}

func (p *PostgresLocationTable) SelectByDeviceID(deviceID string) ([]data_structures.Location, error) {
	query := `
	SELECT 
		location.locationID,
		location.locationType,
		location.locationReferenceID,
		location.zoneCode,
		location.countyCode,
		location.latitude,
		location.longitude,
		location.locationName,
		location.created,
		locationOptions.option,
		locationOptions.optionType
	FROM location
	JOIN locationOptions ON location.locationID = locationOptions.locationID
	WHERE locationReferenceID = $1 AND locationType = 2`

	rows, err := p.db.Query(query, deviceID)
	if err != nil {
		return nil, err
	}

	return p.scanRows(rows)
}

func (p *PostgresLocationTable) SelectByCodes(codes []string) ([]data_structures.Location, error) {
	var userNotifications []data_structures.Location
	for _, code := range codes {
		query := `
		SELECT 
			location.locationID,
			location.locationType,
			location.locationReferenceID,
			location.zoneCode,
			location.countyCode,
			location.latitude,
			location.longitude,
			location.locationName,
			location.created,
			locationOptions.option,
			locationOptions.optionType
		FROM location
		JOIN locationOptions 
		    ON location.locationID = locationOptions.locationID
		WHERE zoneCode = $1 OR countyCode = $1`

		rows, err := p.db.Query(query, code)
		if err != nil {
			return nil, err
		}

		userNotificationsToAppend, err := p.scanRows(rows)
		if err != nil {
			return nil, err
		}

		userNotifications = append(userNotifications, userNotificationsToAppend...)
	}

	return userNotifications, nil
}

func (p *PostgresLocationTable) SelectNotificationsWithConvectiveOutlook() ([]data_structures.Location, error) {
	query := `
	SELECT 
		location.locationID,
		location.locationType,
		location.locationReferenceID,
		location.zoneCode,
		location.countyCode,
		location.latitude,
		location.longitude,
		location.locationName,
		location.created,
		locationOptions.option,
		locationOptions.optionType
	FROM (
	    SELECT 
	        locationID
		FROM locationoptions AS innerLocationOptions
		WHERE innerLocationOptions.optionType = 1
	) AS mesoscaleLocations
	JOIN location 
	    ON mesoscaleLocations.locationID = location.locationID
	JOIN locationOptions 
	    ON mesoscaleLocations.locationID = locationOptions.locationID`

	rows, err := p.db.Query(query)
	if err != nil {
		return nil, err
	}

	userNotifications, err := p.scanRows(rows)
	if err != nil {
		return nil, err
	}

	return userNotifications, nil
}

// SelectNotificationsWithMDNotifications Selects all of the notifications that want mesoscale discussion notifications.
// Note this does not fill out AlertOptions or SPCOptions in the returned UserNotifications struct
func (p *PostgresLocationTable) SelectNotificationsWithMDNotifications() ([]data_structures.Location, error) {
	query := `
	SELECT 
		location.locationID,
		location.locationType,
		location.locationReferenceID,
		location.zoneCode,
		location.countyCode,
		location.latitude,
		location.longitude,
		location.locationName,
		location.created,
		locationOptions.option,
		locationOptions.optionType
	FROM (
	    SELECT 
	        locationID
		FROM locationoptions as innerLocationOptions
		WHERE innerLocationOptions.optionType = 2 && innerLocationOptions.option = 'true'
	) as mesoscaleLocations
	JOIN location on mesoscaleLocations.locationID = location.locationID
	JOIN locationOptions ON mesoscaleLocations.locationID = locationOptions.locationID`

	rows, err := p.db.Query(query)
	if err != nil {
		return nil, err
	}

	userNotifications, err := p.scanRows(rows)
	if err != nil {
		return nil, err
	}

	return userNotifications, nil
}

func (p *PostgresLocationTable) Delete(locationID string) error {
	ctx := context.Background()
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = p.delete(tx, locationID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (p *PostgresLocationTable) delete(transaction *sql.Tx, locationID string) error {
	query := `DELETE FROM location WHERE locationID = $1`
	optionsQuery := `DELETE FROM locationoptions WHERE locationID = $1`

	_, err := transaction.Exec(query, locationID)
	if err != nil {
		return err
	}

	_, err = transaction.Exec(optionsQuery, locationID)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresLocationTable) Update(location data_structures.Location) error {
	ctx := context.Background()
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = p.delete(tx, location.LocationID)
	if err != nil {
		return err
	}

	err = p.insert(tx, location)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (p *PostgresLocationTable) scanRows(rows *sql.Rows) ([]data_structures.Location, error) {
	locations := make(map[string]*data_structures.Location)
	for rows.Next() {
		var locationID string
		var locationType data_structures.LocationType
		var locationReferenceID string
		var zoneCode string
		var countyCode string
		var latitude float64
		var longitude float64
		var locationName string
		var created time.Time
		var option string
		var optionType LocationOptionType

		err := rows.Scan(
			&locationID,
			&locationType,
			&locationReferenceID,
			&zoneCode,
			&countyCode,
			&latitude,
			&longitude,
			&locationName,
			&created,
			&option,
			&optionType,
		)
		if err != nil {
			return nil, err
		}

		if locations[locationID] == nil {
			locations[locationID] = &data_structures.Location{
				LocationID:          locationID,
				LocationType:        locationType,
				LocationReferenceID: locationReferenceID,
				ZoneCode:            zoneCode,
				CountyCode:          countyCode,
				Latitude:            latitude,
				Longitude:           longitude,
				LocationName:        locationName,
				Created:             created,
			}
		}

		switch optionType {
		case LocationOptionType_AlertOption:
			locations[locationID].AlertOptions = append(locations[locationID].AlertOptions, golang.AlertType(option))
		case LocationOptionType_ConvectiveOutlookOption:
			locations[locationID].ConvectiveOutlookOptions = append(locations[locationID].ConvectiveOutlookOptions, golang.ConvectiveOutlookType(option))
		case LocationOptionType_MesoscaleDiscussionNotifications:
			boolOption, err := strconv.ParseBool(option)
			if err != nil {
				return nil, err
			}
			locations[locationID].MesoscaleDiscussionNotifications = boolOption
		}
	}

	var locationsArray []data_structures.Location
	for _, location := range locations {
		locationsArray = append(locationsArray, *location)
	}

	return locationsArray, nil
}
