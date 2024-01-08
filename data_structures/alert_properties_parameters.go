package data_structures

import (
	"github.com/cmeyer18/weather-common/v2/data_structures/internal/processing"
)

type AlertPropertiesParameters struct {
	AWIPSIdentifier []string `json:"AWIPSidentifier"`
	WMOIdentifier   []string `json:"WMOidentifier"`
	NWSHeadline     []string `json:"NWSheadline"`
	BlockChannel    []string `json:"BLOCKCHANNEL"`
	VTEC            []string `json:"VTEC"`
	//EventEndingTime   []time.Time `json:"SAME"`
	ExpiredReferences []string `json:"expiredReferences"`
}

func ParseAlertPropertiesParameters(data interface{}) AlertPropertiesParameters {
	mappedData := data.(map[string]interface{})

	AWIPSidentifier := processing.ProcessStringArray(mappedData["AWIPSidentifier"])
	WMOIdentifier := processing.ProcessStringArray(mappedData["WMOidentifier"])
	NWSHeadline := processing.ProcessStringArray(mappedData["NWSheadline"])
	BlockChannel := processing.ProcessStringArray(mappedData["BLOCKCHANNEL"])
	VTEC := processing.ProcessStringArray(mappedData["VTEC"])
	//eventEndingTime := processing.ProcessStringArray(mappedData["eventEndingTime"])
	ExpiredReferences := processing.ProcessStringArray(mappedData["expiredReferences"])

	return AlertPropertiesParameters{
		AWIPSIdentifier: AWIPSidentifier,
		WMOIdentifier:   WMOIdentifier,
		NWSHeadline:     NWSHeadline,
		BlockChannel:    BlockChannel,
		VTEC:            VTEC,
		//EventEndingTime: nil,
		ExpiredReferences: ExpiredReferences,
	}
}
