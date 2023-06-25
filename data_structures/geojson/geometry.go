package geojson

import (
	"fmt"
)

type Geometry struct {
	Type         string     `json:"type" bson:"type"`
	Polygon      *Polygon   `json:"Polygon" bson:"Polygon"`
	MultiPolygon []*Polygon `json:"MultiPolygon" bson:"MultiPolygon"`
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
			return nil, fmt.Errorf("unable to process geometry polygons %v", polygon)
		}

		parsedPolygon, err := parsePolygon(polygon)
		if err != nil {
			return nil, err
		}

		g.Polygon = parsedPolygon
		break
	case "MultiPolygon":
		parsedPolygons, ok := geometry["coordinates"].([]interface{})
		if !ok {
			return nil, fmt.Errorf("unable to process geometry polygons %v", parsedPolygons)
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
