package data_structures

import (
	"time"

	"github.com/cmeyer18/weather-common/v4/data_structures/geojson"
)

// Deprecated: use AlertV2
type Alert struct {
	ID         string            `json:"id"`
	Type       string            `json:"type"`
	Geometry   *geojson.Geometry `json:"geometry"`
	Properties AlertProperties   `json:"properties"`
}

// Deprecated: use AlertV2
type AlertProperties struct {
	AtID          string                       `json:"@id"`
	Type          string                       `json:"@type"`
	ID            string                       `json:"id"`
	AreaDesc      string                       `json:"areaDesc"`
	Geocode       *AlertPropertiesGeocode      `json:"geocode"`
	AffectedZones []string                     `json:"affectedZones"`
	References    []*AlertPropertiesReferences `json:"references"`
	Sent          time.Time                    `json:"sent"`
	Effective     time.Time                    `json:"effective"`
	Onset         time.Time                    `json:"onset"`
	Expires       time.Time                    `json:"expires"`
	Ends          time.Time                    `json:"ends"`
	Status        string                       `json:"status"`
	MessageType   string                       `json:"messageType"`
	Category      string                       `json:"category"`
	Severity      string                       `json:"severity"`
	Certainty     string                       `json:"certainty"`
	Urgency       string                       `json:"urgency"`
	Event         string                       `json:"event"`
	Sender        string                       `json:"sender"`
	SenderName    string                       `json:"senderName"`
	Headline      string                       `json:"headline"`
	Description   string                       `json:"description"`
	Instruction   string                       `json:"instruction"`
	Response      string                       `json:"response"`
	Parameters    AlertPropertiesParameters    `json:"parameters"`
}

// Deprecated: use AlertV2
type AlertPropertiesGeocode struct {
	SAME []string `json:"SAME"`
	UGC  []string `json:"UGC"`
}

// Deprecated: use AlertV2
type AlertPropertiesParameters struct {
	AWIPSIdentifier   []string `json:"AWIPSidentifier"`
	WMOIdentifier     []string `json:"WMOidentifier"`
	NWSHeadline       []string `json:"NWSheadline"`
	BlockChannel      []string `json:"BLOCKCHANNEL"`
	VTEC              []string `json:"VTEC"`
	ExpiredReferences []string `json:"expiredReferences"`
}

// Deprecated: use AlertV2
type AlertPropertiesReferences struct {
	AtID       string    `json:"@id"`
	Identifier string    `json:"identifier"`
	Sender     string    `json:"sender"`
	Sent       time.Time `json:"sent"`
}

func (a *Alert) GetListOfZones() []string {
	return a.Properties.Geocode.UGC
}
