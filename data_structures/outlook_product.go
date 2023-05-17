package data_structures

import "math"

type Point struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Geometry struct {
	Type        string          `json:"type"`
	Coordinates [][][][]float64 `json:"coordinates"`
	Polygons    []PolygonShape  `json:"polygons,omitempty"`
}

type SPCFeature struct {
	Type       string   `json:"type"`
	Geometry   Geometry `json:"geometry"`
	Properties struct {
		Dn     int    `json:"DN"`
		Valid  string `json:"VALID"`
		Expire string `json:"EXPIRE"`
		Issue  string `json:"ISSUE"`
		Label  string `json:"LABEL"`
		Label2 string `json:"LABEL2"`
		Stroke string `json:"stroke"`
		Fill   string `json:"fill"`
	} `json:"properties"`
}

type OutlookProduct struct {
	Type       string       `json:"type"`
	SPCFeature []SPCFeature `json:"features"`
}

type PolygonShape struct {
	Points []Point `json:"points,omitempty"`
}

func NewPolygonShape(points []Point) *PolygonShape {
	return &PolygonShape{Points: points}
}

func (p *PolygonShape) Add(point Point) {
	p.Points = append(p.Points, point)
}

func (p *PolygonShape) IsClosed() bool {
	return len(p.Points) >= 3
}

func (p *PolygonShape) Contains(point Point) bool {
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

func (p *PolygonShape) intersectsWithRaycast(point Point, start Point, end Point) bool {
	// Always ensure that the the first point
	// has a y coordinate that is less than the second point
	if start.Longitude > end.Longitude {

		// Switch the points if otherwise.
		start, end = end, start

	}

	// Move the point's y coordinate
	// outside of the bounds of the testing region
	// so we can start drawing a ray
	for point.Longitude == start.Longitude || point.Longitude == end.Longitude {
		newLng := math.Nextafter(point.Longitude, math.Inf(1))
		point = Point{point.Latitude, newLng}
	}

	// If we are outside of the polygon, indicate so.
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
