package strutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandInt64(t *testing.T) {
	assert.True(t, RandUint64(5) < 100000)
	assert.True(t, RandUint64(5) < 100000)
	assert.True(t, RandUint64(5) < 100000)
	assert.True(t, RandUint64(5) < 100000)
	assert.True(t, RandUint64(5) < 100000)
	assert.True(t, RandUint64(5) < 100000)
	assert.True(t, RandUint64(5) < 100000)
}

func TestRandNumbers(t *testing.T) {
	assert.Equal(t, 4, len(RandNumbers(4)))
	assert.Equal(t, 4, len(RandNumbers(4)))
	assert.Equal(t, 4, len(RandNumbers(4)))
	assert.Equal(t, 6, len(RandNumbers(6)))
	assert.Equal(t, 10, len(RandNumbers(10)))
}

func TestRandString(t *testing.T) {
	assert.Equal(t, 4, len(RandString(4)))
	assert.Equal(t, 4, len(RandString(4)))
	assert.Equal(t, 4, len(RandString(4)))
	assert.Equal(t, 6, len(RandString(6)))
	assert.Equal(t, 10, len(RandString(10)))
}

func TestRandPrintable(t *testing.T) {
	assert.Equal(t, 4, len(RandPrintable(4)))
	assert.Equal(t, 4, len(RandPrintable(4)))
	assert.Equal(t, 4, len(RandPrintable(4)))
	assert.Equal(t, 6, len(RandPrintable(6)))
	assert.Equal(t, 10, len(RandPrintable(10)))
}

func TestRandLowCased(t *testing.T) {
	assert.Equal(t, 4, len(RandLowCased(4)))
	assert.Equal(t, 4, len(RandLowCased(4)))
	assert.Equal(t, 4, len(RandLowCased(4)))
	assert.Equal(t, 6, len(RandLowCased(6)))
	assert.Equal(t, 10, len(RandLowCased(10)))
}

func TestRandHash(t *testing.T) {
	assert.Equal(t, 4, len(RandHash(4)))
	assert.Equal(t, 4, len(RandHash(4)))
	assert.Equal(t, 4, len(RandHash(4)))
	assert.Equal(t, 6, len(RandHash(6)))
	assert.Equal(t, 10, len(RandHash(10)))
}
