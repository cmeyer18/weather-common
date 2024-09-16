package geojson_v2

import (
	"encoding/json"
	"fmt"
)

type Point struct {
	Latitude  float64
	Longitude float64
}

func (p *Point) MarshalJSON() ([]byte, error) {
	return json.Marshal([]float64{p.Longitude, p.Latitude})
}

func (p *Point) UnmarshalJSON(data []byte) error {
	var coordinates []float64
	if err := json.Unmarshal(data, &coordinates); err != nil {
		return err
	}

	return p.decodePoint(coordinates)
}

func (p *Point) encodePoint() []float64 {
	return []float64{p.Longitude, p.Latitude}
}

func (p *Point) decodePoint(point []float64) error {
	if len(point) != 2 {
		return fmt.Errorf("invalid number of coordinates for Point: %v", point)
	}
	p.Longitude = point[0]
	p.Latitude = point[1]

	return nil
}
