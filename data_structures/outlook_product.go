package data_structures

import (
	"encoding/json"
	"errors"
	"github.com/cmeyer18/weather-common/data_structures/geojson"
	"github.com/cmeyer18/weather-common/generative/golang"
	"time"
)

type SPCOutlook struct {
	Type        string                `json:"type" bson:"type"`
	OutlookType golang.SPCOutlookType `json:"outlookType" bson:"outlookType"`
	Features    []SPCFeature          `json:"features" bson:"features"`
}

type SPCFeature struct {
	Type       string               `json:"type" bson:"type"`
	Geometry   *geojson.Geometry    `json:"geometry" bson:"features"`
	Properties SPCFeatureProperties `json:"properties" bson:"properties"`
}

type SPCFeatureProperties struct {
	DN     int       `json:"DN" bson:"DN"`
	Valid  time.Time `json:"VALID" bson:"VALID"`
	Expire time.Time `json:"EXPIRE" bson:"EXPIRE"`
	Issue  time.Time `json:"ISSUE" bson:"ISSUE"`
	Label  string    `json:"LABEL" bson:"LABEL"`
	Label2 string    `json:"LABEL2" bson:"LABEL2"`
	Stroke string    `json:"stroke" bson:"stroke"`
	Fill   string    `json:"fill" bson:"fill"`
}

func ParseSPCOutlook(data []byte, outlookType golang.SPCOutlookType) (*SPCOutlook, error) {
	var jsonSPCOutlook struct {
		Type     string `json:"type"`
		Features []struct {
			Type       string                 `json:"type"`
			Geometry   map[string]interface{} `json:"geometry"`
			Properties struct {
				Dn     int    `json:"DN"`
				Valid  string `json:"VALID"`
				Expire string `json:"EXPIRE"`
				Issue  string `json:"ISSUE"`
				Label  string `json:"LABEL"`
				Label2 string `json:"LABEL2"`
				Stroke string `json:"stroke"`
				Fill   string `json:"fill"`
			} `json:"properties"`
		} `json:"features"`
	}

	err := json.Unmarshal(data, &jsonSPCOutlook)
	if err != nil {
		return nil, err
	}

	if len(jsonSPCOutlook.Features) == 0 {
		return nil, errors.New("no features found in outlook")
	}

	var parsedFeatures []SPCFeature
	for _, feature := range jsonSPCOutlook.Features {
		parsedGeometry, err := geojson.ParseGeometry(feature.Geometry)
		if err != nil {
			return nil, err
		}

		validTime, err := time.Parse("200601021504", feature.Properties.Issue)
		if err != nil {
			return nil, err
		}

		expireTime, err := time.Parse("200601021504", feature.Properties.Issue)
		if err != nil {
			return nil, err
		}

		issuedTime, err := time.Parse("200601021504", feature.Properties.Issue)
		if err != nil {
			return nil, err
		}

		parsedProperties := SPCFeatureProperties{
			DN:     feature.Properties.Dn,
			Valid:  validTime,
			Expire: expireTime,
			Issue:  issuedTime,
			Label:  feature.Properties.Label,
			Label2: feature.Properties.Label2,
			Stroke: feature.Properties.Stroke,
			Fill:   feature.Properties.Fill,
		}

		parsedFeature := SPCFeature{
			Type:       feature.Type,
			Geometry:   parsedGeometry,
			Properties: parsedProperties,
		}

		parsedFeatures = append(parsedFeatures, parsedFeature)
	}

	outlook := SPCOutlook{
		Type:        jsonSPCOutlook.Type,
		OutlookType: outlookType,
		Features:    parsedFeatures,
	}

	return &outlook, nil
}
