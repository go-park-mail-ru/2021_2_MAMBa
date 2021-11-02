package cast

import (
	"encoding/binary"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestToString(t *testing.T) {
	targetString := "привет"
	bytes := []uint8(targetString)
	actual := ToString(bytes)
	assert.Equal(t, targetString, actual)
}

func TestToUint64(t *testing.T) {
	targetUint := uint64(420)
	bytes := make([]uint8, 8)
	binary.BigEndian.PutUint64(bytes, targetUint)
	actual := ToUint64(bytes)
	assert.Equal(t, targetUint, actual)
}

func TestToUint32(t *testing.T) {
	targetUint := uint32(420)
	bytes := make([]uint8, 4)
	binary.BigEndian.PutUint32(bytes, targetUint)
	actual := ToUint32(bytes)
	assert.Equal(t, targetUint, actual)
}

func TestToFloat64(t *testing.T) {
	targetFloat := 10.0
	bytes := make([]uint8, 8)
	binary.BigEndian.PutUint64(bytes, math.Float64bits(targetFloat))
	actual := ToFloat64(bytes)
	assert.Equal(t, targetFloat, actual)
}
