package data_structures

import (
	"github.com/cmeyer18/weather-common/v2/data_structures/testdata"
	"github.com/cmeyer18/weather-common/v2/generative/golang"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_ParseOutlookProduct(t *testing.T) {
	bytes := testdata.OpenFile(t, "testdata/outlookproduct.txt")

	spcOutlook, err := ParseSPCOutlook(bytes, golang.Day1Categorical)
	assert.NoError(t, err)
	assert.NotNil(t, spcOutlook)

	assert.Equal(t, golang.Day1Categorical, spcOutlook.OutlookType)
	assert.Equal(t, "FeatureCollection", spcOutlook.Type)
	assert.NotNil(t, spcOutlook.Features)
	assert.Len(t, spcOutlook.Features, 4)

	// Validate TSTM
	assert.Equal(t, "Feature", spcOutlook.Features[0].Type)

	assert.NotNil(t, spcOutlook.Features[0].Properties)
	assert.Equal(t, "TSTM", spcOutlook.Features[0].Properties.Label)
	assert.Equal(t, "General Thunderstorms Risk", spcOutlook.Features[0].Properties.Label2)
	assert.Equal(t, "#C1E9C1", spcOutlook.Features[0].Properties.Fill)
	assert.Equal(t, "#55BB55", spcOutlook.Features[0].Properties.Stroke)
	assert.Equal(t, time.Date(2023, time.July, 30, 16, 25, 0, 0, time.UTC), spcOutlook.Features[0].Properties.Issue)
	assert.Equal(t, time.Date(2023, time.July, 30, 16, 30, 0, 0, time.UTC), spcOutlook.Features[0].Properties.Valid)
	assert.Equal(t, time.Date(2023, time.July, 31, 12, 0, 0, 0, time.UTC), spcOutlook.Features[0].Properties.Expire)
	assert.Equal(t, 2, spcOutlook.Features[0].Properties.DN)

	assert.NotNil(t, spcOutlook.Features[0].Geometry)
	assert.Equal(t, "MultiPolygon", spcOutlook.Features[0].Geometry.Type)
	assert.Nil(t, spcOutlook.Features[0].Geometry.Polygon)
	assert.NotNil(t, spcOutlook.Features[0].Geometry.MultiPolygon)
	assert.Len(t, spcOutlook.Features[0].Geometry.MultiPolygon, 4)

	// Validate MRGL
	assert.Equal(t, "Feature", spcOutlook.Features[1].Type)

	assert.NotNil(t, spcOutlook.Features[1].Properties)
	assert.Equal(t, "MRGL", spcOutlook.Features[1].Properties.Label)
	assert.Equal(t, "Marginal Risk", spcOutlook.Features[1].Properties.Label2)
	assert.Equal(t, "#66A366", spcOutlook.Features[1].Properties.Fill)
	assert.Equal(t, "#005500", spcOutlook.Features[1].Properties.Stroke)
	assert.Equal(t, time.Date(2023, time.July, 30, 16, 25, 0, 0, time.UTC), spcOutlook.Features[1].Properties.Issue)
	assert.Equal(t, time.Date(2023, time.July, 30, 16, 30, 0, 0, time.UTC), spcOutlook.Features[1].Properties.Valid)
	assert.Equal(t, time.Date(2023, time.July, 31, 12, 0, 0, 0, time.UTC), spcOutlook.Features[1].Properties.Expire)
	assert.Equal(t, 3, spcOutlook.Features[1].Properties.DN)

	assert.NotNil(t, spcOutlook.Features[1].Geometry)
	assert.Equal(t, "MultiPolygon", spcOutlook.Features[1].Geometry.Type)
	assert.Nil(t, spcOutlook.Features[1].Geometry.Polygon)
	assert.NotNil(t, spcOutlook.Features[1].Geometry.MultiPolygon)
	assert.Len(t, spcOutlook.Features[1].Geometry.MultiPolygon, 2)

	// Validate SLGT
	assert.Equal(t, "Feature", spcOutlook.Features[2].Type)

	assert.NotNil(t, spcOutlook.Features[2].Properties)
	assert.Equal(t, "SLGT", spcOutlook.Features[2].Properties.Label)
	assert.Equal(t, "Slight Risk", spcOutlook.Features[2].Properties.Label2)
	assert.Equal(t, "#FFE066", spcOutlook.Features[2].Properties.Fill)
	assert.Equal(t, "#DDAA00", spcOutlook.Features[2].Properties.Stroke)
	assert.Equal(t, time.Date(2023, time.July, 30, 16, 25, 0, 0, time.UTC), spcOutlook.Features[2].Properties.Issue)
	assert.Equal(t, time.Date(2023, time.July, 30, 16, 30, 0, 0, time.UTC), spcOutlook.Features[2].Properties.Valid)
	assert.Equal(t, time.Date(2023, time.July, 31, 12, 0, 0, 0, time.UTC), spcOutlook.Features[2].Properties.Expire)
	assert.Equal(t, 4, spcOutlook.Features[2].Properties.DN)

	assert.NotNil(t, spcOutlook.Features[2].Geometry)
	assert.Equal(t, "Polygon", spcOutlook.Features[2].Geometry.Type)
	assert.NotNil(t, spcOutlook.Features[2].Geometry.Polygon)
	assert.Nil(t, spcOutlook.Features[2].Geometry.MultiPolygon)

	// Validate ENH
	assert.Equal(t, "Feature", spcOutlook.Features[1].Type)

	assert.NotNil(t, spcOutlook.Features[3].Properties)
	assert.Equal(t, "ENH", spcOutlook.Features[3].Properties.Label)
	assert.Equal(t, "Enhanced Risk", spcOutlook.Features[3].Properties.Label2)
	assert.Equal(t, "#FFA366", spcOutlook.Features[3].Properties.Fill)
	assert.Equal(t, "#FF6600", spcOutlook.Features[3].Properties.Stroke)
	assert.Equal(t, time.Date(2023, time.July, 30, 16, 25, 0, 0, time.UTC), spcOutlook.Features[3].Properties.Issue)
	assert.Equal(t, time.Date(2023, time.July, 30, 16, 30, 0, 0, time.UTC), spcOutlook.Features[3].Properties.Valid)
	assert.Equal(t, time.Date(2023, time.July, 31, 12, 0, 0, 0, time.UTC), spcOutlook.Features[3].Properties.Expire)
	assert.Equal(t, 5, spcOutlook.Features[3].Properties.DN)

	assert.NotNil(t, spcOutlook.Features[3].Geometry)
	assert.Equal(t, "Polygon", spcOutlook.Features[3].Geometry.Type)
	assert.NotNil(t, spcOutlook.Features[3].Geometry.Polygon)
	assert.Nil(t, spcOutlook.Features[3].Geometry.MultiPolygon)
}
