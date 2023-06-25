package data_structures

import (
	"github.com/cmeyer18/weather-common/generative/golang"
	"time"
)

type UserNotification struct {
	ID               string                  `json:"id" bson:"id"`
	UserID           string                  `json:"userid" bson:"userid"`
	ZoneCode         string                  `json:"zonecode" bson:"zonecode"`
	CountyCode       string                  `json:"countycode" bson:"countycode"`
	CreationTime     time.Time               `json:"creationtime" bson:"creationtime"`
	Lat              float64                 `json:"lat" bson:"lat"`
	Lng              float64                 `json:"lng" bson:"lng"`
	FormattedAddress string                  `json:"formattedaddress" bson:"formattedaddress"`
	APNKey           string                  `json:"apnKey" bson:"apnKey"`
	LocationName     string                  `json:"locationName" bson:"locationName"`
	SPCOptions       []golang.SPCOutlookType `json:"spcOptions" bson:"spcOptions"`
	AlertOptions     []golang.AlertType      `json:"alertOptions" bson:"alertOptions"`
}
