package data_structures

import (
	"github.com/cmeyer18/weather-common/v4/data_structures/geojsonv2"
)

type MesoscaleDiscussionV2 struct {
	ID       string
	Number   int
	Year     int
	Geometry *geojson.GeometryV2
	RawText  string
}
