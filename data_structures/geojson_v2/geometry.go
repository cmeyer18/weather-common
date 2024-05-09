package geojson_v2

import (
	"encoding/json"
	"fmt"
)

type Geometry struct {
	Type         string
	Point        *Point
	MultiPoint   *MultiPoint
	Polygon      *Polygon
	MultiPolygon *MultiPolygon
}

func (g *Geometry) MarshalJSON() ([]byte, error) {
	if g == nil || g.Type == "" {
		return []byte("null"), nil
	}

	switch g.Type {
	case "Point":
		return json.Marshal(struct {
			Type        string `json:"type"`
			Coordinates *Point `json:"coordinates"`
		}{
			Type:        g.Type,
			Coordinates: g.Point,
		})
	case "MultiPoint":
		return json.Marshal(struct {
			Type        string      `json:"type"`
			Coordinates *MultiPoint `json:"coordinates"`
		}{
			Type:        g.Type,
			Coordinates: g.MultiPoint,
		})
	case "Polygon":
		return json.Marshal(struct {
			Type        string   `json:"type"`
			Coordinates *Polygon `json:"coordinates"`
		}{
			Type:        g.Type,
			Coordinates: g.Polygon,
		})
	case "MultiPolygon":
		return json.Marshal(struct {
			Type        string        `json:"type"`
			Coordinates *MultiPolygon `json:"coordinates"`
		}{
			Type:        g.Type,
			Coordinates: g.MultiPolygon,
		})
	case "GeometryCollection":
		return []byte("null"), nil
	default:
		return nil, fmt.Errorf("unsupported geometry type: %s", g.Type)
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

	err = json.Unmarshal(raw["type"], &g.Type)
	if err != nil {
		return err
	}

	coordinates, ok := raw["coordinates"]
	if !ok {
		return fmt.Errorf("missing coordinates for Point geometry")
	}

	switch g.Type {
	case "Point":
		err := json.Unmarshal(coordinates, g.Point)
		if err != nil {
			return err
		}
	case "MultiPoint":
		err := json.Unmarshal(coordinates, g.MultiPoint)
		if err != nil {
			return err
		}
	case "Polygon":
		err := json.Unmarshal(coordinates, g.Polygon)
		if err != nil {
			return err
		}
	case "MultiPolygon":
		err := json.Unmarshal(coordinates, g.MultiPolygon)
		if err != nil {
			return err
		}
	case "GeometryCollection":
	default:
		return fmt.Errorf("unsupported geometry type: %s", g.Type)
	}

	return nil
}
