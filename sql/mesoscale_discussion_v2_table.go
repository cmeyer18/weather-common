package sql

import (
	"database/sql"
	"encoding/json"

	"github.com/cmeyer18/weather-common/v4/data_structures"
)

var _ IMesoscaleDiscussionV2Table = (*PostgresMesoscaleDiscussionV2Table)(nil)

type IMesoscaleDiscussionV2Table interface {
	Insert(md data_structures.MesoscaleDiscussionV2) error

	Select(mdNumber int, year int) (*data_structures.MesoscaleDiscussionV2, error)

	SelectById(id string) (*data_structures.MesoscaleDiscussionV2, error)

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
		INSERT INTO mesoscaleDiscussionV2 (id, number, year, geometry, rawText) 
		VALUES (
		$1, 
		$2, 
		$3,
		CASE 
			WHEN $4::TEXT IS NULL OR $4::TEXT = '' OR  jsonb_typeof($4::JSONB) = 'null' THEN NULL 
			ELSE ST_GeomFromGeoJSON($4::JSONB) 
		END,
		$5)`)
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

	_, err = statement.Exec(md.ID, md.Number, md.Year, marshalledGeometryBytes, md.RawText)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresMesoscaleDiscussionV2Table) Select(year, mdNumber int) (*data_structures.MesoscaleDiscussionV2, error) {
	statement, err := p.db.Prepare(`SELECT id, number, year, geometry::JSONB, rawText FROM mesoscaleDiscussionV2 WHERE year = $1 AND mdNumber = $2`)
	if err != nil {
		return nil, err
	}
	defer statement.Close()
	row := statement.QueryRow(year, mdNumber)

	md := data_structures.MesoscaleDiscussionV2{}
	var marshalledGeometry []byte
	err = row.Scan(
		&md.ID,
		&md.Number,
		&md.Year,
		&marshalledGeometry,
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

func (p *PostgresMesoscaleDiscussionV2Table) SelectById(id string) (*data_structures.MesoscaleDiscussionV2, error) {
	statement, err := p.db.Prepare(`SELECT id, number, year, geometry::JSONB, rawText FROM mesoscaleDiscussionV2 WHERE id = $1`)
	if err != nil {
		return nil, err
	}
	defer statement.Close()
	row := statement.QueryRow(id)

	md := data_structures.MesoscaleDiscussionV2{}
	var marshalledGeometry []byte
	err = row.Scan(
		&md.ID,
		&md.Number,
		&md.Year,
		&marshalledGeometry,
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
	statement, err := p.db.Prepare(`SELECT number FROM mesoscaleDiscussionV2 WHERE year = $1`)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	rows, err := statement.Query(year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mdInTable := make(map[int]bool)
	for rows.Next() {
		var md int
		err := rows.Scan(&md)
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
	statement, err := p.db.Prepare(`DELETE FROM mesoscaleDiscussionV2 WHERE year = $1 AND number = $2`)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(year, mdNumber)
	if err != nil {
		return err
	}

	return nil
}
