package data_structures

import (
	"github.com/cmeyer18/weather-common/v3/generative/golang"
	"time"
)

type UserNotification struct {
	NotificationId           string                         `json:"notificationId"`
	UserID                   string                         `json:"userid"`
	ZoneCode                 string                         `json:"zonecode"`
	CountyCode               string                         `json:"countycode"`
	CreationTime             time.Time                      `json:"creationtime"`
	Lat                      float64                        `json:"lat"`
	Lng                      float64                        `json:"lng"`
	FormattedAddress         string                         `json:"formattedaddress"`
	APNKey                   string                         `json:"apnKey"`
	LocationName             string                         `json:"locationName"`
	ConvectiveOutlookOptions []golang.ConvectiveOutlookType `json:"convectiveOutlookOptions"`
	AlertOptions             []golang.AlertType             `json:"alertOptions"`
}
