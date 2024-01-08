package geojson

import (
	"errors"
	"fmt"
)

type MultiPoint struct {
	Points []*Point `json:"points,omitempty"`
}

func parseMultiPoint(multiPoint interface{}) (*MultiPoint, error) {
	rawMultiPoint, ok := multiPoint.([]interface{})
	if !ok {
		return nil, fmt.Errorf("not a valid MultiPoint, got %v", multiPoint)
	}

	if len(rawMultiPoint) == 0 {
		return nil, errors.New("MultiPoint contains no points")
	}

	var points []*Point
	for _, rawPoint := range rawMultiPoint {
		parsedPoint, err := parsePoint(rawPoint)
		if err != nil {
			return nil, err
		}
		points = append(points, parsedPoint)
	}

	p := MultiPoint{
		Points: points,
	}

	return &p, nil
}
