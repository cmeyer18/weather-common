package data_structures

type SPCOutlook struct {
	Type     string       `json:"type"`
	Day      string       `json:"day"`
	Features []SPCFeature `json:"features"`
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
