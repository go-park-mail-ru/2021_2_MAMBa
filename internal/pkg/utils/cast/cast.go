package cast

import (
	"encoding/binary"
	"github.com/jackc/pgtype"
	"math"
	"time"
)

func ToString(src []byte) string {
	return string(src)
}

func ToUint64(src []byte) uint64 {
	return binary.BigEndian.Uint64(src)
}

func ToFloat64(src []byte) float64 {
	return math.Float64frombits(binary.BigEndian.Uint64(src))
}

func ToUint32(src []byte) uint32 {
	return binary.BigEndian.Uint32(src)
}

func ToTime(src []byte) (time.Time, error) {
	timeBuffer := pgtype.Timestamp{}
	err := timeBuffer.DecodeBinary(nil, src)
	return timeBuffer.Time, err
}

func ToDate(src []byte) (time.Time, error) {
	timeBuffer := pgtype.Date{}
	err := timeBuffer.DecodeBinary(nil, src)
	return timeBuffer.Time, err
}
