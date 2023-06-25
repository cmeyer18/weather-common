package geojson

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_parsePoint(t *testing.T) {
	data := []interface{}{12.0, 13.0}

	point, err := parsePoint(data)
	assert.NoError(t, err)
	assert.Equal(t, point.Latitude, 13.0)
	assert.Equal(t, point.Longitude, 12.0)
}

func Test_parsePoint_EmptyData(t *testing.T) {
	var data []interface{}

	point, err := parsePoint(data)
	assert.Error(t, err)
	assert.Nil(t, point)
}
