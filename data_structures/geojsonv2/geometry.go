package geojson

import (
	"encoding/json"
	"fmt"
)

type GeometryV2 struct {
	Type         string
	Point        *PointV2
	MultiPoint   *MultiPointV2
	Polygon      *PolygonV2
	MultiPolygon *MultiPolygonV2
}

func (g *GeometryV2) MarshalJSON() ([]byte, error) {
	if g == nil || g.Type == "" {
		return []byte("null"), nil
	}

	switch g.Type {
	case "Point":
		return json.Marshal(struct {
			Type        string   `json:"type"`
			Coordinates *PointV2 `json:"coordinates"`
		}{
			Type:        g.Type,
			Coordinates: g.Point,
		})
	case "MultiPoint":
		return json.Marshal(struct {
			Type        string        `json:"type"`
			Coordinates *MultiPointV2 `json:"coordinates"`
		}{
			Type:        g.Type,
			Coordinates: g.MultiPoint,
		})
	case "Polygon":
		return json.Marshal(struct {
			Type        string     `json:"type"`
			Coordinates *PolygonV2 `json:"coordinates"`
		}{
			Type:        g.Type,
			Coordinates: g.Polygon,
		})
	case "MultiPolygon":
		return json.Marshal(struct {
			Type        string          `json:"type"`
			Coordinates *MultiPolygonV2 `json:"coordinates"`
		}{
			Type:        g.Type,
			Coordinates: g.MultiPolygon,
		})
	case "GeometryCollection":
		println("cdm return end")
		return []byte("null"), nil
	default:
		return nil, fmt.Errorf("unsupported geometry type: %s", g.Type)
	}
}

// UnmarshalJSON implements custom JSON unmarshaling for the GeometryV2 struct.
func (g *GeometryV2) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == "" {
		return nil
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if err := json.Unmarshal(raw["type"], &g.Type); err != nil {
		return err
	}

	switch g.Type {
	case "Point":
		if val, ok := raw["coordinates"]; ok {
			g.Point = &PointV2{}
			if err := json.Unmarshal(val, g.Point); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("missing coordinates for Point geometry")
		}
	case "MultiPoint":
		if val, ok := raw["coordinates"]; ok {
			g.MultiPoint = &MultiPointV2{}
			if err := json.Unmarshal(val, g.MultiPoint); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("missing coordinates for MultiPoint geometry")
		}
	case "Polygon":
		if val, ok := raw["coordinates"]; ok {
			g.Polygon = &PolygonV2{}
			if err := json.Unmarshal(val, g.Polygon); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("missing coordinates for Polygon geometry")
		}
	case "MultiPolygon":
		if val, ok := raw["coordinates"]; ok {
			g.MultiPolygon = &MultiPolygonV2{}
			if err := json.Unmarshal(val, g.MultiPolygon); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("missing coordinates for MultiPolygon geometry")
		}
	case "GeometryCollection":
	default:
		return fmt.Errorf("unsupported geometry type: %s", g.Type)
	}

	return nil
}
