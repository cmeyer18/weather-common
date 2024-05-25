package geojson

// Deprecated: use MultiPointV2
type MultiPoint struct {
	Points []*Point `json:"points,omitempty"`
}
