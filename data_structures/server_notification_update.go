package data_structures

// Deprecated: use NotificationUpdateV2
type NotificationUpdate struct {
	Alert               *Alert               `json:"alert,omitempty"`
	ConvectiveOutlook   *ConvectiveOutlook   `json:"convectiveOutlook,omitempty"`
	MesoscaleDiscussion *MesoscaleDiscussion `json:"mesoscaleDiscussion,omitempty"`
}
