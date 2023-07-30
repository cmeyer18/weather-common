package geojson

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_parseEmptyGeometry(t *testing.T) {
	// Any empty map should not throw an error, but will pass an empty object back.
	var data map[string]interface{}

	geometry, err := ParseGeometry(data)
	assert.NoError(t, err)
	assert.Nil(t, geometry)
}

func Test_parseGeometryPolygon(t *testing.T) {
	var data map[string]interface{}

	data = make(map[string]interface{})
	data["type"] = "Polygon"

	// Prepare a 3 coordinate structure
	var coordinateData []interface{}
	coordinateData = append(coordinateData, makeCoordinate(0.0, 0.0))
	coordinateData = append(coordinateData, makeCoordinate(1.0, 1.0))
	coordinateData = append(coordinateData, makeCoordinate(2.0, 2.0))

	// Prepare a simple polygon
	var polygonData []interface{}
	polygonData = append(polygonData, coordinateData)

	data["coordinates"] = polygonData

	geometry, err := ParseGeometry(data)
	assert.NoError(t, err)
	assert.Equal(t, geometry.Type, "Polygon")
	assert.NotNil(t, geometry.Polygon)
}

func Test_parseGeometryMultiPolygon(t *testing.T) {
	var data map[string]interface{}

	data = make(map[string]interface{})
	data["type"] = "MultiPolygon"

	// Prepare a square structure
	var multiPoint1Data []interface{}
	multiPoint1Data = append(multiPoint1Data, makeCoordinate(0.0, 0.0))
	multiPoint1Data = append(multiPoint1Data, makeCoordinate(0.0, 10.0))
	multiPoint1Data = append(multiPoint1Data, makeCoordinate(10.0, 10.0))
	multiPoint1Data = append(multiPoint1Data, makeCoordinate(10.0, 0.0))

	// Prepare a simple polygon
	var polygon1Data []interface{}
	polygon1Data = append(polygon1Data, multiPoint1Data)

	// Prepare a square structure
	var multiPoint2Data []interface{}
	multiPoint2Data = append(multiPoint2Data, makeCoordinate(10.0, 10.0))
	multiPoint2Data = append(multiPoint2Data, makeCoordinate(10.0, 110.0))
	multiPoint2Data = append(multiPoint2Data, makeCoordinate(110.0, 110.0))
	multiPoint2Data = append(multiPoint2Data, makeCoordinate(110.0, 10.0))

	// Prepare a simple polygon
	var polygon2Data []interface{}
	polygon2Data = append(polygon2Data, multiPoint2Data)

	// Prepare a multi polygon
	var multiPolygonData []interface{}
	multiPolygonData = append(multiPolygonData, polygon1Data)
	multiPolygonData = append(multiPolygonData, polygon2Data)

	data["coordinates"] = multiPolygonData

	geometry, err := ParseGeometry(data)
	assert.NoError(t, err)
	assert.Equal(t, geometry.Type, "MultiPolygon")
	assert.Nil(t, geometry.Polygon)
	assert.NotNil(t, geometry.MultiPolygon)
	assert.Len(t, geometry.MultiPolygon, 2)
}
