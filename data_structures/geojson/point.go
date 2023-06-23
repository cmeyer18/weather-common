package geojson

import (
	"fmt"
)

type Point struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func parsePoint(point interface{}) (*Point, error) {
	rawPoint, ok := point.([]interface{})
	if !ok {
		return nil, fmt.Errorf("not a valid point, got %v", rawPoint)
	}

	if len(rawPoint) != 2 {
		return nil, fmt.Errorf("coordinate array not length of 2, got %v", rawPoint)
	}

	longitude := rawPoint[0].(float64)
	latitude := rawPoint[1].(float64)

	p := Point{}
	p.Latitude = latitude
	p.Longitude = longitude
	return &p, nil
}
