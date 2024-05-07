package data_structures

import "github.com/cmeyer18/weather-common/v4/data_structures/geojson"

type MesoscaleDiscussion struct {
	ID           string
	MDNumber     int
	Year         int
	AffectedArea *geojson.Polygon
	RawText      string
}
