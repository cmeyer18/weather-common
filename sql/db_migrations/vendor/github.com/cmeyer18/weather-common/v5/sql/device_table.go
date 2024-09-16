package sql

import (
	"database/sql"
	"errors"
	"strconv"

	"github.com/cmeyer18/weather-common/v5/data_structures"
	"github.com/cmeyer18/weather-common/v5/sql/internal/common_tables"
)

var _ IDeviceTable = (*PostgresDeviceTable)(nil)

type IDeviceTable interface {
	common_tables.IIdTable[data_structures.Device]

	SelectByUser(userId string) ([]data_structures.Device, error)

	UpdateApnsToken(id, apnsToken string) error
}

type PostgresDeviceTable struct {
	db *sql.DB
}

func NewPostgresDeviceTable(db *sql.DB) PostgresDeviceTable {
	return PostgresDeviceTable{
		db: db,
	}
}

func (p PostgresDeviceTable) Insert(device data_structures.Device) error {
	//language=SQL
	query := `INSERT INTO device (id, userId, apnsToken) VALUES ($1, $2, $3)`

	_, err := p.db.Exec(query, device.DeviceId, device.UserId, device.APNSToken)
	if err != nil {
		return err
	}

	return nil
}

func (p PostgresDeviceTable) Select(id string) (*data_structures.Device, error) {
	query := `SELECT id, userId, apnsToken FROM device WHERE id = $1`

	row := p.db.QueryRow(query, id)

	device := data_structures.Device{}
	err := row.Scan(
		&device.DeviceId,
		&device.UserId,
		&device.APNSToken,
	)
	if err != nil {
		return nil, err
	}

	return &device, nil
}

func (p PostgresDeviceTable) SelectByUser(userId string) ([]data_structures.Device, error) {
	query := `SELECT id, userId, apnsToken FROM device WHERE userId = $1`

	rows, err := p.db.Query(query, userId)
	if err != nil {
		return nil, err
	}

	var devices []data_structures.Device
	for rows.Next() {
		device := data_structures.Device{}
		err := rows.Scan(
			&device.DeviceId,
			&device.UserId,
			&device.APNSToken,
		)
		if err != nil {
			return nil, err
		}

		devices = append(devices, device)
	}

	return devices, nil
}

func (p PostgresDeviceTable) Delete(id string) error {
	query := `DELETE FROM device WHERE id = $1`

	exec, err := p.db.Exec(query, id)
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

func (p PostgresDeviceTable) UpdateApnsToken(id, apnsToken string) error {
	//language=SQL
	query := `UPDATE device SET apnsToken = $2 WHERE id = ($1)`

	_, err := p.db.Exec(query, id, apnsToken)
	if err != nil {
		return err
	}

	return nil
}
