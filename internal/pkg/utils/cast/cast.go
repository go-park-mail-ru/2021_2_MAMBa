package cast

import (
	"encoding/binary"
	"github.com/jackc/pgtype"
	"math"
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

func TimestampToString(src []byte) (string, error) {
	timeBuffer := pgtype.Timestamp{}
	err := timeBuffer.DecodeBinary(nil, src)
	timeString := timeBuffer.Time.Format("02.01.2006")
	if timeString == "01.01.0001" {
		return "", err
	}
	return timeString, err
}

func DateToString(src []byte) (string, error) {
	timeBuffer := pgtype.Date{}
	err := timeBuffer.DecodeBinary(nil, src)
	timeString := timeBuffer.Time.Format("02.01.2006")
	if timeString == "01.01.0001" {
		return "", err
	}
	return timeString, err
}
