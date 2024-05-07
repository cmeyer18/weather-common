package sql

import (
	"database/sql"
	"encoding/json"
	"regexp"

	"github.com/cmeyer18/weather-common/v4/data_structures"
)

var _ IMesoscaleDiscussionV2Table = (*PostgresMesoscaleDiscussionV2Table)(nil)

type IMesoscaleDiscussionV2Table interface {
	Insert(md data_structures.MesoscaleDiscussionV2) error

	Select(mdNumber int, year int) (*data_structures.MesoscaleDiscussionV2, error)

	SelectMDNotInTable(year int, mdsToCheck map[int]bool) ([]int, error)

	Delete(year, mdNumber int) error
}

type PostgresMesoscaleDiscussionV2Table struct {
	db *sql.DB
}

func NewPostgresMesoscaleDiscussionV2Table(db *sql.DB) PostgresMesoscaleDiscussionV2Table {
	return PostgresMesoscaleDiscussionV2Table{
		db: db,
	}
}

func (p *PostgresMesoscaleDiscussionV2Table) Insert(md data_structures.MesoscaleDiscussionV2) error {
	//language=SQL
	statement, err := p.db.Prepare(`
		INSERT INTO mesoscaleDiscussionV2 (number, year, geometry, rawText) 
		VALUES (
		$1, 
		$2, 
		CASE 
			WHEN $3::TEXT IS NULL OR $3::TEXT = '' OR  jsonb_typeof($3::JSONB) = 'null' THEN NULL 
			ELSE ST_GeomFromGeoJSON($3::JSONB) 
		END,
		$4)`)
	if err != nil {
		return err
	}
	defer statement.Close()

	var marshalledGeometryBytes []byte
	if md.Geometry != nil {
		marshalledGeometryBytes, err = json.Marshal(md.Geometry)
		if err != nil {
			return err
		}
	}

	// Clean up the json, SQL doesn't like these escape characters.
	pattern := regexp.MustCompile(`\\+`)
	unescapedString := pattern.ReplaceAllString(string(marshalledGeometryBytes), "")

	_, err = statement.Exec(md.Number, md.Year, unescapedString, md.RawText)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresMesoscaleDiscussionV2Table) Select(year, mdNumber int) (*data_structures.MesoscaleDiscussionV2, error) {
	statement, err := p.db.Prepare(`SELECT mdNumber, year, affectedArea, rawText FROM mesoscaleDiscussionV2 WHERE year = $1 AND mdNumber = $2`)
	if err != nil {
		return nil, err
	}
	defer statement.Close()
	row := statement.QueryRow(year, mdNumber)

	md := data_structures.MesoscaleDiscussionV2{}
	var marshalledGeometry []byte
	err = row.Scan(
		&md.Number,
		&md.Year,
		marshalledGeometry,
		&md.RawText,
	)
	if err != nil {
		return nil, err
	}

	if !(string(marshalledGeometry) == "" || string(marshalledGeometry) == `""` || string(marshalledGeometry) == "null") {
		err = json.Unmarshal(marshalledGeometry, &md.Geometry)
		if err != nil {
			return nil, err
		}
	}


	return &md, nil
}

func (p *PostgresMesoscaleDiscussionV2Table) SelectMDNotInTable(year int, mdsToCheck map[int]bool) ([]int, error) {
	query := `SELECT number FROM mesoscaleDiscussionV2 WHERE year = $1`

	row, err := p.db.Query(query, year)
	if err != nil {
		return nil, err
	}

	mdInTable := make(map[int]bool)
	for row.Next() {
		var md int
		err := row.Scan(&md)
		if err != nil {
			return nil, err
		}
		mdInTable[md] = true
	}

	var mdsNotInTable []int
	for md := range mdsToCheck {
		if !mdInTable[md] {
			mdsNotInTable = append(mdsNotInTable, md)
		}
	}

	return mdsNotInTable, nil
}

func (p *PostgresMesoscaleDiscussionV2Table) Delete(year, mdNumber int) error {
	query := `DELETE FROM mesoscaleDiscussionV2 WHERE year = $1 AND number = $2`

	_, err := p.db.Exec(query, year, mdNumber)
	if err != nil {
		return err
	}

	return nil
}
