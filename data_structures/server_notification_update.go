package data_structures

type NotificationUpdate struct {
	Alert      Feature           `json:"alert,omitempty"`
	SPCOutlook SPCOutlookProduct `json:"spcoutlook,omitempty"`
}
