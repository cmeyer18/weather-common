package data_structures

type NotificationUpdate struct {
	Alert               *Alert                `json:"alert,omitempty"`
	AlertV2             *AlertV2              `json:"alertV2,omitempty"`
	ConvectiveOutlook   *ConvectiveOutlook    `json:"convectiveOutlook,omitempty"`
	ConvectiveOutlookV2 []ConvectiveOutlookV2 `json:"convectiveOutlookV2,omitempty"`
	MesoscaleDiscussion *MesoscaleDiscussion  `json:"mesoscaleDiscussion,omitempty"`
}
