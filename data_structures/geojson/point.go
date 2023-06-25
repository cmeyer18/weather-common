package geojson

import (
	"fmt"
)

type Point struct {
	Latitude  float64 `json:"lat" bson:"lat"`
	Longitude float64 `json:"lon" bson:"lon"`
}

func parsePoint(point interface{}) (*Point, error) {
	rawPoint, ok := point.([]interface{})
	if !ok {
		return nil, fmt.Errorf("not a valid point, got %v", rawPoint)
	}

	if len(rawPoint) != 2 {
		return nil, fmt.Errorf("array not length of 2, got %v", rawPoint)
	}

	longitude := rawPoint[0].(float64)
	latitude := rawPoint[1].(float64)

	p := Point{
		Latitude:  latitude,
		Longitude: longitude,
	}
	return &p, nil
}
