package data_structures

type NotificationUpdateV2 struct {
	Id               string           `json:"id"`
	NotificationType NotificationType `json:"notificationType"`
}

type NotificationType string

const (
	AlertType               NotificationType = "alert"
	ConvectiveOutlookType   NotificationType = "convectiveOutlook"
	MesoscaleDiscussionType NotificationType = "mesoscaleDiscussion"
)
