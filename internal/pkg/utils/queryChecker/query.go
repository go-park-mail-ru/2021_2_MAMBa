package queryChecker

import (
	"net/http"
	"strconv"
)

func CheckIsIn(w http.ResponseWriter, r *http.Request, queryString string, defaultValue int, returnError error) (int, error) {
	valueString, isIn := r.URL.Query()[queryString]
	if isIn {
		value, err := strconv.Atoi(valueString[0])
		if err != nil || value < 0 {
			http.Error(w, returnError.Error(), http.StatusBadRequest)
			return defaultValue, returnError
		}
		return value, nil
	}
	return defaultValue, nil
}

func CheckIsIn64(w http.ResponseWriter, r *http.Request, queryString string, defaultValue uint64, returnError error) (uint64, error) {
	valueString, isIn := r.URL.Query()[queryString]
	if isIn {
		value, err := strconv.ParseUint(valueString[0], 10, 64)
		if err != nil || value < 0 {
			http.Error(w, returnError.Error(), http.StatusBadRequest)
			return defaultValue, returnError
		}
		return value, nil
	}
	return defaultValue, nil
}

func CheckIsInFloat64(w http.ResponseWriter, r *http.Request, queryString string, defaultValue float64, returnError error) (float64, error) {
	valueString, isIn := r.URL.Query()[queryString]
	if isIn {
		value, err := strconv.ParseFloat(valueString[0], 64)
		if err != nil || value < 0 {
			http.Error(w, returnError.Error(), http.StatusBadRequest)
			return defaultValue, returnError
		}
		return value, nil
	}
	return defaultValue, nil
}
