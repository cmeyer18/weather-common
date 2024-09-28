package data_structures

import (
	"time"

	"github.com/cmeyer18/weather-common/v6/data_structures/geojson_v2"
)

type AlertV2 struct {
	ID            string                    `json:"id"`
	Type          string                    `json:"type"`
	Geometry      *geojson_v2.Geometry      `json:"geometry"`
	AreaDesc      string                    `json:"areaDesc"`
	Geocode       *AlertPropertiesGeocodeV2 `json:"geocode"`
	AffectedZones []string                  `json:"affectedZones"`
	References    []string                  `json:"references"`
	Sent          time.Time                 `json:"sent"`
	Effective     time.Time                 `json:"effective"`
	Onset         time.Time                 `json:"onset"`
	Expires       time.Time                 `json:"expires"`
	Ends          time.Time                 `json:"ends"`
	Status        string                    `json:"status"`
	MessageType   string                    `json:"messageType"`
	Category      string                    `json:"category"`
	Severity      string                    `json:"severity"`
	Certainty     string                    `json:"certainty"`
	Urgency       string                    `json:"urgency"`
	Event         string                    `json:"event"`
	Sender        string                    `json:"sender"`
	SenderName    string                    `json:"senderName"`
	Headline      string                    `json:"headline"`
	Description   string                    `json:"description"`
	Instruction   string                    `json:"instruction"`
	Response      string                    `json:"response"`
	Parameters    map[string]interface{}    `json:"parameters"`
}

type AlertPropertiesGeocodeV2 struct {
	SAME []string `json:"SAME"`
	UGC  []string `json:"UGC"`
}
