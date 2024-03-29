package sql

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/cmeyer18/weather-common/v4/data_structures"
	"github.com/cmeyer18/weather-common/v4/generative/golang"
)

var _ IConvectiveOutlookTable = (*PostgresConvectiveOutlookTable)(nil)

type IConvectiveOutlookTable interface {
	Insert(outlook data_structures.ConvectiveOutlook) error

	Select(publishedTime time.Time, outlookType golang.ConvectiveOutlookType) (*data_structures.ConvectiveOutlook, error)

	SelectLatest(outlookType golang.ConvectiveOutlookType) (*data_structures.ConvectiveOutlook, error)
}

type PostgresConvectiveOutlookTable struct {
	db *sql.DB
}

func NewPostgresConvectiveOutlookTable(db *sql.DB) PostgresConvectiveOutlookTable {
	return PostgresConvectiveOutlookTable{
		db: db,
	}
}

func (p *PostgresConvectiveOutlookTable) Insert(outlook data_structures.ConvectiveOutlook) error {
	//language=SQL
	query := `INSERT INTO convectiveOutlookTable (outlookType, outlook) VALUES ($1, $2)`
	marshalledOutlook, err := json.Marshal(outlook)
	if err != nil {
		return err
	}

	_, err = p.db.Exec(query, outlook.OutlookType, marshalledOutlook)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresConvectiveOutlookTable) Select(publishedTime time.Time, outlookType golang.ConvectiveOutlookType) (*data_structures.ConvectiveOutlook, error) {
	query := `SELECT outlook FROM convectiveOutlookTable WHERE outlookType = $1`

	row := p.db.QueryRow(query, publishedTime, string(outlookType))

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

func (p *PostgresConvectiveOutlookTable) SelectLatest(outlookType golang.ConvectiveOutlookType) (*data_structures.ConvectiveOutlook, error) {
	query := `SELECT outlook FROM convectiveOutlookTable WHERE outlooktype = $1 ORDER BY (outlook->'features'->0->'properties'->>'VALID')::timestamptz DESC LIMIT 1`

	row := p.db.QueryRow(query, string(outlookType))

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
