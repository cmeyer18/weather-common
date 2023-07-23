package geojson

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_parseMultipoint(t *testing.T) {
	data := []byte("[-93.787988,32.392335]")

	multiPoint, err := parseMultiPoint(data)
	assert.NoError(t, err)
	assert.NotNil(t, len(multiPoint.Points), 1)
}
