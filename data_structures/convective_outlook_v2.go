package data_structures

import (
	"time"

	"github.com/cmeyer18/weather-common/v5/data_structures/geojson_v2"
	"github.com/cmeyer18/weather-common/v5/generative/golang"
)

type ConvectiveOutlookV2 struct {
	ID          string                       `json:"id"`
	OutlookType golang.ConvectiveOutlookType `json:"outlookType"`
	Geometry    *geojson_v2.Geometry         `json:"geometry"`
	DN          int                          `json:"dn"`
	Valid       time.Time                    `json:"valid"`
	Expires     time.Time                    `json:"expires"`
	Issued      time.Time                    `json:"issued"`
	Label       string                       `json:"label"`
	Label2      string                       `json:"label2"`
	Stroke      string                       `json:"stroke"`
	Fill        string                       `json:"fill"`
}
