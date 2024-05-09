package data_structures

import (
	"github.com/cmeyer18/weather-common/v4/data_structures/geojson"
)

type MesoscaleDiscussionV2 struct {
	ID       string            `json:"id"`
	Number   int               `json:"number"`
	Year     int               `json:"year"`
	Geometry *geojson.Geometry `json:"geometry"`
	RawText  string            `json:"rawText"`
}
