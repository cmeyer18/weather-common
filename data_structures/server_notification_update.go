package data_structures

type NotificationUpdate struct {
	Alert               *Alert                `json:"alert,omitempty"`
	ConvectiveOutlook   *ConvectiveOutlook    `json:"convectiveOutlook,omitempty"`
	MesoscaleDiscussion *MesoscaleDiscussion  `json:"mesoscaleDiscussion,omitempty"`
}

type NotificationUpdateV2 struct {
	Id string `json:"id"`
	NotificationType NotificationType `json:"notificationType"`
}

type NotificationType string

const (
	AlertType NotificationType = "alert"
	ConvectiveOutlookType NotificationType = "convectiveOutlook"
	MesoscaleDiscussionType NotificationType = "mesoscaleDiscussion"
)

