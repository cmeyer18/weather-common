package geojson

/*
	Copyright (c) 2013 Kelly Dunn
	Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
	The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

import (
	"errors"
	"fmt"
	"math"
)

type Polygon struct {
	Points []*Point `json:"points,omitempty"`
}

func parsePolygon(polygon interface{}) (*Polygon, error) {
	rawPolygon, ok := polygon.([]interface{})
	if !ok {
		return nil, fmt.Errorf("not a valid polygon, got %v", rawPolygon)
	}

	if len(rawPolygon) == 0 {
		return nil, errors.New("polygon contains no points")
	}

	var points []*Point
	for _, point := range rawPolygon {
		parsedPoint, err := parsePoint(point)
		if err != nil {
			return nil, err
		}
		points = append(points, parsedPoint)
	}

	p := Polygon{}
	p.Points = points

	return &p, nil
}

func NewPolygonShape(points []*Point) *Polygon {
	return &Polygon{Points: points}
}

func (p *Polygon) Add(point *Point) {
	p.Points = append(p.Points, point)
}

func (p *Polygon) IsClosed() bool {
	return len(p.Points) >= 3
}

func (p *Polygon) Contains(point *Point) bool {
	if !p.IsClosed() {
		return false
	}

	start := len(p.Points) - 1
	end := 0

	contains := p.intersectsWithRaycast(point, p.Points[start], p.Points[end])

	for i := 1; i < len(p.Points); i++ {
		if p.intersectsWithRaycast(point, p.Points[i-1], p.Points[i]) {
			contains = !contains
		}
	}

	return contains
}

func (p *Polygon) intersectsWithRaycast(point *Point, start *Point, end *Point) bool {
	// Always ensure that the first point has a y coordinate that is less than the second point
	if start.Longitude > end.Longitude {
		// Switch the points if otherwise.
		start = end
		end = start
	}

	// Move the point's y coordinate outside the bounds of the testing region so we can start drawing a ray
	for point.Longitude == start.Longitude || point.Longitude == end.Longitude {
		newLng := math.Nextafter(point.Longitude, math.Inf(1))
		point = &Point{point.Latitude, newLng}
	}

	// If we are outside the polygon, indicate so.
	if point.Longitude < start.Longitude || point.Longitude > end.Longitude {
		return false
	}

	if start.Latitude > end.Latitude {
		if point.Latitude > start.Latitude {
			return false
		}
		if point.Latitude < end.Latitude {
			return true
		}

	} else {
		if point.Latitude > end.Latitude {
			return false
		}
		if point.Latitude < start.Latitude {
			return true
		}
	}

	raySlope := (point.Longitude - start.Longitude) / (point.Latitude - start.Latitude)
	diagSlope := (end.Longitude - start.Longitude) / (end.Latitude - start.Latitude)

	return raySlope >= diagSlope
}
