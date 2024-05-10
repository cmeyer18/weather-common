package geojson_v2

import (
	"encoding/json"
	"fmt"
)

type Geometry struct {
	Point              *Point
	MultiPoint         *MultiPoint
	Polygon            *Polygon
	MultiPolygon       *MultiPolygon
	GeometryCollection *GeometryCollection
}

func (g *Geometry) MarshalJSON() ([]byte, error) {
	if g == nil {
		return []byte("null"), nil
	}

	if g.Point != nil {
		return json.Marshal(struct {
			Type        string `json:"type"`
			Coordinates *Point `json:"coordinates"`
		}{
			Type:        "Point",
			Coordinates: g.Point,
		})
	} else if g.MultiPoint != nil {
		return json.Marshal(struct {
			Type        string      `json:"type"`
			Coordinates *MultiPoint `json:"coordinates"`
		}{
			Type:        "MultiPoint",
			Coordinates: g.MultiPoint,
		})
	} else if g.Polygon != nil {
		return json.Marshal(struct {
			Type        string   `json:"type"`
			Coordinates *Polygon `json:"coordinates"`
		}{
			Type:        "Polygon",
			Coordinates: g.Polygon,
		})
	} else if g.MultiPolygon != nil {
		return json.Marshal(struct {
			Type        string        `json:"type"`
			Coordinates *MultiPolygon `json:"coordinates"`
		}{
			Type:        "MultiPolygon",
			Coordinates: g.MultiPolygon,
		})
	} else {
		return []byte("null"), nil
	}
}

func (g *Geometry) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == "" || string(data) == `""` {
		return nil
	}

	var raw map[string]json.RawMessage
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	var geometryType string
	err = json.Unmarshal(raw["type"], &geometryType)
	if err != nil {
		return err
	}

	coordinates, _ := raw["coordinates"]

	switch geometryType {
	case "Point":
		err := json.Unmarshal(coordinates, &g.Point)
		if err != nil {
			return err
		}
	case "MultiPoint":
		err := json.Unmarshal(coordinates, &g.MultiPoint)
		if err != nil {
			return err
		}
	case "Polygon":
		err := json.Unmarshal(coordinates, &g.Polygon)
		if err != nil {
			return err
		}
	case "MultiPolygon":
		err := json.Unmarshal(coordinates, &g.MultiPolygon)
		if err != nil {
			return err
		}
	case "GeometryCollection":
		g.GeometryCollection = &GeometryCollection{}
	default:
		return fmt.Errorf("unsupported geometry type: %s", geometryType)
	}

	return nil
}
