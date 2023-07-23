package data_structures

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const dataJson = `
	{
		"id": "https://api.weather.gov/alerts/urn:oid:2.49.0.1.840.0.dee430c9241c51b190b43d2bf35cc82b1c355f7f.001.1",
		"type": "Feature",
		"geometry": {
			"type": "Polygon",
			"coordinates": [
				[
					[
						-104.40000000000001,
						47.579999999999998
					],
					[
						-105.36,
						48.309999999999995
					],
					[
						-104.7,
						48.629999999999995
					],
					[
						-104.05,
						48.239999999999995
					],
					[
						-104.05,
						47.749999999999993
					],
					[
						-104.40000000000001,
						47.579999999999998
					]
				]
			]
		},
		"properties": {
			"@id": "https://api.weather.gov/alerts/urn:oid:2.49.0.1.840.0.dee430c9241c51b190b43d2bf35cc82b1c355f7f.001.1",
			"@type": "wx:Alert",
			"id": "urn:oid:2.49.0.1.840.0.dee430c9241c51b190b43d2bf35cc82b1c355f7f.001.1",
			"areaDesc": "Richland, MT; Roosevelt, MT; Sheridan, MT",
			"geocode": {
				"SAME": [
					"030083",
					"030085",
					"030091"
				],
				"UGC": [
					"MTC083",
					"MTC085",
					"MTC091"
				]
			},
			"affectedZones": [
				"https://api.weather.gov/zones/county/MTC083",
				"https://api.weather.gov/zones/county/MTC085",
				"https://api.weather.gov/zones/county/MTC091"
			],
			"references": [],
			"sent": "2023-07-22T23:24:00-06:00",
			"effective": "2023-07-22T23:24:00-06:00",
			"onset": "2023-07-22T23:24:00-06:00",
			"expires": "2023-07-23T00:30:00-06:00",
			"ends": "2023-07-23T00:30:00-06:00",
			"status": "Actual",
			"messageType": "Alert",
			"category": "Met",
			"severity": "Severe",
			"certainty": "Observed",
			"urgency": "Immediate",
			"event": "Severe Thunderstorm Warning",
			"sender": "w-nws.webmaster@noaa.gov",
			"senderName": "NWS Glasgow MT",
			"headline": "Severe Thunderstorm Warning issued July 22 at 11:24PM MDT until July 23 at 12:30AM MDT by NWS Glasgow MT",
			"description": "The National Weather Service in Glasgow has issued a\n\n* Severe Thunderstorm Warning for...\nCentral Roosevelt County in northeastern Montana...\nSouthwestern Sheridan County in northeastern Montana...\nNortheastern Richland County in northeastern Montana...\n\n* Until 1230 AM MDT.\n\n* At 1124 PM MDT, a severe thunderstorm was located 11 miles west of\nFroid, or 14 miles northwest of Culbertson, moving southeast at 45\nmph.\n\nHAZARD...Ping pong ball size hail and 60 mph wind gusts.\n\nSOURCE...Radar indicated.\n\nIMPACT...People and animals outdoors will be injured. Expect hail\ndamage to roofs, siding, windows, and vehicles. Expect\nwind damage to roofs, siding, and trees.\n\n* Locations impacted include...\nSidney, Culbertson, Fairview, Brockton, Medicine Lake, Bainville,\nFroid, Snowden, Fort Kipp, Homestead, Nohly, Wooley, Sprole and\nMccabe.",
			"instruction": "For your protection move to an interior room on the lowest floor of a\nbuilding.",
			"response": "Shelter",
			"parameters": {
				"AWIPSidentifier": [
					"SVRGGW"
				],
				"WMOidentifier": [
					"WUUS55 KGGW 230524"
				],
				"eventMotionDescription": [
					"2023-07-23T05:24:00-00:00...storm...314DEG...40KT...48.31,-104.73"
				],
				"windThreat": [
					"RADAR INDICATED"
				],
				"maxWindGust": [
					"60 MPH"
				],
				"hailThreat": [
					"RADAR INDICATED"
				],
				"maxHailSize": [
					"1.50"
				],
				"BLOCKCHANNEL": [
					"EAS",
					"NWEM",
					"CMAS"
				],
				"EAS-ORG": [
					"WXR"
				],
				"VTEC": [
					"/O.NEW.KGGW.SV.W.0050.230723T0524Z-230723T0630Z/"
				],
				"eventEndingTime": [
					"2023-07-23T06:30:00+00:00"
				]
			}
		}
	}
`

