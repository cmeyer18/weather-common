package data_migrations

import (
	"database/sql"
	"encoding/json"

	"github.com/cmeyer18/weather-common/v4/data_structures"
	"github.com/cmeyer18/weather-common/v4/data_structures/geojson"
	"github.com/cmeyer18/weather-common/v4/data_structures/geojson_v2"
	table "github.com/cmeyer18/weather-common/v4/sql"
)

type Migrator struct {
	DB *sql.DB
}

func (m *Migrator) MigrateAlerts() error {
	alertV2Table := table.NewPostgresAlertV2Table(m.DB)

	statement, err := m.DB.Prepare(`SELECT payload FROM alerts`)
	if err != nil {
		return err
	}
	defer statement.Close()

	rows, err := statement.Query()
	if err != nil {
		return err
	}

	for rows.Next() {
		var alert data_structures.Alert
		var rawAlert []byte
		err := rows.Scan(&rawAlert)
		if err != nil {
			return err
		}

		err = json.Unmarshal(rawAlert, &alert)
		if err != nil {
			return err
		}

		exists, err := alertV2Table.Exists(alert.ID)
		if exists {
			continue
		}

		alertV2 := GetAlertV2(alert)

		err = alertV2Table.Insert(alertV2)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetAlertV2(a data_structures.Alert) data_structures.AlertV2 {
	var geocode *data_structures.AlertPropertiesGeocodeV2
	if a.Properties.Geocode != nil {
		geocode = &data_structures.AlertPropertiesGeocodeV2{}
		if a.Properties.Geocode.UGC != nil {
			geocode.UGC = make([]string, len(a.Properties.Geocode.UGC))
			for i, ugc := range a.Properties.Geocode.UGC {
				geocode.UGC[i] = ugc
			}
		}

		if a.Properties.Geocode.SAME != nil {
			geocode.SAME = make([]string, len(a.Properties.Geocode.SAME))
			for i, same := range a.Properties.Geocode.SAME {
				geocode.SAME[i] = same
			}
		}
	}

	var referenceIds []string
	for _, reference := range a.Properties.References {
		referenceIds = append(referenceIds, reference.AtID)
	}

	parameters := make(map[string]interface{})
	parameters["AWIPSidentifier"] = a.Properties.Parameters.AWIPSIdentifier
	parameters["WMOidentifier"] = a.Properties.Parameters.WMOIdentifier
	parameters["NWSheadline"] = a.Properties.Parameters.AWIPSIdentifier
	parameters["BLOCKCHANNEL"] = a.Properties.Parameters.BlockChannel
	parameters["VTEC"] = a.Properties.Parameters.VTEC
	parameters["expiredReferences"] = a.Properties.Parameters.ExpiredReferences

	return data_structures.AlertV2{
		ID:            a.ID,
		Type:          a.Type,
		Geometry:      GetGeometryV2(a.Geometry),
		AreaDesc:      a.Properties.AreaDesc,
		Geocode:       geocode,
		AffectedZones: a.Properties.AffectedZones,
		References:    referenceIds,
		Sent:          a.Properties.Sent,
		Effective:     a.Properties.Effective,
		Onset:         a.Properties.Onset,
		Expires:       a.Properties.Expires,
		Ends:          a.Properties.Expires,
		Status:        a.Properties.Status,
		MessageType:   a.Properties.MessageType,
		Category:      a.Properties.Category,
		Severity:      a.Properties.Severity,
		Certainty:     a.Properties.Certainty,
		Urgency:       a.Properties.Urgency,
		Event:         a.Properties.Event,
		Sender:        a.Properties.Sender,
		SenderName:    a.Properties.SenderName,
		Headline:      a.Properties.Headline,
		Description:   a.Properties.Description,
		Instruction:   a.Properties.Instruction,
		Response:      a.Properties.Response,
		Parameters:    parameters,
	}
}

func GetGeometryV2(g *geojson.Geometry) *geojson_v2.Geometry {
	if g == nil {
		return nil
	}

	if g.Polygon != nil {
		return &geojson_v2.Geometry{
			Polygon: GetPolygonV2(g.Polygon),
		}
	} else if g.MultiPolygon != nil {
		v2 := &geojson_v2.Geometry{
			MultiPolygon: &geojson_v2.MultiPolygon{
				Polygons: make([]*geojson_v2.Polygon, len(g.MultiPolygon)),
			},
		}

		for i, polygon := range g.MultiPolygon {
			v2.MultiPolygon.Polygons[i] = GetPolygonV2(polygon)
		}

		return v2
	} else {
		return nil
	}
}

func GetPolygonV2(p *geojson.Polygon) *geojson_v2.Polygon {
	if p.OuterPath == nil && p.InnerPaths == nil {
		return nil
	}

	outerpathV2 := GetMultiPointV2(p.OuterPath)

	v2 := geojson_v2.Polygon{
		OuterPath:  &outerpathV2,
		InnerPaths: make([]*geojson_v2.MultiPoint, len(p.InnerPaths)),
	}

	for i, innerPath := range p.InnerPaths {
		innerPathV2 := GetMultiPointV2(innerPath)
		v2.InnerPaths[i] = &innerPathV2
	}

	return &v2
}

func GetMultiPointV2(m *geojson.MultiPoint) geojson_v2.MultiPoint {
	v2 := geojson_v2.MultiPoint{
		Points: make([]*geojson_v2.Point, len(m.Points)),
	}
	for i, point := range m.Points {
		v2.Points[i] = &geojson_v2.Point{
			Longitude: point.Longitude,
			Latitude:  point.Latitude,
		}
	}

	return v2
}
