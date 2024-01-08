package data_structures

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/cmeyer18/weather-common/v2/data_structures/geojson"
)

type Alert struct {
	ID         string            `json:"id"`
	Type       string            `json:"type"`
	Geometry   *geojson.Geometry `json:"geometry"`
	Properties AlertProperties   `json:"properties"`
}

func (a *Alert) GetListOfZones() []string {
	return a.Properties.Geocode.UGC
}

func (a *Alert) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Alert) Scan(data interface{}) error {
	dataBytes, ok := data.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return a.ParseAlert(dataBytes)
}

func (a *Alert) ParseAlert(data []byte) error {
	if string(data) == "null" || string(data) == `""` {
		return nil
	}

	var jsonAlert struct {
		ID         string                 `json:"id"`
		Type       string                 `json:"type"`
		Geometry   map[string]interface{} `json:"geometry"`
		Properties map[string]interface{} `json:"properties"`
	}

	if err := json.Unmarshal(data, &jsonAlert); err != nil {
		return err
	}

	alertGeometry, err := geojson.ParseGeometry(jsonAlert.Geometry)
	if err != nil {
		return err
	}

	alertProperties, err := ParseAlertProperties(jsonAlert.Properties)
	if err != nil {
		return err
	}

	a.ID = jsonAlert.ID
	a.Type = jsonAlert.Type
	a.Geometry = alertGeometry
	a.Properties = *alertProperties

	return nil
}
