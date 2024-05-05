package geojson

import (
	"encoding/json"
)

type MultiPointV2 struct {
	Points []*PointV2
}

// MarshalJSON implements custom JSON marshaling for the MultiPointV2 struct.
func (mp *MultiPointV2) MarshalJSON() ([]byte, error) {
	return json.Marshal(mp.encodeMultiPoint())
}

// UnmarshalJSON implements custom JSON unmarshaling for the MultiPointV2 struct.
func (mp *MultiPointV2) UnmarshalJSON(data []byte) error {
	var coordinates [][]float64
	if err := json.Unmarshal(data, &coordinates); err != nil {
		return err
	}

	return mp.decodeMultiPoint(coordinates)
}

func (mp *MultiPointV2) encodeMultiPoint() [][]float64 {
	coordinates := make([][]float64, len(mp.Points))
	for i, point := range mp.Points {
		coordinates[i] = point.encodePoint()
	}
	return coordinates
}

func (mp *MultiPointV2) decodeMultiPoint(multiPoint [][]float64) error {
	points := make([]*PointV2, len(multiPoint))
	for i, point := range multiPoint {
		pt := &PointV2{}
		err := pt.decodePoint(point)
		if err != nil {
			return err
		}

		points[i] = pt
	}
	mp.Points = points

	return nil
}
