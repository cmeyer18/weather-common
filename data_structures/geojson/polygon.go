package geojson

// Deprecated: use PolygonV2, not using go based location parsing
type Polygon struct {
	OuterPath  *MultiPoint   `json:"outerPath"`
	InnerPaths []*MultiPoint `json:"innerPaths,omitempty"`
}

// Deprecated:
func NewPolygonShape(outerPath *MultiPoint, innerPaths []*MultiPoint) *Polygon {
	return &Polygon{OuterPath: outerPath, InnerPaths: innerPaths}
}

// Deprecated:
func (p *Polygon) ContainsPoint(point *Point) bool {
	return p.containedInOuterPath(point) && !p.containedInInnerPaths(point)
}

func (p *Polygon) containedInInnerPaths(point *Point) bool {
	for _, interPath := range p.InnerPaths {
		if isPointInPolygon(point, interPath) {
			return true
		}
	}
	return false
}

func (p *Polygon) containedInOuterPath(point *Point) bool {
	return isPointInPolygon(point, p.OuterPath)
}

func isPointInPolygon(point *Point, multiPoint *MultiPoint) bool {
	// Check if the point is one of the multiPoint points
	for _, polyPoint := range multiPoint.Points {
		if polyPoint.Latitude == point.Latitude && polyPoint.Longitude == point.Longitude {
			return true
		}
	}

	// Check if the point is on the edge of the multiPoint
	for i := 0; i < len(multiPoint.Points); i++ {
		currPoint := multiPoint.Points[i]
		nextPoint := multiPoint.Points[(i+1)%len(multiPoint.Points)]

		if isPointOnSegment(point, currPoint, nextPoint) {
			return true
		}
	}

	// Apply the Winding number algorithm to check if the point is inside the multiPoint
	wn := 0
	for i := 0; i < len(multiPoint.Points); i++ {
		currPoint := multiPoint.Points[i]
		nextPoint := multiPoint.Points[(i+1)%len(multiPoint.Points)]

		if currPoint.Latitude <= point.Latitude {
			if nextPoint.Latitude > point.Latitude && isPointLeftOfSegment(point, currPoint, nextPoint) {
				wn++
			}
		} else {
			if nextPoint.Latitude <= point.Latitude && isPointRightOfSegment(point, currPoint, nextPoint) {
				wn--
			}
		}
	}

	return wn != 0
}

func isPointOnSegment(point, segmentStart, segmentEnd *Point) bool {
	// Check if the point is collinear with the segment
	if crossProduct(point, segmentStart, segmentEnd) != 0 {
		return false
	}

	// Check if the point is within the segment's bounding box
	if point.Latitude >= min(segmentStart.Latitude, segmentEnd.Latitude) &&
		point.Latitude <= max(segmentStart.Latitude, segmentEnd.Latitude) &&
		point.Longitude >= min(segmentStart.Longitude, segmentEnd.Longitude) &&
		point.Longitude <= max(segmentStart.Longitude, segmentEnd.Longitude) {
		return true
	}

	return false
}

func isPointLeftOfSegment(point, segmentStart, segmentEnd *Point) bool {
	return crossProduct(point, segmentStart, segmentEnd) > 0
}

func isPointRightOfSegment(point, segmentStart, segmentEnd *Point) bool {
	return crossProduct(point, segmentStart, segmentEnd) < 0
}

func crossProduct(pointA, pointB, pointC *Point) float64 {
	return (pointB.Latitude-pointA.Latitude)*(pointC.Longitude-pointA.Longitude) -
		(pointB.Longitude-pointA.Longitude)*(pointC.Latitude-pointA.Latitude)
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
