package data_structures

import "time"

type AlertPropertiesReferences struct {
	AtID       string    `json:"@id"`
	Identifier string    `json:"identifier"`
	Sender     string    `json:"sender"`
	Sent       time.Time `json:"sent"`
}

func ParseAlertPropertiesReferences(data interface{}) (*AlertPropertiesReferences, error) {
	mappedData := data.(map[string]interface{})

	sentTimeRaw := mappedData["sent"].(string)
	var sentTime time.Time
	var err error
	if sentTimeRaw != "" {
		sentTime, err = time.Parse(time.RFC3339, sentTimeRaw)
		if err != nil {
			return nil, err
		}
	}

	return &AlertPropertiesReferences{
		AtID:       mappedData["@id"].(string),
		Identifier: mappedData["identifier"].(string),
		Sender:     mappedData["sender"].(string),
		Sent:       sentTime,
	}, nil
}
