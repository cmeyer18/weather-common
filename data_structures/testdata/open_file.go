package testdata

import (
	"bufio"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func OpenFile(t *testing.T, fileName string) []byte {
	file, err := os.Open(fileName)
	assert.NoError(t, err)

	// Get the file size
	stat, err := file.Stat()
	assert.NoError(t, err)

	bs := make([]byte, stat.Size())

	_, err = bufio.NewReader(file).Read(bs)
	assert.NoError(t, err)
	return bs
}
