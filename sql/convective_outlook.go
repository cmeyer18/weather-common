package sql

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/cmeyer18/weather-common/v3/data_structures"
	"github.com/cmeyer18/weather-common/v3/generative/golang"
	"time"
)

var _ PostgresTable[data_structures.ConvectiveOutlook] = (*PostgresConvectiveOutlookTable)(nil)

type PostgresConvectiveOutlookTable struct {
	conn Connection
}

func NewConvectiveOutlookTable(conn Connection) PostgresConvectiveOutlookTable {
	return PostgresConvectiveOutlookTable{
		conn: conn,
	}
}

func (p *PostgresConvectiveOutlookTable) Init() error {
	//language=SQL
	query := `CREATE TABLE IF NOT EXISTS convectiveOutlookTable(
		outlookType   varchar(225),
		outlook		  jsonb
	)`

	err := p.conn.AddTable(query)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresConvectiveOutlookTable) Create(outlook *data_structures.ConvectiveOutlook) error {
	//language=SQL
	query := `INSERT INTO convectiveOutlookTable (outlookType, outlook) VALUES ($1, $2)`
	marshalledOutlook, err := json.Marshal(outlook)
	if err != nil {
		return err
	}

	_, err = p.conn.db.Exec(query, outlook.OutlookType, marshalledOutlook)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresConvectiveOutlookTable) Find(publishedTime time.Time, outlookType golang.ConvectiveOutlookType) (*data_structures.ConvectiveOutlook, error) {
	query := `SELECT outlook FROM convectiveOutlookTable WHERE outlookType = $1`

	row := p.conn.db.QueryRow(query, publishedTime, string(outlookType))

	var rawOutlook []byte
	err := row.Scan(&rawOutlook)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var outlook data_structures.ConvectiveOutlook
	err = json.Unmarshal(rawOutlook, &outlook)
	if err != nil {
		return nil, err
	}

	return &outlook, nil
}

func (p *PostgresConvectiveOutlookTable) FindLatest(outlookType golang.ConvectiveOutlookType) (*data_structures.ConvectiveOutlook, error) {
	query := `SELECT outlook FROM convectiveOutlookTable WHERE outlooktype = $1 ORDER BY (outlook->'features'->0->'properties'->>'VALID')::timestamptz DESC LIMIT 1`

	row := p.conn.db.QueryRow(query, string(outlookType))

	var rawOutlook []byte
	err := row.Scan(&rawOutlook)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var outlook data_structures.ConvectiveOutlook
	err = json.Unmarshal(rawOutlook, &outlook)
	if err != nil {
		return nil, err
	}

	return &outlook, nil
}
