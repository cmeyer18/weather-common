package sql

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/cmeyer18/weather-common/v4/data_structures"
	"github.com/cmeyer18/weather-common/v4/data_structures/geojson"
)

var _ IPostgresMesoscaleDiscussionTable = (*PostgresMesoscaleDiscussionTable)(nil)

type IPostgresMesoscaleDiscussionTable interface {
	Insert(md data_structures.MesoscaleDiscussion) error

	Select(mdNumber int, year int) (*data_structures.MesoscaleDiscussion, error)

	SelectMDNotInTable(year int, mdsToCheck map[int]bool) ([]int, error)
}

type PostgresMesoscaleDiscussionTable struct {
	db *sql.DB
}

func NewPostgresMesoscaleDicussionTable(db *sql.DB) PostgresMesoscaleDiscussionTable {
	return PostgresMesoscaleDiscussionTable{
		db: db,
	}
}

func (p *PostgresMesoscaleDiscussionTable) Insert(md data_structures.MesoscaleDiscussion) error {
	//language=SQL
	query := `INSERT INTO mesoscaleDiscussion (mdNumber, year, affectedArea, rawText) VALUES ($1, $2, $3, $4)`

	affectedArea, err := json.Marshal(md.AffectedArea)
	if err != nil {
		return err
	}

	_, err = p.db.Exec(query, md.MDNumber, md.Year, affectedArea, md.RawText)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresMesoscaleDiscussionTable) Select(year, mdNumber int) (*data_structures.MesoscaleDiscussion, error) {
	query := `SELECT mdNumber, year, affectedArea, rawText FROM mesoscaleDiscussion WHERE year = $1 AND mdNumber = $2`

	row := p.db.QueryRow(query, year, mdNumber)

	md := data_structures.MesoscaleDiscussion{}

	var rawAffectedArea []byte
	err := row.Scan(
		&md.MDNumber,
		&md.Year,
		rawAffectedArea,
		&md.RawText,
	)
	if err != nil {
		return nil, err
	}

	var polygon *geojson.Polygon
	err = json.Unmarshal(rawAffectedArea, &polygon)
	if err != nil {
		return nil, err
	}

	md.AffectedArea = polygon

	return &md, nil
}

func (p *PostgresMesoscaleDiscussionTable) SelectMDNotInTable(year int, mdsToCheck map[int]bool) ([]int, error) {
	query := `SELECT mdNumber FROM mesoscaleDiscussion WHERE year = $1`

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

func (p *PostgresMesoscaleDiscussionTable) Delete(year, mdNumber int) error {
	query := `DELETE FROM mesoscaleDiscussion WHERE year = $1 AND mdNumber = $2`

	exec, err := p.db.Exec(query, year, mdNumber)
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
