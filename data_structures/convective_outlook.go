package data_structures

import (
	"github.com/cmeyer18/weather-common/v3/data_structures/geojson"
	"github.com/cmeyer18/weather-common/v3/generative/golang"
	"time"
)

type ConvectiveOutlook struct {
	Type          string                       `json:"type"`
	PublishedTime time.Time                    `json:"publishedTime"`
	OutlookType   golang.ConvectiveOutlookType `json:"outlookType"`
	Features      []ConvectiveOutlookFeature   `json:"features"`
}

type ConvectiveOutlookFeature struct {
	Type       string                             `json:"type"`
	Geometry   *geojson.Geometry                  `json:"geometry"`
	Properties ConvectiveOutlookFeatureProperties `json:"properties"`
}

type ConvectiveOutlookFeatureProperties struct {
	DN     int       `json:"DN"`
	Valid  time.Time `json:"VALID"`
	Expire time.Time `json:"EXPIRE"`
	Issue  time.Time `json:"ISSUE"`
	Label  string    `json:"LABEL"`
	Label2 string    `json:"LABEL2"`
	Stroke string    `json:"stroke"`
	Fill   string    `json:"fill"`
}
