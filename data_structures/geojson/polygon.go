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
	OuterPath  *MultiPoint   `json:"outerPath" bson:"outerPath"`
	InnerPaths []*MultiPoint `json:"innerPaths,omitempty" bson:"innerPaths"`
}

func parsePolygon(polygon interface{}) (*Polygon, error) {
	rawMultiPoints, ok := polygon.([]interface{})
	if !ok {
		return nil, fmt.Errorf("not a valid MultiPoint, got %v", rawMultiPoints)
	}

	if len(rawMultiPoints) == 0 {
		return nil, errors.New("MultiPoint contains no Points ")
	}

	var outerPath *MultiPoint
	var innerPaths []*MultiPoint
	for i, point := range rawMultiPoints {
		parsedMultiPoint, err := parseMultiPoint(point)
		if err != nil {
			return nil, err
		}

		if i == 0 {
			outerPath = parsedMultiPoint
		} else {
			innerPaths = append(innerPaths, parsedMultiPoint)
		}
	}

	p := Polygon{
		OuterPath:  outerPath,
		InnerPaths: innerPaths,
	}

	return &p, nil
}

func NewPolygonShape(outerPath *MultiPoint, innerPaths []*MultiPoint) *Polygon {
	return &Polygon{OuterPath: outerPath, InnerPaths: innerPaths}
}

func (p *Polygon) ContainsPoint(point *Point) bool {
	return p.containedInOuterPath(point) && !p.containedInInnerPaths(point)
}

func (p *Polygon) contains(point *Point, multiPoint *MultiPoint) bool {
	start := len(multiPoint.Points) - 1
	end := 0

	contains := p.intersectsWithRaycast(point, multiPoint.Points[start], multiPoint.Points[end])

	for i := 1; i < len(multiPoint.Points); i++ {
		if p.intersectsWithRaycast(point, multiPoint.Points[i-1], multiPoint.Points[i]) {
			contains = !contains
		}
	}

	return contains
}

func (p *Polygon) containedInInnerPaths(point *Point) bool {
	for _, interPath := range p.InnerPaths {
		if p.contains(point, interPath) {
			return true
		}
	}
	return false
}

func (p *Polygon) containedInOuterPath(point *Point) bool {
	return p.contains(point, p.OuterPath)
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
