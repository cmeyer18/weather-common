package geojson

import "encoding/json"

// MultiPolygonV2 represents a GeoJSON MultiPolygon geometry.
type MultiPolygonV2 struct {
	Polygons []*PolygonV2 `json:"polygons"`
}

func (mpg *MultiPolygonV2) MarshalJSON() ([]byte, error) {
	coordinates := make([][][][]float64, len(mpg.Polygons))
	for i, polygon := range mpg.Polygons {
		coordinates[i] = polygon.encodePolygon()
	}
	return json.Marshal(coordinates)
}

func (mpg *MultiPolygonV2) UnmarshalJSON(data []byte) error {
	var coordinates [][][][]float64
	if err := json.Unmarshal(data, &coordinates); err != nil {
		return err
	}
	return mpg.decode(coordinates)
}

func (mpg *MultiPolygonV2) decode(coordinates [][][][]float64) error {
	polygons := make([]*PolygonV2, len(coordinates))
	for i, polygonCoords := range coordinates {
		polygon := &PolygonV2{}
		err := polygon.decodePolygon(polygonCoords)
		if err != nil {
			return err
		}

		polygons[i] = polygon
	}
	mpg.Polygons = polygons
	return nil
}
