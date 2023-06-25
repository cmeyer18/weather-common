package data_structures

import (
	"encoding/json"
	"github.com/cmeyer18/weather-common/data_structures/geojson"
	"time"
)

type Alert struct {
	ID         string            `json:"id" bson:"id"`
	Type       string            `json:"type" bson:"type"`
	Geometry   *geojson.Geometry `json:"geometry" bson:"geometry"`
	Properties AlertProperties   `json:"properties" bson:"properties"`
}

type AlertProperties struct {
	AtID          string                      `json:"@id" bson:"@id"`
	Type          string                      `json:"@type" bson:"@type"`
	ID            string                      `json:"id" bson:"id"`
	AreaDesc      string                      `json:"areaDesc" bson:"areaDesc"`
	Geocode       AlertPropertiesGeocode      `json:"geocode" bson:"geocode"`
	AffectedZones []string                    `json:"affectedZones" bson:"affectedZones"`
	References    []AlertPropertiesReferences `json:"references" bson:"references"`
	Sent          time.Time                   `json:"sent" bson:"sent"`
	Effective     time.Time                   `json:"effective" bson:"effective"`
	Onset         time.Time                   `json:"onset" bson:"onset"`
	Expires       time.Time                   `json:"expires" bson:"expires"`
	Ends          time.Time                   `json:"ends" bson:"ends"`
	Status        string                      `json:"status" bson:"status"`
	MessageType   string                      `json:"messageType" bson:"messageType"`
	Category      string                      `json:"category" bson:"category"`
	Severity      string                      `json:"severity" bson:"severity"`
	Certainty     string                      `json:"certainty" bson:"certainty"`
	Urgency       string                      `json:"urgency" bson:"urgency"`
	Event         string                      `json:"event" bson:"event"`
	Sender        string                      `json:"sender" bson:"sender"`
	SenderName    string                      `json:"senderName" bson:"senderName"`
	Headline      string                      `json:"headline" bson:"headline"`
	Description   string                      `json:"description" bson:"description"`
	Instruction   string                      `json:"instruction" bson:"instruction"`
	Response      string                      `json:"response" bson:"response"`
	Parameters    AlertPropertiesParameters   `json:"parameters" bson:"parameters"`
}

type AlertPropertiesGeocode struct {
	SAME []string `json:"SAME" bson:"SAME"`
	UGC  []string `json:"UGC" bson:"UGC"`
}

type AlertPropertiesParameters struct {
	AWIPSIdentifier   []string    `json:"AWIPSidentifier" bson:"AWIPSidentifier"`
	WMOIdentifier     []string    `json:"WMOidentifier" bson:"WMOidentifier"`
	NWSHeadline       []string    `json:"NWSheadline" bson:"NWSheadline"`
	BlockChannel      []string    `json:"BLOCKCHANNEL" bson:"BLOCKCHANNEL"`
	VTEC              []string    `json:"VTEC" bson:"VTEC"`
	EventEndingTime   []time.Time `json:"eventEndingTime" bson:"eventEndingTime"`
	ExpiredReferences []string    `json:"expiredReferences" bson:"expiredReferences"`
}

type AlertPropertiesReferences struct {
	ID         string    `json:"@id" bson:"@id"`
	Identifier string    `json:"identifier" bson:"identifier"`
	Sender     string    `json:"sender" bson:"sender"`
	Sent       time.Time `json:"sent" bson:"sent"`
}

func (a *Alert) GetListOfZones() []string {
	return a.Properties.Geocode.UGC
}