func Test_ParseAlert(t *testing.T) {
	data := []byte(dataJson)

	alert, err := ParseAlert(data)
	assert.NoError(t, err)
	assert.NotNil(t, alert)

	assert.Equal(t, "https://api.weather.gov/alerts/urn:oid:2.49.0.1.840.0.dee430c9241c51b190b43d2bf35cc82b1c355f7f.001.1", alert.ID)
	assert.Equal(t, "Feature", alert.Type)

	assert.NotNil(t, alert.Geometry)
	assert.Equal(t, "Polygon", alert.Geometry.Type)
	assert.NotNil(t, alert.Geometry.Polygon)
	assert.Nil(t, alert.Geometry.MultiPolygon)

	assert.NotNil(t, alert.Properties)
	assert.Equal(t, "https://api.weather.gov/alerts/urn:oid:2.49.0.1.840.0.dee430c9241c51b190b43d2bf35cc82b1c355f7f.001.1", alert.Properties.AtID)
	assert.Equal(t, "wx:Alert", alert.Properties.Type)
	assert.Equal(t, "urn:oid:2.49.0.1.840.0.dee430c9241c51b190b43d2bf35cc82b1c355f7f.001.1", alert.Properties.ID)
	assert.Equal(t, "Richland, MT; Roosevelt, MT; Sheridan, MT", alert.Properties.AreaDesc)

	assert.NotNil(t, alert.Properties.Geocode)
	assert.Equal(t, []string{"030083", "030085", "030091"}, alert.Properties.Geocode.SAME)
	assert.Equal(t, []string{"MTC083", "MTC085", "MTC091"}, alert.Properties.Geocode.UGC)

	assert.Equal(t, []string{
		"https://api.weather.gov/zones/county/MTC083",
		"https://api.weather.gov/zones/county/MTC085",
		"https://api.weather.gov/zones/county/MTC091",
	}, alert.Properties.AffectedZones)

	assert.NotNil(t, "", alert.Properties.References)

	assert.NotNil(t, alert.Properties.Sent)
	assert.NotNil(t, alert.Properties.Effective)
	assert.NotNil(t, alert.Properties.Onset)
	assert.NotNil(t, alert.Properties.Expires)
	assert.NotNil(t, alert.Properties.Ends)
	assert.Equal(t, "Actual", alert.Properties.Status)
	assert.Equal(t, "Alert", alert.Properties.MessageType)
	assert.Equal(t, "Met", alert.Properties.Category)
	assert.Equal(t, "Severe", alert.Properties.Severity)
	assert.Equal(t, "Observed", alert.Properties.Certainty)
	assert.Equal(t, "Immediate", alert.Properties.Urgency)
	assert.Equal(t, "Severe Thunderstorm Warning", alert.Properties.Event)
	assert.Equal(t, "w-nws.webmaster@noaa.gov", alert.Properties.Sender)
	assert.Equal(t, "NWS Glasgow MT", alert.Properties.SenderName)
	assert.Equal(t, "Severe Thunderstorm Warning issued July 22 at 11:24PM MDT until July 23 at 12:30AM MDT by NWS Glasgow MT", alert.Properties.Headline)
	assert.Equal(t, "The National Weather Service in Glasgow has issued a\n\n* Severe Thunderstorm Warning for...\nCentral Roosevelt County in northeastern Montana...\nSouthwestern Sheridan County in northeastern Montana...\nNortheastern Richland County in northeastern Montana...\n\n* Until 1230 AM MDT.\n\n* At 1124 PM MDT, a severe thunderstorm was located 11 miles west of\nFroid, or 14 miles northwest of Culbertson, moving southeast at 45\nmph.\n\nHAZARD...Ping pong ball size hail and 60 mph wind gusts.\n\nSOURCE...Radar indicated.\n\nIMPACT...People and animals outdoors will be injured. Expect hail\ndamage to roofs, siding, windows, and vehicles. Expect\nwind damage to roofs, siding, and trees.\n\n* Locations impacted include...\nSidney, Culbertson, Fairview, Brockton, Medicine Lake, Bainville,\nFroid, Snowden, Fort Kipp, Homestead, Nohly, Wooley, Sprole and\nMccabe.", alert.Properties.Description)
	assert.Equal(t, "For your protection move to an interior room on the lowest floor of a\nbuilding.", alert.Properties.Instruction)
	assert.Equal(t, "Shelter", alert.Properties.Response)

	// Anything under here is optional, and may not be defined.
	assert.NotNil(t, alert.Properties.Parameters)
}
