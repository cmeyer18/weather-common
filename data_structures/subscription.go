package data_structures

type Subscriber struct {
	ID  string `json:"id" bson:"id"`
	URL string `json:"url" bson:"url"`
}
