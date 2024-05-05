package geojson

import (
	"encoding/json"
	"fmt"
)

type PolygonV2 struct {
	OuterPath  *MultiPointV2
	InnerPaths []*MultiPointV2
}

// MarshalJSON implements custom JSON marshaling for the PolygonV2 struct.
func (pg *PolygonV2) MarshalJSON() ([]byte, error) {
	return json.Marshal(pg.encodePolygon())
}

// UnmarshalJSON implements custom JSON unmarshaling for the PolygonV2 struct.
func (pg *PolygonV2) UnmarshalJSON(data []byte) error {
	var coordinates [][][]float64
	if err := json.Unmarshal(data, &coordinates); err != nil {
		return err
	}
	if len(coordinates) < 1 {
		return fmt.Errorf("invalid number of paths for Polygon: %v", coordinates)
	}

	return pg.decodePolygon(coordinates)
}

func (pg *PolygonV2) encodePolygon() [][][]float64 {
	coordinates := make([][][]float64, 1+len(pg.InnerPaths))
	coordinates[0] = pg.OuterPath.encodeMultiPoint()
	for i, innerPath := range pg.InnerPaths {
		coordinates[i+1] = innerPath.encodeMultiPoint()
	}
	return coordinates
}

func (pg *PolygonV2) decodePolygon(polygon [][][]float64) error {
	outerPath := &MultiPointV2{}
	err := outerPath.decodeMultiPoint(polygon[0])
	if err != nil {
		return err
	}

	innerPaths := make([]*MultiPointV2, len(polygon)-1)
	for i := 1; i < len(polygon); i++ {
		innerPath := &MultiPointV2{}
		err := innerPath.decodeMultiPoint(polygon[i])
		if err != nil {
			return err
		}
		innerPaths[i-1] = innerPath
	}
	pg.OuterPath = outerPath
	pg.InnerPaths = innerPaths

	return nil
}
