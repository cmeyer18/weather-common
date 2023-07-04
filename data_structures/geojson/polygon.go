package geojson

import (
	"errors"
	"fmt"
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

func (p *Polygon) containedInInnerPaths(point *Point) bool {
	for _, interPath := range p.InnerPaths {
		if pointInMultiPoint(point, interPath) {
			return true
		}
	}
	return false
}

func (p *Polygon) containedInOuterPath(point *Point) bool {
	return pointInMultiPoint(point, p.OuterPath)
}

func pointInMultiPoint(point *Point, multiPoint *MultiPoint) bool {
	numPoints := len(multiPoint.Points)
	if numPoints < 3 {
		return false
	}

	// Check if the point is one of the multiPoint points
	for _, p := range multiPoint.Points {
		if point.Latitude == p.Latitude && point.Longitude == p.Longitude {
			return true
		}
	}

	// Check if the point is on the edge of the multiPoint
	for i := 0; i < numPoints; i++ {
		j := (i + 1) % numPoints
		p1 := multiPoint.Points[i]
		p2 := multiPoint.Points[j]
		if pointOnLine(point, p1, p2) {
			return true
		}
	}

	// Perform raycasting algorithm to determine if the point is inside the multiPoint
	inside := false
	for i := 0; i < numPoints; i++ {
		j := (i + 1) % numPoints
		p1 := multiPoint.Points[i]
		p2 := multiPoint.Points[j]
		if intersectsRay(point, p1, p2) {
			inside = !inside
		}
	}

	return inside
}

func pointOnLine(point, lineStart, lineEnd *Point) bool {
	crossProduct := (point.Latitude-lineStart.Latitude)*(lineEnd.Longitude-lineStart.Longitude) - (point.Longitude-lineStart.Longitude)*(lineEnd.Latitude-lineStart.Latitude)
	if crossProduct != 0 {
		return false
	}

	dotProduct := (point.Longitude-lineStart.Longitude)*(lineEnd.Longitude-lineStart.Longitude) + (point.Latitude-lineStart.Latitude)*(lineEnd.Latitude-lineStart.Latitude)
	if dotProduct < 0 {
		return false
	}

	squaredLength := (lineEnd.Longitude-lineStart.Longitude)*(lineEnd.Longitude-lineStart.Longitude) + (lineEnd.Latitude-lineStart.Latitude)*(lineEnd.Latitude-lineStart.Latitude)
	if dotProduct > squaredLength {
		return false
	}

	return true
}

func intersectsRay(point, lineStart, lineEnd *Point) bool {
	if point.Longitude > max(lineStart.Longitude, lineEnd.Longitude) {
		return false
	}

	if point.Longitude < min(lineStart.Longitude, lineEnd.Longitude) {
		return false
	}

	if point.Latitude > max(lineStart.Latitude, lineEnd.Latitude) {
		return false
	}

	if lineStart.Longitude == lineEnd.Longitude {
		if point.Longitude == lineStart.Longitude {
			return true
		} else {
			return false
		}
	}

	xIntersection := lineStart.Longitude + (point.Latitude-lineStart.Latitude)*(lineEnd.Longitude-lineStart.Longitude)/(lineEnd.Latitude-lineStart.Latitude)
	if point.Longitude < xIntersection {
		return true
	} else {
		return false
	}
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
