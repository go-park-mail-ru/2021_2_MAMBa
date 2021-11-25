package cast

import (
	authRPC "2021_2_MAMBa/internal/pkg/sessions/delivery/grpc"
	"encoding/binary"
	"encoding/json"
	"github.com/jackc/pgtype"
	"math"
	"net/http"
)

func CookieToRq(request *http.Request, id uint64) authRPC.Request {
	cookie, _ := request.Cookie("session-name")
	if cookie == nil {
		return authRPC.Request{ID: 0}
	}
	return authRPC.Request{
		Name:     cookie.Name,
		Value:    cookie.Value,
		Path:     cookie.Path,
		Domain:   cookie.Domain,
		MaxAge:   int64(cookie.MaxAge),
		Secure:   cookie.Secure,
		HttpOnly: cookie.HttpOnly,
		SameSite: int64(cookie.SameSite),
		Raw:      cookie.Raw,
		Unparsed: cookie.Unparsed,
		ID:       id,
	}
}

type JsonErr struct {
	Error string `json:"error"`
}

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

func DateToStringUnderscore(src []byte) (string, error) {
	timeBuffer := pgtype.Date{}
	err := timeBuffer.DecodeBinary(nil, src)
	timeString := timeBuffer.Time.Format("2006-01-02")
	if timeString == "0001-01-01" {
		return "", err
	}
	return timeString, err
}

func StringToJson(src string) []byte {
	res, _ := json.Marshal(src)
	return res
}

func ErrorToJson(src string) []byte {
	res, _ := json.Marshal(JsonErr{Error: src})
	return res
}
