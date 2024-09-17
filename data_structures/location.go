package data_structures

import (
	"time"

	"github.com/cmeyer18/weather-common/v5/generative/golang"
)

type Location struct {
	LocationID                       string
	LocationType                     LocationType
	LocationReferenceID              string
	ZoneCode                         string
	CountyCode                       string
	Created                          time.Time
	Latitude                         float64
	Longitude                        float64
	LocationName                     string
	ConvectiveOutlookOptions         []golang.ConvectiveOutlookType
	AlertOptions                     []golang.AlertType
	MesoscaleDiscussionNotifications bool
}

type LocationType int8

const (
	LocationType_Unknown        LocationType = 0
	LocationType_UserLocation   LocationType = 1
	LocationType_DeviceLocaiton LocationType = 2
)
