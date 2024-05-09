package geojson_v2

import (
	"encoding/json"
	"fmt"
)

type Polygon struct {
	OuterPath  *MultiPoint
	InnerPaths []*MultiPoint
}

func (pg *Polygon) MarshalJSON() ([]byte, error) {
	return json.Marshal(pg.encodePolygon())
}

func (pg *Polygon) UnmarshalJSON(data []byte) error {
	var coordinates [][][]float64
	if err := json.Unmarshal(data, &coordinates); err != nil {
		return err
	}
	if len(coordinates) < 1 {
		return fmt.Errorf("invalid number of paths for Polygon: %v", coordinates)
	}

	return pg.decodePolygon(coordinates)
}

func (pg *Polygon) encodePolygon() [][][]float64 {
	coordinates := make([][][]float64, 1+len(pg.InnerPaths))
	coordinates[0] = pg.OuterPath.encodeMultiPoint()
	for i, innerPath := range pg.InnerPaths {
		coordinates[i+1] = innerPath.encodeMultiPoint()
	}
	return coordinates
}

func (pg *Polygon) decodePolygon(polygon [][][]float64) error {
	outerPath := &MultiPoint{}
	err := outerPath.decodeMultiPoint(polygon[0])
	if err != nil {
		return err
	}

	innerPaths := make([]*MultiPoint, len(polygon)-1)
	for i := 1; i < len(polygon); i++ {
		innerPath := &MultiPoint{}
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
