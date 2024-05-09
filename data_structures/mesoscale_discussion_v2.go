package data_structures

import (
	"github.com/cmeyer18/weather-common/v4/data_structures/geojsonv2"
)

type MesoscaleDiscussionV2 struct {
	ID       string              `json:"id"`
	Number   int                 `json:"number"`
	Year     int                 `json:"year"`
	Geometry *geojson.GeometryV2 `json:"geometry"`
	RawText  string              `json:"rawText"`
}
