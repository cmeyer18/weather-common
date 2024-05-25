package data_structures

import (
	"time"

	"github.com/cmeyer18/weather-common/v5/data_structures/geojson_v2"
)

type MesoscaleDiscussionV2 struct {
	ID                         string               `json:"id"`
	Number                     int                  `json:"number"`
	Year                       int                  `json:"year"`
	Geometry                   *geojson_v2.Geometry `json:"geometry"`
	RawText                    string               `json:"rawText"`
	ProbabilityOfWatchIssuance *int                 `json:"probabilityOfWatchIssuance"`
	Effective                  *time.Time           `json:"effective"`
	Expires                    *time.Time           `json:"expires"`
}
