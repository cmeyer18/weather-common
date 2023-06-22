package data_structures

import "time"

type SPCOutlook struct {
	Type string `json:"type"`
	// Deprecated: Use OutlookProduct instead
	Day            string       `json:"day,omitempty"`
	OutlookProduct string       `json:"outlookProduct,omitempty"`
	IssuedTime     time.Time    `json:"issuedTime,omitempty"`
	Features       []SPCFeature `json:"features"`
}

type SPCFeature struct {
	Type       string               `json:"type"`
	Geometry   SPCFeatureGeometry   `json:"geometry"`
	Properties SPCFeatureProperties `json:"properties"`
}

type SPCFeatureGeometry struct {
	Type        string          `json:"type"`
	Coordinates [][][][]float64 `json:"coordinates"`
	Polygons    []PolygonShape  `json:"polygons,omitempty"`
}

type SPCFeatureProperties struct {
	Dn     int    `json:"DN"`
	Valid  string `json:"VALID"`
	Expire string `json:"EXPIRE"`
	Issue  string `json:"ISSUE"`
	Label  string `json:"LABEL"`
	Label2 string `json:"LABEL2"`
	Stroke string `json:"stroke"`
	Fill   string `json:"fill"`
}
