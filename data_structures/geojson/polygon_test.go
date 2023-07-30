package geojson

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_parsePolygon(t *testing.T) {
	var data []interface{}
	data = append(data, []interface{}{
		makeCoordinate(1.2, 2.2),
		makeCoordinate(1.4, 5.5),
		makeCoordinate(1.23, 1.234),
	})

	point, err := parsePolygon(data)
	assert.NoError(t, err)
	assert.NotNil(t, point.OuterPath)
	assert.Nil(t, point.InnerPaths)
}

func Test_parsePolygon_EmptyData(t *testing.T) {
	var data map[string]interface{}

	point, err := parsePolygon(data)
	assert.Error(t, err)
	assert.Nil(t, point)
}

func TestPolygon_ContainsPoint(t *testing.T) {
	polygon := Polygon{
		OuterPath: &MultiPoint{
			Points: []*Point{
				{
					Latitude:  0,
					Longitude: 0,
				},
				{
					Latitude:  4,
					Longitude: 0,
				},
				{
					Latitude:  4,
					Longitude: 4,
				},
				{
					Latitude:  0,
					Longitude: 4,
				},
			},
		},
		InnerPaths: []*MultiPoint{
			{
				Points: []*Point{
					{
						Latitude:  1,
						Longitude: 1,
					},
					{
						Latitude:  1,
						Longitude: 2,
					},
					{
						Latitude:  2,
						Longitude: 2,
					},
					{
						Latitude:  2,
						Longitude: 1,
					},
				},
			},
		},
	}

	// Add item in the polygon
	assert.True(t, polygon.ContainsPoint(&Point{Latitude: 2.5, Longitude: 2.5}))

	// On edge of inner path
	assert.False(t, polygon.ContainsPoint(&Point{Latitude: 2, Longitude: 2}))

	// Edge Points
	assert.True(t, polygon.ContainsPoint(&Point{Latitude: 0, Longitude: 0}))
	assert.True(t, polygon.ContainsPoint(&Point{Latitude: 4, Longitude: 4}))

	// Edges
	assert.True(t, polygon.ContainsPoint(&Point{Latitude: 4, Longitude: 2}))
	assert.True(t, polygon.ContainsPoint(&Point{Latitude: 2, Longitude: 4}))

	// Outside
	assert.False(t, polygon.ContainsPoint(&Point{Latitude: 40, Longitude: 40}))

	// In a hole
	assert.False(t, polygon.ContainsPoint(&Point{Latitude: 1.5, Longitude: 1.5}))
}
