package geojson

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_parseEmptyMultipoint(t *testing.T) {
	// Prepare a 0 coordinate structure, this is invalid. Must have at least 3 points
	var data []interface{}

	multiPoint, err := parseMultiPoint(data)
	assert.Error(t, err)
	assert.Nil(t, multiPoint)
}

func Test_parseMultipoint(t *testing.T) {
	// Prepare a 3 coordinate structure
	var data []interface{}
	data = append(data, makeCoordinate(0.0, 0.0))
	data = append(data, makeCoordinate(1.0, 1.0))
	data = append(data, makeCoordinate(2.0, 1.0))

	multiPoint, err := parseMultiPoint(data)
	assert.NoError(t, err)
	assert.Len(t, multiPoint.Points, 3)
}
