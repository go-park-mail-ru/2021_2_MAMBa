package middlewares

import (
	"2021_2_MAMBa/internal/pkg/utils/log"
	"fmt"
	"net/http"
)

var allowedOrigins = map[string]struct{}{
	"http://localhost":           {},
	"http://localhost:3001":      {},
	"http://89.208.198.137":      {},
	"http://89.208.198.137:3001": {},
	"http://film4u.club":         {},
	"http://film4u.club:3001":    {},
	"http://film4u.club:3000":    {},
	"http://film4u.club:9090":    {},
	"http://film4u.club:9100":    {},
	"":                           {}, // Для дебага локально

	"https://89.208.198.137":      {},
	"https://89.208.198.137:3001": {},
	"https://film4u.club":         {},
	"https://film4u.club:3001":    {},
}

func CORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		_, isIn := allowedOrigins[origin]
		if isIn {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			log.Warn(fmt.Sprintf("unknown origin: \"%s\"", origin))
			http.Error(w, "Access denied", http.StatusForbidden)
		}
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Accept-language, Content-Type, Content-Language, Content-Encoding")
		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}
