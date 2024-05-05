package sql

import (
	"context"
	"database/sql"
	"encoding/json"
	"regexp"
	"time"

	"github.com/cmeyer18/weather-common/v4/data_structures"
	"github.com/cmeyer18/weather-common/v4/generative/golang"
)

var _ IConvectiveOutlookTableV2 = (*PostgresConvectiveOutlookTableV2)(nil)

type IConvectiveOutlookTableV2 interface {
	Insert(outlooks []data_structures.ConvectiveOutlookV2) error

	Select(publishedTime time.Time, outlookType golang.ConvectiveOutlookType) ([]data_structures.ConvectiveOutlookV2, error)

	SelectLatest(outlookType golang.ConvectiveOutlookType) ([]data_structures.ConvectiveOutlookV2, error)
}

type PostgresConvectiveOutlookTableV2 struct {
	db *sql.DB
}

func NewPostgresConvectiveOutlookTableV2(db *sql.DB) PostgresConvectiveOutlookTableV2 {
	return PostgresConvectiveOutlookTableV2{
		db: db,
	}
}

func (p *PostgresConvectiveOutlookTableV2) Insert(outlooks []data_structures.ConvectiveOutlookV2) error {
	ctx := context.Background()
	tx, err := p.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	for _, outlook := range outlooks {
		//language=SQL
		statement, err := p.db.Prepare(`
		INSERT INTO convectiveOutlookV2(outlookType, geometry, dn, issued, expires, valid, label, label2, stroke, fill) 
		VALUES (
			$1, 
			CASE 
				WHEN $2::TEXT IS NULL OR $2::TEXT = '' OR  jsonb_typeof($2::JSONB) = 'null' THEN NULL 
				ELSE ST_GeomFromGeoJSON($2::JSONB) 
    		END,
			$3, $4, $5, $6, $7, $8, $9, $10)`)
		if err != nil {
			return err
		}

		var marshalledGeometryBytes []byte
		if outlook.Geometry != nil {
			marshalledGeometryBytes, err = json.Marshal(outlook.Geometry)
			if err != nil {
				return err
			}
		}

		// Clean up the json, SQL doesn't like these escape characters.
		pattern := regexp.MustCompile(`\\+`)
		unescapedString := pattern.ReplaceAllString(string(marshalledGeometryBytes), "")

		_, err = statement.Exec(string(outlook.OutlookType), unescapedString, outlook.DN, outlook.Issued, outlook.Expires, outlook.Valid, outlook.Label, outlook.Label2, outlook.Stroke, outlook.Fill)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresConvectiveOutlookTableV2) Select(issuedTime time.Time, outlookType golang.ConvectiveOutlookType) ([]data_structures.ConvectiveOutlookV2, error) {
	statement, err := p.db.Prepare(`SELECT outlookType, geometry::JSONB, dn, issued, expires, valid, label, label2, stroke, fill FROM convectiveOutlookV2 WHERE $1 = issued AND $2 = outlookType`)
	if err != nil {
		return nil, err
	}

	rows, err := statement.Query(issuedTime, string(outlookType))
	if err != nil {
		return nil, err
	}

	return p.processConvectiveOutlooks(rows)
}

func (p *PostgresConvectiveOutlookTableV2) SelectLatest(outlookType golang.ConvectiveOutlookType) ([]data_structures.ConvectiveOutlookV2, error) {
	statement, err := p.db.Prepare(`
	SELECT outlookType, geometry::JSONB, dn, issued, expires, valid, label, label2, stroke, fill 
	FROM convectiveOutlookV2 
	WHERE $1 = outlookType AND valid = (
		SELECT MAX(valid)
		FROM convectiveOutlookV2
		WHERE $1 = outlookType
	);`)
	if err != nil {
		return nil, err
	}

	rows, err := statement.Query(string(outlookType))
	if err != nil {
		return nil, err
	}

	return p.processConvectiveOutlooks(rows)
}

func (p *PostgresConvectiveOutlookTableV2) processConvectiveOutlooks(rows *sql.Rows) ([]data_structures.ConvectiveOutlookV2, error) {
	var outlooks []data_structures.ConvectiveOutlookV2
	for rows.Next() {
		outlook := data_structures.ConvectiveOutlookV2{}
		var marshalledGeometry []byte
		var outlookType string

		err := rows.Scan(&outlookType, &marshalledGeometry, &outlook.DN, &outlook.Issued, &outlook.Expires, &outlook.Valid, &outlook.Label, &outlook.Label2, &outlook.Stroke, &outlook.Fill)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(marshalledGeometry, &outlook.Geometry)
		if err != nil {
			return nil, err
		}

		outlook.OutlookType = golang.ConvectiveOutlookType(outlookType)
		outlooks = append(outlooks, outlook)
	}

	return outlooks, nil
}
