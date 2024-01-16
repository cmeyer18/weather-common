package geojson

type MultiPoint struct {
	Points []*Point `json:"points,omitempty"`
}
