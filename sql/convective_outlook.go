package sql

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/cmeyer18/weather-common/v2/data_structures"
	"github.com/cmeyer18/weather-common/v2/generative/golang"
	"time"
)

var _ PostgresTable[data_structures.ConvectiveOutlook] = (*ConvectiveOutlookTable)(nil)

type ConvectiveOutlookTable struct {
	conn Connection
}

func NewConvectiveOutlookTable(conn Connection) ConvectiveOutlookTable {
	return ConvectiveOutlookTable{
		conn: conn,
	}
}

func (p *ConvectiveOutlookTable) Init() error {
	//language=SQL
	query := `CREATE TABLE IF NOT EXISTS convectiveOutlookTable(
		publishedTime timestamptz,
		outlookType   varchar(225),
		outlook		  jsonb
	)`

	err := p.conn.AddTable(query)
	if err != nil {
		return err
	}

	return nil
}

func (p *ConvectiveOutlookTable) Create(outlook *data_structures.ConvectiveOutlook) error {
	//language=SQL
	query := `INSERT INTO convectiveOutlookTable (publishedTime, outlookType, outlook) VALUES ($1, $2, $3)`
	marshalledOutlook, err := json.Marshal(outlook)
	if err != nil {
		return err
	}

	_, err = p.conn.db.Exec(query, outlook.PublishedTime, outlook.OutlookType, marshalledOutlook)
	if err != nil {
		return err
	}

	return nil
}

func (p *ConvectiveOutlookTable) Find(publishedTime time.Time, outlookType golang.ConvectiveOutlookType) (*data_structures.ConvectiveOutlook, error) {
	query := `SELECT outlook FROM convectiveOutlookTable WHERE publishedTime = $1 AND outlookType = $2`

	row := p.conn.db.QueryRow(query, publishedTime, string(outlookType))

	var rawOutlook []byte
	err := row.Scan(&rawOutlook)
	if err != nil {
		return nil, err
	}

	outlook, err := data_structures.ParseConvectiveOutlook(rawOutlook, outlookType)
	if err != nil {
		return nil, err
	}

	return outlook, nil
}

func (p *ConvectiveOutlookTable) FindLatest(outlookType golang.ConvectiveOutlookType) (*data_structures.ConvectiveOutlook, error) {
	query := `SELECT publishedTime, outlook FROM convectiveOutlookTable WHERE outlookType = $1 ORDER BY publishedTime DESC LIMIT 1`

	row := p.conn.db.QueryRow(query, string(outlookType))

	var publishedTime time.Time
	var rawOutlook []byte
	err := row.Scan(&publishedTime, &rawOutlook)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	outlook, err := data_structures.ParseConvectiveOutlook(rawOutlook, outlookType)
	if err != nil {
		return nil, err
	}

	return outlook, nil
}
