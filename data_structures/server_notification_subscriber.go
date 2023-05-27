package data_structures

type ServerNotificationSubscriber struct {
	ID  string `json:"id" bson:"id"`
	URL string `json:"url" bson:"url"`
}
