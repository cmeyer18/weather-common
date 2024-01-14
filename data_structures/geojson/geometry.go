package geojson

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Geometry struct {
	Type         string     `json:"type"`
	Polygon      *Polygon   `json:"Polygon"`
	MultiPolygon []*Polygon `json:"MultiPolygon"`
}

func (g *Geometry) Value() (driver.Value, error) {
	if g.Type == "" {
		return "", nil
	}

	return json.Marshal(g)
}

func (g *Geometry) Scan(value interface{}) error {
	byteData := value.([]byte)

	var unmarshalledData map[string]interface{}
	err := json.Unmarshal(byteData, &unmarshalledData)
	if err != nil {
		return err
	}

	geometry, err := ParseGeometry(unmarshalledData)
	if err != nil {
		return err
	}

	g.Type = geometry.Type
	g.Polygon = geometry.Polygon
	g.MultiPolygon = geometry.MultiPolygon
	return nil
}

func ParseGeometry(geometry map[string]interface{}) (*Geometry, error) {
	if len(geometry) == 0 {
		return nil, nil
	}

	parsedType, ok := geometry["type"].(string)
	if !ok {
		return nil, fmt.Errorf("unable to parse type on geometry %v", geometry)
	}

	g := Geometry{}
	g.Type = parsedType

	switch parsedType {
	case "Polygon":
		polygon, ok := geometry["coordinates"]
		if !ok {
			return nil, fmt.Errorf("unable to process Polygon %v", polygon)
		}

		parsedPolygon, err := parsePolygon(polygon)
		if err != nil {
			return nil, err
		}

		g.Polygon = parsedPolygon
		break
	case "MultiPolygon":
		// TODO check if empty
		parsedPolygons, ok := geometry["coordinates"].([]interface{})
		if !ok {
			return nil, fmt.Errorf("unable to process MultiPolygon %v", parsedPolygons)
		}

		var polygons []*Polygon
		for _, polygon := range parsedPolygons {
			parsedPolygon, err := parsePolygon(polygon)
			if err != nil {
				return nil, err
			}

			polygons = append(polygons, parsedPolygon)
		}

		g.MultiPolygon = polygons
		break
	case "GeometryCollection":
		break
	default:
		return nil, fmt.Errorf("unknown type %v", parsedType)
	}

	return &g, nil
}
