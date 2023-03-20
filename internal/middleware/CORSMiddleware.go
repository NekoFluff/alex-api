package middleware

import (
	"fmt"
	"net/http"
	"strings"
)

func CORSMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		fmt.Println("ORIGIN", origin)
		reqURL := r.Header.Get("Request URL")
		fmt.Println("REQ URL", reqURL)

		allowedOrigins := []string{"http://localhost:3000", "https://nekofluff.github.io", "https://bdo-craft-profit.herokuapp.com/"}

		for _, allowedOrigin := range allowedOrigins {
			if strings.Contains(origin, allowedOrigin) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST, GET, DELETE")
			}
		}

		next(w, r)
	}
}
