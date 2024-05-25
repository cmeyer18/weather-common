package data_structures

import "github.com/cmeyer18/weather-common/v4/data_structures/geojson"

// Deprecated: use MesoscaleDiscussionV2
type MesoscaleDiscussion struct {
	MDNumber     int
	Year         int
	AffectedArea *geojson.Polygon
	RawText      string
}
