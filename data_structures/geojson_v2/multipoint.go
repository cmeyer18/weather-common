package geojson_v2

import (
	"encoding/json"
)

type MultiPoint struct {
	Points []*Point
}

func (mp *MultiPoint) MarshalJSON() ([]byte, error) {
	return json.Marshal(mp.encodeMultiPoint())
}

func (mp *MultiPoint) UnmarshalJSON(data []byte) error {
	var coordinates [][]float64
	if err := json.Unmarshal(data, &coordinates); err != nil {
		return err
	}

	return mp.decodeMultiPoint(coordinates)
}

func (mp *MultiPoint) encodeMultiPoint() [][]float64 {
	coordinates := make([][]float64, len(mp.Points))
	for i, point := range mp.Points {
		coordinates[i] = point.encodePoint()
	}
	return coordinates
}

func (mp *MultiPoint) decodeMultiPoint(multiPoint [][]float64) error {
	points := make([]*Point, len(multiPoint))
	for i, point := range multiPoint {
		pt := &Point{}
		err := pt.decodePoint(point)
		if err != nil {
			return err
		}

		points[i] = pt
	}
	mp.Points = points

	return nil
}
