package data_structures

import (
	"encoding/json"
	"errors"
	"github.com/cmeyer18/weather-common/v2/data_structures/geojson"
	"github.com/cmeyer18/weather-common/v2/generative/golang"
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

func ParseConvectiveOutlook(data []byte, outlookType golang.ConvectiveOutlookType) (*ConvectiveOutlook, error) {
	var jsonConvectiveOutlook struct {
		Type     string        `json:"type"`
		Features []interface{} `json:"features"`
	}

	err := json.Unmarshal(data, &jsonConvectiveOutlook)
	if err != nil {
		return nil, err
	}

	// If no features, something is not right
	if len(jsonConvectiveOutlook.Features) == 0 {
		return nil, errors.New(`convective outlooks should have at least one feature`)
	}

	features, err := ParseConvectiveOutlookFeatures(jsonConvectiveOutlook.Features)
	if err != nil {
		return nil, err
	}

	co := &ConvectiveOutlook{
		OutlookType: outlookType,
		Type:        jsonConvectiveOutlook.Type,
		Features:    features,
	}

	return co, nil
}

func ParseConvectiveOutlookFeatures(features []interface{}) ([]ConvectiveOutlookFeature, error) {
	if len(features) == 0 {
		return nil, errors.New("no features found in outlook")
	}

	var parsedFeatures []ConvectiveOutlookFeature
	for _, feature := range features {
		featureMap := feature.(map[string]interface{})

		parsedGeometry, err := geojson.ParseGeometry(featureMap["geometry"].(map[string]interface{}))
		if err != nil {
			return nil, err
		}

		propertiesMap := featureMap["properties"].(map[string]interface{})

		validTime, err := time.Parse("200601021504", propertiesMap["VALID"].(string))
		if err != nil {
			err = nil
			validTime, err = time.Parse(time.RFC3339, propertiesMap["VALID"].(string))
			if err != nil {
				return nil, err
			}

		}

		expireTime, err := time.Parse("200601021504", propertiesMap["EXPIRE"].(string))
		if err != nil {
			err = nil
			expireTime, err = time.Parse(time.RFC3339, propertiesMap["VALID"].(string))
			if err != nil {
				return nil, err
			}
		}

		issuedTime, err := time.Parse("200601021504", propertiesMap["ISSUE"].(string))
		if err != nil {
			err = nil
			issuedTime, err = time.Parse(time.RFC3339, propertiesMap["VALID"].(string))
			if err != nil {
				return nil, err
			}
		}

		dn := int(propertiesMap["DN"].(float64))

		parsedProperties := ConvectiveOutlookFeatureProperties{
			DN:     dn,
			Valid:  validTime,
			Expire: expireTime,
			Issue:  issuedTime,
			Label:  propertiesMap["LABEL"].(string),
			Label2: propertiesMap["LABEL2"].(string),
			Stroke: propertiesMap["stroke"].(string),
			Fill:   propertiesMap["fill"].(string),
		}

		parsedFeature := ConvectiveOutlookFeature{
			Type:       featureMap["type"].(string),
			Geometry:   parsedGeometry,
			Properties: parsedProperties,
		}

		parsedFeatures = append(parsedFeatures, parsedFeature)
	}

	return parsedFeatures, nil
}
