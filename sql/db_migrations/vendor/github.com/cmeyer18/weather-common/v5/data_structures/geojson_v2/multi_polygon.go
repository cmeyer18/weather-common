package geojson_v2

import "encoding/json"

type MultiPolygon struct {
	Polygons []*Polygon `json:"polygons"`
}

func (mpg *MultiPolygon) MarshalJSON() ([]byte, error) {
	coordinates := make([][][][]float64, len(mpg.Polygons))
	for i, polygon := range mpg.Polygons {
		coordinates[i] = polygon.encodePolygon()
	}
	return json.Marshal(coordinates)
}

func (mpg *MultiPolygon) UnmarshalJSON(data []byte) error {
	var coordinates [][][][]float64
	if err := json.Unmarshal(data, &coordinates); err != nil {
		return err
	}
	return mpg.decode(coordinates)
}

func (mpg *MultiPolygon) decode(coordinates [][][][]float64) error {
	polygons := make([]*Polygon, len(coordinates))
	for i, polygonCoords := range coordinates {
		polygon := &Polygon{}
		err := polygon.decodePolygon(polygonCoords)
		if err != nil {
			return err
		}

		polygons[i] = polygon
	}
	mpg.Polygons = polygons
	return nil
}
