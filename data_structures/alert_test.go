package data_structures

import (
	"github.com/cmeyer18/weather-common/v2/data_structures/testdata"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ParseAlert(t *testing.T) {
	bytes := testdata.OpenFile(t, "testdata/alert.txt")

	alert, err := ParseAlert(bytes)
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
