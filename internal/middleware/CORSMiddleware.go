package middleware

import "net/http"

func CORSMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		allowedOrigins := []string{"http://localhost:3000", "https://nekofluff.github.io"}

		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				w.Header().Set("Access-Control-Allow-Origin", "")
				w.Header().Set("Access-Control-Allow-Methods", "POST, GET, DELETE, OPTIONS")
			}
		}

		next(w, r)
	}
}
