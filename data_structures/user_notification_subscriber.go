package data_structures

import "time"

type UserNotificationSubscriber struct {
	ID               string    `bson:"id"`
	UserID           string    `bson:"userid"`
	ZoneCode         string    `bson:"zonecode"`
	CountyCode       string    `bson:"countycode"`
	CreationTime     time.Time `bson:"creationtime"`
	Lat              string    `bson:"lat"`
	Lng              string    `bson:"lng"`
	FormattedAddress string    `bson:"formattedaddress"`
}