func ParseAlert(data []byte) (*Alert, error) {
	if string(data) == "null" || string(data) == `""` {
		return nil, nil
	}

	var jsonAlert struct {
		ID         string                 `json:"id"`
		Type       string                 `json:"type"`
		Geometry   map[string]interface{} `json:"geometry"`
		Properties struct {
			AtID     string `json:"@id"`
			Type     string `json:"@type"`
			ID       string `json:"id"`
			AreaDesc string `json:"areaDesc"`
			Geocode  struct {
				Same []string `json:"SAME"`
				Ugc  []string `json:"UGC"`
			} `json:"geocode"`
			AffectedZones []string `json:"affectedZones"`
			References    []struct {
				ID         string `json:"@id"`
				Identifier string `json:"identifier"`
				Sender     string `json:"sender"`
				Sent       string `json:"sent"`
			} `json:"references"`
			Sent        string `json:"sent"`
			Effective   string `json:"effective"`
			Onset       string `json:"onset"`
			Expires     string `json:"expires"`
			Ends        string `json:"ends"`
			Status      string `json:"status"`
			MessageType string `json:"messageType"`
			Category    string `json:"category"`
			Severity    string `json:"severity"`
			Certainty   string `json:"certainty"`
			Urgency     string `json:"urgency"`
			Event       string `json:"event"`
			Sender      string `json:"sender"`
			SenderName  string `json:"senderName"`
			Headline    string `json:"headline"`
			Description string `json:"description"`
			Instruction string `json:"instruction"`
			Response    string `json:"response"`
			Parameters  struct {
				AWIPSidentifier   []string    `json:"AWIPSidentifier"`
				WMOidentifier     []string    `json:"WMOidentifier"`
				NWSheadline       []string    `json:"NWSheadline"`
				Blockchannel      []string    `json:"BLOCKCHANNEL"`
				EASORG            []string    `json:"EAS-ORG"`
				Vtec              []string    `json:"VTEC"`
				EventEndingTime   []time.Time `json:"eventEndingTime"`
				ExpiredReferences []string    `json:"expiredReferences"`
			} `json:"parameters"`
		} `json:"properties"`
	}

	if err := json.Unmarshal(data, &jsonAlert); err != nil {
		return nil, err
	}

	alertGeometry, err := geojson.ParseGeometry(jsonAlert.Geometry)
	if err != nil {
		return nil, err
	}

	alertPropertiesGeocode := AlertPropertiesGeocode{
		UGC:  jsonAlert.Properties.Geocode.Ugc,
		SAME: jsonAlert.Properties.Geocode.Same,
	}

	var alertReferences []AlertPropertiesReferences
	for _, reference := range jsonAlert.Properties.References {
		sentTime, err := time.Parse(time.RFC3339, reference.Sent)
		if err != nil {
			return nil, err
		}

		alertReferences = append(alertReferences, AlertPropertiesReferences{
			ID:         reference.ID,
			Identifier: reference.Identifier,
			Sender:     reference.Sender,
			Sent:       sentTime,
		})
	}

	var sentTime time.Time
	if jsonAlert.Properties.Sent != "" {
		sentTime, err = time.Parse(time.RFC3339, jsonAlert.Properties.Sent)
		if err != nil {
			return nil, err
		}
	}

	var effectiveTime time.Time
	if jsonAlert.Properties.Effective != "" {
		effectiveTime, err = time.Parse(time.RFC3339, jsonAlert.Properties.Effective)
		if err != nil {
			return nil, err
		}
	}

	var onsetTime time.Time
	if jsonAlert.Properties.Onset != "" {
		onsetTime, err = time.Parse(time.RFC3339, jsonAlert.Properties.Onset)
		if err != nil {
			return nil, err
		}
	}

	var expiresTime time.Time
	if jsonAlert.Properties.Expires != "" {
		expiresTime, err = time.Parse(time.RFC3339, jsonAlert.Properties.Expires)
		if err != nil {
			return nil, err
		}
	}

	var endsTime time.Time
	if jsonAlert.Properties.Ends != "" {
		endsTime, err = time.Parse(time.RFC3339, jsonAlert.Properties.Ends)
		if err != nil {
			return nil, err
		}
	}

	alertPropertiesParameters := AlertPropertiesParameters{
		AWIPSIdentifier:   jsonAlert.Properties.Parameters.AWIPSidentifier,
		WMOIdentifier:     jsonAlert.Properties.Parameters.WMOidentifier,
		NWSHeadline:       jsonAlert.Properties.Parameters.NWSheadline,
		BlockChannel:      jsonAlert.Properties.Parameters.Blockchannel,
		VTEC:              jsonAlert.Properties.Parameters.Vtec,
		EventEndingTime:   jsonAlert.Properties.Parameters.EventEndingTime,
		ExpiredReferences: jsonAlert.Properties.Parameters.ExpiredReferences,
	}

	alertProperties := AlertProperties{
		AtID:          jsonAlert.Properties.AtID,
		Type:          jsonAlert.Properties.Type,
		ID:            jsonAlert.Properties.ID,
		AreaDesc:      jsonAlert.Properties.AreaDesc,
		Geocode:       alertPropertiesGeocode,
		AffectedZones: jsonAlert.Properties.AffectedZones,
		References:    alertReferences,
		Sent:          sentTime,
		Effective:     effectiveTime,
		Onset:         onsetTime,
		Expires:       expiresTime,
		Ends:          endsTime,
		Status:        jsonAlert.Properties.Status,
		MessageType:   jsonAlert.Properties.MessageType,
		Category:      jsonAlert.Properties.Category,
		Severity:      jsonAlert.Properties.Severity,
		Certainty:     jsonAlert.Properties.Certainty,
		Urgency:       jsonAlert.Properties.Urgency,
		Event:         jsonAlert.Properties.Event,
		Sender:        jsonAlert.Properties.Sender,
		SenderName:    jsonAlert.Properties.SenderName,
		Headline:      jsonAlert.Properties.Headline,
		Description:   jsonAlert.Properties.Description,
		Instruction:   jsonAlert.Properties.Instruction,
		Response:      jsonAlert.Properties.Response,
		Parameters:    alertPropertiesParameters,
	}

	a := Alert{
		ID:         jsonAlert.ID,
		Properties: alertProperties,
		Type:       jsonAlert.Type,
		Geometry:   alertGeometry,
	}

	return &a, nil
}
