package data_structures

type NotificationUpdate struct {
	Alert      *Alert             `json:"alert,omitempty"`
	SPCOutlook *SPCOutlookProduct `json:"spcoutlook,omitempty"`
}
