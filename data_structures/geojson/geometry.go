package geojson

// Deprecated: use GeometryV2
type Geometry struct {
	Type         string     `json:"type"`
	Polygon      *Polygon   `json:"Polygon,omitempty"`
	MultiPolygon []*Polygon `json:"MultiPolygon,omitempty"`
}
