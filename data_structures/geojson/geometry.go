package geojson

type Geometry struct {
	Type         string     `json:"type"`
	Polygon      *Polygon   `json:"Polygon,omitempty"`
	MultiPolygon []*Polygon `json:"MultiPolygon,omitempty"`
}
