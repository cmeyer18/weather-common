package data_structures

import (
	"github.com/cmeyer18/weather-common/v2/data_structures/internal/processing"
	"time"
)

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

func ParseAlertProperties(data map[string]interface{}) (*AlertProperties, error) {
	alertProperties := &AlertProperties{}

	alertProperties.Geocode = ParseAlertPropertiesGeocode(data["geocode"])

	alertProperties.AffectedZones = processing.ProcessStringArray(data["affectedZones"])

	rawReferences := data["references"].([]interface{})
	for _, reference := range rawReferences {
		processedReference, err := ParseAlertPropertiesReferences(reference)
		if err != nil {
			return nil, err
		}

		alertProperties.References = append(alertProperties.References, processedReference)
	}

	sentTimeRaw := data["sent"]
	var err error
	if sentTimeRaw != nil {
		alertProperties.Sent, err = time.Parse(time.RFC3339, sentTimeRaw.(string))
		if err != nil {
			return nil, err
		}
	}

	effectiveTimeRaw := data["effective"]
	if effectiveTimeRaw != nil {
		alertProperties.Effective, err = time.Parse(time.RFC3339, effectiveTimeRaw.(string))
		if err != nil {
			return nil, err
		}
	}

	onsetTimeRaw := data["onset"]
	if onsetTimeRaw != nil {
		alertProperties.Onset, err = time.Parse(time.RFC3339, onsetTimeRaw.(string))
		if err != nil {
			return nil, err
		}
	}

	expiresTimeRaw := data["expires"]
	if expiresTimeRaw != nil {
		alertProperties.Expires, err = time.Parse(time.RFC3339, expiresTimeRaw.(string))
		if err != nil {
			return nil, err
		}
	}

	endsTimeRaw := data["ends"]
	if endsTimeRaw != nil {
		alertProperties.Ends, err = time.Parse(time.RFC3339, endsTimeRaw.(string))
		if err != nil {
			return nil, err
		}
	}

	alertProperties.Parameters = ParseAlertPropertiesParameters(data["parameters"])

	if data["@id"] != nil {
		alertProperties.AtID = data["@id"].(string)
	}

	if data["@type"] != nil {
		alertProperties.Type = data["@type"].(string)
	}

	if data["id"] != nil {
		alertProperties.ID = data["id"].(string)
	}

	if data["areaDesc"] != nil {
		alertProperties.AreaDesc = data["areaDesc"].(string)
	}

	if data["status"] != nil {
		alertProperties.Status = data["status"].(string)
	}

	if data["messageType"] != nil {
		alertProperties.MessageType = data["messageType"].(string)
	}

	if data["category"] != nil {
		alertProperties.Category = data["category"].(string)
	}

	if data["severity"] != nil {
		alertProperties.Severity = data["severity"].(string)
	}

	if data["certainty"] != nil {
		alertProperties.Certainty = data["certainty"].(string)
	}

	if data["urgency"] != nil {
		alertProperties.Urgency = data["urgency"].(string)
	}

	if data["event"] != nil {
		alertProperties.Event = data["event"].(string)
	}

	if data["sender"] != nil {
		alertProperties.Sender = data["sender"].(string)
	}

	if data["senderName"] != nil {
		alertProperties.SenderName = data["senderName"].(string)
	}

	if data["headline"] != nil {
		alertProperties.Headline = data["headline"].(string)
	}

	if data["description"] != nil {
		alertProperties.Description = data["description"].(string)
	}

	if data["instruction"] != nil {
		alertProperties.Instruction = data["instruction"].(string)
	}

	if data["response"] != nil {
		alertProperties.Response = data["response"].(string)
	}

	return alertProperties, nil
}
