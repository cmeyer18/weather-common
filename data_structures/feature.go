package data_structures

import "time"

type Feature struct {
	ID         string      `json:"id" bson:"id"`
	Type       string      `json:"type" bson:"type"`
	Geometry   interface{} `json:"geometry" bson:"geometry"`
	Properties struct {
		AtID     string `json:"@id" bson:"atid"`
		Type     string `json:"@type" bson:"type"`
		ID       string `json:"id" bson:"id"`
		AreaDesc string `json:"areaDesc" bson:"areadesc"`
		Geocode  struct {
			SAME []string `json:"SAME" bson:"same"`
			UGC  []string `json:"UGC" bson:"ugc"`
		} `json:"geocode" bson:"geocode"`
		AffectedZones []string      `json:"affectedZones" bson:"affectedzones"`
		References    []interface{} `json:"references" bson:"references"`
		Sent          string        `json:"sent" bson:"sent"`
		Effective     string        `json:"effective" bson:"effective"`
		Onset         string        `json:"onset" bson:"onset"`
		Expires       string        `json:"expires" bson:"expires"`
		Ends          string        `json:"ends" bson:"ends"`
		Status        string        `json:"status" bson:"status"`
		MessageType   string        `json:"messageType" bson:"messagetype"`
		Category      string        `json:"category" bson:"category"`
		Severity      string        `json:"severity" bson:"severity"`
		Certainty     string        `json:"certainty" bson:"certainty"`
		Urgency       string        `json:"urgency" bson:"urgency"`
		Event         string        `json:"event" bson:"event"`
		Sender        string        `json:"sender" bson:"sender"`
		SenderName    string        `json:"senderName" bson:"sendername"`
		Headline      string        `json:"headline" bson:"headline"`
		Description   string        `json:"description" bson:"description"`
		Instruction   string        `json:"instruction" bson:"instruction"`
		Response      string        `json:"response" bson:"response"`
		Parameters    struct {
			AWIPSidentifier   []string    `json:"AWIPSidentifier" bson:"awipsidentifier"`
			WMOidentifier     []string    `json:"WMOidentifier" bson:"wmoidentifier"`
			NWSheadline       []string    `json:"NWSheadline" bson:"nwsheadline"`
			BLOCKCHANNEL      []string    `json:"BLOCKCHANNEL" bson:"blockchannel"`
			VTEC              []string    `json:"VTEC" bson:"vtec"`
			EventEndingTime   []time.Time `json:"eventEndingTime" bson:"eventendingtime"`
			ExpiredReferences []string    `json:"expiredReferences" bson:"expiredreferences"`
		} `json:"parameters" bson:"parameters"`
	} `json:"properties" bson:"properties"`
}

func (f *Feature) GetListOfZones() []string {
	var zones []string
	for _, zone := range f.Properties.Geocode.UGC {
		zones = append(zones, zone)
	}
	return zones
}
