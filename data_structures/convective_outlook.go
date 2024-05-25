package data_structures

import (
	"time"

	"github.com/cmeyer18/weather-common/v4/data_structures/geojson"
	"github.com/cmeyer18/weather-common/v4/generative/golang"
)

// Deprecated: use ConvectiveOutlookV2
type ConvectiveOutlook struct {
	Type          string                       `json:"type"`
	PublishedTime time.Time                    `json:"publishedTime"`
	OutlookType   golang.ConvectiveOutlookType `json:"outlookType"`
	Features      []ConvectiveOutlookFeature   `json:"features"`
}

// Deprecated: use ConvectiveOutlookV2
type ConvectiveOutlookFeature struct {
	Type       string                             `json:"type"`
	Geometry   *geojson.Geometry                  `json:"geometry"`
	Properties ConvectiveOutlookFeatureProperties `json:"properties"`
}

// Deprecated: use ConvectiveOutlookV2
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
