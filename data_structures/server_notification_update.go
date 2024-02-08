package data_structures

type NotificationUpdate struct {
	Alert               *Alert               `json:"alert,omitempty"`
	ConvectiveOutlook   *ConvectiveOutlook   `json:"convectiveOutlook,omitempty"`
	MesoscaleDiscussion *MesoscaleDiscussion `json:"mesoscaleDiscussion,omitempty"`
}
