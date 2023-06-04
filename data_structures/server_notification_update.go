package data_structures

type NotificationUpdate struct {
	Alert      *Alert      `json:"alert,omitempty"`
	SPCOutlook *SPCOutlook `json:"spcoutlook,omitempty"`
}
