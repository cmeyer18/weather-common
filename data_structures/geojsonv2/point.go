package geojson

import (
	"encoding/json"
	"fmt"
)

type PointV2 struct {
	Latitude  float64
	Longitude float64
}

func (p *PointV2) MarshalJSON() ([]byte, error) {
	return json.Marshal([]float64{p.Longitude, p.Latitude})
}

func (p *PointV2) UnmarshalJSON(data []byte) error {
	var coordinates []float64
	if err := json.Unmarshal(data, &coordinates); err != nil {
		return err
	}

	return p.decodePoint(coordinates)
}

func (p *PointV2) encodePoint() []float64 {
	return []float64{p.Longitude, p.Latitude}
}

func (p *PointV2) decodePoint(point []float64) error {
	if len(point) != 2 {
		return fmt.Errorf("invalid number of coordinates for Point: %v", point)
	}
	p.Longitude = point[0]
	p.Latitude = point[1]

	return nil
}
