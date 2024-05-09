package data_structures

import (
	"time"

	geojson2 "github.com/cmeyer18/weather-common/v4/data_structures/geojsonv2"
	"github.com/cmeyer18/weather-common/v4/generative/golang"
)

type ConvectiveOutlookV2 struct {
	ID          string                       `json:"id"`
	OutlookType golang.ConvectiveOutlookType `json:"outlookType"`
	Geometry    *geojson2.GeometryV2         `json:"geometry"`
	DN          int                          `json:"dn"`
	Valid       time.Time                    `json:"valid"`
	Expires     time.Time                    `json:"expires"`
	Issued      time.Time                    `json:"issued"`
	Label       string                       `json:"label"`
	Label2      string                       `json:"label2"`
	Stroke      string                       `json:"stroke"`
	Fill        string                       `json:"fill"`
}
