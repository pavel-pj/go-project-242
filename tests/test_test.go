package tests

import (
	pathsize "hexlet-path-size/path_size"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTeat(t *testing.T) {
	// Пример использования твоего пакета
	result, err := pathsize.GetSize(".", false, false, false)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)

	b := 10
	assert.Equal(t, b, 10)
}
