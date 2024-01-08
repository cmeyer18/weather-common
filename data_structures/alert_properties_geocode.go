package data_structures

import "github.com/cmeyer18/weather-common/v2/data_structures/internal/processing"

type AlertPropertiesGeocode struct {
	SAME []string `json:"SAME"`
	UGC  []string `json:"UGC"`
}

func ParseAlertPropertiesGeocode(data interface{}) *AlertPropertiesGeocode {
	mappedData := data.(map[string]interface{})

	SAME := processing.ProcessStringArray(mappedData["SAME"])
	UGC := processing.ProcessStringArray(mappedData["UGC"])

	return &AlertPropertiesGeocode{
		SAME: SAME,
		UGC:  UGC,
	}
}
